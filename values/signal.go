package values

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
	"syscall"
)

// Signal represents an system signal.
// @inline
type Signal syscall.Signal

const (
	NOOP   = Signal(0x0)
	ABRT   = Signal(0x6)
	ALRM   = Signal(0xe)
	BUS    = Signal(0x7)
	CHLD   = Signal(0x11)
	CONT   = Signal(0x12)
	FPE    = Signal(0x8)
	HUP    = Signal(0x1)
	ILL    = Signal(0x4)
	INT    = Signal(0x2)
	IO     = Signal(0x1d)
	KILL   = Signal(0x9)
	PIPE   = Signal(0xd)
	PROF   = Signal(0x1b)
	PWR    = Signal(0x1e)
	QUIT   = Signal(0x3)
	SEGV   = Signal(0xb)
	STKFLT = Signal(0x10)
	STOP   = Signal(0x13)
	SYS    = Signal(0x1f)
	TERM   = Signal(0xf)
	TRAP   = Signal(0x5)
	TSTP   = Signal(0x14)
	TTIN   = Signal(0x15)
	TTOU   = Signal(0x16)
	URG    = Signal(0x17)
	USR1   = Signal(0xa)
	USR2   = Signal(0xc)
	VTALRM = Signal(0x1a)
	WINCH  = Signal(0x1c)
	XCPU   = Signal(0x18)
	XFSZ   = Signal(0x19)
)

// SignalToName contains all possible variants of Signal and their name
var SignalToName = map[Signal]string{
	NOOP:   "NOOP",
	ABRT:   "ABRT",
	ALRM:   "ALRM",
	BUS:    "BUS",
	CHLD:   "CHLD",
	CONT:   "CONT",
	FPE:    "FPE",
	HUP:    "HUP",
	ILL:    "ILL",
	INT:    "INT",
	IO:     "IO",
	KILL:   "KILL",
	PIPE:   "PIPE",
	PROF:   "PROF",
	PWR:    "PWR",
	QUIT:   "QUIT",
	SEGV:   "SEGV",
	STKFLT: "STKFLT",
	STOP:   "STOP",
	SYS:    "SYS",
	TERM:   "TERM",
	TRAP:   "TRAP",
	TSTP:   "TSTP",
	TTIN:   "TTIN",
	TTOU:   "TTOU",
	URG:    "URG",
	USR1:   "USR1",
	USR2:   "USR2",
	VTALRM: "VTALRM",
	WINCH:  "WINCH",
	XCPU:   "XCPU",
	XFSZ:   "XFSZ",
}

// SignalToName contains all possible variants of Signal names and their values
var NameToSignal = map[string]Signal{}

func init() {
	for value, name := range SignalToName {
		NameToSignal[name] = value
	}
}

func (instance Signal) String() string {
	result, err := instance.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

// CheckedString is like String but return also an optional error if there are some
// validation errors.
func (instance Signal) CheckedString() (string, error) {
	if name, ok := SignalToName[instance]; ok {
		return name, nil
	}
	return "", errors.New("Illegal signal: %d", instance)
}

// Set the given string to current object from a string.
// Return an error object if there are some problems while transforming the string.
func (instance *Signal) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		candidate := Signal(valueAsInt)
		if _, ok := SignalToName[candidate]; ok {
			(*instance) = candidate
			return nil
		}
		return errors.New("Illegal signal: " + value)
	}
	valueAsUpperCase := strings.ToUpper(value)
	if candidate, ok := NameToSignal[valueAsUpperCase]; ok {
		(*instance) = candidate
		return nil
	}
	if candidate, ok := NameToSignal["SIG"+valueAsUpperCase]; ok {
		(*instance) = candidate
		return nil
	}
	return errors.New("Illegal signal: " + value)
}

// MarshalYAML is used until yaml marshalling. Do not call directly.
func (instance Signal) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call directly.
func (instance *Signal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call directly.
func (instance Signal) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call directly.
func (instance *Signal) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate do validate action on this object and return an error object if any.
func (instance Signal) Validate() error {
	_, err := instance.CheckedString()
	return err
}
