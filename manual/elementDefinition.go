package main

import (
    "bytes"
    "fmt"
"strings"
)

type ElementDefinition struct {
    identifier Identifier
    comment    string
}

func (instance ElementDefinition) Identifier() Identifier {
    return instance.identifier
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

func newElementDefinition(identifier Identifier, comment string) *ElementDefinition {
    return &ElementDefinition{
        identifier: identifier,
        comment: comment,
    }
}

func (instance ElementDefinition) String() string {
    buf := new(bytes.Buffer)
    fmt.Fprintf(buf, "%s %s // %s", instance.Identifier().AsTargetIdentifier(), instance.TypeName(), strings.Replace(instance.comment, "\n", " - ", -1))
    return buf.String()
}

