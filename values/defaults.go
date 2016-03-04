package values

import (
	"reflect"
)

// SetDefaultsTo takes defaults from given defaults map and apply it to given to object.
func SetDefaultsTo(defaults map[string]interface{}, to interface{}) interface{} {
	toValue := reflect.ValueOf(to)
	if toValue.Kind() == reflect.Ptr {
		toValue = toValue.Elem()
	}
	for candidateName, candidate := range defaults {
		candidateValue := toValue.FieldByName(candidateName)
		candidateValue.Set(reflect.ValueOf(candidate))
	}
	return to
}

// IsDefaultValue is checking if the given fieldName in the given defaults map has the same value as the given one.
func IsDefaultValue(defaults map[string]interface{}, fieldName string, value interface{}) bool {
	if defaultValue, ok := defaults[fieldName]; ok {
		return reflect.DeepEqual(defaultValue, value)
	}
	return false
}

// IsDefaultReflectValue is checking if the given field in the given defaults map has the same value as the given one.
func IsDefaultReflectValue(defaults map[string]interface{}, field reflect.StructField, value reflect.Value) bool {
	return IsDefaultValue(defaults, field.Name, value.Interface())
}
