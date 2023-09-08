package utils

import (
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
)

func StructToString(t interface{}) (string, error) {
	if reflect.TypeOf(t).Kind() != reflect.Struct {
		return "", errors.New("payload is not struct")
	}
	bytes, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	reg, err := regexp.Compile(`("([a-zA-Z0-9_]+)":null,?)|"`)
	if err != nil {
		return "", err
	}
	str := reg.ReplaceAllString(string(bytes), "")
	return str, nil
}
