package main

import (
	"github.com/echocat/caretakerd/errors"
	"go/parser"
	"go/token"
	"github.com/echocat/caretakerd/logger"
	"regexp"
	"go/ast"
	"go/types"
	"go/build"
	"runtime"
	"path/filepath"
	"fmt"
	"strings"
	"reflect"
)

var extractPropertyPattern = regexp.MustCompile("(?m)^\\s*@([a-zA-Z][a-zA-Z0-9]*)\\s+(.*)\\s*(:?\r\n|\n)")

type Api struct {
	Project Project
}

type PosEnabled interface {
	Pos() token.Pos
}

type parsedPackage struct {
	sourceFiles map[string]*ast.File
	pkg         *types.Package
	fileSet     *token.FileSet
}

func (instance *parsedPackage) FileFor(object PosEnabled) (*ast.File, error) {
	tokenFile := instance.fileSet.File(object.Pos())
	if tokenFile == nil {
		return nil, errors.New("Package %v does not contain object %v.", instance.pkg.Path(), object)
	}
	if file, ok := instance.sourceFiles[tokenFile.Name()]; ok {
		return file, nil
	}
	return nil, errors.New("Package %v does not contain file %v.", instance.pkg.Path(), tokenFile.Name())
}

func (instance *parsedPackage) CommentTextFor(object PosEnabled) (string, error) {
	comment, err := instance.CommentFor(object)
	if err != nil {
		return "", err
	}
	if comment != nil {
		return comment.Text(), nil
	}
	return "", nil
}

func (instance *parsedPackage) CommentFor(object PosEnabled) (*ast.CommentGroup, error) {
	file, err := instance.FileFor(object)
	if err != nil {
		return nil, err
	}
	pos := object.Pos()
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if spec.Pos() == pos {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						LOGGER.Log(logger.Info, "\t\t\t\t %v <-> %v: %v", reflect.TypeOf(object), reflect.TypeOf(spec), typeSpec)
						if typeSpec.Comment == nil && len(genDecl.Specs) == 1 {
							return genDecl.Doc, nil
						} else {
							return typeSpec.Comment, nil
						}
					} else if valueSpec, ok := spec.(*ast.ValueSpec); ok {
						if valueSpec.Comment == nil && len(genDecl.Specs) == 1 {
							return genDecl.Doc, nil
						} else {
							return valueSpec.Comment, nil
						}
					}
				} else {
					if sSpec, ok := spec.(*ast.TypeSpec); ok {
						if strctType, ok := sSpec.Type.(*ast.StructType); ok {
							for _, field := range strctType.Fields.List {
								if field.Pos() == pos {
									return field.Doc, nil
								}
							}
						}
					}
				}
			}
		}
	}
	return nil, nil
}

type extractionTask struct {
	info                       *types.Info
	project                    Project
	packageNameToParsedPackage map[string]*parsedPackage
	context                    *build.Context
}

func (instance *extractionTask) findDeclFor(object PosEnabled) (*ast.Decl, error) {
	return nil, nil
}

func (instance *extractionTask) parsePackage(packageName string) (*parsedPackage, error) {
	result, ok := instance.packageNameToParsedPackage[packageName]
	if instance.packageNameToParsedPackage == nil {
		instance.packageNameToParsedPackage = map[string]*parsedPackage{}
	}
	if !ok {
		sourceFiles := []*ast.File{}
		contextPackage, err := instance.context.Import(packageName, "", build.ImportComment)
		if err != nil {
			return nil, errors.New("Could not import package %v.", packageName).CausedBy(err)
		}
		result = &parsedPackage{
			sourceFiles: map[string]*ast.File{},
		}
		result.fileSet = token.NewFileSet()
		for _, goFile := range contextPackage.GoFiles {
			sourceFilename := fmt.Sprintf("%v%c%v", contextPackage.Dir, filepath.Separator, goFile)
			sourceFile, err := parser.ParseFile(result.fileSet, sourceFilename, nil, parser.ParseComments)
			if err != nil {
				return nil, errors.New("Could not parse source file %v.", sourceFilename).CausedBy(err)
			}
			sourceFiles = append(sourceFiles, sourceFile)
			result.sourceFiles[sourceFilename] = sourceFile
		}
		typesConfig := types.Config{
			Importer: instance,
			FakeImportC: true,
			DisableUnusedImportCheck: true,
			IgnoreFuncBodies: true,
		}
		if contextPackage.Dir == instance.project.SrcRootPath || strings.HasPrefix(contextPackage.Dir, instance.project.SrcRootPath + string([]byte{filepath.Separator})) {
			LOGGER.Log(logger.Info, "Check package %v...", packageName)
		}
		pkg, err := typesConfig.Check(packageName, result.fileSet, sourceFiles, instance.info)
		if err != nil {
			return nil, errors.New("Could not check package %v.", packageName).CausedBy(err)
		}
		result.pkg = pkg
		instance.packageNameToParsedPackage[packageName] = result
	}
	return result, nil
}

func (instance *extractionTask) Import(packageName string) (*types.Package, error) {
	pp, err := instance.parsePackage(packageName)
	if err != nil {
		return nil, err
	}
	return pp.pkg, nil
}

func ExtractApiFrom(project Project) (*Api, error) {
	api := &Api{
		Project: project,
	}

	et := &extractionTask{
		info: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		},
		project: project,
		context: &build.Context{
			GOARCH: runtime.GOARCH,
			GOOS: runtime.GOOS,
			GOROOT: GOROOT,
			GOPATH: GOPATH,
			Compiler: runtime.Compiler,
		},
	}
	pp, err := et.parsePackage(project.RootPackage + "/logger")
	if err != nil {
		return nil, err
	}
	LOGGER.Log(logger.Info, "Package: %v", pp.pkg.Path())
	scope := pp.pkg.Scope()
	for _, name := range scope.Names() {
		element := scope.Lookup(name)
		eUnderlying := element.Type().Underlying()
		file, err := pp.FileFor(element)
		if err != nil {
			return nil, err
		}
		if _, ok := element.(*types.TypeName); ok {
			if _, ok := eUnderlying.(*types.Basic); ok {
				comment, err := pp.CommentTextFor(element)
				if err != nil {
					return nil, err
				}
				LOGGER.Log(logger.Info, "  %v.%v: %v", pp.pkg.Path(), name, len(comment))
			} else if eStruct, ok := eUnderlying.(*types.Struct); ok {
				comment, err := pp.CommentTextFor(element)
				if err != nil {
					return nil, err
				}
				LOGGER.Log(logger.Info, "  %v.%v: %v", pp.pkg.Path(), name, len(comment))
				for n := 0; n < eStruct.NumFields(); n++ {
					field := eStruct.Field(n)
					tag := eStruct.Tag(n)
					comment, err := pp.CommentTextFor(field)
					if err != nil {
						return nil, err
					}
					LOGGER.Log(logger.Info, "  \t\t%v.%v %v: %v", pp.pkg.Path(), field.Name(), tag, len(comment))
				}
			} else {
				LOGGER.Log(logger.Info, "# %v(%v) %v", name, file.Name, reflect.TypeOf(element))
			}
		} else if eConst, ok := element.(*types.Const); ok {
			comment, err := pp.CommentTextFor(eConst)
			if err != nil {
				return nil, err
			}
			LOGGER.Log(logger.Info, "C %v.%v: %v", pp.pkg.Path(), name, comment)
		} else {
			LOGGER.Log(logger.Info, "@ %v(%v): %v", name, file.Name, reflect.TypeOf(element))
		}
	}
	return api, nil
}
