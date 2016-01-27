package stack

import (
	"runtime"
	"path/filepath"
	"fmt"
	"strings"
)

var (
	unknownFunction = string("<unknown function>")
	runtimeMain = string("runtime.main")
	runtimeGoexit = string("runtime.goexit")
)

type StackElement struct {
	File      string
	ShortFile string
	Line      int
	Function  string
	Package   string
	Pc        uintptr
}

func (i StackElement) String() string {
	return fmt.Sprintf("%s.%s(%s:%d)", i.Package, i.Function, i.ShortFile, i.Line)
}

type Stack []StackElement

func (i Stack) String() string {
	result := ""
	for _, element := range i {
		result += "\tat " + element.String() + "\n"
	}
	return result
}

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
				result = append(result, StackElement{
					File: file,
					ShortFile: filepath.Base(file),
					Line:  line,
					Function: functionNameOf(fullFunctionName),
					Package: packageNameOf(fullFunctionName),
					Pc: pc,
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
	if lastDot >= 0 && lastDot + 1 < len(fullFunctionName) {
		return fullFunctionName[lastDot + 1:]
	} else {
		return fullFunctionName
	}
}

func packageNameOf(fullFunctionName string) string {
	lastDot := strings.LastIndex(fullFunctionName, ".")
	if lastDot > 0 {
		return fullFunctionName[:lastDot]
	} else {
		return ""
	}
}

func fullFunctionNameOf(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return unknownFunction
	}
	return fn.Name()
}
