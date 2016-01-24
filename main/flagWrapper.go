package main

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
        value: initialValue,
        explicitSet: false,
    }
}

func (this FlagWrapper) String() string {
    return this.value.String()
}

func (this *FlagWrapper) Set(value string) error {
    err := this.value.Set(value)
    if err != nil {
        return err
    }
    this.explicitSet = true
    return nil
}

func (this FlagWrapper) IsExplicitSet() bool {
    return this.explicitSet
}

func (this FlagWrapper) Value() cli.Generic {
    return this.value
}

func (this FlagWrapper) AssignIfExplicitSet(to interface{}) {
    if this.explicitSet {
        v := reflect.ValueOf(this.value)
        reflect.ValueOf(to).Elem().Set(v.Elem())
    }
}
