package main

import (
    "strings"
    "path"
)

type Identifier struct {
    Package          string
    Name             string
    AtomicName       string

    SourcePackage    string
    SourceName       string
    SourceAtomicName string
}

func (instance Identifier) AsSourceIdentifier() string {
    result := ""
    if len(instance.SourcePackage) > 0 {
        result += instance.SourcePackage + "."
    }
    result += instance.SourceName
    return result
}

func (instance Identifier) AsTargetIdentifier() string {
    result := ""
    if len(instance.Package) > 0 {
        result += instance.Package + "."
    }
    result += instance.Name
    return result
}

func (instance Identifier) String() string {
    return instance.AsSourceIdentifier() + " > " + instance.AsTargetIdentifier()
}

func NewIdentifier(project Project, packageName string, sourceName string, targetName string) Identifier {
    sourcePackage := packageName
    targetPackage := sourcePackage
    if targetName == capitalize(path.Base(sourcePackage)) {
        targetName = "_" + targetName
    }
    if targetName == "Config" {
        targetName = capitalize(path.Base(sourcePackage))
    }

    if sourcePackage == project.RootPackage {
        targetPackage = ""
    } else if strings.HasPrefix(sourcePackage, project.RootPackage) {
        targetPackage = sourcePackage[len(project.RootPackage) + 1:]
    }

    sourceNameParts := strings.SplitAfter(sourceName, ".")
    targetNameParts := strings.SplitAfter(targetName, ".")

    return Identifier{
        SourcePackage: sourcePackage,
        SourceName: sourceName,
        SourceAtomicName: sourceNameParts[len(sourceNameParts) - 1],

        Package: targetPackage,
        Name: targetName,
        AtomicName: targetNameParts[len(targetNameParts) - 1],
    }
}

func capitalize(in string) string {
    l := len(in)
    if l <= 0 {
        return ""
    } else if l == 1 {
        return strings.ToUpper(in)
    } else {
        return strings.ToUpper(in[0:1]) + in[1:]
    }
}
