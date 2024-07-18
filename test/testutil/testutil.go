package testutil

import (
	"errors"
	"reflect"
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
