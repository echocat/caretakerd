package main

type Definition interface {
    Identifier() Identifier
    Comment() string
    TypeName() string
    IsTopLevel() bool
}

type DefinitionWithDefaultValue interface {
    Default() *string
}

type DefinitionWithChildren interface {
    Children() []Definition
    AddChild(child Definition)
}
