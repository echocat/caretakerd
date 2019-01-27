package main

import (
	"fmt"
	"github.com/alecthomas/kingpin"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	_ = app.Command("build", "executes builds for the project").
		Action(func(*kingpin.ParseContext) error {
			build(branch, commit)
			return nil
		})
)

func build(branch, commit string) {
	buildBinaries(branch, commit)
}

func buildBinaries(branch, commit string) {
	for _, t := range targets {
		buildBinary(branch, commit, t, false)
		buildManual(branch, t)
	}
}

func buildBinary(branch, commit string, t target, forTesting bool) {
	ldFlags := buildLdFlagsFor(branch, commit, forTesting)
	outputName := t.executable()
	must(os.MkdirAll(filepath.Dir(outputName), 0755))
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch, "GO111MODULE=on")
	}, os.Stderr, os.Stdout, "go", "build", "-ldflags", ldFlags, "-o", outputName, "./main")
}

func buildLdFlagsFor(branch, commit string, forTesting bool) string {
	testPrefix := ""
	testSuffix := ""
	if forTesting {
		testPrefix = "TEST"
		testSuffix = "TEST"
	}
	return fmt.Sprintf(" -X github.com/echocat/caretakerd/app.version=%s%s%s", testPrefix, branch, testSuffix) +
		fmt.Sprintf(" -X github.com/echocat/caretakerd/app.revision=%s%s%s", testPrefix, commit, testSuffix) +
		fmt.Sprintf(" -X github.com/echocat/caretakerd/app.compiled=%s", startTime.Format("2006-01-02T15:04:05Z"))
}

func buildManual(branch string, t target) {
	outputName := "var/manual-builder" + t.executableExtension()
	must(os.MkdirAll(filepath.Dir(outputName), 0755))

	execute("go", "build", "-o", outputName, "./manual")

	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch, "GO111MODULE=on")
	}, os.Stderr, os.Stdout, outputName, branch, t.manual())
}
