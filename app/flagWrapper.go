package app

import (
	"github.com/urfave/cli"
	"reflect"
)

// FlagWrapper wraps generic command line flags.
type FlagWrapper struct {
	value       cli.Generic
	explicitSet bool
}

// NewFlagWrapper creates a new instance of FlagWrapper
func NewFlagWrapper(initialValue cli.Generic) *FlagWrapper {
	return &FlagWrapper{
		value:       initialValue,
		explicitSet: false,
	}
}

func (instance FlagWrapper) String() string {
	return instance.value.String()
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *FlagWrapper) Set(value string) error {
	err := instance.value.Set(value)
	if err != nil {
		return err
	}
	instance.explicitSet = true
	return nil
}

// IsExplicitSet returns true if someone called the Set(string) method before.
func (instance FlagWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

// Value returns the wrapped generic command line flag.
func (instance FlagWrapper) Value() cli.Generic {
	return instance.value
}

// AssignIfExplicitSet will assign the wrapped value to given object
// if it was explicit set before.
func (instance FlagWrapper) AssignIfExplicitSet(to interface{}) {
	if instance.explicitSet {
		v := reflect.ValueOf(instance.value)
		reflect.ValueOf(to).Elem().Set(v.Elem())
	}
}
