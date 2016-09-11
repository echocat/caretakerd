package main

// EnumDefinition represents a Definition that is an enum.
type EnumDefinition struct {
	id       IDType
	comment  string
	children []Definition
}

// ID returns the ID of this Definition.
func (instance EnumDefinition) ID() IDType {
	return instance.id
}

// Description returns the description of this Definition.
func (instance EnumDefinition) Description() string {
	return instance.comment
}

// TypeName returns the type name of this Definition.
func (instance EnumDefinition) TypeName() string {
	return "enum"
}

// IsTopLevel returns "true" if this element is a top level Definition.
func (instance EnumDefinition) IsTopLevel() bool {
	return true
}

// Children returns every children of this Definition.
func (instance *EnumDefinition) Children() []Definition {
	return instance.children
}

// AddChild adds a child to this Definition.
func (instance *EnumDefinition) AddChild(child Definition) {
	instance.children = append(instance.children, child)
}

func newEnumDefinition(id IDType, comment string) *EnumDefinition {
	return &EnumDefinition{
		id:       id,
		comment:  comment,
		children: []Definition{},
	}
}

func (instance EnumDefinition) String() string {
	return FormatDefinition(&instance)
}
