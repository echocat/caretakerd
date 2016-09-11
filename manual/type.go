package main

import (
	"fmt"
	"github.com/echocat/caretakerd/panics"
	"reflect"
	"strings"
)

// Type is a type definition in the source code.
type Type interface {
	String() string
}

// IDType is a type definition in the source code which directly references a type by name/id.
type IDType struct {
	Package   string
	Name      string
	Primitive bool
}

func (instance IDType) String() string {
	result := ""
	if len(instance.Package) > 0 {
		result += instance.Package + "."
	}
	result += instance.Name
	return result
}

// NewIDType creates a new instance of IDType
func NewIDType(packageName string, name string, primitive bool) IDType {
	return IDType{
		Package:   packageName,
		Name:      name,
		Primitive: primitive,
	}
}

// MapType is a type definition in the source code which references the map type that contains other types.
type MapType struct {
	Key   Type
	Value Type
}

func (instance MapType) String() string {
	return "map[" + instance.Key.String() + "]" + instance.Value.String()
}

// NewMapType creates a new instance of MapType
func NewMapType(key Type, value Type) MapType {
	return MapType{
		Key:   key,
		Value: value,
	}
}

// ArrayType is a type definition in the source code which references the array type that contains another type.
type ArrayType struct {
	Value Type
}

func (instance ArrayType) String() string {
	return "[]" + instance.Value.String()
}

// NewArrayType creates a new instance of ArrayType.
func NewArrayType(value Type) ArrayType {
	return ArrayType{
		Value: value,
	}
}

// PointerType is a type definition in the source code which references the pointer type that contains another type.
type PointerType struct {
	Value Type
}

func (instance PointerType) String() string {
	return "*" + instance.Value.String()
}

// NewPointerType creates a new instance of PointerType.
func NewPointerType(value Type) PointerType {
	return PointerType{
		Value: value,
	}
}

// ParseType parses the given plain string and transform it to a Type.
func ParseType(plain string) Type {
	if strings.HasPrefix(plain, "map[") {
		inBracketCount := 1
		for i := 4; i < len(plain); i++ {
			c := plain[i]
			if c == '[' {
				inBracketCount++
			} else if c == ']' {
				inBracketCount--
			}
			if inBracketCount <= 0 {
				key := ParseType(plain[4:i])
				value := ParseType(plain[i+1:])
				return NewMapType(key, value)
			}
		}
		panic(fmt.Sprintf("Unexpexted end of type: '%v'", plain))
	} else if strings.HasPrefix(plain, "[]") {
		value := ParseType(plain[2:])
		return NewArrayType(value)
	} else if strings.HasPrefix(plain, "*") {
		value := ParseType(plain[1:])
		return NewPointerType(value)
	}
	lastDot := strings.LastIndex(plain, ".")
	if lastDot <= 0 || len(plain) <= lastDot+1 {
		return NewIDType("", plain, true)
	}
	packageName := plain[:lastDot]
	name := plain[lastDot+1:]
	return NewIDType(packageName, name, false)
}

// ExtractAllIDTypesFrom extracts all referenced ID types from the given type as an array.
func ExtractAllIDTypesFrom(t Type) []IDType {
	if idType, ok := t.(IDType); ok {
		return []IDType{idType}
	} else if arrayType, ok := t.(ArrayType); ok {
		return ExtractAllIDTypesFrom(arrayType.Value)
	} else if pointerType, ok := t.(PointerType); ok {
		return ExtractAllIDTypesFrom(pointerType.Value)
	} else if mapType, ok := t.(MapType); ok {
		result := []IDType{}
		result = append(result, ExtractAllIDTypesFrom(mapType.Key)...)
		result = append(result, ExtractAllIDTypesFrom(mapType.Value)...)
		return result
	}
	panic(panics.New("Unknown type %v.", reflect.TypeOf(t)))
}

// ExtractValueIDType extracts the value ID types from the given type.
func ExtractValueIDType(t Type) IDType {
	if idType, ok := t.(IDType); ok {
		return idType
	} else if arrayType, ok := t.(ArrayType); ok {
		return ExtractValueIDType(arrayType.Value)
	} else if pointerType, ok := t.(PointerType); ok {
		return ExtractValueIDType(pointerType.Value)
	} else if mapType, ok := t.(MapType); ok {
		return ExtractValueIDType(mapType.Value)
	}
	panic(panics.New("Unknown type %v.", reflect.TypeOf(t)))
}
