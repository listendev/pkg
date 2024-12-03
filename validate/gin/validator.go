package ginval

import (
	"reflect"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/listendev/pkg/validate"
)

type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

func New() *Validator {
	return &Validator{
		validate: validate.Singleton,
	}
}

func (v *Validator) lazyinit() {
	v.once.Do(func() {
		v.validate = validate.Singleton
	})
}

func (v *Validator) Engine() any {
	v.lazyinit()

	return v.validate
}

func (v *Validator) validateStruct(obj any) error {
	v.lazyinit()

	return v.validate.Struct(obj)
}

func (v *Validator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return eventuallyMany(v.validateStruct(obj))
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)

		for i := range count {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}

		return validateRet
	default:
		return nil
	}
}

func eventuallyMany(err error) error {
	if err == nil {
		return nil
	}

	validationErrors, _ := err.(validate.ValidationError)
	// TODO: check the input not having `validate.ValidationErrors` type
	if len(validationErrors) == 0 {
		return nil
	}

	ret := make(binding.SliceValidationError, 0)
	for _, ve := range validationErrors {
		ret = append(ret, ve)
	}

	return ret
}
