package main

type EnumDefinition struct {
	identifier Identifier
	comment    string
	children   []Definition
}

func (instance EnumDefinition) Identifier() Identifier {
	return instance.identifier
}

func (instance EnumDefinition) Comment() string {
	return instance.comment
}

func (instance EnumDefinition) TypeName() string {
	return "enum"
}

func (instance EnumDefinition) IsTopLevel() bool {
	return true
}

func (instance *EnumDefinition) Children() []Definition {
	return instance.children
}

func (instance *EnumDefinition) AddChild(child Definition) {
	instance.children = append(instance.children, child)
}

func newEnumDefinition(identifier Identifier, comment string) *EnumDefinition {
	return &EnumDefinition{
		identifier: identifier,
		comment: comment,
		children: []Definition{},
	}
}

func (instance EnumDefinition) String() string {
	return FormatDefinition(&instance)
}

