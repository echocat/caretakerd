package main

import (
	"fmt"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
	"os"
	"regexp"
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

func main() {
	defer panicHandler()
	srcRootPath := getSrcRootPath()
	project, err := DeterminateProject(srcRootPath)
	if err != nil {
		panic(err)
	}
	LOGGER.Log(logger.Info, "Root package: %v", project.RootPackage)
	LOGGER.Log(logger.Info, "Source root path: %v", project.SrcRootPath)

	definitions, err := ParseDefinitions(project)
	if err != nil {
		panic(err)
	}
	pd, err := PickDefinitionsFrom(definitions, NewIdType(project.RootPackage, "Config", false))
	if err != nil {
		panic(err)
	}

	renderer, err := NewRendererFor(project, pd)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("target/manual.html", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0)
	if err != nil {
		panic(err)
	}
	err = renderer.Execute(file)
	if err != nil {
		panic(err)
	}

	err = file.Close()
	if err != nil {
		panic(err)
	}
}
