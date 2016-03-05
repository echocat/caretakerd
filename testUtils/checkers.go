package testUtils

import (
	"github.com/echocat/caretakerd/values"
	"gopkg.in/check.v1"
	"reflect"
)

type isEmptyChecker struct {
	*check.CheckerInfo
}

var IsEmpty check.Checker = &isEmptyChecker{
	&check.CheckerInfo{
		Name:   "IsEmpty",
		Params: []string{"value"},
	},
}

func (checker *isEmptyChecker) Check(params []interface{}, names []string) (result bool, error string) {
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
