package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

// GOPATH points the the current GOPATH.
var GOPATH = os.Getenv("GOPATH")

// GOROOT points the the current GOROOT.
var GOROOT = runtime.GOROOT()
var GOROOTSRC, _ = filepath.Abs(GOROOT + "/src/")

// Project represents a Go project and its sources.
type Project struct {
	GoSrcPath   string
	SrcRootPath string
	RootPackage string
}

func (instance Project) String() string {
	b, _ := json.MarshalIndent(instance, "", "   ")
	return string(b)
}

// DeterminateProject determinate the Project for the given package name and returns it.
func DeterminateProject(packageName string) (Project, error) {
	cleanGoSrcPath, err := filepath.Abs(GOPATH + "/src")
	if err != nil {
		return Project{}, err
	}
	cleanSrcRootPath, err := filepath.Abs(".")
	if err != nil {
		return Project{}, err
	}

	return Project{
		GoSrcPath:   cleanGoSrcPath,
		SrcRootPath: cleanSrcRootPath,
		RootPackage: packageName,
	}, nil
}
