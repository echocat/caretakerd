package main

import (
    "bytes"
    "fmt"
"strings"
)

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

func (instance ObjectDefinition) Children() []Definition {
    return instance.children
}

func (instance *ObjectDefinition) AdChild(child Definition) {
    instance.children = append(instance.children, child)
}

func newObjectDefinition(identifier Identifier, comment string, children ... Definition) *ObjectDefinition {
    return &ObjectDefinition{
        identifier: identifier,
        comment: comment,
        children: children,
    }
}

func (instance ObjectDefinition) String() string {
    buf := new(bytes.Buffer)
    fmt.Fprintf(buf, "%s %s // %s", instance.Identifier().AsTargetIdentifier(), instance.TypeName(), strings.Replace(instance.comment, "\n", " - ", -1))
    for _, child := range instance.children {
        fmt.Fprintf(buf, "\n\t%v", child)
    }
    return buf.String()
}

