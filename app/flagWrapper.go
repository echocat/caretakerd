package app

import (
	"github.com/codegangsta/cli"
	"reflect"
)

type FlagWrapper struct {
	value       cli.Generic
	explicitSet bool
}

func NewFlagWrapper(initialValue cli.Generic) *FlagWrapper {
	return &FlagWrapper{
		value:       initialValue,
		explicitSet: false,
	}
}

func (instance FlagWrapper) String() string {
	return instance.value.String()
}

func (instance *FlagWrapper) Set(value string) error {
	err := instance.value.Set(value)
	if err != nil {
		return err
	}
	instance.explicitSet = true
	return nil
}

func (instance FlagWrapper) IsExplicitSet() bool {
	return instance.explicitSet
}

func (instance FlagWrapper) Value() cli.Generic {
	return instance.value
}

func (instance FlagWrapper) AssignIfExplicitSet(to interface{}) {
	if instance.explicitSet {
		v := reflect.ValueOf(instance.value)
		reflect.ValueOf(to).Elem().Set(v.Elem())
	}
}
