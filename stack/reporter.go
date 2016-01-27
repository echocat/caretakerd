package stack

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type MessageEnabled interface {
	Message() string
}

type ErrorEnabled interface {
	Error() string
}

type StackEnabled interface {
	Stack() Stack
}

type CauseEnabled interface {
	Cause() interface{}
}

func ErrorMessageFor(what error) string {
	message := ""
	var current interface{}
	current = what
	for i := 0; current != nil; i++ {
		if m, ok := current.(MessageEnabled); ok {
			message += m.Message()
		} else if m, ok := current.(ErrorEnabled); ok {
			message += m.Error()
		} else {
			message += fmt.Sprintf("%v", current)
		}
		if m, ok := current.(CauseEnabled); ok {
			current = m.Cause()
		} else {
			break
		}
		if i >= 0 {
			message += " Caused by: "
		}
	}
	return message
}

func StringOf(what interface{}, framesToSkip int) string {
	buf := new(bytes.Buffer)
	Print(what, buf, framesToSkip+1)
	return buf.String()
}

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
		if m, ok := current.(StackEnabled); ok {
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
