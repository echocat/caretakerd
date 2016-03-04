package main

// SimpleDefinition represents a simple Definition.
type SimpleDefinition struct {
	id        IDType
	valueType Type
	comment   string
	inlined   bool
}

// ID returns the ID of this Definition.
func (instance SimpleDefinition) ID() IDType {
	return instance.id
}

// ValueType returns the value type of this Definition.
func (instance SimpleDefinition) ValueType() Type {
	return instance.valueType
}

// Description returns the description of this Definition.
func (instance SimpleDefinition) Description() string {
	return instance.comment
}

// Inlined returns true if this Definition should be inlined.
func (instance SimpleDefinition) Inlined() bool {
	return instance.inlined
}

// TypeName returns the type name of this Definition.
func (instance SimpleDefinition) TypeName() string {
	return "simple"
}

// IsTopLevel returns true if this element is a top level Definition.
func (instance SimpleDefinition) IsTopLevel() bool {
	return true
}

func newSimpleDefinition(id IDType, valueType Type, comment string, inlined bool) *SimpleDefinition {
	return &SimpleDefinition{
		id:        id,
		valueType: valueType,
		comment:   comment,
		inlined:   inlined,
	}
}

func (instance SimpleDefinition) String() string {
	return FormatDefinition(&instance)
}
