package main

type EnumDefinition struct {
	id       IdType
	comment  string
	children []Definition
}

func (instance EnumDefinition) Id() IdType {
	return instance.id
}

func (instance EnumDefinition) Description() string {
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

func newEnumDefinition(id IdType, comment string) *EnumDefinition {
	return &EnumDefinition{
		id:       id,
		comment:  comment,
		children: []Definition{},
	}
}

func (instance EnumDefinition) String() string {
	return FormatDefinition(&instance)
}
