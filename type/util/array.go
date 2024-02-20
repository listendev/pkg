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
