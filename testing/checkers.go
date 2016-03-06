package testing

import (
	"fmt"
	"github.com/echocat/caretakerd/values"
	"gopkg.in/check.v1"
	"reflect"
	"regexp"
)

type hasStringMethod interface {
	String() string
}

type throwsPanicThatMatches struct {
	*check.CheckerInfo
}

var ThrowsPanicThatMatches check.Checker = &throwsPanicThatMatches{
	&check.CheckerInfo{
		Name:   "ThrowsPanicThatMatches",
		Params: []string{"action", "regex"},
	},
}

func (checker *throwsPanicThatMatches) Check(params []interface{}, names []string) (result bool, errorMessage string) {
	if len(params) != 2 {
		panic("Illegal number of parameters.")
	}
	action, ok := params[0].(func())
	if !ok {
		return false, "There is no <action> function with signature func() provided."
	}
	plainRegex, ok := params[1].(string)
	if !ok {
		return false, "There is no <regex> of type string provided."
	}
	regex, err := regexp.Compile(plainRegex)
	if err != nil {
		return false, fmt.Sprintf("Can't compile regex '%s'. Got: %s", plainRegex, err.Error())
	}
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = regex.MatchString(err.Error())
			} else if rs, ok := r.(hasStringMethod); ok {
				result = regex.MatchString(rs.String())
			} else if rs, ok := r.(string); ok {
				result = regex.MatchString(rs)
			} else {
				panic(r)
			}
		}
	}()
	action()
	return false, ""
}

type isEmptyChecker struct {
	*check.CheckerInfo
}

var IsEmpty check.Checker = &isEmptyChecker{
	&check.CheckerInfo{
		Name:   "IsEmpty",
		Params: []string{"value"},
	},
}

func (checker *isEmptyChecker) Check(params []interface{}, names []string) (bool, string) {
	if len(params) != 1 {
		panic("Illegal number of parameters.")
	}
	param := params[0]
	if asString, ok := param.(string); ok {
		return len(asString) == 0, ""
	}
	if asString, ok := param.(values.String); ok {
		return len(asString) == 0, ""
	}
	pv := reflect.ValueOf(param)
	switch pv.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		fallthrough
	case reflect.Map:
		return pv.Len() == 0, ""
	}
	return false, "No compatible value."
}

type isLessThan struct {
	*check.CheckerInfo
}

var IsLessThan check.Checker = &isLessThan{
	&check.CheckerInfo{
		Name:   "IsLessThan",
		Params: []string{"obtained", "compareTo"},
	},
}

func (checker *isLessThan) Check(params []interface{}, names []string) (bool, string) {
	if len(params) != 2 {
		panic("Illegal number of parameters.")
	}
	obtained := params[0]
	compareTo := params[1]
	obtainedType := reflect.TypeOf(obtained)
	compareToType := reflect.TypeOf(compareTo)
	if !reflect.DeepEqual(obtainedType, compareToType) {
		return false, fmt.Sprintf("'obtained' type not equal to the type to 'compareTo' type.")
	}
	if casted, ok := obtained.(uint8); ok {
		return casted < compareTo.(uint8), ""
	}
	if casted, ok := obtained.(uint16); ok {
		return casted < compareTo.(uint16), ""
	}
	if casted, ok := obtained.(uint32); ok {
		return casted < compareTo.(uint32), ""
	}
	if casted, ok := obtained.(uint64); ok {
		return casted < compareTo.(uint64), ""
	}
	if casted, ok := obtained.(int8); ok {
		return casted < compareTo.(int8), ""
	}
	if casted, ok := obtained.(int16); ok {
		return casted < compareTo.(int16), ""
	}
	if casted, ok := obtained.(int32); ok {
		return casted < compareTo.(int32), ""
	}
	if casted, ok := obtained.(int64); ok {
		return casted < compareTo.(int64), ""
	}
	if casted, ok := obtained.(float32); ok {
		return casted < compareTo.(float32), ""
	}
	if casted, ok := obtained.(float64); ok {
		return casted < compareTo.(float64), ""
	}
	if casted, ok := obtained.(int); ok {
		return casted < compareTo.(int), ""
	}
	if casted, ok := obtained.(uint); ok {
		return casted < compareTo.(uint), ""
	}
	if casted, ok := obtained.(uintptr); ok {
		return casted < compareTo.(uintptr), ""
	}
	if casted, ok := obtained.(byte); ok {
		return casted < compareTo.(byte), ""
	}
	if casted, ok := obtained.(rune); ok {
		return casted < compareTo.(rune), ""
	}
	return false, "No compatible type."
}

type isLessThanOrEqualTo struct {
	*check.CheckerInfo
}

var IsLessThanOrEqualTo check.Checker = &isLessThanOrEqualTo{
	&check.CheckerInfo{
		Name:   "IsLessThanOrEqualTo",
		Params: []string{"obtained", "compareTo"},
	},
}

func (checker *isLessThanOrEqualTo) Check(params []interface{}, names []string) (bool, string) {
	result, err := IsLessThan.Check(params, names)
	if result || err != "" {
		return result, err
	}
	return check.Equals.Check(params, names)
}

type isLargerThan struct {
	*check.CheckerInfo
}

var IsLargerThan check.Checker = &isLargerThan{
	&check.CheckerInfo{
		Name:   "IsLargerThan",
		Params: []string{"obtained", "compareTo"},
	},
}

func (checker *isLargerThan) Check(params []interface{}, names []string) (bool, string) {
	if len(params) != 2 {
		panic("Illegal number of parameters.")
	}
	obtained := params[0]
	compareTo := params[1]
	obtainedType := reflect.TypeOf(obtained)
	compareToType := reflect.TypeOf(compareTo)
	if !reflect.DeepEqual(obtainedType, compareToType) {
		return false, fmt.Sprintf("'obtained' type not equal to the type to 'compareTo' type.")
	}
	if casted, ok := obtained.(uint8); ok {
		return casted > compareTo.(uint8), ""
	}
	if casted, ok := obtained.(uint16); ok {
		return casted > compareTo.(uint16), ""
	}
	if casted, ok := obtained.(uint32); ok {
		return casted > compareTo.(uint32), ""
	}
	if casted, ok := obtained.(uint64); ok {
		return casted > compareTo.(uint64), ""
	}
	if casted, ok := obtained.(int8); ok {
		return casted > compareTo.(int8), ""
	}
	if casted, ok := obtained.(int16); ok {
		return casted > compareTo.(int16), ""
	}
	if casted, ok := obtained.(int32); ok {
		return casted > compareTo.(int32), ""
	}
	if casted, ok := obtained.(int64); ok {
		return casted > compareTo.(int64), ""
	}
	if casted, ok := obtained.(float32); ok {
		return casted > compareTo.(float32), ""
	}
	if casted, ok := obtained.(float64); ok {
		return casted > compareTo.(float64), ""
	}
	if casted, ok := obtained.(int); ok {
		return casted > compareTo.(int), ""
	}
	if casted, ok := obtained.(uint); ok {
		return casted > compareTo.(uint), ""
	}
	if casted, ok := obtained.(uintptr); ok {
		return casted > compareTo.(uintptr), ""
	}
	if casted, ok := obtained.(byte); ok {
		return casted > compareTo.(byte), ""
	}
	if casted, ok := obtained.(rune); ok {
		return casted > compareTo.(rune), ""
	}
	return false, "No compatible type."
}

type isLargerThanOrEqualTo struct {
	*check.CheckerInfo
}

var IsLargerThanOrEqualTo check.Checker = &isLargerThanOrEqualTo{
	&check.CheckerInfo{
		Name:   "IsLargerThanOrEqualTo",
		Params: []string{"obtained", "compareTo"},
	},
}

func (checker *isLargerThanOrEqualTo) Check(params []interface{}, names []string) (bool, string) {
	result, err := IsLargerThan.Check(params, names)
	if result || err != "" {
		return result, err
	}
	return check.Equals.Check(params, names)
}
