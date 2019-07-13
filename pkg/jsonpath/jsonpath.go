package jsonpath

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PaesslerAG/jsonpath"
	"reflect"
	"strings"
)

// https://github.com/steinfletcher/apitest-jsonpath/blob/master/jsonpath.go

// Contains is a convenience function to assert that a jsonpath expression extracts a value in an array
func Contains(bytes []byte, expression string, expected interface{}) (bool, error) {
	value, err := jsonPath(bytes, expression)
	if err != nil {
		return false, err
	}

	ok, found := includesElement(value, expected)
	if !ok {
		return false, errors.New(fmt.Sprintf("\"%s\" could not be applied builtin len()", expected))
	}
	if !found {
		return false, errors.New(fmt.Sprintf("\"%s\" does not contain \"%s\"", expected, value))
	}
	return true, nil
}

// Equal is a convenience function to assert that a jsonpath expression extracts a value
func Equal(bytes []byte, expression string, expected interface{}) (bool, error) {
	value, err := jsonPath(bytes, expression)
	if err != nil {
		return false, err
	}

	if !objectsAreEqual(value, expected) {
		return false, errors.New(fmt.Sprintf("\"%s\" not equal to \"%s\"", value, expected))
	}
	return true, nil
}

func jsonPath(b []byte, expression string) (interface{}, error) {
	v := interface{}(nil)
	err := json.Unmarshal(b, &v)
	if err != nil {
		return nil, err
	}

	value, err := jsonpath.Get(expression, v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// courtesy of github.com/stretchr/testify
func includesElement(list interface{}, element interface{}) (ok, found bool) {
	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)
	defer func() {
		if e := recover(); e != nil {
			ok = false
			found = false
		}
	}()

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if objectsAreEqual(mapKeys[i].Interface(), element) {
				return true, true
			}
		}
		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if objectsAreEqual(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}
	return true, false
}

func objectsAreEqual(expected, actual interface{}) bool {
	if expected == nil || actual == nil {
		return expected == actual
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, actual)
	}

	act, ok := actual.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}
