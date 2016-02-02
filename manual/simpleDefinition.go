package main

type SimpleDefinition struct {
	id        IdType
	valueType Type
	comment   string
	inlined   bool
}

func (instance SimpleDefinition) Id() IdType {
	return instance.id
}

func (instance SimpleDefinition) ValueType() Type {
	return instance.valueType
}

func (instance SimpleDefinition) Description() string {
	return instance.comment
}

func (instance SimpleDefinition) Inlined() bool {
	return instance.inlined
}

func (instance SimpleDefinition) TypeName() string {
	return "simple"
}

func (instance SimpleDefinition) IsTopLevel() bool {
	return true
}

func newSimpleDefinition(id IdType, valueType Type, comment string, inlined bool) *SimpleDefinition {
	return &SimpleDefinition{
		id: id,
		valueType: valueType,
		comment: comment,
		inlined: inlined,
	}
}

func (instance SimpleDefinition) String() string {
	return FormatDefinition(&instance)
}
