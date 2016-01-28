package main

import (
    "bytes"
    "fmt"
)

type Definitions struct {
    project              Project
    sourceToDefinition map[string]Definition
    targetToDefinition map[string]Definition
}

func NewDefinitions(project Project) *Definitions {
    return &Definitions{
        project: project,
        sourceToDefinition: map[string]Definition{},
        targetToDefinition: map[string]Definition{},
    }
}

func (instance *Definitions) NewSimpleDefinition(packageName string, name string, comment string) *SimpleDefinition {
    identifier := instance.newIdentifier(packageName, name, name)
    definition := newSimpleDefinition(identifier, comment)
    instance.add(definition)
    return definition
}

func (instance *Definitions) NewObjectDefinition(packageName string, name string, comment string, children ... Definition) *ObjectDefinition {
    identifier := instance.newIdentifier(packageName, name, name)
    definition := newObjectDefinition(identifier, comment, children...)
    instance.add(definition)
    return definition
}

func (instance *Definitions) NewEnumDefinition(packageName string, name string, comment string, children ... Definition) *EnumDefinition {
    identifier := instance.newIdentifier(packageName, name, name)
    definition := newEnumDefinition(identifier, comment, children...)
    instance.add(definition)
    return definition
}

func (instance *Definitions) NewPropertyDefinition(parent *ObjectDefinition, sourceName string, targetName string, comment string, def *string) *PropertyDefinition {
    identifier := instance.newIdentifier(
        parent.identifier.SourcePackage,
        parent.identifier.SourceName + "." + sourceName,
        parent.identifier.Name + "." + targetName,
    )
    definition := newPropertyDefinition(identifier, comment, def)
    parent.AdChild(definition)
    instance.add(definition)
    return definition
}

func (instance *Definitions) NewElementDefinition(parent *EnumDefinition, packageName string, sourceName string, targetName string, comment string) *ElementDefinition {
    identifier := instance.newIdentifier(
        parent.identifier.SourcePackage,
        parent.identifier.SourceName + "." + sourceName,
        parent.identifier.Name + "." + targetName,
    )
    definition := newElementDefinition(identifier, comment)
    parent.AdChild(definition)
    instance.add(definition)
    return definition
}

func (instance *Definitions) newIdentifier(packageName string, sourceName string, targetName string) Identifier {
    return NewIdentifier(instance.project, packageName, sourceName, targetName)
}

func (instance *Definitions) add(definition Definition) {
    identifier := definition.Identifier()
    instance.sourceToDefinition[identifier.AsSourceIdentifier()] = definition
    instance.targetToDefinition[identifier.AsTargetIdentifier()] = definition
}

func (instance *Definitions) AllTopLevel() []Definition {
    result := []Definition{}
    for _, definition := range instance.sourceToDefinition {
        if definition.IsTopLevel() {
            result = append(result, definition)
        }
    }
    return result
}

func (instance *Definitions) GetBy(identifier Identifier) Definition {
    return instance.sourceToDefinition[identifier.AsSourceIdentifier()]
}

func (instance *Definitions) GetBySource(identifier string) Definition {
    return instance.sourceToDefinition[identifier]
}

func (instance *Definitions) GetByTarget(identifier string) Definition {
    return instance.targetToDefinition[identifier]
}

func (instance Definitions) String() string {
    buf := new(bytes.Buffer)
    for _, definition := range instance.sourceToDefinition {
        if definition.IsTopLevel() {
            fmt.Fprintf(buf, "%v\n", definition)
        }
    }
    return buf.String()
}
