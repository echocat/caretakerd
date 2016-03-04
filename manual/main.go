package main

import (
	"fmt"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	"io/ioutil"
	"os"
	"path/filepath"
)

var log, _ = logger.NewLogger(logger.Config{
	Level:    logger.Info,
	Filename: "console",
	Pattern:  "%d{YYYY-MM-DD HH:mm:ss} [%-5.5p] %m%n%P{%m}",
}, "manual", sync.NewSyncGroup())

func panicHandler() {
	if r := recover(); r != nil {
		log.LogProblem(r, logger.Fatal, "There is an unrecoverable problem occured.")
		os.Exit(2)
	}
}

func getSrcRootPath() string {
	if len(os.Args) < 2 || len(os.Args[1]) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v <package> <output>\n", os.Args[0])
		os.Exit(1)
	}
	return os.Args[1]
}

func main() {
	defer panicHandler()
	srcRootPath := getSrcRootPath()
	project, err := DeterminateProject(srcRootPath)
	if err != nil {
		panic(err)
	}
	log.Log(logger.Info, "Root package: %v", project.RootPackage)
	log.Log(logger.Info, "Source root path: %v", project.SrcRootPath)

	definitions, err := ParseDefinitions(project)
	if err != nil {
		panic(err)
	}
	pd, err := PickDefinitionsFrom(definitions, NewIdType(project.RootPackage, "Config", false))
	if err != nil {
		panic(err)
	}

	apps := app.NewApps()

	renderer, err := NewRendererFor(project, pd, apps)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 3 || len(os.Args[2]) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v <package> <output>\n", os.Args[0])
		os.Exit(1)
	}

	content, err := renderer.Execute()
	if err != nil {
		panic(err)
	}

	plainFile := os.Args[2]
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
