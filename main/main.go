package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/app"
	"github.com/echocat/caretakerd/panics"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var executableNamePattern = regexp.MustCompile("(?:^|" + regexp.QuoteMeta(string(os.PathSeparator)) + ")" + caretakerd.BaseName + "(d|ctl)(?:$|[\\.\\-\\_].*$)")

func main() {
	defer panics.DefaultPanicHandler()
	a := app.NewAppFor(runtime.GOOS, getExecutableType())

	kingpin.MustParse(a.Parse(os.Args[1:]))
}

func getExecutableType() app.ExecutableType {
	executable := strings.ToLower(filepath.Base(os.Args[0]))
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
