package main

import (
	"bytes"
	"fmt"
)

type Definitions struct {
	project                Project
	identifierToDefinition map[string]Definition
}

func NewDefinitions(project Project) *Definitions {
	return &Definitions{
		project: project,
		identifierToDefinition: map[string]Definition{},
	}
}

func (instance *Definitions) NewSimpleDefinition(packageName string, name string, comment string) *SimpleDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newSimpleDefinition(identifier, comment)
	instance.add(definition)
	return definition
}

func (instance *Definitions) NewObjectDefinition(packageName string, name string, comment string) *ObjectDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newObjectDefinition(identifier, comment)
	instance.add(definition)
	return definition
}

func (instance *Definitions) NewEnumDefinition(packageName string, name string, comment string) *EnumDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newEnumDefinition(identifier, comment)
	instance.add(definition)
	return definition
}

func (instance *Definitions) NewPropertyDefinition(parent *ObjectDefinition, name string, key string, comment string, def *string) *PropertyDefinition {
	identifier := instance.newIdentifierWithParent(parent, name)
	definition := newPropertyDefinition(identifier, key, comment, def)
	parent.AddChild(definition)
	instance.add(definition)
	return definition
}

func (instance *Definitions) NewElementDefinition(parent *EnumDefinition, name string, key string, comment string) *ElementDefinition {
	identifier := instance.newIdentifierWithParent(parent, name)
	definition := newElementDefinition(identifier, key, comment)
	parent.AddChild(definition)
	instance.add(definition)
	return definition
}

func (instance *Definitions) newIdentifier(packageName string, name string) Identifier {
	return NewIdentifier(packageName, name)
}

func (instance *Definitions) newIdentifierWithParent(parent Definition, name string) Identifier {
	parentIdentifier := parent.Identifier()
	return instance.newIdentifier(
		parentIdentifier.Package,
		parentIdentifier.Name + "." + name,
	)
}

func (instance *Definitions) add(definition Definition) {
	identifier := definition.Identifier()
	instance.identifierToDefinition[identifier.String()] = definition
}

func (instance *Definitions) AllTopLevel() []Definition {
	result := []Definition{}
	for _, definition := range instance.identifierToDefinition {
		if definition.IsTopLevel() {
			result = append(result, definition)
		}
	}
	return result
}

func (instance *Definitions) GetBy(identifier Identifier) Definition {
	return instance.GetByPlain(identifier.String())
}

func (instance *Definitions) GetByPlain(identifier string) Definition {
	return instance.identifierToDefinition[identifier]
}

func (instance Definitions) String() string {
	buf := new(bytes.Buffer)
	for _, definition := range instance.identifierToDefinition {
		if definition.IsTopLevel() {
			fmt.Fprintf(buf, "%v\n", definition)
		}
	}
	return buf.String()
}
