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
		stack: stack.CaptureStack(1),
	}
	return result
}

func (this Error) String() string {
	return reflect.TypeOf(this).Name() + ": " + this.Message()
}

func (this Error) Error() string {
	return stack.ErrorMessageFor(this)
}

func (this Error) Message() string {
	return this.message
}

func (this Error) Cause() interface{} {
	return this.cause
}

func (this Error) Stack() stack.Stack {
	return this.stack
}

func (this Error) CausedBy(what interface{}) Error {
	this.cause = what
	return this
}
