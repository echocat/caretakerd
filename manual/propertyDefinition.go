package main

type PropertyDefinition struct {
	identifier Identifier
	key        string
	valueType  Identifier
	comment    string
	def        *string
}

func (instance PropertyDefinition) Identifier() Identifier {
	return instance.identifier
}

func (instance PropertyDefinition) Key() string {
	return instance.key
}

func (instance PropertyDefinition) ValueType() Identifier {
	return instance.valueType
}

func (instance PropertyDefinition) Comment() string {
	return instance.comment
}

func (instance PropertyDefinition) TypeName() string {
	return "property"
}

func (instance PropertyDefinition) IsTopLevel() bool {
	return false
}

func (instance PropertyDefinition) DefaultValue() *string {
	return instance.def
}

func newPropertyDefinition(identifier Identifier, key string, comment string, def *string) *PropertyDefinition {
	return &PropertyDefinition{
		identifier: identifier,
		key: key,
		comment: comment,
		def: def,
	}
}

func (instance PropertyDefinition) String() string {
	return FormatDefinition(&instance)
}
