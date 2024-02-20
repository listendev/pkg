package typeutil

import "reflect"

// Dectect whether the concrete underlying value of the given input is one or more
// Kinds of value.
func IsKind(in interface{}, kinds ...reflect.Kind) bool {
	var inT reflect.Type

	if v, ok := in.(reflect.Value); ok && v.IsValid() {
		inT = v.Type()
	} else if v, ok := in.(reflect.Type); ok {
		inT = v
	} else {
		in = ResolveValue(in)
		inT = reflect.TypeOf(in)
	}

	if inT == nil {
		return false
	}

	for _, k := range kinds {
		if inT.Kind() == k {
			return true
		}
	}

	return false
}

func ResolveValue(in interface{}) interface{} {
	var inV reflect.Value

	if vV, ok := in.(reflect.Value); ok {
		inV = vV
	} else {
		inV = reflect.ValueOf(in)
	}

	if inV.IsValid() {
		if inT := inV.Type(); inT == nil {
			return nil
		}

		switch inV.Kind() {
		case reflect.Ptr, reflect.Interface:
			return ResolveValue(inV.Elem())
		}

		in = inV.Interface()
	}

	return in
}
