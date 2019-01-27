package main

import (
	"fmt"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var log, _ = logger.NewLogger(logger.Config{
	Level:    logger.Info,
	Filename: "console",
	Pattern:  "%d{YYYY-MM-DD HH:mm:ss} [%-5.5p] %m%n%P{%m}",
}, "manual", sync.NewGroup())

func panicHandler() {
	if r := recover(); r != nil {
		log.LogProblem(r, logger.Fatal, "There is an unrecoverable problem occurred.")
		os.Exit(2)
	}
}

func main() {
	if len(os.Args) < 2 || len(os.Args[1]) <= 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %v <version> <output>\n", os.Args[0])
		os.Exit(1)
	}
	version := os.Args[1]
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}
	plainFile := os.Args[2]

	defer panicHandler()
	project, err := DeterminateProject("github.com/echocat/caretakerd")
	if err != nil {
		panic(err)
	}
	log.Log(logger.Info, "Root package: %v", project.RootPackage)
	log.Log(logger.Info, "Source root path: %v", project.SrcRootPath)

	definitions, err := ParseDefinitions(project)
	if err != nil {
		panic(err)
	}
	pd, err := PickDefinitionsFrom(definitions, NewIDType(project.RootPackage, "Config", false))
	if err != nil {
		panic(err)
	}

	apps := app.NewApps()

	renderer, err := NewRendererFor(version, project, pd, apps)
	if err != nil {
		panic(err)
	}

	content, err := renderer.Execute()
	if err != nil {
		panic(err)
	}

	file, err := filepath.Abs(plainFile)
	if err != nil {
		panic(err)
	}
	directory := filepath.Dir(file)
	if err := os.MkdirAll(directory, 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(file, []byte(content), 0655); err != nil {
		panic(err)
	}
}
