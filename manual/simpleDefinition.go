package main

type SimpleDefinition struct {
	id      IdType
	comment string
}

func (instance SimpleDefinition) Id() IdType {
	return instance.id
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

func newSimpleDefinition(id IdType, comment string) *SimpleDefinition {
	return &SimpleDefinition{
		id: id,
		comment: comment,
	}
}

func (instance SimpleDefinition) String() string {
	return FormatDefinition(&instance)
}
