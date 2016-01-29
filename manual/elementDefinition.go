package main

type ElementDefinition struct {
	identifier Identifier
	valueType  Identifier
	key        string
	comment    string
}

func (instance ElementDefinition) Identifier() Identifier {
	return instance.identifier
}

func (instance ElementDefinition) Key() string {
	return instance.key
}

func (instance ElementDefinition) ValueType() Identifier {
	return instance.valueType
}

func (instance ElementDefinition) Comment() string {
	return instance.comment
}

func (instance ElementDefinition) TypeName() string {
	return "element"
}

func (instance ElementDefinition) IsTopLevel() bool {
	return false
}

func newElementDefinition(identifier Identifier, key string, comment string) *ElementDefinition {
	return &ElementDefinition{
		identifier: identifier,
		key: key,
		comment: comment,
	}
}

func (instance ElementDefinition) String() string {
	return FormatDefinition(&instance)
}

