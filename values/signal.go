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
	// NOOP represents the system signal NOOP
	NOOP = Signal(0x0)
	// ABRT represents the system signal ABRT
	ABRT = Signal(0x6)
	// ALRM represents the system signal ALRM
	ALRM = Signal(0xe)
	// BUS represents the system signal BUS
	BUS = Signal(0x7)
	// CHLD represents the system signal CHLD
	CHLD = Signal(0x11)
	// CONT represents the system signal CONT
	CONT = Signal(0x12)
	// FPE represents the system signal FPE
	FPE = Signal(0x8)
	// HUP represents the system signal HUP
	HUP = Signal(0x1)
	// ILL represents the system signal ILL
	ILL = Signal(0x4)
	// INT represents the system signal INT
	INT = Signal(0x2)
	// IO represents the system signal IO
	IO = Signal(0x1d)
	// KILL represents the system signal KILL
	KILL = Signal(0x9)
	// PIPE represents the system signal PIPE
	PIPE = Signal(0xd)
	// PROF represents the system signal PROF
	PROF = Signal(0x1b)
	// PWR represents the system signal PWR
	PWR = Signal(0x1e)
	// QUIT represents the system signal QUIT
	QUIT = Signal(0x3)
	// SEGV represents the system signal SEGV
	SEGV = Signal(0xb)
	// STKFLT represents the system signal STKFLT
	STKFLT = Signal(0x10)
	// STOP represents the system signal STOP
	STOP = Signal(0x13)
	// SYS represents the system signal SYS
	SYS = Signal(0x1f)
	// TERM represents the system signal TERM
	TERM = Signal(0xf)
	// TRAP represents the system signal TRAP
	TRAP = Signal(0x5)
	// TSTP represents the system signal TSTP
	TSTP = Signal(0x14)
	// TTIN represents the system signal TTIN
	TTIN = Signal(0x15)
	// TTOU represents the system signal TTOU
	TTOU = Signal(0x16)
	// URG represents the system signal URG
	URG = Signal(0x17)
	// USR1 represents the system signal USR1
	USR1 = Signal(0xa)
	// USR2 represents the system signal USR2
	USR2 = Signal(0xc)
	// VTALRM represents the system signal VTALRM
	VTALRM = Signal(0x1a)
	// WINCH represents the system signal WINCH
	WINCH = Signal(0x1c)
	// XCPU represents the system signal XCPU
	XCPU = Signal(0x18)
	// XFSZ represents the system signal XFSZ
	XFSZ = Signal(0x19)
)

// SignalToName contains all possible variants of Signal and their names
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

// NameToSignal contains all possible variants of Signal names and their values
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

// CheckedString is like String but also returns an optional error if there are any
// validation errors.
func (instance Signal) CheckedString() (string, error) {
	if name, ok := SignalToName[instance]; ok {
		return name, nil
	}
	return "", errors.New("Illegal signal: %d", instance)
}

// Set sets the given string to current object from a string.
// Returns an error object if there are any problems while transforming the string.
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

// MarshalYAML is used until yaml marshalling. Do not call this method directly.
func (instance Signal) MarshalYAML() (interface{}, error) {
	return instance.CheckedString()
}

// UnmarshalYAML is used until yaml unmarshalling. Do not call this method directly.
func (instance *Signal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return instance.Set(value)
}

// MarshalJSON is used until json marshalling. Do not call this method directly.
func (instance Signal) MarshalJSON() ([]byte, error) {
	s, err := instance.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

// UnmarshalJSON is used until json unmarshalling. Do not call this method directly.
func (instance *Signal) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return instance.Set(value)
}

// Validate validates actions on this object and returns an error object if there are any.
func (instance Signal) Validate() error {
	_, err := instance.CheckedString()
	return err
}
