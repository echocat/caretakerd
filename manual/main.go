package main

import (
	"fmt"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	"os"
	"text/template"
	"strings"
	"path"
	"github.com/russross/blackfriday"
	"regexp"
	"io/ioutil"
)

var headerPrefix = regexp.MustCompile("(?m)^([\\* 0-9\\.]*)#")

var LOGGER, _ = logger.NewLogger(logger.Config{
	Level:    logger.Info,
	Filename: "console",
	Pattern:  "%d{YYYY-MM-DD HH:mm:ss} [%-5.5p] %m%n%P{%m}",
}, "manual", sync.NewSyncGroup())

func panicHandler() {
	if r := recover(); r != nil {
		LOGGER.LogProblem(r, logger.Info, "There is an unrecoverable problem occured.")
		os.Exit(2)
	}
}

func getSrcRootPath() string {
	if len(os.Args) < 2 || len(os.Args[1]) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v <package>\n", os.Args[0])
		os.Exit(1)
	}
	return os.Args[1]
}

func (project Project) transformIdType (idType IdType) string {
	if len(idType.Package) == 0 {
		return idType.Name
	}
	name := idType.Name
	if name == "Config" {
		name = capitalize(path.Base(idType.Package))
	} else if name == capitalize(path.Base(idType.Package)) {
		name = "_" + name
	}

	if idType.Package == project.RootPackage {
		return name
	}
	p := idType.Package
	if strings.HasPrefix(idType.Package, project.RootPackage + "/") {
		p = p[len(project.RootPackage) + 1:]
	}
	return p + "." + name
}

func (project Project) transformValueType(valueType Type) string {
	if idType, ok := valueType.(IdType); ok {
		return project.transformIdType(idType)
	}
	return "--error--"
}

func main() {
	defer panicHandler()
	srcRootPath := getSrcRootPath()
	project, err := DeterminateProject(srcRootPath)
	if err != nil {
		panic(err)
	}
	LOGGER.Log(logger.Info, "Root package: %v", project.RootPackage)
	LOGGER.Log(logger.Info, "Source root path: %v", project.SrcRootPath)

	templateFunction := template.FuncMap{
		"transformIdType": project.transformIdType,
		"transformValueType": project.transformValueType,
		"markdown": func(code string) string {
			markdown := headerPrefix.ReplaceAllString(code, "$1####")
			html := blackfriday.MarkdownCommon([]byte(markdown))
			return string(html)
		},
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
		"includeTemplateCss": func() (string, error) {
			content, err := ioutil.ReadFile("manual/template.css")
			if err != nil {
				return "", err
			}
			return string(content), err
		},
		"anchorFor": func(prefix string, definition Definition, suffix string) string {
			id := ""
			if len(prefix) > 0 {
				id += prefix + "."
			}
			id += project.transformIdType(definition.Id())
			if len(suffix) > 0 {
				id += "." + suffix
			}
			return "<a id=\"" + id + "\" class=\"anchor\" href=\"#" + id + "\" aria-hidden=\"true\">" +
				"<span aria-hidden=\"true\" class=\"octicon octicon-link\"></span>" +
				"</a>"
		},
	}

	definitions, err := ParseDefinitions(project)
	if err != nil {
		panic(err)
	}
	pd, err := PickDefinitionsFrom(definitions, NewIdType(project.RootPackage, "Config", false))
	if err != nil {
		panic(err)
	}
	//	for _, definition := range pd.TopLevelDefinitions {
	//		LOGGER.Log(logger.Info, "%v", definition)
	//	}

	plainTemplate, err := ioutil.ReadFile("manual/template.html")
	if err != nil {
		panic(err)
	}
	file, err := os.OpenFile("target/manual.html", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0)
	if err != nil {
		panic(err)
	}
	tmpl, err := template.New("moo").Funcs(templateFunction).Parse(string(plainTemplate))
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(file, pd)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func capitalize(what string) string {
	l := len(what)
	if l <= 0 {
		return ""
	} else if l == 1 {
		return strings.ToUpper(what)
	} else {
		return strings.ToUpper(what[0:1]) + what[1:]
	}
}
