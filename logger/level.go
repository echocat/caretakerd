package logger

import (
	"strconv"
	"strings"
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
)

// @id Level
// @type enum
//
// ## Description
//
// Represents a level for logging with a {@ref Logger}
type Level int

const (
// @id debug
//
// Used for debugging proposes. This level is only required you something goes wrong and you need more information.
	Debug Level = 200

// @id info
//
// This is the regular level. Every normal message will be logged with this level.
	Info Level = 300

// @id warning
//
// If a problem appears but the program is still able to continue its work, this level is used.
	Warning Level = 400

// @id error
//
// If a problem appears and the program is not longer able to continue its work, this level is used.
	Error Level = 500

// @id fatal
//
// This level is used on dramatic problems.
	Fatal Level = 600
)

var AllLevels []Level = []Level{
	Debug,
	Info,
	Warning,
	Error,
	Fatal,
}

func (i Level) String() string {
	result, err := i.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

func (i Level) CheckedString() (string, error) {
	switch i {
	case Debug:
		return "debug", nil
	case Info:
		return "info", nil
	case Warning:
		return "warning", nil
	case Error:
		return "error", nil
	case Fatal:
		return "fatal", nil
	}
	return strconv.Itoa(int(i)), nil
}

func (i Level) DisplayForLogging() string {
	if i == Warning {
		return "WARN"
	} else {
		return strings.ToUpper(i.String())
	}
}

func (i *Level) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllLevels {
			if int(candidate) == valueAsInt {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal level: " + value)
	} else {
		lowerValue := strings.ToLower(value)
		switch lowerValue {
		case "warn":
			*i = Warning
			return nil
		case "err":
			*i = Error
			return nil
		}
		for _, candidate := range AllLevels {
			if candidate.String() == lowerValue {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal level: " + value)
	}
}

func (i Level) IsIndicatingProblem() bool {
	return i == Warning || i == Error || i == Fatal
}

func (i Level) MarshalYAML() (interface{}, error) {
	return i.CheckedString()
}

func (i *Level) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Level) MarshalJSON() ([]byte, error) {
	s, err := i.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (i *Level) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Level) Validate() error {
	_, err := i.CheckedString()
	return err
}
