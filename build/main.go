package main

import (
	"github.com/alecthomas/kingpin/v2"
	"os"
)

var (
	app = kingpin.New("build", "helps to build caretakerd").
		Interspersed(false)

	branch = "snapshot"
	commit = "unknown"
)

func init() {
	app.Flag("branch", "something like either main, v1.2.3 or snapshot-feature-foo").
		Required().
		Envar("GITHUB_REF_NAME").
		StringVar(&branch)
	app.Flag("commit", "something like 463e189796d5e96a7b605ab51985458faf8fd0d4").
		Required().
		Envar("GITHUB_SHA").
		StringVar(&commit)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}
