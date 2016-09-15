package keyStore

import (
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/values"
	"strconv"
	"strings"
)

var defaults = map[string]interface{}{
	"Type":    Generated,
	"PemFile": values.String(""),
	"Hints":   values.String("algorithm:`rsa` bits:`1024`"),
	"CaFile":  values.String(""),
}

// # Description
//
// Defines the keyStore of caretakerd.
type Config struct {
	// @default generated
	//
	// Defines the type of the instance keyStore.
	Type Type `json:"type" yaml:"type"`

	// @default ""
	//
	// Defines the pemFile which contains the key and certificate to be used.
	// This has to be of type PEM and has to contain the certificate and private key.
	// Currently only private keys of type RSA are supported.
	//
	// This property is only evaluated and required if {@ref #Type type} is set to
	// {@ref .Type#FromFile fromFile}.
	PemFile values.String `json:"pemFile,omitempty" yaml:"pemFile"`

	// @default "algorithm:`rsa` bits:`1024`"
	//
	// Defines some hints, for example to store in the format ``[<key:`value`>...]``.
	// Possible hints are:
	//
	// * ``algorithm``: Algorithm to be used to create new keys. Currently only ``rsa`` is supported.
	// * ``bits``: Number of bits to create a new key with.
	Hints values.String `json:"hints,omitempty" yaml:"hints"`

	// @default ""
	//
	// File where trusted certificates are stored in. This has to be in PEM format.
	CaFile values.String `json:"caFile,omitempty" yaml:"caFile"`
}

// NewConfig creates a new instance of Config.
func NewConfig() Config {
	result := Config{}
	values.SetDefaultsTo(defaults, &result)
	return result
}

// Validate validates an action on this object and returns an error object if there are any.
func (instance Config) Validate() error {
	err := instance.Type.Validate()
	if err == nil {
		err = instance.validateRequireStringOrNotValue(instance.PemFile, "pemFile", instance.Type.IsTakingFilename)
	}
	if err == nil {
		err = instance.validateStringOnlyAllowedValue(instance.CaFile, "caFile", instance.Type.IsConsumingCAFile)
	}
	if err == nil {
		algorithm := instance.GetHintsArgument("algorithm")
		if len(algorithm) > 0 && strings.ToLower(algorithm) != "rsa" {
			err = errors.New("Unsupported algorithm: %s", algorithm)
		}
	}
	return err
}

func (instance Config) validateRequireStringOrNotValue(value values.String, fieldName string, isAllowedMethod func() bool) error {
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

func (instance Config) validateStringOnlyAllowedValue(value values.String, fieldName string, isAllowedMethod func() bool) error {
	if !isAllowedMethod() && !value.IsEmpty() {
		return errors.New("There is no %s allowed for type %v.", fieldName, instance.Type)
	}
	return nil
}

// GetHintsArgument returns hints argument content for the given key.
// If there is no hint for this key and empty string is returned.
func (instance Config) GetHintsArgument(key string) string {
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
