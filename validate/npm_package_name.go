package validate

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/leodido/go-npmpackagename"
)

func isNpmPackageName(fl validator.FieldLevel) bool {
	field := fl.Field()
	// Do you want strict validation or not?
	strict := false
	if fl.Param() == "strict" {
		strict = true
	}

	if field.Kind() == reflect.String {
		valid, warnings, err := npmpackagename.Validate([]byte(field.String()))
		if strict {
			if err == nil && len(warnings) == 0 {
				return valid
			}
		} else {
			if err == nil {
				return valid
			}
		}

		return false
	}

	panic(fmt.Sprintf("bad field type: %T", field.Interface()))
}
