package utils

import "github.com/spf13/cast"

// cast type
func CastType(value, valueType string) interface{} {
	if valueType == "float64" {
		return cast.ToFloat64(value)
	} else if valueType == "string" {
		return value
	} else if valueType == "int" {
		return cast.ToInt(value)
	}
	return value
}
