package main

import (
	"fmt"
	"github.com/echocat/caretakerd"
	"path/filepath"
	"runtime"
)

var (
	currentTarget = target{os: runtime.GOOS, arch: runtime.GOARCH}
	linuxAmd64    = target{os: "linux", arch: "amd64"}
	targets       = []target{
		{os: "windows", arch: "amd64"},
		{os: "windows", arch: "386"},
		{os: "darwin", arch: "amd64"},
		{os: "darwin", arch: "386"},
		linuxAmd64,
		{os: "linux", arch: "386"},
	}
)

type target struct {
	os   string
	arch string
}

func (instance target) outputName() string {
	return fmt.Sprintf(caretakerd.DaemonName+"-%s-%s", instance.os, instance.arch)
}

func (instance target) executable() string {
	return filepath.Join("var", "executables", instance.outputName()+instance.executableExtension())
}

func (instance target) manual() string {
	return filepath.Join("var", "manuals", instance.outputName()+".html")
}

func (instance target) archive() string {
	return filepath.Join("var", "dist", instance.outputName()+instance.archiveExtension())
}

func (instance target) executableExtension() string {
	if instance.os == "windows" {
		return ".exe"
	}
	return ""
}

func (instance target) archiveExtension() string {
	if instance.os == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}
