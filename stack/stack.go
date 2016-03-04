package stack

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	unknownFunction = string("<unknown function>")
	runtimeMain     = string("runtime.main")
	runtimeGoexit   = string("runtime.goexit")
)

// Element represents an element from the whole stack trace.
type Element struct {
	File      string
	ShortFile string
	Line      int
	Function  string
	Package   string
	Pc        uintptr
}

func (i Element) String() string {
	return fmt.Sprintf("%s.%s(%s:%d)", i.Package, i.Function, i.ShortFile, i.Line)
}

// Stack represents the whole stack trace with a couple of Elements.
type Stack []Element

func (i Stack) String() string {
	result := ""
	for _, element := range i {
		result += "\tat " + element.String() + "\n"
	}
	return result
}

// CaptureStack creates a new stack capture of the current stack.
// It is possible to cut off the returned stack with framesToSkip.
func CaptureStack(framesToSkip int) Stack {
	result := Stack{}
	valid := true
	for i := framesToSkip + 1; valid; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if ok {
			fullFunctionName := fullFunctionNameOf(pc)
			if fullFunctionName == runtimeMain || fullFunctionName == runtimeGoexit {
				valid = false
			} else {
				result = append(result, Element{
					File:      file,
					ShortFile: filepath.Base(file),
					Line:      line,
					Function:  functionNameOf(fullFunctionName),
					Package:   packageNameOf(fullFunctionName),
					Pc:        pc,
				})
			}
		} else {
			valid = false
		}
	}
	return result
}

func functionNameOf(fullFunctionName string) string {
	lastDot := strings.LastIndex(fullFunctionName, ".")
	if lastDot >= 0 && lastDot+1 < len(fullFunctionName) {
		return fullFunctionName[lastDot+1:]
	}
	return fullFunctionName
}

func packageNameOf(fullFunctionName string) string {
	lastDot := strings.LastIndex(fullFunctionName, ".")
	if lastDot > 0 {
		return fullFunctionName[:lastDot]
	}
	return ""
}

func fullFunctionNameOf(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return unknownFunction
	}
	return fn.Name()
}
