package main

type PropertyDefinition struct {
	id        IdType
	key       string
	valueType Type
	comment   string
	def       *string
}

func (instance PropertyDefinition) Id() IdType {
	return instance.id
}

func (instance PropertyDefinition) Key() string {
	return instance.key
}

func (instance PropertyDefinition) ValueType() Type {
	return instance.valueType
}

func (instance PropertyDefinition) Description() string {
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

func newPropertyDefinition(id IdType, key string, valueType Type, comment string, def *string) *PropertyDefinition {
	return &PropertyDefinition{
		id: id,
		key: key,
		valueType: valueType,
		comment: comment,
		def: def,
	}
}

func (instance PropertyDefinition) String() string {
	return FormatDefinition(&instance)
}
