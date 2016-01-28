package main

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/panics"
	"go/ast"
	"go/token"
)

type Element struct {
	Id            string
	Type          Type
	Default       string
	Documentation string
	Filename      string
	Comment       *ast.CommentGroup
	Position      token.Position
}

type Type int

const (
	Simple = 0
	Struct = 1
	Enum = 2
)

func (instance *Type) Set(value string) error {
	switch value {
	case "simple":
		(*instance) = Simple
	case "struct":
		(*instance) = Struct
	case "enum":
		(*instance) = Enum
	default:
		return errors.New("Unknown element type: %v", value)
	}
	return nil
}

func (instance Type) String() string {
	switch instance {
	case Simple:
		return "simple"
	case Struct:
		return "struct"
	case Enum:
		return "enum"
	}
	panic(panics.New("Unknown element type: %v", instance))
}