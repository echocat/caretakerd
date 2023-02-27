package main

import (
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"os"
	"os/exec"
	"strings"
)

var (
	_ = app.Command("test", "executes tests for the project").
		Action(func(*kingpin.ParseContext) error {
			test(branch, commit)
			return nil
		})
)

func test(branch, commit string) {
	testGoCode(currentTarget)

	buildBinary(branch, commit, currentTarget, true)
	testBinary(branch, commit, currentTarget)
	buildManual(branch, &currentTarget)
}

func testGoCode(t target) {
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GOOS="+t.os, "GOARCH="+t.arch)
	}, os.Stderr, os.Stdout,
		"go", "test",
		"-v",
		//"-race", // TODO! We should improve that later.
		"-covermode", "atomic",
		"-coverprofile", "profile.cov",
		"./...",
	)
}

func testBinary(branch, commit string, t target) {
	testBinaryByExpectingResponse(t, `Version:      TEST`+branch+`TEST`, t.executable(), "version")
	testBinaryByExpectingResponse(t, `Git revision: TEST`+commit+`TEST`, t.executable(), "version")
}

func testBinaryByExpectingResponse(t target, expectedPartOfResponse string, args ...string) {
	cmd := append([]string{t.executable()}, args...)
	response := executeAndRecord(args...)
	if !strings.Contains(response, expectedPartOfResponse) {
		panic(fmt.Sprintf("Command failed [%s]\nResponse should contain: %s\nBut response was: %s",
			quoteAllIfNeeded(cmd...), expectedPartOfResponse, response))
	}
}
