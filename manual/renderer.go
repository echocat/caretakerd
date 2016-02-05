package main

import (
	"html/template"
	"io/ioutil"
	"strings"
	"path"
	"bytes"
	"regexp"
	"path/filepath"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/html"
	"github.com/russross/blackfriday"
	"github.com/echocat/caretakerd/errors"
	"github.com/codegangsta/cli"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd"
)

var headerPrefixPattern = regexp.MustCompile("(?m)^([\\* 0-9\\.]*)#")
var excerptFromCommentExtractionPattern = regexp.MustCompile("(?s)^(.*?(?:\\.\\s|$))")
var removeHtmlTags = regexp.MustCompile("(?sm)<[^>]+>")
var findLeadingWhitespaces = regexp.MustCompile("(?m)^( +)")
var refPropertyPattern = regexp.MustCompile("{@ref +([^\\}\\s]+)\\s*([^\\}]*)}")
var titlePropertyPattern = regexp.MustCompile("(?m)^#\\s*@title\\s+(.*)\\s*(:?\r\n|\n)")
var commandsInAppPattern = regexp.MustCompile("(?s)(COMMANDS:\n)(.*?)(\n[\t ]*\n|$)")
var commandLinePattern = regexp.MustCompile("(?m)^( +|\t)([a-z][a-zA-Z0-9]+)(.*)$")

type Describeable interface {
	Id() IdType
	Description() string
}

type Renderer struct {
	Template                    *Template
	IdTemplate                  *Template
	PointerTemplate             *Template
	ArrayTemplate               *Template
	MapTemplate                 *Template
	DefinitionStructureTemplate *Template

	Functions                   template.FuncMap
	Project                     Project
	PickedDefinitions           *PickedDefinitions
	Apps                        map[app.ExecutableType]*cli.App

	Name                        string
	Version                     string
	Description                 string
	Url                         string
}

func (instance *Renderer) Execute() (template.HTML, error) {
	return instance.Template.Execute(instance)
}

type Template struct {
	tmpl *template.Template
	name string
}

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
		KeepWhitespace:  true,
	})
	m.AddFunc("text/javascript", js.Minify)
	if err := m.Minify("text/html", compressed, uncompressedReader); err != nil {
		return "", err
	}
	return template.HTML(compressed.String()), nil
}

func NewRendererFor(project Project, pickedDefinitions *PickedDefinitions, apps map[app.ExecutableType]*cli.App) (*Renderer, error) {
	renderer := &Renderer{
		Project: project,
		PickedDefinitions: pickedDefinitions,
		Apps: apps,
		Name: caretakerd.DAEMON_NAME,
		Version: caretakerd.VERSION,
		Description: caretakerd.DESCRIPTION,
		Url: caretakerd.URL,
	}
	renderer.Functions = newFunctionsFor(renderer)

	var err error
	renderer.Template, err = parseTemplate(project, "templates/root", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.IdTemplate, err = parseTemplate(project, "templates/idType", renderer.Functions)
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
		"transformIdType": renderer.transformIdType,
		"getDisplayIdOf": renderer.getDisplayIdOf,
		"renderValueType": renderer.renderValueType,
		"renderMarkdown": renderer.renderMarkdown,
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
			tmpl, err := parseTemplate(renderer.Project, "includes/" + name, renderer.Functions)
			if err != nil {
				return "", err
			}
			return tmpl.Execute(data)
		},
		"includeJavaScript": func(name string) (template.JS, error) {
			content, err := ioutil.ReadFile(renderer.Project.SrcRootPath + "/manual/templates/scripts/" + name + ".js")
			if err != nil {
				return "", err
			}
			return template.JS(content), err
		},
		"includeCss": func(name string) (template.CSS, error) {
			content, err := ioutil.ReadFile(renderer.Project.SrcRootPath + "/manual/templates/styles/" + name + ".css")
			if err != nil {
				return "", err
			}
			return template.CSS(content), err
		},
		"includeMarkdown": func(name string, headerTypeStart int, headerIdPrefix string, data interface{}) (template.HTML, error) {
			source := renderer.Project.SrcRootPath + "/manual/includes/" + name + ".md"
			content, err := ioutil.ReadFile(source)
			if err != nil {
				return "", err
			}
			html, err := renderer.renderMarkdownWithContext(string(content), nil, headerTypeStart, headerIdPrefix)
			if err != nil {
				return "", err
			}
			tmpl, err := template.New(source).Funcs(renderer.Functions).Parse(string(html))
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
			content, err := ioutil.ReadFile(renderer.Project.SrcRootPath + "/LICENSE")
			if err != nil {
				return "", err
			}
			return string(content), err
		},
		"includeAppUsageOf": func(executableType app.ExecutableType, app *cli.App) template.HTML {
			app.HelpName = executableType.String()
			buf := new(bytes.Buffer)
			cli.HelpPrinter(buf, cli.AppHelpTemplate, app)
			content := commandsInAppPattern.ReplaceAllStringFunc(buf.String(), func(what string) string {
				match := commandsInAppPattern.FindStringSubmatch(what)
				content := commandLinePattern.ReplaceAllStringFunc(match[2], func(subWhat string) string {
					subMatch := commandLinePattern.FindStringSubmatch(subWhat)
					return subMatch[1] + "<a href=\"#commands." + executableType.String() + "." + subMatch[2] + "\">" + subMatch[2] + "</a>" + subMatch[3]
				})
				return match[1] + content + match[3]
			})
			content = findLeadingWhitespaces.ReplaceAllStringFunc(content, func(spaces string) string {
				return strings.Replace(spaces, " ", "<span class=\"indent\"></span>", -1)
			})
			content = strings.Replace(content, "\n", "<br>", -1)
			return template.HTML(content)
		},
		"includeCommandUsageOf": func(command *cli.Command) string {
			if len(command.HelpName) <= 0 {
				command.HelpName = command.Name
			}
			buf := new(bytes.Buffer)
			cli.HelpPrinter(buf, cli.CommandHelpTemplate, command)
			return buf.String()
		},
		"collectExamples": renderer.collectExamples,
		"transformElementHtmlId": renderer.transformElementHtmlId,
		"renderDefinitionStructure": renderer.renderDefinitionStructure,
	}
}

func (instance *Renderer) transformIdType(id IdType) string {
	if len(id.Package) == 0 {
		return id.Name
	}
	name := id.Name
	suffix := ""
	lastHash := strings.LastIndex(id.Name, "#")
	if lastHash > 0 && len(id.Name) > lastHash + 1 {
		name = id.Name[:lastHash]
		suffix = "." + id.Name[lastHash + 1:]
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
	if strings.HasPrefix(id.Package, project.RootPackage + "/") {
		p = p[len(project.RootPackage) + 1:]
	}
	return p + "." + name + suffix
}

func (instance *Renderer) getDisplayIdOf(describeable Describeable) string {
	id := describeable.Id()
	if withKey, ok := describeable.(WithKey); ok {
		lastHash := strings.LastIndex(id.Name, "#")
		if lastHash > 0 {
			id.Name = id.Name[:lastHash] + "#" + withKey.Key()
		}
	}
	return instance.transformIdType(id)
}

func (instance *Renderer) isMapType(t Type) bool {
	if _, ok := t.(MapType); ok {
		return true
	} else if idType, ok := t.(IdType); ok {
		inlined := instance.PickedDefinitions.FindInlinedFor(idType)
		if inlined != nil && inlined.Inlined() {
			return instance.isMapType(inlined.ValueType())
		}
	}
	return false
}

func (instance *Renderer) renderValueType(t Type) (template.HTML, error) {
	if idType, ok := t.(IdType); ok {
		inlined := instance.PickedDefinitions.FindInlinedFor(idType)
		if inlined != nil && inlined.Inlined() {
			return instance.renderValueType(inlined.ValueType())
		} else {
			return instance.IdTemplate.Execute(idType)
		}
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

func (instance *Renderer) transformElementHtmlId(definition Definition) (string, error) {
	id := instance.getDisplayIdOf(definition)
	return "configuration.dataType." + id, nil
}

func (instance *Renderer) extractExcerptFrom(definition Definition, headerTypeStart int, headerIdPrefix string) (template.HTML, error) {
	excerpt := definition.Description()
	match := excerptFromCommentExtractionPattern.FindStringSubmatch(excerpt)
	if match != nil && len(match) == 2 {
		excerpt = match[1]
	}
	excerpt = strings.Replace(excerpt, "\r", "", -1)
	excerpt = strings.Replace(excerpt, "\n", " ", -1)
	excerpt = strings.TrimSpace(excerpt)
	excerptHtml, err := instance.renderMarkdownWithContext(excerpt, definition, headerTypeStart, headerIdPrefix)
	if err != nil {
		return "", err
	}
	excerpt = removeHtmlTags.ReplaceAllString(string(excerptHtml), "")
	return template.HTML(excerpt), nil
}

type RenderDefinitionProperty struct {
	Definition   *PropertyDefinition
	MapKeyMarker string
	Id           IdType
	Excerpt      template.HTML
}

func (instance *Renderer) renderDefinitionStructure(level int, id IdType, headerTypeStart int, headerIdPrefix string) (template.HTML, error) {
	definition, err := instance.PickedDefinitions.GetSourceElementBy(id)
	if err != nil {
		return "", err
	}
	if definition == nil {
		return "", nil
	}
	if objectDefinition, ok := definition.(*ObjectDefinition); ok {
		properties := []RenderDefinitionProperty{}
		for _, child := range objectDefinition.Children() {
			propertyDefinition := child.(*PropertyDefinition)
			id := ExtractValueIdType(propertyDefinition.ValueType())
			inlined := instance.PickedDefinitions.FindInlinedFor(id)
			for inlined != nil && inlined.Inlined() {
				id = ExtractValueIdType(inlined.ValueType())
				inlined = instance.PickedDefinitions.FindInlinedFor(id)
			}
			excerpt, err := instance.extractExcerptFrom(propertyDefinition, headerTypeStart, headerIdPrefix)
			if err != nil {
				return "", err
			}
			renderDefinitionProperty := RenderDefinitionProperty{
				Definition: propertyDefinition,
				Id: id,
				Excerpt: excerpt,
			}
			if instance.isMapType(propertyDefinition.ValueType()) {
				renderDefinitionProperty.MapKeyMarker = instance.singular(propertyDefinition.Key()) + " name"
			}
			properties = append(properties, renderDefinitionProperty)
		}
		indentContent := "<span class=\"tabIndent\"></span>"
		indent := template.HTML(strings.Repeat(indentContent, level))
		nextIndent := template.HTML(strings.Repeat(indentContent, level + 1))
		object := map[string]interface{}{
			"object": objectDefinition,
			"properties": properties,
			"level": level,
			"nextLevel": level + 1,
			"nextNextLevel": level + 2,
			"indent": indent,
			"nextIndent": nextIndent,
			"headerTypeStart": headerTypeStart,
			"headerIdPrefix": headerIdPrefix,
		}
		html, err := instance.DefinitionStructureTemplate.Execute(object)
		if err != nil {
			return "", err
		}
		if level == 0 {
			html = template.HTML(strings.TrimSpace(string(html)))
		}
		return html, nil
	} else {
		return "", nil
	}
}

func (instance *Renderer) renderMarkdown(of Describeable, headerTypeStart int, headerIdPrefix string) (template.HTML, error) {
	return instance.renderMarkdownWithContext(of.Description(), of, headerTypeStart, headerIdPrefix)
}

func (instance *Renderer) renderMarkdownWithContext(markup string, context Describeable, headerTypeStart int, headerIdPrefix string) (template.HTML, error) {
	var err error

	markup = headerPrefixPattern.ReplaceAllString(markup, "$1" + strings.Repeat("#", headerTypeStart))
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
			targetType := instance.getDisplayIdOf(element)
			display := strings.TrimSpace(match[2])
			if len(display) <= 0 {
				display = targetType
			}
			return "[``" + display + "``](#configuration.dataType." + targetType + ")"
		} else {
			err = errors.New("Unknonwn reference: %v", ref)
			return markup
		}
	})
	if err != nil {
		return "", err
	}
	prefix := ""
	if len(headerIdPrefix) > 0 {
		prefix = headerIdPrefix + "."
	}
	renderer := blackfriday.HtmlRendererWithParameters(blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES,
		"",
		"",
		blackfriday.HtmlRendererParameters{
			HeaderIDPrefix: prefix,
		},
	)
	html := blackfriday.MarkdownOptions([]byte(markup), renderer, blackfriday.Options{
		Extensions:  blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
			blackfriday.EXTENSION_TABLES |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_SPACE_HEADERS |
			blackfriday.EXTENSION_HEADER_IDS |
			blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
			blackfriday.EXTENSION_DEFINITION_LISTS |
			blackfriday.EXTENSION_AUTO_HEADER_IDS,
	})
	return template.HTML(strings.TrimSpace(string(html))), nil
}

func (instance *Renderer) resolveRef(ref string, context Describeable) IdType {
	if context != nil && strings.HasPrefix(ref, "#") {
		name := ref[1:]
		contextId := context.Id()
		lastDotIfContextId := strings.LastIndex(contextId.Name, "#")
		if lastDotIfContextId > 0 {
			name = contextId.Name[:lastDotIfContextId] + "#" + name
		}
		return IdType{
			Package: contextId.Package,
			Name: name,
		}
	} else if context != nil && strings.HasPrefix(ref, ".") {
		contextId := context.Id()
		return IdType{
			Package: contextId.Package,
			Name: ref[1:],
		}
	} else {
		t := ParseType(ref)
		if idType, ok := t.(IdType); ok {
			return idType
		} else {
			return IdType{Name: ref}
		}
	}
}

type Example struct {
	Id          string
	Title       string
	CodeType    string
	CodeContent string
}

func (instance *Renderer) collectExamples() ([]Example, error) {
	examplesSources, err := filepath.Glob(instance.Project.SrcRootPath + "/manual/examples/*.yaml")
	if err != nil {
		return []Example{}, err
	}
	examples := []Example{}
	for _, exampleSource := range examplesSources {
		bytes, err := ioutil.ReadFile(exampleSource)
		if err != nil {
			return nil, err
		}
		content, title, id := instance.extractTitleFrom(string(bytes), exampleSource)
		examples = append(examples, Example{
			Id: "configuration.examples." + id,
			Title: title,
			CodeType: "yaml",
			CodeContent: content,
		})
	}
	return examples, nil
}

func (instance *Renderer) extractTitleFrom(source string, filename string) (string, string, string) {
	id := filepath.Base(filename)
	ext := filepath.Ext(id)
	if len(ext) > 0 {
		id = id[:len(id) - len(ext)]
	}
	match := titlePropertyPattern.FindStringSubmatch(source)
	if match != nil && len(match) > 0 {
		return titlePropertyPattern.ReplaceAllString(source, ""), match[1], id
	}

	return source, id, id
}

func parseTemplate(project Project, name string, functions template.FuncMap) (*Template, error) {
	source := project.SrcRootPath + "/manual/" + name + ".html"
	bytes, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(source).Funcs(functions).Parse(string(bytes))
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
		return what[:len(what) - 1]
	}
	return what
}
