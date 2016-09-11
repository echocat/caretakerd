package main

// ElementDefinition represents a Definition that is an element.
type ElementDefinition struct {
	id        IDType
	valueType Type
	key       string
	comment   string
}

// ID returns the ID of this Definition.
func (instance ElementDefinition) ID() IDType {
	return instance.id
}

// Key returns the key of this Definition.
func (instance ElementDefinition) Key() string {
	return instance.key
}

// ValueType returns the value type of this Definition.
func (instance ElementDefinition) ValueType() Type {
	return instance.valueType
}

// Description returns the description of this Definition.
func (instance ElementDefinition) Description() string {
	return instance.comment
}

// TypeName returns the type name of this Definition.
func (instance ElementDefinition) TypeName() string {
	return "element"
}

// IsTopLevel returns "true" if this element is a top level Definition.
func (instance ElementDefinition) IsTopLevel() bool {
	return false
}

func newElementDefinition(id IDType, key string, valueType Type, comment string) *ElementDefinition {
	return &ElementDefinition{
		id:        id,
		key:       key,
		valueType: valueType,
		comment:   comment,
	}
}

func (instance ElementDefinition) String() string {
	return FormatDefinition(&instance)
}
