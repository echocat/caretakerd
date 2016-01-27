package errors

import (
	"fmt"
	"github.com/echocat/caretakerd/stack"
	"reflect"
)

type Error struct {
	message string
	cause   interface{}
	stack   stack.Stack
}

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

func (instance Error) Message() string {
	return instance.message
}

func (instance Error) Cause() interface{} {
	return instance.cause
}

func (instance Error) Stack() stack.Stack {
	return instance.stack
}

func (instance Error) CausedBy(what interface{}) Error {
	instance.cause = what
	return instance
}
