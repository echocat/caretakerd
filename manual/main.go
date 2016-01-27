package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"fmt"
	"os"
	"github.com/echocat/caretakerd/logger"
	"github.com/echocat/caretakerd/sync"
)

var LOGGER, _ = logger.NewLogger(logger.Config{
	Level: logger.Info,
	Filename: "console",
	Pattern: "%d{YYYY-MM-DD HH:mm:ss} [%-5.5p] %m%n%P{%m}",
}, "manual", sync.NewSyncGroup())

func panicHandler() {
	if r := recover(); r != nil {
		LOGGER.LogProblem(r, logger.Fatal, "There is an unrecoverable problem occured.")
		os.Exit(2)
	}
}

func getSrcRootPath() string {
	if len(os.Args) < 2 || len(os.Args[1]) <= 0 {
		fmt.Fprintf(os.Stderr, "Usage: %v <project source root path>\n", os.Args[0])
		os.Exit(1)
	}
	return os.Args[1]
}

func main() {
	defer panicHandler()
	srcRootPath := getSrcRootPath()
	project := DeterminateProject(srcRootPath)
	LOGGER.Log(logger.Info, "Root package: %v", project.RootPackage)
	LOGGER.Log(logger.Info, "Source root path: %v", project.SrcRootPath)

	bytes, err := ioutil.ReadFile("manual/docs/configuration/examples.md")
	if err != nil {
		panic(err)
	}
	content := blackfriday.MarkdownCommon(bytes)
	err = ioutil.WriteFile("targets/test.html", content, 0)
	if err != nil {
		panic(err)
	}

}
