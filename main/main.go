package main

import (
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd/panics"
	"os"
	"regexp"
	"strings"
)

var executableNamePattern = regexp.MustCompile("(?:^|" + regexp.QuoteMeta(string(os.PathSeparator)) + ")" + caretakerd.BaseName + "(d|ctl)(?:$|[\\.\\-\\_].*$)")

func main() {
	defer panics.DefaultPanicHandler()
	app := app.NewAppFor(getExecutableType())

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(1)
	}
}

func getExecutableType() app.ExecutableType {
	executable := strings.ToLower(os.Args[0])
	match := executableNamePattern.FindStringSubmatch(executable)
	if match != nil && len(match) == 2 {
		switch match[1] {
		case "d":
			return app.Daemon
		case "ctl":
			return app.Control
		}
	}
	return app.Generic
}
