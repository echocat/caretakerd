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
	if len(os.Args) < 4 || len(os.Args[1]) <= 0 || len(os.Args[2]) <= 0 || len(os.Args[3]) <= 0 {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %v <version> <platform> <output>\n", os.Args[0])
		os.Exit(1)
	}
	version := os.Args[1]
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}
	platform := os.Args[2]
	plainFile := os.Args[3]

	defer panicHandler()
	project, err := DeterminateProject("github.com/echocat/caretakerd")
	if err != nil {
		panic(err)
	}
	log.Log(logger.Debug, "Build manual for package=%v, path=%v, platform=%s", project.RootPackage, project.SrcRootPath, platform)

	definitions, err := ParseDefinitions(project)
	if err != nil {
		panic(err)
	}
	pd, err := PickDefinitionsFrom(definitions, NewIDType(project.RootPackage, "Config", false))
	if err != nil {
		panic(err)
	}

	apps := app.NewAppsFor(platform)

	renderer, err := NewRendererFor(platform, version, project, pd, apps)
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
