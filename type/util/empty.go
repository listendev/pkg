package typeutil

import (
	"fmt"
	"reflect"
	"strings"
)

func IsEmpty(value interface{}) bool {
	valueV := reflect.ValueOf(value)

	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	switch valueV.Kind() {
	case reflect.Struct:
		if IsZero(value) {
			return true
		}

	case reflect.Array, reflect.Slice:
		if valueV.Len() == 0 {
			return true
		}
		for i := 0; i < valueV.Len(); i++ {
			if indexV := valueV.Index(i); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
				return false
			}
		}

		return true

	case reflect.Map:
		if valueV.Len() == 0 {
			return true
		}
		for _, keyV := range valueV.MapKeys() {
			if indexV := valueV.MapIndex(keyV); indexV.IsValid() && !IsEmpty(indexV.Interface()) {
				return false
			}
		}

		return true

	case reflect.Chan:
		if valueV.Len() == 0 {
			return true
		}

	case reflect.String:
		if len(strings.TrimSpace(fmt.Sprintf("%v", value))) == 0 {
			return true
		}
	}

	return valueV.IsZero()
}
