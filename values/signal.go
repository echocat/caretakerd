package values

import (
	"strconv"
	"strings"
	"github.com/echocat/caretakerd/errors"
	"syscall"
	"time"
	"encoding/json"
)

// @id Signal
// @type enum
//
// ## Description
//
// This represents a system signal.
type Signal syscall.Signal

const (
// @id NOOP
	NOOP = Signal(0x0)
// @id ABRT
	ABRT = Signal(0x6)
// @id ALRM
	ALRM = Signal(0xe)
// @id BUS
	BUS = Signal(0x7)
// @id CHLD
	CHLD = Signal(0x11)
// @id CONT
	CONT = Signal(0x12)
// @id FPE
	FPE = Signal(0x8)
// @id HUP
	HUP = Signal(0x1)
// @id ILL
	ILL = Signal(0x4)
// @id INT
	INT = Signal(0x2)
// @id IO
	IO = Signal(0x1d)
// @id KILL
	KILL = Signal(0x9)
// @id PIPE
	PIPE = Signal(0xd)
// @id PROF
	PROF = Signal(0x1b)
// @id PWR
	PWR = Signal(0x1e)
// @id QUIT
	QUIT = Signal(0x3)
// @id SEGV
	SEGV = Signal(0xb)
// @id STKFLT
	STKFLT = Signal(0x10)
// @id STOP
	STOP = Signal(0x13)
// @id SYS
	SYS = Signal(0x1f)
// @id TERM
	TERM = Signal(0xf)
// @id TRAP
	TRAP = Signal(0x5)
// @id TSTP
	TSTP = Signal(0x14)
// @id TTIN
	TTIN = Signal(0x15)
// @id TTOU
	TTOU = Signal(0x16)
// @id URG
	URG = Signal(0x17)
// @id USR1
	USR1 = Signal(0xa)
// @id USR2
	USR2 = Signal(0xc)
// @id VTALRM
	VTALRM = Signal(0x1a)
// @id WINCH
	WINCH = Signal(0x1c)
// @id XCPU
	XCPU = Signal(0x18)
// @id XFSZ
	XFSZ = Signal(0x19)
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
			if candidateAsString == lowerValue || "sig" + candidateAsString == lowerValue {
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

const (
	lastSendSignalThreshold = 500 * time.Millisecond
)

var (
	lastSendSignals = map[Signal]time.Time{}
)

func IsHandlingOfSignalIgnoreable(what Signal) bool {
	timeout, ok := lastSendSignals[what]
	return ok && timeout.After(time.Now())
}

func RecordSendSignal(what Signal) {
	lastSendSignals[what] = time.Now().Add(lastSendSignalThreshold)
	if (what == INT) {
		lastSendSignals[TERM] = time.Now().Add(lastSendSignalThreshold)
	}
	if (what == TERM) {
		lastSendSignals[INT] = time.Now().Add(lastSendSignalThreshold)
	}
}

