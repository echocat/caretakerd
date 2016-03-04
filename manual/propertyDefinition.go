package main

// PropertyDefinition represents a Definition that is a property.
type PropertyDefinition struct {
	id        IDType
	key       string
	valueType Type
	comment   string
	def       *string
}

// Id returns the ID of this Definition.
func (instance PropertyDefinition) Id() IDType {
	return instance.id
}

// Key returns the key of this Definition.
func (instance PropertyDefinition) Key() string {
	return instance.key
}

// ValueType returns the value type of this Definition.
func (instance PropertyDefinition) ValueType() Type {
	return instance.valueType
}

// Description returns the description of this Definition.
func (instance PropertyDefinition) Description() string {
	return instance.comment
}

// TypeName returns the type name of this Definition.
func (instance PropertyDefinition) TypeName() string {
	return "property"
}

// IsTopLevel returns true if this element is a top level Definition.
func (instance PropertyDefinition) IsTopLevel() bool {
	return false
}

// DefaultValue returns the default value of this Definition.
func (instance PropertyDefinition) DefaultValue() *string {
	return instance.def
}

func newPropertyDefinition(id IDType, key string, valueType Type, comment string, def *string) *PropertyDefinition {
	return &PropertyDefinition{
		id:        id,
		key:       key,
		valueType: valueType,
		comment:   comment,
		def:       def,
	}
}

func (instance PropertyDefinition) String() string {
	return FormatDefinition(&instance)
}
