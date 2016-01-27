package main

import (
	"os"
	"runtime"
	"path/filepath"
	"github.com/echocat/caretakerd/panics"
	"strings"
	"encoding/json"
)

var GOPATH = os.Getenv("GOPATH")
var GOROOT = runtime.GOROOT()

type Project struct {
	GoSrcPath   string
	SrcRootPath string
	RootPackage string
}

func (instance Project) String() string {
	b, _ := json.MarshalIndent(instance, "", "   ")
	return string(b)
}

func DeterminateProject(srcRootPath string) Project {
	result := determinateProjectIn(GOPATH + "/src", srcRootPath)
	if result == nil {
		result = determinateProjectIn(GOROOT + "/src", srcRootPath)
	}
	if result == nil {
		if len(GOPATH) > 0 {
			panics.New("'%v' is not a subpath of GOROOT(%v). Hint: But environment variable GOPATH is not set.", srcRootPath, GOROOT).Throw()
		} else {
			panics.New("'%v' is neither a subpath of GOPATH(%v) nor GOROOT(%v).", srcRootPath, GOPATH, GOROOT).Throw()
		}
	}
	return *result
}

func determinateProjectIn(goSrcPath string, srcRootPath string) *Project {
	cleanGoPath, err := filepath.Abs(goSrcPath)
	if err != nil {
		return nil
	}
	cleanSrcRootPath, err := filepath.Abs(srcRootPath)
	if err != nil {
		panics.New("Could not make srcRootPath '%v' absolute.", srcRootPath).CausedBy(err).Throw()
	}
	if strings.HasPrefix(cleanSrcRootPath, cleanGoPath) && len(cleanSrcRootPath) + 1 > len(cleanGoPath) {
		rootPackage := cleanSrcRootPath[len(cleanGoPath) + 1:]
		return &Project{
			GoSrcPath: cleanGoPath,
			SrcRootPath: cleanSrcRootPath,
			RootPackage: rootPackage,
		}
	}
	return nil
}

