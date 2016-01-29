package main

type SimpleDefinition struct {
	identifier Identifier
	comment    string
}

func (instance SimpleDefinition) Identifier() Identifier {
	return instance.identifier
}

func (instance SimpleDefinition) Comment() string {
	return instance.comment
}

func (instance SimpleDefinition) TypeName() string {
	return "simple"
}

func (instance SimpleDefinition) IsTopLevel() bool {
	return true
}

func newSimpleDefinition(identifier Identifier, comment string) *SimpleDefinition {
	return &SimpleDefinition{
		identifier: identifier,
		comment: comment,
	}
}

func (instance SimpleDefinition) String() string {
	return FormatDefinition(&instance)
}
