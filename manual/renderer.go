package main

import (
	"html/template"
	"io/ioutil"
	"github.com/russross/blackfriday"
	"io"
	"strings"
	"path"
	"bytes"
)

type Renderer struct {
	Template                    *template.Template
	IdTemplate                  *template.Template
	PointerTemplate             *template.Template
	ArrayTemplate               *template.Template
	MapTemplate                 *template.Template
	DataTypeAnchorTemplate      *template.Template
	DefinitionStructureTemplate *template.Template
	HeaderTemplate              *template.Template

	Functions                   template.FuncMap
	Project                     Project
	PickedDefinitions           *PickedDefinitions
}

func (instance *Renderer) Execute(writer io.Writer) error {
	return instance.Template.ExecuteTemplate(writer, "manual/template.html", instance.PickedDefinitions)
}

func NewRendererFor(project Project, pickedDefinitions *PickedDefinitions) (*Renderer, error) {
	renderer := &Renderer{
		Project: project,
		PickedDefinitions: pickedDefinitions,
	}
	renderer.Functions = newFunctionsFor(renderer)

	var err error
	renderer.Template, err = parseTemplate("template", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.IdTemplate, err = parseTemplate("template.idType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.PointerTemplate, err = parseTemplate("template.pointerType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.ArrayTemplate, err = parseTemplate("template.arrayType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.MapTemplate, err = parseTemplate("template.mapType", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.DataTypeAnchorTemplate, err = parseTemplate("template.dataTypeAnchor", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.DefinitionStructureTemplate, err = parseTemplate("template.definitionStructure", renderer.Functions)
	if err != nil {
		return nil, err
	}
	renderer.HeaderTemplate, err = parseTemplate("template.header", renderer.Functions)
	if err != nil {
		return nil, err
	}

	return renderer, nil
}

func newFunctionsFor(renderer *Renderer) template.FuncMap {
	return template.FuncMap{
		"transformIdType": renderer.transformIdType,
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
		"includeTemplateCss": func() (template.CSS, error) {
			content, err := ioutil.ReadFile("manual/template.css")
			if err != nil {
				return "", err
			}
			return template.CSS(content), err
		},
		"renderDataTypeAnchor": renderer.renderDataTypeAnchor,
		"renderDefinitionStructure": renderer.renderDefinitionStructure,
		"header": renderer.header,
	}
}

func (instance *Renderer) transformIdType(id IdType) string {
	if len(id.Package) == 0 {
		return id.Name
	}
	name := id.Name
	if name == "Config" {
		name = instance.capitalize(path.Base(id.Package))
	} else if name == instance.capitalize(path.Base(id.Package)) {
		name = "_" + name
	}

	project := instance.Project
	if id.Package == project.RootPackage {
		return name
	}
	p := id.Package
	if strings.HasPrefix(id.Package, project.RootPackage + "/") {
		p = p[len(project.RootPackage) + 1:]
	}
	return p + "." + name
}

func (instance *Renderer) renderValueType(t Type) (template.HTML, error) {
	buf := new(bytes.Buffer)
	if idType, ok := t.(IdType); ok {
		inlined := instance.PickedDefinitions.FindInlinedFor(idType)
		if inlined != nil && inlined.Inlined() {
			return instance.renderValueType(inlined.ValueType())
		} else {
			err := instance.IdTemplate.Execute(buf, idType)
			if err != nil {
				return "", err
			}
		}
	} else if arrayType, ok := t.(ArrayType); ok {
		err := instance.ArrayTemplate.Execute(buf, arrayType)
		if err != nil {
			return "", err
		}
	} else if pointerType, ok := t.(PointerType); ok {
		err := instance.PointerTemplate.Execute(buf, pointerType)
		if err != nil {
			return "", err
		}
	} else if mapType, ok := t.(MapType); ok {
		err := instance.MapTemplate.Execute(buf, mapType)
		if err != nil {
			return "", err
		}
	}
	return template.HTML(buf.String()), nil
}

func (instance *Renderer) renderDataTypeAnchor(definition Definition, suffix string) (template.HTML, error) {
	id := ""
	id += instance.transformIdType(definition.Id())
	if len(suffix) > 0 {
		id += "." + suffix
	}

	buf := new(bytes.Buffer)
	err := instance.DataTypeAnchorTemplate.Execute(buf, id)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

type RenderDefinitionProperty struct {
	Definition *PropertyDefinition
	Id         IdType
}

func (instance *Renderer) renderDefinitionStructure(level int, id IdType) (template.HTML, error) {
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
			properties = append(properties, RenderDefinitionProperty{
				Definition: propertyDefinition,
				Id: id,
			})
		}
		indent := ""
		for i := 0; i < level; i++ {
			indent += "    "
		}
		object := map[string]interface{}{
			"object": objectDefinition,
			"properties": properties,
			"level": level,
			"nextLevel": level + 1,
			"indent": indent,
		}
		buf := new(bytes.Buffer)
		err := instance.DefinitionStructureTemplate.Execute(buf, object)
		if err != nil {
			return "", err
		}
		html := buf.String()
		if level == 0 {
			html = strings.TrimSpace(html)
		}
		return template.HTML(html), nil
	} else {
		return "", nil
	}
}

func (instance *Renderer) header(level int, id string, css string, content string) (template.HTML, error) {
	buf := new(bytes.Buffer)
	object := map[string]interface{}{
		"level": level,
		"id": id,
		"css": css,
		"content": content,
	}
	err := instance.HeaderTemplate.Execute(buf, object)
	if err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}

func (instance *Renderer) renderMarkdown(markup string) (template.HTML, error) {
	escapedMarkup := headerPrefix.ReplaceAllString(markup, "$1####")
	html := blackfriday.MarkdownCommon([]byte(escapedMarkup))
	return template.HTML(html), nil
}

func parseTemplate(name string, functions template.FuncMap) (*template.Template, error) {
	source := "manual/" + name + ".html"
	bytes, err := ioutil.ReadFile(source)
	if err != nil {
		return nil, err
	}
	result, err := template.New(source).Funcs(functions).Parse(string(bytes))
	if err != nil {
		return nil, err
	}
	return result, nil
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
