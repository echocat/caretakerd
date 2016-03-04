package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Definition represents a definition of an element or whatever that could be configured.
type Definition interface {
	Id() IDType
	Description() string
	TypeName() string
	IsTopLevel() bool
}

// WithDefaultValue is a Definition with a default value.
type WithDefaultValue interface {
	DefaultValue() *string
}

// WithChildren is a Definition with children definitions.
type WithChildren interface {
	Children() []Definition
	AddChild(child Definition)
}

// WithKey is a Definition with a key.
type WithKey interface {
	Key() string
}

// WithValueType is a Definition with a value type.
type WithValueType interface {
	ValueType() Type
}

// WithInlinedMarker is a Definition which could be marked as inlined.
type WithInlinedMarker interface {
	Inlined() bool
	ValueType() Type
}

// FormatDefinition formats a given definition to be printed for logging purposes.
func FormatDefinition(definition Definition) string {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v %s", definition.Id(), definition.TypeName())
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
	comment := definition.Description()
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
