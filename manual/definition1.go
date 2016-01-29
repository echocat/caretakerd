package main

import (
	"bytes"
	"fmt"
	"strings"
)

type Definition interface {
	Identifier() Identifier
	Comment() string
	TypeName() string
	IsTopLevel() bool
}

type WithDefaultValue interface {
	DefaultValue() *string
}

type WithChildren interface {
	Children() []Definition
	AddChild(child Definition)
}

type WithKey interface {
	Key() string
}

type WithValueType interface {
	ValueType() Identifier
}

func FormatDefinition(definition Definition) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v %s", definition.Identifier(), definition.TypeName())
	if key, ok := definition.(WithKey); ok {
		fmt.Fprintf(buf, ": %s", key.Key())
		if valueType, ok := definition.(WithValueType); ok {
			fmt.Fprintf(buf, " %v", valueType.ValueType())
		}
	} else if valueType, ok := definition.(WithValueType); ok {
		fmt.Fprintf(buf, ": %v", valueType.ValueType())
	}
	if defaultValue, ok := definition.(WithDefaultValue); ok {
		def := defaultValue.DefaultValue()
		if def != nil {
			fmt.Fprintf(buf, " = %v", *def)
		}
	}
	comment := definition.Comment()
	if len(comment) > 0 {
		fmt.Fprintf(buf, " // %s", strings.Replace(comment, "\n", " - ", -1))
	}
	if children, ok := definition.(WithChildren); ok {
		for _, child := range children.Children() {
			fmt.Fprintf(buf, "\n\t\t%v", child)

		}
	}
	return buf.String()
}
