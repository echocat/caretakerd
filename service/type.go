package service

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// # Description
//
// Identifies the differed ways how caretakerd handles services.
type Type int

const (
	// @id onDemand
	//
	// The service is not automatically started by caretakerd.
	// You have to use [``caretakerctl``](#commands.caretakerctl) or to execute a RPC call from another service
	// to start it.
	//
	// This service will be automatically stopped if the {@ref #Master master} was also stopped.
	OnDemand Type = 0
	// @id autoStart
	//
	// This services is automatically started by caretakerd.
	//
	// This service will be automatically stopped if the {@ref #Master master} was also stopped.
	AutoStart Type = 1
	// @id master
	//
	// This services is automatically started by caretakerd and influence all other services.
	//
	// > **Important:** There have to be exact one of all services specified as ``master``.
	//
	// Every other service will live and die together with the ``master``.
	Master Type = 2
)

var AllTypes = []Type{
	OnDemand,
	AutoStart,
	Master,
}

func (instance Type) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Type) CheckedString() (string, error) {
	switch instance {
	case OnDemand:
		return "onDemand", nil
	case AutoStart:
		return "autoStart", nil
	case Master:
		return "master", nil
	}
	return "", errors.New("Illegal type: %d", instance)
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal type: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllTypes {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal type: " + value)
	}
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance Type) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

func (instance Type) IsAutoStartable() bool {
	switch instance {
	case Master:
		fallthrough
	case AutoStart:
		return true
	}
	return false
}

// Validate do validate action on this object and return an error object if any.
func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
