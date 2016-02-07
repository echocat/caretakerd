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
	OnDemand  Type = 0
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
	Master    Type = 2
)

var AllTypes = []Type{
	OnDemand,
	AutoStart,
	Master,
}

func (i Type) String() string {
	result, err := i.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

func (i Type) CheckedString() (string, error) {
	switch i {
	case OnDemand:
		return "onDemand", nil
	case AutoStart:
		return "autoStart", nil
	case Master:
		return "master", nil
	}
	return "", errors.New("Illegal type: %d", i)
}

func (i *Type) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllTypes {
			if int(candidate) == valueAsInt {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal type: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		for _, candidate := range AllTypes {
			if strings.ToLower(candidate.String()) == lowerValue {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal type: " + value)
	}
}

func (i Type) MarshalYAML() (interface{}, error) {
	return i.CheckedString()
}

func (i *Type) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Type) MarshalJSON() ([]byte, error) {
	s, err := i.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (i *Type) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Type) IsAutoStartable() bool {
	switch i {
	case Master:
		fallthrough
	case AutoStart:
		return true
	}
	return false
}

func (i Type) Validate() error {
	_, err := i.CheckedString()
	return err
}
