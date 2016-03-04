package main

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"os"
	"path/filepath"
	"runtime"
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

func DeterminateProject(packageName string) (Project, error) {
	result, err := determinateProjectIn(GOPATH+"/src", packageName)
	if err != nil {
		return Project{}, err
	}
	if result == nil {
		result, err = determinateProjectIn(GOROOT+"/src", packageName)
		if err != nil {
			return Project{}, err
		}
	}
	if result == nil {
		if len(GOPATH) <= 0 {
			return Project{}, errors.New("'%v' is not contained in GOROOT(%v). Hint: Environment variable GOPATH is not set which could contain the package.", packageName, GOROOT)
		} else {
			return Project{}, errors.New("'%v' is neither a contained in GOPATH(%v) nor GOROOT(%v).", packageName, GOPATH, GOROOT)
		}
	}
	return *result, nil
}

func determinateProjectIn(goSrcPath string, packageName string) (*Project, error) {
	cleanGoSrcPath, err := filepath.Abs(goSrcPath)
	if err != nil {
		return nil, err
	}
	cleanSrcRootPath, err := filepath.Abs(goSrcPath + "/" + packageName)
	if err != nil {
		return nil, err
	}
	fileInfo, err := os.Stat(cleanSrcRootPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	if !fileInfo.IsDir() {
		return nil, nil
	}
	return &Project{
		GoSrcPath:   cleanGoSrcPath,
		SrcRootPath: cleanSrcRootPath,
		RootPackage: packageName,
	}, nil
}
