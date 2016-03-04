package values

import (
	"encoding/json"
	"github.com/echocat/caretakerd/errors"
	"strconv"
	"strings"
	"syscall"
)

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

var AllSignals = []Signal{
	NOOP,
	ABRT,
	ALRM,
	BUS,
	CHLD,
	CONT,
	FPE,
	HUP,
	ILL,
	INT,
	IO,
	KILL,
	PIPE,
	PROF,
	PWR,
	QUIT,
	SEGV,
	STKFLT,
	STOP,
	SYS,
	TERM,
	TRAP,
	TSTP,
	TTIN,
	TTOU,
	URG,
	USR1,
	USR2,
	VTALRM,
	WINCH,
	XCPU,
	XFSZ,
}

func (i Signal) String() string {
	result, err := i.CheckedString()
	if err != nil {
		panic(err)
	}
	return result
}

func (i Signal) CheckedString() (string, error) {
	switch i {
	case NOOP:
		return "NOOP", nil
	case ABRT:
		return "ABRT", nil
	case ALRM:
		return "ALRM", nil
	case BUS:
		return "BUS", nil
	case CHLD:
		return "CHLD", nil
	case CONT:
		return "CONT", nil
	case FPE:
		return "FPE", nil
	case HUP:
		return "HUP", nil
	case ILL:
		return "ILL", nil
	case INT:
		return "INT", nil
	case IO:
		return "IO", nil
	case KILL:
		return "KILL", nil
	case PIPE:
		return "PIPE", nil
	case PROF:
		return "PROF", nil
	case PWR:
		return "PWR", nil
	case QUIT:
		return "QUIT", nil
	case SEGV:
		return "SEGV", nil
	case STKFLT:
		return "STKFLT", nil
	case STOP:
		return "STOP", nil
	case SYS:
		return "SYS", nil
	case TERM:
		return "TERM", nil
	case TRAP:
		return "TRAP", nil
	case TSTP:
		return "TSTP", nil
	case TTIN:
		return "TTIN", nil
	case TTOU:
		return "TTOU", nil
	case URG:
		return "URG", nil
	case USR1:
		return "USR1", nil
	case USR2:
		return "USR2", nil
	case VTALRM:
		return "VTALRM", nil
	case WINCH:
		return "WINCH", nil
	case XCPU:
		return "XCPU", nil
	case XFSZ:
		return "XFSZ", nil
	}
	return "", errors.New("Illegal signal: %d", i)
}

func (i *Signal) Set(value string) error {
	if valueAsInt, err := strconv.Atoi(value); err == nil {
		for _, candidate := range AllSignals {
			if int(candidate) == valueAsInt {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal signal: " + value)
	} else {
		lowerValue := strings.ToUpper(value)
		for _, candidate := range AllSignals {
			candidateAsString := strings.ToUpper(candidate.String())
			if candidateAsString == lowerValue || "sig"+candidateAsString == lowerValue {
				(*i) = candidate
				return nil
			}
		}
		return errors.New("Illegal signal: " + value)
	}
}

func (i Signal) MarshalYAML() (interface{}, error) {
	return i.CheckedString()
}

func (i *Signal) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var value string
	if err := unmarshal(&value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Signal) MarshalJSON() ([]byte, error) {
	s, err := i.CheckedString()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(s)
}

func (i *Signal) UnmarshalJSON(b []byte) error {
	var value string
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}
	return i.Set(value)
}

func (i Signal) Validate() error {
	_, err := i.CheckedString()
	return err
}
