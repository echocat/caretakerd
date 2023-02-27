package main

import (
	"bytes"
	"github.com/alecthomas/kingpin/v2"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd/errors"
	"github.com/russross/blackfriday/v2"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	lineBreakCorrectPattern             = regexp.MustCompile(`(?m)\r\n`)
	headerPrefixPattern                 = regexp.MustCompile(`(?m)^([* 0-9.]*)#`)
	excerptFromCommentExtractionPattern = regexp.MustCompile(`(?s)^(.*?(?:\.\s|$))`)
	removeHTMLTags                      = regexp.MustCompile(`(?sm)<[^>]+>`)
	refPropertyPattern                  = regexp.MustCompile(`{@ref +([^}\s]+)\s*([^}]*)}`)
	titlePropertyPattern                = regexp.MustCompile(`(?m)^#\s*@title\s+(.*)\s*(:?\r\n|\n)`)
	windowsEnvarPattern                 = regexp.MustCompile(`(?m)%([a-zA-Z0-9_.]+)%`)
	otherEnvarPattern                   = regexp.MustCompile(`(?m)$([a-zA-Z0-9_.]+)`)
)

// Describable represents an object that describes itself and has an ID.
type Describable interface {
	ID() IDType
	Description() string
}

// Renderer is an object that renders a manual to HTML.
type Renderer struct {
	Template                    *Template
	IDTemplate                  *Template
	PointerTemplate             *Template
	ArrayTemplate               *Template
	MapTemplate                 *Template
	DefinitionStructureTemplate *Template
	Platform                    string

	Functions         template.FuncMap
	Project           Project
	PickedDefinitions *PickedDefinitions
	Apps              map[app.ExecutableType]*kingpin.Application

	Name        string
	Version     string
	Description string
	URL         string
}

// Execute executes the rendering.
func (instance *Renderer) Execute() (template.HTML, error) {
	return instance.Template.Execute(instance)
}

// Template represents an HTML template with its name.
type Template struct {
	tmpl *template.Template
	name string
}

// Execute executes the rendering of this template with the given data.
func (instance *Template) Execute(data interface{}) (template.HTML, error) {
	uncompressed := new(bytes.Buffer)
	err := instance.tmpl.ExecuteTemplate(uncompressed, instance.name, data)
	if err != nil {
		return "", err
	}
	compressed := new(bytes.Buffer)
	uncompressedReader := strings.NewReader(uncompressed.String())
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.Add("text/html", &html.Minifier{
		KeepWhitespace: true,
	})
	m.AddFunc("text/javascript", js.Minify)
	if err := m.Minify("text/html", compressed, uncompressedReader); err != nil {
		return "", err
	}
	return template.HTML(compressed.String()), nil
}

// NewRendererFor creates a new renderer for the given parameters.
func NewRendererFor(platform, version string, project Project, pickedDefinitions *PickedDefinitions, apps map[app.ExecutableType]*kingpin.Application) (*Renderer, error) {
	renderer := &Renderer{
		Project:           project,
		PickedDefinitions: pickedDefinitions,
		Apps:              apps,
		Name:              caretakerd.DaemonName,
		Version:           version,
		Description:       caretakerd.Description,
		URL:               caretakerd.URL,
		Platform:          platform,
	}
	renderer.Functions = newFunctionsFor(renderer)

	var err error
	renderer.Template, err = parseTemplate(project, "templates/root", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.IDTemplate, err = parseTemplate(project, "templates/idType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.PointerTemplate, err = parseTemplate(project, "templates/pointerType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.ArrayTemplate, err = parseTemplate(project, "templates/arrayType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.MapTemplate, err = parseTemplate(project, "templates/mapType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.DefinitionStructureTemplate, err = parseTemplate(project, "templates/definitionStructure", renderer.Functions)
	if err != nil {
		return nil, err
	}
	return renderer, nil
}

func newFunctionsFor(renderer *Renderer) template.FuncMap {
	return template.FuncMap{
		"transformIdType": renderer.transformIDType,
		"getDisplayIdOf":  renderer.getDisplayIDOf,
		"renderValueType": renderer.renderValueType,
		"renderMarkdown":  renderer.renderMarkdown,
		"toSimple": func(definition Definition) *SimpleDefinition {
			if result, ok := definition.(*SimpleDefinition); ok {
				return result
			}
			return nil
		},
		"toObject": func(definition Definition) *ObjectDefinition {
			if result, ok := definition.(*ObjectDefinition); ok {
				return result
			}
			return nil
		},
		"toEnum": func(definition Definition) *EnumDefinition {
			if result, ok := definition.(*EnumDefinition); ok {
				return result
			}
			return nil
		},
		"toProperty": func(definition Definition) *PropertyDefinition {
			if result, ok := definition.(*PropertyDefinition); ok {
				return result
			}
			return nil
		},
		"toElement": func(definition Definition) *ElementDefinition {
			if result, ok := definition.(*ElementDefinition); ok {
				return result
			}
			return nil
		},
		"include": func(name string, data interface{}) (template.HTML, error) {
			tmpl, err := parseTemplate(renderer.Project, "includes/"+name, renderer.Functions)
			if err != nil {
				return "", err
			}
			return tmpl.Execute(data)
		},
		"includeJavaScript": func(name string) (template.JS, error) {
			content, err := os.ReadFile(renderer.Project.SrcRootPath + "/manual/templates/scripts/" + name + ".js")
			if err != nil {
				return "", err
			}
			return template.JS(content), err
		},
		"includeCss": func(name string) (template.CSS, error) {
			content, err := os.ReadFile(renderer.Project.SrcRootPath + "/manual/templates/styles/" + name + ".css")
			if err != nil {
				return "", err
			}
			return template.CSS(content), err
		},
		"includeMarkdown": func(name string, headerTypeStart int, headerIdPrefix string, data interface{}) (template.HTML, error) {
			source := renderer.Project.SrcRootPath + "/manual/includes/" + name + ".md"
			content, err := os.ReadFile(source)
			if err != nil {
				return "", err
			}
			rhtml, err := renderer.renderMarkdownWithContext(string(content), nil, headerTypeStart, headerIdPrefix)
			if err != nil {
				return "", err
			}
			tmpl, err := template.New(source).Funcs(renderer.Functions).Parse(string(rhtml))
			if err != nil {
				return "", err
			}
			buf := new(bytes.Buffer)
			err = tmpl.ExecuteTemplate(buf, source, data)
			if err != nil {
				return "", err
			}
			return template.HTML(buf.String()), nil
		},
		"includeLicense": func() (string, error) {
			content, err := os.ReadFile(renderer.Project.SrcRootPath + "/LICENSE")
			if err != nil {
				return "", err
			}
			return string(content), err
		},
		"includeAppUsageOf": func(executableType app.ExecutableType, a *kingpin.Application) (string, error) {
			buf := new(bytes.Buffer)
			a.Name = executableType.String()
			a.UsageWriter(buf)
			context, err := a.ParseContext([]string{})
			if err != nil {
				return "", err
			}
			if err := a.UsageForContextWithTemplate(context, 2, kingpin.LongHelpTemplate); err != nil {
				return "", err
			}
			content := strings.TrimSpace(buf.String())
			content = renderer.replaceUsageEnvVarDisplaysIfNeeded(content)
			return content, nil
		},
		"collectExamples":           renderer.collectExamples,
		"transformElementHtmlId":    renderer.transformElementHTMLID,
		"renderDefinitionStructure": renderer.renderDefinitionStructure,
	}
}

func (instance *Renderer) replaceUsageEnvVarDisplaysIfNeeded(content string) string {
	if instance.Platform == "windows" {
		return otherEnvarPattern.ReplaceAllString(content, "%$1%")
	}
	return windowsEnvarPattern.ReplaceAllString(content, "$$$1")
}

func (instance *Renderer) transformIDType(id IDType) string {
	if len(id.Package) == 0 {
		return id.Name
	}
	name := id.Name
	suffix := ""
	lastHash := strings.LastIndex(id.Name, "#")
	if lastHash > 0 && len(id.Name) > lastHash+1 {
		name = id.Name[:lastHash]
		suffix = "." + id.Name[lastHash+1:]
	}

	if name == "Config" {
		name = instance.capitalize(path.Base(id.Package))
	} else if name == instance.capitalize(path.Base(id.Package)) {
		name = "_" + name
	}

	project := instance.Project
	if id.Package == project.RootPackage {
		return name + suffix
	}
	p := id.Package
	if strings.HasPrefix(id.Package, project.RootPackage+"/") {
		p = p[len(project.RootPackage)+1:]
	}
	return p + "." + name + suffix
}

func (instance *Renderer) getDisplayIDOf(describeable Describable) string {
	id := describeable.ID()
	if withKey, ok := describeable.(WithKey); ok {
		lastHash := strings.LastIndex(id.Name, "#")
		if lastHash > 0 {
			id.Name = id.Name[:lastHash] + "#" + withKey.Key()
		}
	}
	return instance.transformIDType(id)
}

func (instance *Renderer) isMapType(t Type) bool {
	if _, ok := t.(MapType); ok {
		return true
	} else if idType, ok := t.(IDType); ok {
		inlined := instance.PickedDefinitions.FindInlinedFor(idType)
		if inlined != nil && inlined.Inlined() {
			return instance.isMapType(inlined.ValueType())
		}
	}
	return false
}

func (instance *Renderer) renderValueType(t Type) (template.HTML, error) {
	if idType, ok := t.(IDType); ok {
		inlined := instance.PickedDefinitions.FindInlinedFor(idType)
		if inlined != nil && inlined.Inlined() {
			return instance.renderValueType(inlined.ValueType())
		}
		return instance.IDTemplate.Execute(idType)
	} else if arrayType, ok := t.(ArrayType); ok {
		return instance.ArrayTemplate.Execute(arrayType)
	} else if pointerType, ok := t.(PointerType); ok {
		return instance.PointerTemplate.Execute(pointerType)
	} else if mapType, ok := t.(MapType); ok {
		return instance.MapTemplate.Execute(mapType)
	} else {
		return "", errors.New("Unknown type: %v", t)
	}
}

func (instance *Renderer) transformElementHTMLID(definition Definition) (string, error) {
	id := instance.getDisplayIDOf(definition)
	return "configuration.dataType." + id, nil
}

func (instance *Renderer) extractExcerptFrom(definition Definition, headerTypeStart int, headerIDPrefix string) (template.HTML, error) {
	excerpt := definition.Description()
	match := excerptFromCommentExtractionPattern.FindStringSubmatch(excerpt)
	if len(match) == 2 {
		excerpt = match[1]
	}
	excerpt = strings.Replace(excerpt, "\r", "", -1)
	excerpt = strings.Replace(excerpt, "\n", " ", -1)
	excerpt = strings.TrimSpace(excerpt)
	excerptHTML, err := instance.renderMarkdownWithContext(excerpt, definition, headerTypeStart, headerIDPrefix)
	if err != nil {
		return "", err
	}
	excerpt = removeHTMLTags.ReplaceAllString(string(excerptHTML), "")
	return template.HTML(excerpt), nil
}

type renderDefinitionProperty struct {
	Definition   *PropertyDefinition
	MapKeyMarker string
	ID           IDType
	Excerpt      template.HTML
}

func (instance *Renderer) renderDefinitionStructure(level int, id IDType, headerTypeStart int, headerIDPrefix string) (template.HTML, error) {
	definition, err := instance.PickedDefinitions.GetSourceElementBy(id)
	if err != nil {
		return "", err
	}
	if definition == nil {
		return "", nil
	}
	if objectDefinition, ok := definition.(*ObjectDefinition); ok {
		var properties []renderDefinitionProperty
		for _, child := range objectDefinition.Children() {
			propertyDefinition := child.(*PropertyDefinition)
			id := ExtractValueIDType(propertyDefinition.ValueType())
			inlined := instance.PickedDefinitions.FindInlinedFor(id)
			for inlined != nil && inlined.Inlined() {
				id = ExtractValueIDType(inlined.ValueType())
				inlined = instance.PickedDefinitions.FindInlinedFor(id)
			}
			excerpt, err := instance.extractExcerptFrom(propertyDefinition, headerTypeStart, headerIDPrefix)
			if err != nil {
				return "", err
			}
			renderDefinitionProperty := renderDefinitionProperty{
				Definition: propertyDefinition,
				ID:         id,
				Excerpt:    excerpt,
			}
			if instance.isMapType(propertyDefinition.ValueType()) {
				renderDefinitionProperty.MapKeyMarker = instance.singular(propertyDefinition.Key()) + " name"
			}
			properties = append(properties, renderDefinitionProperty)
		}
		indentContent := "<span class=\"tabIndent\"></span>"
		indent := template.HTML(strings.Repeat(indentContent, level))
		nextIndent := template.HTML(strings.Repeat(indentContent, level+1))
		object := map[string]interface{}{
			"object":          objectDefinition,
			"properties":      properties,
			"level":           level,
			"nextLevel":       level + 1,
			"nextNextLevel":   level + 2,
			"indent":          indent,
			"nextIndent":      nextIndent,
			"headerTypeStart": headerTypeStart,
			"headerIdPrefix":  headerIDPrefix,
		}
		rhtml, err := instance.DefinitionStructureTemplate.Execute(object)
		if err != nil {
			return "", err
		}
		if level == 0 {
			rhtml = template.HTML(strings.TrimSpace(string(rhtml)))
		}
		return rhtml, nil
	}
	return "", nil
}

func (instance *Renderer) renderMarkdown(of Describable, headerTypeStart int, headerIDPrefix string) (template.HTML, error) {
	return instance.renderMarkdownWithContext(of.Description(), of, headerTypeStart, headerIDPrefix)
}

func (instance *Renderer) renderMarkdownWithContext(markup string, context Describable, headerTypeStart int, headerIDPrefix string) (template.HTML, error) {
	var err error

	markup = lineBreakCorrectPattern.ReplaceAllString(markup, "\n")
	markup = headerPrefixPattern.ReplaceAllString(markup, "$1"+strings.Repeat("#", headerTypeStart))
	markup = refPropertyPattern.ReplaceAllStringFunc(markup, func(inline string) string {
		match := refPropertyPattern.FindStringSubmatch(inline)
		ref := match[1]
		idType := instance.resolveRef(ref, context)
		element, pErr := instance.PickedDefinitions.GetSourceElementBy(idType)
		if pErr != nil {
			err = pErr
			return inline
		}
		if element != nil {
			targetType := instance.getDisplayIDOf(element)
			display := strings.TrimSpace(match[2])
			if len(display) <= 0 {
				display = targetType
			}
			return "[``" + display + "``](#configuration.dataType." + targetType + ")"
		}
		err = errors.New("Unknonwn reference: %v", ref)
		return markup
	})
	if err != nil {
		return "", err
	}
	prefix := ""
	if len(headerIDPrefix) > 0 {
		prefix = headerIDPrefix + "."
	}
	renderer := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
		HeadingIDPrefix: prefix,
		Flags: blackfriday.UseXHTML |
			blackfriday.Smartypants |
			blackfriday.SmartypantsFractions |
			blackfriday.SmartypantsDashes |
			blackfriday.SmartypantsLatexDashes,
	})
	rhtml := blackfriday.Run([]byte(markup),
		blackfriday.WithRenderer(renderer),
		blackfriday.WithExtensions(blackfriday.NoIntraEmphasis|
			blackfriday.Tables|
			blackfriday.FencedCode|
			blackfriday.Autolink|
			blackfriday.Strikethrough|
			blackfriday.SpaceHeadings|
			blackfriday.HeadingIDs|
			blackfriday.BackslashLineBreak|
			blackfriday.DefinitionLists|
			blackfriday.AutoHeadingIDs),
	)
	return template.HTML(strings.TrimSpace(string(rhtml))), nil
}

func (instance *Renderer) resolveRef(ref string, context Describable) IDType {
	if context != nil && strings.HasPrefix(ref, "#") {
		name := ref[1:]
		contextID := context.ID()
		lastDotIfContextID := strings.LastIndex(contextID.Name, "#")
		if lastDotIfContextID > 0 {
			name = contextID.Name[:lastDotIfContextID] + "#" + name
		}
		return IDType{
			Package: contextID.Package,
			Name:    name,
		}
	} else if context != nil && strings.HasPrefix(ref, ".") {
		contextID := context.ID()
		return IDType{
			Package: contextID.Package,
			Name:    ref[1:],
		}
	} else {
		t := ParseType(ref)
		if idType, ok := t.(IDType); ok {
			return idType
		}
		return IDType{Name: ref}
	}
}

type example struct {
	ID          string
	Title       string
	CodeType    string
	CodeContent string
}

func (instance *Renderer) collectExamples() ([]example, error) {
	examplesSources, err := filepath.Glob(instance.Project.SrcRootPath + "/manual/examples/*.yaml")
	if err != nil {
		return []example{}, err
	}
	var examples []example
	for _, exampleSource := range examplesSources {
		b, err := os.ReadFile(exampleSource)
		if err != nil {
			return nil, err
		}
		content, title, id := instance.extractTitleFrom(string(b), exampleSource)
		examples = append(examples, example{
			ID:          "configuration.examples." + id,
			Title:       title,
			CodeType:    "yaml",
			CodeContent: content,
		})
	}
	return examples, nil
}

func (instance *Renderer) extractTitleFrom(source string, filename string) (string, string, string) {
	id := filepath.Base(filename)
	ext := filepath.Ext(id)
	if len(ext) > 0 {
		id = id[:len(id)-len(ext)]
	}
	match := titlePropertyPattern.FindStringSubmatch(source)
	if len(match) > 0 {
		return titlePropertyPattern.ReplaceAllString(source, ""), match[1], id
	}

	return source, id, id
}

func parseTemplate(project Project, name string, functions template.FuncMap) (*Template, error) {
	source := project.SrcRootPath + "/manual/" + name + ".html"
	b, err := os.ReadFile(source)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(source).Funcs(functions).Parse(string(b))
	if err != nil {
		return nil, err
	}
	return &Template{
		tmpl: tmpl,
		name: source,
	}, nil
}

func (instance *Renderer) capitalize(what string) string {
	l := len(what)
	if l <= 0 {
		return ""
	} else if l == 1 {
		return strings.ToUpper(what)
	} else {
		return strings.ToUpper(what[0:1]) + what[1:]
	}
}

func (instance *Renderer) singular(what string) string {
	if strings.HasPrefix(what, "s") {
		return what[:len(what)-1]
	}
	return what
}
