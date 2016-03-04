package main

// ObjectDefinition represents a Definition that is an object/struct.
type ObjectDefinition struct {
	id       IDType
	comment  string
	children []Definition
}

// ID returns the ID of this Definition.
func (instance ObjectDefinition) ID() IDType {
	return instance.id
}

// Description returns the description of this Definition.
func (instance ObjectDefinition) Description() string {
	return instance.comment
}

// TypeName returns the type name of this Definition.
func (instance ObjectDefinition) TypeName() string {
	return "object"
}

// IsTopLevel returns true if this element is a top level Definition.
func (instance ObjectDefinition) IsTopLevel() bool {
	return true
}

// Children returns every children of this Definition.
func (instance *ObjectDefinition) Children() []Definition {
	return instance.children
}

// AddChild adds a child to this Definition.
func (instance *ObjectDefinition) AddChild(child Definition) {
	instance.children = append(instance.children, child)
}

func newObjectDefinition(id IDType, comment string) *ObjectDefinition {
	return &ObjectDefinition{
		id:       id,
		comment:  comment,
		children: []Definition{},
	}
}

func (instance ObjectDefinition) String() string {
	return FormatDefinition(&instance)
}
