package keyStore

import (
	"github.com/echocat/caretakerd/errors"
	. "github.com/echocat/caretakerd/values"
	"strconv"
	"strings"
)

var defaults = map[string]interface{}{
	"Type":    Generated,
	"PemFile": String(""),
	"Hints":   String("algorithm:`rsa` bits:`1024`"),
	"CaFile":  String(""),
}

// # Description
//
// Defines the keyStore of caretakerd.
type Config struct {
	// @default generated
	//
	// Defines the type of instance keyStore.
	Type Type `json:"type" yaml:"type"`

	// @default ""
	//
	// Defines the pemFile which contains the key and certificate to use.
	// This have to be of type PEM and have to contain the certificate and private key.
	// Currently the private key is only supported of type RSA.
	//
	// This property is only evaluated and required if {@ref Type} is set to
	// {@ref github.com/echocat/caretakerd/keyStore.Type#FromFile}.
	PemFile String `json:"pemFile,omitempty" yaml:"pemFile"`

	// @default "algorithm:`rsa` bits:`1024`"
	//
	// Defines some hints for instance store in format ``[<key:`value`>...]``.
	// Possible hints are:
	//
	// * ``algorithm``: Algorithm to use for creation of new keys. Currently only ``rsa`` is supported.
	// * ``bits``: Number of bits to create a new key with.
	Hints String `json:"hints,omitempty" yaml:"hints"`

	// @default ""
	//
	// File where trusted certificates are stored in. This have to be in PEM format.
	CaFile String `json:"caFile,omitempty" yaml:"caFile"`
}

func NewConfig() Config {
	result := Config{}
	SetDefaultsTo(defaults, &result)
	return result
}

func (instance Config) Validate() error {
	err := instance.Type.Validate()
	if err == nil {
		err = instance.validateRequireStringOrNotValue(instance.PemFile, "pemFile", instance.Type.IsTakingFilename)
	}
	if err == nil {
		err = instance.validateStringOnlyAllowedValue(instance.CaFile, "caFile", instance.Type.IsConsumingCaFile)
	}
	if err == nil {
		algorithm := instance.GetKeyArgument("algorithm")
		if len(algorithm) > 0 && strings.ToLower(algorithm) != "rsa" {
			err = errors.New("Unsupported algorithm: %s", algorithm)
		}
	}
	return err
}

func (instance Config) validateRequireStringOrNotValue(value String, fieldName string, isAllowedMethod func() bool) error {
	if isAllowedMethod() {
		if value.IsEmpty() {
			return errors.New("There is no %s set for type %v.", fieldName, instance.Type)
		}
	} else {
		if !value.IsEmpty() {
			return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
		}
	}
	return nil
}

func (instance Config) validateStringOnlyAllowedValue(value String, fieldName string, isAllowedMethod func() bool) error {
	if !isAllowedMethod() && !value.IsEmpty() {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}

func (instance Config) GetKeyArgument(key string) string {
	arguments := instance.Hints
	for arguments != "" {
		i := 0
		for i < len(arguments) && arguments[i] == ' ' {
			i++
		}
		arguments = arguments[i:]
		if arguments == "" {
			break
		}
		i = 0
		for i < len(arguments) && arguments[i] > ' ' && arguments[i] != ':' && arguments[i] != '`' && arguments[i] != 0x7f {
			i++
		}
		if i == 0 || i+1 >= len(arguments) || arguments[i] != ':' || arguments[i+1] != '`' {
			break
		}
		name := string(arguments[:i])
		arguments = arguments[i+1:]

		i = 1
		for i < len(arguments) && arguments[i] != '`' {
			if arguments[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(arguments) {
			break
		}
		qvalue := string(arguments[:i+1])
		arguments = arguments[i+1:]

		if key == name {
			value, err := strconv.Unquote(qvalue)
			if err != nil {
				break
			}
			return value
		}
	}
	return ""
}
