package stack

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// MessageEnabled represents an object that has a Message.
type MessageEnabled interface {
	Message() string
}

// ErrorEnabled represents an object that has an Error message.
type ErrorEnabled interface {
	Error() string
}

// HasStack represents an object that has a Stack.
type HasStack interface {
	Stack() Stack
}

// CauseEnabled represents an object that could have a Cause.
type CauseEnabled interface {
	Cause() interface{}
}

// ErrorMessageFor creates an error message for given error.
func ErrorMessageFor(what error) string {
	message := ""
	var current interface{}
	current = what
	for i := 0; current != nil; i++ {
		if m, ok := current.(MessageEnabled); ok {
			message += messageFor(i, m.Message())
		} else if m, ok := current.(ErrorEnabled); ok {
			message += messageFor(i, m.Error())
		} else {
			message += messageFor(i, fmt.Sprintf("%v", current))
		}
		if m, ok := current.(CauseEnabled); ok {
			current = m.Cause()
		} else {
			break
		}
	}
	return message
}

func messageFor(i int, message string) string {
	result := ""
	if i > 0 {
		result = "\n\tCaused by: "
	}
	result += message
	return result
}

// StringOf converts the given problem object (panic, error, ...) to a string which is readable.
func StringOf(what interface{}, framesToSkip int) string {
	buf := new(bytes.Buffer)
	Print(what, buf, framesToSkip+1)
	return buf.String()
}

// Print prints the given problem object (panic, error, ...) to a readable version to writer.
func Print(what interface{}, to io.Writer, framesToSkip int) {
	prefix := ""
	var current interface{}
	current = what
	for i := 0; current != nil; i++ {
		message := typeToString(current) + ": "
		if m, ok := current.(MessageEnabled); ok {
			message += m.Message()
		} else if m, ok := current.(ErrorEnabled); ok {
			message += m.Error()
		} else {
			message += fmt.Sprintf("%v", current)
		}
		fmt.Fprint(to, prefix, message, "\n")
		if m, ok := current.(HasStack); ok {
			fmt.Fprintf(to, "%v", m.Stack())
		} else if i == 0 {
			stack := CaptureStack(1 + framesToSkip)
			fmt.Fprintf(to, "%v", stack)
		}
		if m, ok := current.(CauseEnabled); ok {
			current = m.Cause()
		} else {
			break
		}
		if i >= 0 {
			prefix = "Caused by: "
		}
	}
}

func typeToString(of interface{}) string {
	t := reflect.TypeOf(of)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.String()
}
