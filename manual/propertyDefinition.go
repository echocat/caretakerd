package main

import (
    "bytes"
    "fmt"
    "strings"
)

type PropertyDefinition struct {
    identifier Identifier
    comment    string
    def        *string
}

func (instance PropertyDefinition) Identifier() Identifier {
    return instance.identifier
}

func (instance PropertyDefinition) Comment() string {
    return instance.comment
}

func (instance PropertyDefinition) TypeName() string {
    return "property"
}

func (instance PropertyDefinition) IsTopLevel() bool {
    return false
}

func (instance PropertyDefinition) Default() *string {
    return instance.def
}

func newPropertyDefinition(identifier Identifier, comment string, def *string) *PropertyDefinition {
    return &PropertyDefinition{
        identifier: identifier,
        comment: comment,
        def: def,
    }
}

func (instance PropertyDefinition) String() string {
    buf := new(bytes.Buffer)
    fmt.Fprintf(buf, "%s %s = %v // %s", instance.Identifier().AsTargetIdentifier(), instance.TypeName(), instance.def, strings.Replace(instance.comment, "\n", " - ", -1))
    return buf.String()
}
