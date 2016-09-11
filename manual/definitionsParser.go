package main

import (
	"github.com/echocat/caretakerd/errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"go/importer"
	"os"
	"github.com/echocat/caretakerd/system"
	"github.com/echocat/caretakerd/panics"
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
	defaultImporter            types.Importer
	info                       *types.Info
	project                    Project
	packageNameToParsedPackage map[string]*parsedPackage
	context                    *build.Context
	importer                   types.ImporterFrom
	definitions                *Definitions
	typesConfig                types.Config
}

func (instance *extractionTask) findDeclFor(object posEnabled) (*ast.Decl, error) {
	return nil, nil
}

func (instance *extractionTask) Import(packageName string) (*types.Package, error) {
	panic("Should not be called.")
}

func (instance *extractionTask) packageExistsInGoRootSrc(packageName string) bool {
	targetPath := filepath.Join(GOROOTSRC, packageName)
	info, err := os.Lstat(targetPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func (instance *extractionTask) ImportFrom(packageName, packageSource string, mode types.ImportMode) (*types.Package, error) {
	pp, has := instance.packageNameToParsedPackage[packageName]
	if has {
		return pp.pkg, nil
	}

	if instance.packageExistsInGoRootSrc(packageName) {
		pkg, err := instance.defaultImporter.(types.ImporterFrom).ImportFrom(packageName, packageSource, mode)
		if err == nil {
			instance.packageNameToParsedPackage[packageName] = &parsedPackage{
				fileSet: token.NewFileSet(),
				sourceFiles: make(map[string]*ast.File),
				pkg: pkg,
			}
		}
		return pkg, err
	}

	buildPkg, err := build.Import(packageName, packageSource, build.ImportComment)
	if err != nil {
		return nil, err
	}

	pp = &parsedPackage{
		fileSet: token.NewFileSet(),
		sourceFiles: make(map[string]*ast.File),
	}

	var astFiles []*ast.File
	astFiles, pp.sourceFiles, err = parseAstFiles(pp.fileSet, buildPkg.Dir, buildPkg.GoFiles)
	if err != nil {
		return nil, err
	}

	pp.pkg, err = instance.typesConfig.Check(buildPkg.ImportPath, pp.fileSet, astFiles, instance.info)
	if err != nil {
		return nil, err
	}

	instance.packageNameToParsedPackage[packageName] = pp
	return pp.pkg, nil
}

// ParseAstFiles is a shortcut to parse files from a directory into a set of ast.Files.
func parseAstFiles(fset *token.FileSet, dir string, files []string) (astFiles []*ast.File, sourceFiles map[string]*ast.File, err error) {
	sourceFiles = make(map[string]*ast.File)
	for _, filename := range files {
		var afile *ast.File
		fullFilename := filepath.Join(dir, filename)
		afile, err = parser.ParseFile(fset, fullFilename, nil, parser.ParseComments)
		if err != nil {
			return
		}
		sourceFiles[fullFilename] = afile
		astFiles = append(astFiles, afile)
	}
	return
}

// ParseDefinitions parses every Definitions from the given project and returns it.
// If there is any error it will be returned and the Definitions are nil.
func ParseDefinitions(project Project) (*Definitions, error) {
	et := &extractionTask{
		defaultImporter: importer.Default(),
		info: &types.Info{
			Types: make(map[ast.Expr]types.TypeAndValue),
			Defs:  make(map[*ast.Ident]types.Object),
			Uses:  make(map[*ast.Ident]types.Object),
		},
		packageNameToParsedPackage: make(map[string]*parsedPackage),
		project: project,
		context: &build.Context{
			GOARCH:   runtime.GOARCH,
			GOOS:     runtime.GOOS,
			GOROOT:   GOROOT,
			GOPATH:   GOPATH,
			Compiler: runtime.Compiler,
		},
		definitions: NewDefinitions(project),
		typesConfig: types.Config{
			FakeImportC:              true,
			DisableUnusedImportCheck: true,
			IgnoreFuncBodies:         true,
		},
	}
	et.typesConfig.Importer = et

	exclude1 := project.SrcRootPath + system.PathSeparator + "target"
	err := filepath.Walk(project.SrcRootPath, func(path string, info os.FileInfo, err error) error {
		if info != nil && info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				// Ignore dot files and directories
				return nil
			} else if path == exclude1 || strings.HasPrefix(path, exclude1 + system.PathSeparator) {
				// Do not try to build target directory
				return nil
			} else if path == project.SrcRootPath {
				return et.parsePackageToDefinitions(project.RootPackage, project.SrcRootPath)
			} else if strings.HasPrefix(path, project.GoSrcPath+system.PathSeparator) {
				subPath := path[len(project.GoSrcPath)+1:]
				targetPackage := strings.Replace(subPath, system.PathSeparator, "/", -1)
				err := et.parsePackageToDefinitions(targetPackage, path)
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

	return et.definitions, nil
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

func (instance *extractionTask) parsePackageToDefinitions(packageName string, packageSource string) error {
	_, err := instance.ImportFrom(packageName, packageSource, 0)
	if err != nil {
		return err
	}
	pp, ok := instance.packageNameToParsedPackage[packageName]
	if !ok {
		return errors.New("Parsed but not found!?")
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
																	enumDefinition = instance.definitions.NewEnumDefinition(pp.pkg.Path(), name, comment)
																}
																typeIdentifier := ParseType(eConst.Type().String())
																elementComment, id := extractIDFrom(elementComment, eConst.Name())
																instance.definitions.NewElementDefinition(enumDefinition, eConst.Name(), id, typeIdentifier, elementComment)
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
					instance.definitions.NewSimpleDefinition(pp.pkg.Path(), name, typeIdentifier, comment, inlined)
				}
			} else if eStruct, ok := eUnderlying.(*types.Struct); ok {
				comment, err := pp.commentTextFor(element)
				if err != nil {
					return err
				}
				comment, serializedAs := serializedAs(comment)
				if serializedAs != nil {
					comment, inlined := extractInlinedFrom(comment)
					instance.definitions.NewSimpleDefinition(pp.pkg.Path(), name, serializedAs, comment, inlined)
				} else {
					objectDefinition := instance.definitions.NewObjectDefinition(pp.pkg.Path(), name, comment)
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
						instance.definitions.NewPropertyDefinition(objectDefinition, field.Name(), targetName, typeIdentifier, comment, defValue)
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
