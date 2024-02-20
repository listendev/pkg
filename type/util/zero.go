package typeutil

import "reflect"

func IsZero(value interface{}) bool {
	if value == nil {
		return true
	} else if valueV, ok := value.(reflect.Value); ok && valueV.IsValid() {
		if valueV.CanInterface() {
			value = valueV.Interface()
		}
	}

	return reflect.DeepEqual(
		value,
		reflect.Zero(reflect.TypeOf(value)).Interface(),
	)
}
