package list

import (
	"errors"
	"reflect"
)

func ListIn(list interface{}, element interface{}) (bool, error) {
	sVal := reflect.ValueOf(list)
	kind := sVal.Kind()
	if kind == reflect.Slice || kind == reflect.Array {
		for i := 0; i < sVal.Len(); i++ {
			if sVal.Index(i).Interface() == element {
				return true, nil
			}
		}
		return false, nil
	}
	return false, errors.New("list data type error，is not \"selic \"or \"array\"")
}
