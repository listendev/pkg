package typeutil

import (
	"reflect"

	"github.com/ghetzel/go-stockutil/utils"
)

func IsArray(in interface{}) bool {
	return IsKind(in, utils.SliceTypes...)
}

var SliceTypes = []reflect.Kind{
	reflect.Slice,
	reflect.Array,
}

// RemoveFromSliceAt deletes an element at position `at` from the slice.
//
// Order matters.
func RemoveFromSliceAt[T any](slice []T, at int) []T {
	return append(slice[:at], slice[at+1:]...)
}
