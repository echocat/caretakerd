package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"github.com/alecthomas/kingpin"
	"github.com/echocat/caretakerd"
	"github.com/echocat/caretakerd/logger"
	"io"
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
	buildTarget(branch, commit)
	buildManual(branch, nil)
}

func buildTarget(branch, commit string) {
	for _, t := range targets {
		buildBinary(branch, commit, t, false)
		buildManual(branch, &t)
		buildPackage(t)
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

func buildManual(branch string, t *target) {
	outputName := filepath.Join("var", "dist", caretakerd.DaemonName+".html")
	platform := "linux"
	if t != nil {
		outputName = t.manual()
		platform = t.os
	}
	executeTo(func(cmd *exec.Cmd) {
		cmd.Env = append(os.Environ(), "GO111MODULE=on")
	}, os.Stderr, os.Stdout, "go", "run", "./manual", branch, platform, outputName)
}

func buildPackage(t target) {
	outputName := t.archive()
	log.Log(logger.Info, "Build package: %s", outputName)
	must(os.MkdirAll(filepath.Dir(outputName), 0755))
	switch t.archiveExtension() {
	case ".tar.gz":
		buildTarGzPackage(t, outputName)
	case ".zip":
		buildZipPackage(t, outputName)
	default:
		panic(fmt.Sprintf("Unsupported archive extension: %s", t.archiveExtension()))
	}
}

func buildTarGzPackage(t target, outputName string) {
	f, err := os.Create(outputName)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()
	gw, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer gw.Close()
	tw := tar.NewWriter(gw)
	//noinspection GoUnhandledErrorResult
	defer tw.Close()

	addFileToTar(t.executable(), caretakerd.DaemonName+t.executableExtension(), 0755, tw)
	addLinkToTar(caretakerd.DaemonName+t.executableExtension(), caretakerd.ControlName+t.executableExtension(), 0755, tw)
	addFileToTar(t.manual(), caretakerd.DaemonName+".html", 0644, tw)
}

func addFileToTar(sourceFile string, targetPath string, mode os.FileMode, to *tar.Writer) {
	f, err := os.Open(sourceFile)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()
	fi, err := f.Stat()
	must(err)
	must(to.WriteHeader(&tar.Header{
		Name:    targetPath,
		Size:    fi.Size(),
		Mode:    int64(mode),
		ModTime: fi.ModTime(),
	}))
	_, err = io.Copy(to, f)
	must(err)
}

func addLinkToTar(sourcePath string, targetPath string, mode os.FileMode, to *tar.Writer) {
	must(to.WriteHeader(&tar.Header{
		Typeflag: tar.TypeSymlink,
		Name:     targetPath,
		Linkname: sourcePath,
		Mode:     int64(mode),
	}))
}

func buildZipPackage(t target, outputName string) {
	f, err := os.Create(outputName)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()
	zw := zip.NewWriter(f)
	//noinspection GoUnhandledErrorResult
	defer zw.Close()

	addFileToZip(t.executable(), caretakerd.DaemonName+t.executableExtension(), 0755, zw)
	addFileToZip(t.executable(), caretakerd.ControlName+t.executableExtension(), 0755, zw)
	addFileToZip(t.manual(), caretakerd.DaemonName+".html", 0644, zw)
}

func addFileToZip(sourceFile string, targetPath string, mode os.FileMode, to *zip.Writer) {
	f, err := os.Open(sourceFile)
	must(err)
	//noinspection GoUnhandledErrorResult
	defer f.Close()
	fi, err := f.Stat()
	must(err)
	header, err := zip.FileInfoHeader(fi)
	must(err)
	header.SetMode(mode)
	header.Method = zip.Deflate
	header.Name = targetPath
	ew, err := to.CreateHeader(header)
	must(err)
	_, err = io.Copy(ew, f)
	must(err)
}
