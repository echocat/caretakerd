package main

import (
	"bytes"
	"fmt"
)

// Definitions contains all Definitions of the given project.
type Definitions struct {
	project                Project
	identifierToDefinition map[string]Definition
}

// NewDefinitions creates a new instance of Definitions for the given project.
func NewDefinitions(project Project) *Definitions {
	return &Definitions{
		project:                project,
		identifierToDefinition: map[string]Definition{},
	}
}

// NewSimpleDefinition creates a new SimpleDefinition from the given parameters and append it to the current instance.
func (instance *Definitions) NewSimpleDefinition(packageName string, name string, valueType Type, comment string, inlined bool) *SimpleDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newSimpleDefinition(identifier, valueType, comment, inlined)
	instance.add(definition)
	return definition
}

// NewObjectDefinition creates a new ObjectDefinition from the given parameters and append it to the current instance.
func (instance *Definitions) NewObjectDefinition(packageName string, name string, comment string) *ObjectDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newObjectDefinition(identifier, comment)
	instance.add(definition)
	return definition
}

// NewEnumDefinition creates a new EnumDefinition from the given parameters and append it to the current instance.
func (instance *Definitions) NewEnumDefinition(packageName string, name string, comment string) *EnumDefinition {
	identifier := instance.newIdentifier(packageName, name)
	definition := newEnumDefinition(identifier, comment)
	instance.add(definition)
	return definition
}

// NewPropertyDefinition creates a new PropertyDefinition from the given parameters and append it to the current instance.
func (instance *Definitions) NewPropertyDefinition(parent *ObjectDefinition, name string, key string, valueType Type, comment string, def *string) *PropertyDefinition {
	id := instance.newIDWithParent(parent, name)
	definition := newPropertyDefinition(id, key, valueType, comment, def)
	parent.AddChild(definition)
	instance.add(definition)
	return definition
}

// NewElementDefinition creates a new ElementDefinition from the given parameters and append it to the current instance.
func (instance *Definitions) NewElementDefinition(parent *EnumDefinition, name string, key string, valueType Type, comment string) *ElementDefinition {
	identifier := instance.newIDWithParent(parent, name)
	definition := newElementDefinition(identifier, key, valueType, comment)
	parent.AddChild(definition)
	instance.add(definition)
	return definition
}

func (instance *Definitions) newIdentifier(packageName string, name string) IDType {
	return NewIDType(packageName, name, false)
}

func (instance *Definitions) newIDWithParent(parent Definition, name string) IDType {
	parentIdentifier := parent.ID()
	return instance.newIdentifier(
		parentIdentifier.Package,
		parentIdentifier.Name+"#"+name,
	)
}

func (instance *Definitions) add(definition Definition) {
	identifier := definition.ID()
	instance.identifierToDefinition[identifier.String()] = definition
}

// AllTopLevel returns all definitions that are at the top level and are no children of another Definition.
func (instance *Definitions) AllTopLevel() []Definition {
	result := []Definition{}
	for _, definition := range instance.identifierToDefinition {
		if definition.IsTopLevel() {
			result = append(result, definition)
		}
	}
	return result
}

// GetBy returns a Definition by the given identifier or nil if it does not exist.
func (instance *Definitions) GetBy(identifier IDType) Definition {
	return instance.GetByPlain(identifier.String())
}

// GetByPlain returns a Definition by the given identifier or nil if it does not exist.
func (instance *Definitions) GetByPlain(identifier string) Definition {
	return instance.identifierToDefinition[identifier]
}

func (instance Definitions) String() string {
	buf := new(bytes.Buffer)
	for _, definition := range instance.identifierToDefinition {
		if definition.IsTopLevel() {
			_, _ = fmt.Fprintf(buf, "%v\n", definition)
		}
	}
	return buf.String()
}
