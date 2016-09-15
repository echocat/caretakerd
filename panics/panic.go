package panics

import (
	"fmt"
	"github.com/echocat/caretakerd/stack"
	"os"
	"reflect"
)

// Panic represents a go panic which could have a cause and a call stack.
type Panic struct {
	message string
	cause   interface{}
	stack   stack.Stack
}

// New creates a new Panic instance with the given message.
func New(message string, a ...interface{}) Panic {
	var targetMessage string
	if len(a) == 0 {
		targetMessage = message
	} else {
		targetMessage = fmt.Sprintf(message, a...)
	}
	result := Panic{
		message: targetMessage,
		stack:   stack.CaptureStack(1),
	}
	return result
}

func (i Panic) String() string {
	return reflect.TypeOf(i).Name() + ": " + i.Message()
}

func (i Panic) Error() string {
	return i.Message()
}

// Message returns the message of this panic.
func (i Panic) Message() string {
	return i.message
}

// Cause returns the cause of this panic.
// If the result is nil there is no cause.
func (i Panic) Cause() interface{} {
	return i.cause
}

// Stack returns the call stack of this panic.
func (i Panic) Stack() stack.Stack {
	return i.stack
}

// CausedBy is a shortcut for Stack() for using in builder pattern.
func (i Panic) CausedBy(what interface{}) Panic {
	i.cause = what
	return i
}

// Throw throws this panic.
func (i Panic) Throw() {
	panic(i)
}

// DefaultPanicHandler could be used as panic handler on top method like
//     defer panics.DefaultPanicHandler()
// ...to handle every panic in a better way.
func DefaultPanicHandler() {
	if r := recover(); r != nil {
		stack.Print(r, os.Stderr, 4)
		os.Exit(2)
	}
}

// HandlePanicOfPresent could be used as panic handler on top method like
//     defer panics.HandlePanicOfPresent(2)
// ...to handle every panic in a better way.
func HandlePanicOfPresent(framesToSkip int) {
	if r := recover(); r != nil {
		stack.Print(r, os.Stderr, 4+framesToSkip)
		os.Exit(2)
	}
}
