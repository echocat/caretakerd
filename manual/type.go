package main

import (
	"strings"
	"fmt"
	"github.com/echocat/caretakerd/panics"
	"reflect"
)

type Type interface {
	String() string
}

type IdType struct {
	Package   string
	Name      string
	Primitive bool
}

func (instance IdType) String() string {
	result := ""
	if len(instance.Package) > 0 {
		result += instance.Package + "."
	}
	result += instance.Name
	return result
}

func NewIdType(packageName string, name string, primitive bool) IdType {
	return IdType{
		Package: packageName,
		Name: name,
		Primitive: primitive,
	}
}

type MapType struct {
	Key   Type
	Value Type
}

func (instance MapType) String() string {
	return "map[" + instance.Key.String() + "]" + instance.Value.String()
}

func NewMapType(key Type, value Type) MapType {
	return MapType{
		Key: key,
		Value: value,
	}
}

type ArrayType struct {
	Value Type
}

func (instance ArrayType) String() string {
	return "[]" + instance.Value.String()
}

func NewArrayType(value Type) ArrayType {
	return ArrayType{
		Value: value,
	}
}

type PointerType struct {
	Value Type
}

func (instance PointerType) String() string {
	return "*" + instance.Value.String()
}

func NewPointerType(value Type) PointerType {
	return PointerType{
		Value: value,
	}
}

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
				value := ParseType(plain[i + 1:])
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
	} else {
		lastDot := strings.LastIndex(plain, ".")
		if lastDot <= 0 || len(plain) <= lastDot + 1 {
			return NewIdType("", plain, true)
		}
		packageName := plain[:lastDot]
		name := plain[lastDot + 1:]
		return NewIdType(packageName, name, false)
	}
}

func ExtractAllIdTypesFrom(t Type) []IdType {
	if idType, ok := t.(IdType); ok {
		return []IdType{idType}
	} else if arrayType, ok := t.(ArrayType); ok {
		return ExtractAllIdTypesFrom(arrayType.Value)
	} else if pointerType, ok := t.(PointerType); ok {
		return ExtractAllIdTypesFrom(pointerType.Value)
	} else if mapType, ok := t.(MapType); ok {
		result := []IdType{}
		result = append(result, ExtractAllIdTypesFrom(mapType.Key)...)
		result = append(result, ExtractAllIdTypesFrom(mapType.Value)...)
		return result
	}
	panic(panics.New("Unknown type %v.", reflect.TypeOf(t)))
}