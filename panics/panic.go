package panics

import (
	"fmt"
	"github.com/echocat/caretakerd/stack"
	"os"
	"reflect"
)

type Panic struct {
	message string
	cause   interface{}
	stack   stack.Stack
}

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

func (i Panic) Message() string {
	return i.message
}

func (i Panic) Cause() interface{} {
	return i.cause
}

func (i Panic) Stack() stack.Stack {
	return i.stack
}

func (i Panic) CausedBy(what interface{}) Panic {
	i.cause = what
	return i
}

func (i Panic) Throw() {
	panic(i)
}

func DefaultPanicHandler() {
	if r := recover(); r != nil {
		stack.Print(r, os.Stderr, 4)
		os.Exit(2)
	}
}

func HandlePanicOfPresent(framesToSkip int) {
	if r := recover(); r != nil {
		stack.Print(r, os.Stderr, 4+framesToSkip)
		os.Exit(2)
	}
}
