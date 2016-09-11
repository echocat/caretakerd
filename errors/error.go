package errors

import (
	"fmt"
	"github.com/echocat/caretakerd/stack"
	"reflect"
)

// Error represents a GO error but with more features
// - such as a cause and a call stack.
type Error struct {
	message string
	cause   interface{}
	stack   stack.Stack
}

// New creates a new instance of Error.
func New(message string, a ...interface{}) Error {
	var targetMessage string
	if len(a) == 0 {
		targetMessage = message
	} else {
		targetMessage = fmt.Sprintf(message, a...)
	}
	result := Error{
		message: targetMessage,
		stack:   stack.CaptureStack(1),
	}
	return result
}

func (instance Error) String() string {
	return reflect.TypeOf(instance).Name() + ": " + instance.Message()
}

func (instance Error) Error() string {
	return stack.ErrorMessageFor(instance)
}

// Message queries the error message.
func (instance Error) Message() string {
	return instance.message
}

// Cause queries the cause that causes this error.
// If there is no cause nil is returned.
func (instance Error) Cause() interface{} {
	return instance.cause
}

// Stack queries the stack at the point this Error instance was created.
func (instance Error) Stack() stack.Stack {
	return instance.stack
}

// CausedBy is a shortcut for Cause to use a builder pattern.
func (instance Error) CausedBy(what interface{}) Error {
	instance.cause = what
	return instance
}
