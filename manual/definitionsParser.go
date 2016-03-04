package main

import (
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"github.com/echocat/caretakerd/system"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var extractIDPropertyPattern = regexp.MustCompile("(?m)^\\s*@id\\s+(.*)\\s*(:?\r\n|\n)")
var extractDefaultPropertyPattern = regexp.MustCompile("(?m)^\\s*@default\\s+(.*)\\s*(:?\r\n|\n)")
var extractInlinePropertyPattern = regexp.MustCompile("(?m)^\\s*@inline\\s*(:?\r\n|\n)")
var extractSerializedAsPropertyPattern = regexp.MustCompile("(?m)^\\s*@serializedAs\\s+(.*)\\s*(:?\r\n|\n)")

type posEnabled interface {
	Pos() token.Pos
}

type parsedPackage struct {
	sourceFiles map[string]*ast.File
	pkg         *types.Package
	fileSet     *token.FileSet
}

func (instance *parsedPackage) fileFor(object posEnabled) (*ast.File, error) {
	tokenFile := instance.fileSet.File(object.Pos())
	if tokenFile == nil {
		return nil, errors.New("Package %v does not contain object %v.", instance.pkg.Path(), object)
	}
	if file, ok := instance.sourceFiles[tokenFile.Name()]; ok {
		return file, nil
	}
	return nil, errors.New("Package %v does not contain file %v.", instance.pkg.Path(), tokenFile.Name())
}

func (instance *parsedPackage) commentTextFor(object posEnabled) (string, error) {
	comment, err := instance.commentFor(object)
	if err != nil {
		return "", err
	}
	if comment != nil {
		return comment.Text(), nil
	}
	return "", nil
}

func (instance *parsedPackage) commentFor(object posEnabled) (*ast.CommentGroup, error) {
	file, err := instance.fileFor(object)
	object.Pos().IsValid()
	if err != nil {
		return nil, err
	}
	pos := object.Pos()
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if spec.Pos() == pos {
					if sSpec, ok := spec.(*ast.TypeSpec); ok {
						if sSpec.Comment == nil {
							if sSpec.Doc == nil && len(genDecl.Specs) == 1 {
								return genDecl.Doc, nil
							}
							return sSpec.Doc, nil
						}
						return sSpec.Comment, nil
					} else if sSpec, ok := spec.(*ast.ValueSpec); ok {
						if sSpec.Comment == nil {
							if sSpec.Doc == nil && len(genDecl.Specs) == 1 {
								return genDecl.Doc, nil
							}
							return sSpec.Doc, nil
						}
						return sSpec.Comment, nil
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

func (instance *extractionTask) findDeclFor(object posEnabled) (*ast.Decl, error) {
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
			if _, ok := err.(*build.NoGoError); ok {
				return nil, nil
			}
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
			Importer:                 instance,
			FakeImportC:              true,
			DisableUnusedImportCheck: true,
			IgnoreFuncBodies:         true,
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
	if pp == nil {
		return nil, errors.New("No go package.")
	}
	return pp.pkg, nil
}

// ParseDefinitions parses every Definitions from the given project and return it.
// If there is any error it will be returned and the Definitions are nil.
func ParseDefinitions(project Project) (*Definitions, error) {
	definitions := NewDefinitions(project)

	et := &extractionTask{
		info: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		},
		project: project,
		context: &build.Context{
			GOARCH:   runtime.GOARCH,
			GOOS:     runtime.GOOS,
			GOROOT:   GOROOT,
			GOPATH:   GOPATH,
			Compiler: runtime.Compiler,
		},
	}
	exclude := project.SrcRootPath + system.PATH_SEPARATOR + "target"
	err := filepath.Walk(project.SrcRootPath, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				// Ignore dot files and directories
				return nil
			} else if path == exclude || strings.HasPrefix(path, exclude+system.PATH_SEPARATOR) {
				// Do not try to build target directory
				return nil
			} else if path == project.SrcRootPath {
				return et.parsePackageToDefinitions(project.RootPackage, definitions)
			} else if strings.HasPrefix(path, project.GoSrcPath+system.PATH_SEPARATOR) {
				subPath := path[len(project.GoSrcPath)+1:]
				targetPackage := strings.Replace(subPath, system.PATH_SEPARATOR, "/", -1)
				err := et.parsePackageToDefinitions(targetPackage, definitions)
				if _, ok := err.(*build.NoGoError); ok {
					return nil
				}
				return err
			}
			panics.New("Unexpected path: %v", path).Throw()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return definitions, nil
}

func isBasic(what types.Type) bool {
	if _, ok := what.(*types.Basic); ok {
		return true
	} else if _, ok := what.(*types.Map); ok {
		return true
	} else if _, ok := what.(*types.Slice); ok {
		return true
	}
	return false
}

func (instance *extractionTask) parsePackageToDefinitions(pkg string, definitions *Definitions) error {
	pp, err := instance.parsePackage(pkg)
	if err != nil {
		return err
	}
	if pp == nil {
		return nil
	}
	scope := pp.pkg.Scope()
	for _, name := range scope.Names() {
		element := scope.Lookup(name)
		eUnderlying := element.Type().Underlying()
		if _, ok := element.(*types.TypeName); ok {
			if isBasic(eUnderlying) {
				comment, err := pp.commentTextFor(element)
				if err != nil {
					return err
				}

				file, err := pp.fileFor(element)
				if err != nil {
					return err
				}
				pos := element.Pos()
				decls := file.Decls

				var enumDefinition *EnumDefinition
				for c, decl := range decls {
					if genDecl, ok := decl.(*ast.GenDecl); ok {
						for _, spec := range genDecl.Specs {
							if spec.Pos() == pos {
								for n := c + 1; n < len(decls); n++ {
									nextDecl := decls[n]
									if nextGenDecl, ok := nextDecl.(*ast.GenDecl); ok {
										if nextGenDecl.Tok == token.CONST {
											for _, cSpec := range nextGenDecl.Specs {
												for _, cScopeName := range scope.Names() {
													cScope := scope.Lookup(cScopeName)
													if cScope.Pos() == cSpec.Pos() {
														if eConst, ok := cScope.(*types.Const); ok {
															if eConst.Type().String() == pp.pkg.Path()+"."+name {
																elementComment, err := pp.commentTextFor(eConst)
																if err != nil {
																	return err
																}
																if enumDefinition == nil {
																	_, inlined := extractInlinedFrom(comment)
																	if inlined {
																		break
																	}
																	enumDefinition = definitions.NewEnumDefinition(pp.pkg.Path(), name, comment)
																}
																typeIdentifier := ParseType(eConst.Type().String())
																elementComment, id := extractIDFrom(elementComment, eConst.Name())
																definitions.NewElementDefinition(enumDefinition, eConst.Name(), id, typeIdentifier, elementComment)
															} else {
																break
															}
														}
													}
												}
											}
										} else {
											break
										}
									} else {
										break
									}
								}
							}
						}
					}
				}

				if enumDefinition == nil {
					typeIdentifier := ParseType(eUnderlying.Underlying().String())
					comment, inlined := extractInlinedFrom(comment)
					definitions.NewSimpleDefinition(pp.pkg.Path(), name, typeIdentifier, comment, inlined)
				}
			} else if eStruct, ok := eUnderlying.(*types.Struct); ok {
				comment, err := pp.commentTextFor(element)
				if err != nil {
					return err
				}
				comment, serializedAs := serializedAs(comment)
				if serializedAs != nil {
					comment, inlined := extractInlinedFrom(comment)
					definitions.NewSimpleDefinition(pp.pkg.Path(), name, serializedAs, comment, inlined)
				} else {
					objectDefinition := definitions.NewObjectDefinition(pp.pkg.Path(), name, comment)
					for n := 0; n < eStruct.NumFields(); n++ {
						field := eStruct.Field(n)
						tag := eStruct.Tag(n)
						targetName := fieldNameFor(field.Name(), tag)
						comment, err := pp.commentTextFor(field)
						if err != nil {
							return err
						}
						typeIdentifier := ParseType(field.Type().String())
						comment, defValue := extractDefaultFrom(comment)
						definitions.NewPropertyDefinition(objectDefinition, field.Name(), targetName, typeIdentifier, comment, defValue)
					}
				}
			}
		}
	}
	return nil
}

func fieldNameFor(name string, tag string) string {
	st := reflect.StructTag(tag)
	yaml := st.Get("yaml")
	if len(yaml) > 0 {
		return strings.SplitN(yaml, ",", 2)[0]
	}
	json := st.Get("json")
	if len(json) > 0 {
		return strings.SplitN(json, ",", 2)[0]
	}
	return name
}

func extractIDFrom(comment string, fallbackID string) (string, string) {
	matches := extractIDPropertyPattern.FindAllStringSubmatch(comment, -1)
	if len(matches) > 0 {
		id := strings.TrimSpace(matches[0][1])
		return extractIDPropertyPattern.ReplaceAllString(comment, ""), id
	}
	return comment, fallbackID
}

func extractDefaultFrom(comment string) (string, *string) {
	matches := extractDefaultPropertyPattern.FindAllStringSubmatch(comment, -1)
	if len(matches) > 0 {
		defValue := strings.TrimSpace(matches[0][1])
		return extractDefaultPropertyPattern.ReplaceAllString(comment, ""), &defValue
	}
	return comment, nil
}

func extractInlinedFrom(comment string) (string, bool) {
	matches := extractInlinePropertyPattern.FindAllStringSubmatch(comment, -1)
	if len(matches) > 0 {
		return extractInlinePropertyPattern.ReplaceAllString(comment, ""), true
	}
	return comment, false
}

func serializedAs(comment string) (string, Type) {
	matches := extractSerializedAsPropertyPattern.FindAllStringSubmatch(comment, -1)
	if len(matches) > 0 {
		plainType := strings.TrimSpace(matches[0][1])
		t := ParseType(plainType)
		return extractSerializedAsPropertyPattern.ReplaceAllString(comment, ""), t
	}
	return comment, nil
}
