package main

type ElementDefinition struct {
	id        IdType
	valueType Type
	key       string
	comment   string
}

func (instance ElementDefinition) Id() IdType {
	return instance.id
}

func (instance ElementDefinition) Key() string {
	return instance.key
}

func (instance ElementDefinition) ValueType() Type {
	return instance.valueType
}

func (instance ElementDefinition) Description() string {
	return instance.comment
}

func (instance ElementDefinition) TypeName() string {
	return "element"
}

func (instance ElementDefinition) IsTopLevel() bool {
	return false
}

func newElementDefinition(id IdType, key string, valueType Type, comment string) *ElementDefinition {
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
