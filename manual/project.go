package main

import (
	"encoding/json"
	"go/build"
	"os"
	"path/filepath"
)

// GOPATH points the current GOPATH.
var GOPATH = func() string {
	if v := os.Getenv("GOPATHX"); v != "" {
		return v
	}
	if v := os.Getenv("GOPATH"); v != "" {
		return v
	}
	return build.Default.GOPATH
}()

// GOROOT points the the current GOROOT.
var GOROOT = func() string {
	if v := os.Getenv("GOROOTX"); v != "" {
		return v
	}
	if v := os.Getenv("GOROOT"); v != "" {
		return v
	}
	return build.Default.GOROOT
}()
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
