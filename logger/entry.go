package logger

import (
	"fmt"
	"github.com/echocat/caretakerd/errors"
	"github.com/echocat/caretakerd/stack"
	"github.com/eknkc/dateformat"
	"strconv"
	"time"
)

// Entry represents an entry to be logged.
type Entry struct {
	Time     time.Time
	Message  string
	Priority Level
	Category string
	Stack    stack.Stack
	Uptime   time.Duration
	Problem  interface{}
}

// NewEntry creates a new instance of Entry.
func NewEntry(framesToSkip int, problem interface{}, category string, prioriy Level, time time.Time, message string, uptime time.Duration) Entry {
	return Entry{
		Time:     time,
		Message:  message,
		Priority: prioriy,
		Category: category,
		Stack:    stack.CaptureStack(framesToSkip + 1),
		Uptime:   uptime,
		Problem:  problem,
	}
}

// Format formats the current entry using the given pattern.
func (e Entry) Format(pattern Pattern, framesToSkip int) (string, error) {
	result := []byte{}
	flag := byte(0)
	flagStarted := false
	flagStart := 0
	flagFormat := []byte{}
	flagArgumentsStarted := false
	flagArguments := []byte{}
	for position := 0; position < len(pattern); position++ {
		c := pattern[position]
		if len(flagFormat) > 0 {
			if flagArgumentsStarted {
				if c == '{' {
					return "", NewFormatError(position, "Unexpedted character %c at instance position within flag argument %c.", c, flag)
				} else if c == '}' {
					flagPlainContent, err := e.contentOf(position, flag, string(flagArguments), framesToSkip+1)
					if err != nil {
						return "", err
					}
					//noinspection GoPlaceholderCount
					flagContent := fmt.Sprintf(string(append(flagFormat, 's')), flagPlainContent)
					result = append(result, []byte(flagContent)...)
					flag = 0
					flagFormat = []byte{}
					flagArgumentsStarted = false
					flagArguments = []byte{}
					flagStarted = false
				} else {
					flagArguments = append(flagArguments, c)
				}
			} else if c == '%' {
				if flagStarted {
					return "", NewFormatError(position, "Unexpedted character %c at instance position within flag %c.", c, flag)
				}
				flagFormat = []byte{}
				result = append(result, c)
			} else if c == '*' || c == '.' || c == '-' || c == ' ' || (c >= '0' && c <= '9') {
				flagStarted = true
				flagFormat = append(flagFormat, c)
			} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				flagStarted = true
				flag = c
				if (len(pattern) > position+1) && pattern[position+1] == '{' {
					flagArgumentsStarted = true
					position++
				} else {
					flagPlainContent, err := e.contentOf(position, flag, "", framesToSkip+1)
					if err != nil {
						return "", err
					}
					//noinspection GoPlaceholderCount
					flagContent := fmt.Sprintf(string(append(flagFormat, 's')), flagPlainContent)
					result = append(result, []byte(flagContent)...)
					flag = 0
					flagFormat = []byte{}
					flagStarted = false
				}
			} else {
				return "", NewFormatError(position, "Unexpedted character %c at instance position within flag %c.", c, flag)
			}
		} else if c == '%' {
			flagStart = position
			flagFormat = []byte{c}
		} else {
			result = append(result, c)
		}
	}
	if len(flagFormat) > 0 {
		return "", NewFormatError(flagStart, "Uncompleted flag.")
	}
	return string(result), nil
}

func (e Entry) contentOf(position int, flag byte, arguments string, framesToSkip int) (string, error) {
	switch flag {
	case 'd':
		return dateformat.Format(e.Time, arguments), nil
	case 'm':
		return e.Message, nil
	case 'c':
		return e.cutLeftSideSegmentsOfS(position, flag, e.Category, arguments, isCategorySegmentSeparator)
	case 'F':
		return e.cutLeftSideSegmentsOfS(position, flag, e.Stack[0].File, arguments, isFileSegmentSeparator)
	case 'l':
		return e.Stack[0].String(), nil
	case 'L':
		return strconv.Itoa(e.Stack[0].Line), nil
	case 'C':
		return e.cutLeftSideSegmentsOfS(position, flag, e.Stack[0].Package, arguments, isCategorySegmentSeparator)
	case 'M':
		return e.Stack[0].Function, nil
	case 'p':
		return e.Priority.DisplayForLogging(), nil
	case 'P':
		return e.formatProblemIfNeeded(arguments, framesToSkip+1)
	case 'r':
		return fmt.Sprintf("%d", e.Uptime/time.Millisecond), nil
	case 'n':
		return "\n", nil
	}
	return "", NewFormatError(position, "Unknown flag '%c'.", flag)
}

func (e Entry) cutLeftSideSegmentsOfS(position int, flag byte, in string, maximumAsString string, isSegmentSeparator func(byte) bool) (string, error) {
	if len(maximumAsString) <= 0 {
		return in, nil
	}
	maximum, err := strconv.Atoi(maximumAsString)
	if err != nil || maximum <= 0 {
		return "", NewFormatError(position, "'%s' is not a valid number for argument of flag '%c'.", maximumAsString, flag)
	}
	return e.cutLeftSideSegmentsOf(in, maximum, isSegmentSeparator), nil
}

func (e Entry) cutLeftSideSegmentsOf(in string, maximum int, isSegmentSeparator func(byte) bool) string {
	result := []byte{}
	numberOfSeparators := 0
	for i := len(in) - 1; i >= 0 && numberOfSeparators < maximum; i-- {
		c := in[i]
		if isSegmentSeparator(c) {
			numberOfSeparators++
			if numberOfSeparators < maximum {
				result = append([]byte{c}, result...)
			}
		} else {
			result = append([]byte{c}, result...)
		}
	}
	return string(result)
}

func isCategorySegmentSeparator(c byte) bool {
	return c == '.' || c == '-' || c == '_' || c == '/' || c == '\\'
}

func isFileSegmentSeparator(c byte) bool {
	return c == '/' || c == '\\'
}

func (e Entry) formatProblemIfNeeded(arguments string, framesToSkip int) (string, error) {
	problem := e.Problem
	if problem != nil {
		problemAsString := stack.StringOf(problem, framesToSkip+1)
		subEntry := NewEntry(framesToSkip+1, nil, e.Category, e.Priority, e.Time, problemAsString, e.Uptime)
		pattern := arguments
		if len(pattern) <= 0 {
			pattern = "%n%m"
		}
		return subEntry.Format(Pattern(pattern), framesToSkip+1)
	}
	return "", nil
}

// FormatError represents an error if a given pattern contains wrong arguments.
type FormatError struct {
	Message  string
	Position int
}

// NewFormatError creates a new instance of FormatError.
func NewFormatError(position int, message string, a ...interface{}) FormatError {
	return FormatError{
		Message:  errors.New(message, a...).Message(),
		Position: position,
	}
}

func (e FormatError) Error() string {
	return fmt.Sprintf("At index %d: %s", e.Position, e.Message)
}
