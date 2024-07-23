package testutil

import (
	"errors"
	"reflect"
	"strings"
)

// errIsNil checks if err is nil, including nil pointers.
// A nil pointer makes err == nil return false because
// it's an interface with a non-nil type but nil value.
// Therefore, we must check if the value is nil.
func errIsNil(err error) bool {
	return err == nil || reflect.ValueOf(err).IsNil()
}

// IsSameErr checks if err is the same as wantErr.
func IsSameErr(err, wantErr error) bool {
	if errIsNil(err) && errIsNil(wantErr) {
		return true
	}

	if errIsNil(err) || errIsNil(wantErr) {
		return false
	}

	return errors.Is(err, wantErr)
}

// CompareAsPtr checks if 'val' is equal to 'cmp'.
// It handles both pointers and non-pointers.
// If either 'val' or 'cmp' is a pointer, it converts both to pointers
// before performing the comparison. This approach is useful for cases
// like comparing an empty string to a nil string pointer, where the expected
// result should be true.
func CompareAsPtr[V, C any](val V, cmp C) bool {
	valType := reflect.TypeOf(val).String()
	cmpType := reflect.TypeOf(cmp).String()

	valTypeDisregardingPtr := strings.ReplaceAll(valType, "*", "")
	cmpTypeDisregardingPtr := strings.ReplaceAll(cmpType, "*", "")

	if valTypeDisregardingPtr != cmpTypeDisregardingPtr {
		return false
	}

	var zeroVal V
	var zeroCmp C
	if reflect.DeepEqual(val, zeroVal) && reflect.DeepEqual(cmp, zeroCmp) {
		return true
	}

	valValue := reflect.ValueOf(val)
	cmpValue := reflect.ValueOf(cmp)

	if valValue.Kind() != reflect.Ptr {
		valValue = reflect.New(valValue.Type())
		valValue.Elem().Set(reflect.ValueOf(val))
	}

	if cmpValue.Kind() != reflect.Ptr {
		cmpValue = reflect.New(cmpValue.Type())
		cmpValue.Elem().Set(reflect.ValueOf(cmp))
	}

	return reflect.DeepEqual(valValue.Interface(), cmpValue.Interface())
}
