package app

import (
	"reflect"
)

type Stringer interface {
	Set(value string) error
	String() string
}

// FlagWrapper wraps generic command line flags.
type FlagWrapper struct {
	value       Stringer
	explicitSet bool
}

// NewFlagWrapper creates a new instance of FlagWrapper.
func NewFlagWrapper(initialValue Stringer) *FlagWrapper {
	return &FlagWrapper{
		value:       initialValue,
		explicitSet: false,
	}
}

func (instance FlagWrapper) String() string {
	return instance.value.String()
}

// Set sets the given string to the current object from a string.
// Returns an error object if there are problems while transforming the string.
func (instance *FlagWrapper) Set(value string) error {
	err := instance.value.Set(value)
	if err != nil {
		return err
	}
	instance.explicitSet = true
	return nil
}

// IsExplicitSet returns true if the Set(string) method has been called before.
func (instance FlagWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

// Value returns the wrapped generic command line flag.
func (instance FlagWrapper) Value() Stringer {
	return instance.value
}

// AssignIfExplicitSet assigns the wrapped value to given object
// if it was explicitly set before.
func (instance FlagWrapper) AssignIfExplicitSet(to interface{}) {
	if instance.explicitSet {
		v := reflect.ValueOf(instance.value)
		reflect.ValueOf(to).Elem().Set(v.Elem())
	}
}
