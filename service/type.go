package service

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
)

// # Description
//
// Identifies the different ways caretakerd handles services.
type Type int

const (
	// @id onDemand
	//
	// The service is not automatically started by caretakerd.
	// You have to use [``caretakerctl``](#commands.caretakerctl) or execute an RPC call from another service
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
	// This service is automatically started by caretakerd and influences all other services.
	//
	// > **Important:** One of all available services must be specified as ``master``.
	//
	// Every other service lives and dies together with the ``master``.
	Master Type = 2
)

// AllTypes contains all possible variants of Type.
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

// CheckedString is like String but also returns an optional error if there are any
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
// Returns an error object if there are any problems while transforming the string.
func (instance *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*instance) = candidate
				return nil
			}
		}
		return errors.New("Illegal type: " + value)
	}
	lowerValue := strings.ToLower(value)
	for _, candidate := range AllTypes {
		if strings.ToLower(candidate.String()) == lowerValue {
			(*instance) = candidate
			return nil
		}
	}
	return errors.New("Illegal type: " + value)
}

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Type) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Type) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// IsAutoStartable returns "true" if the given type indicates that the service
// have to be started automatically together with caretakerd.
func (instance Type) IsAutoStartable() bool {
	switch instance {
	case Master:
		fallthrough
	case AutoStart:
		return true
	}
	return false
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Type) Validate() error {
	_, err := instance.CheckedString()
	return err
}
