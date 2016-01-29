package main

type ObjectDefinition struct {
	id       IdType
	comment  string
	children []Definition
}

func (instance ObjectDefinition) Id() IdType {
	return instance.id
}

func (instance ObjectDefinition) Comment() string {
	return instance.comment
}

func (instance ObjectDefinition) TypeName() string {
	return "object"
}

func (instance ObjectDefinition) IsTopLevel() bool {
	return true
}

func (instance *ObjectDefinition) Children() []Definition {
	return instance.children
}

func (instance *ObjectDefinition) AddChild(child Definition) {
	instance.children = append(instance.children, child)
}

func newObjectDefinition(id IdType, comment string) *ObjectDefinition {
	return &ObjectDefinition{
		id: id,
		comment: comment,
		children: []Definition{},
	}
}

func (instance ObjectDefinition) String() string {
	return FormatDefinition(&instance)
}

