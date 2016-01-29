package main

type ObjectDefinition struct {
	identifier Identifier
	comment    string
	children   []Definition
}

func (instance ObjectDefinition) Identifier() Identifier {
	return instance.identifier
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

func newObjectDefinition(identifier Identifier, comment string) *ObjectDefinition {
	return &ObjectDefinition{
		identifier: identifier,
		comment: comment,
		children: []Definition{},
	}
}

func (instance ObjectDefinition) String() string {
	return FormatDefinition(&instance)
}

