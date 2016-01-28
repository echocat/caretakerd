package main

import (
    "bytes"
    "fmt"
"strings"
)

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

func (instance EnumDefinition) Children() []Definition {
    return instance.children
}

func (instance *EnumDefinition) AdChild(child Definition) {
    instance.children = append(instance.children, child)
}

func newEnumDefinition(identifier Identifier, comment string, children ... Definition) *EnumDefinition {
    return &EnumDefinition{
        identifier: identifier,
        comment: comment,
        children: children,
    }
}

func (instance EnumDefinition) String() string {
    buf := new(bytes.Buffer)
    fmt.Fprintf(buf, "%s %s // %s", instance.Identifier().AsTargetIdentifier(), instance.TypeName(), strings.Replace(instance.comment, "\n", " - ", -1))
    for _, child := range instance.children {
        fmt.Fprintf(buf, "\n\t%v", child)
    }
    return buf.String()
}

