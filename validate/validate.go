package validate

import (
	"fmt"
	"reflect"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type ValidationErrors = validator.ValidationErrors

// Singleton is the validator singleton instance.
//
// This way it caches the structs info.
var Singleton *validator.Validate

// Translator is the universal translator for validation errors.
var Translator ut.Translator

func init() {
	Singleton = validator.New()

	// Register a function to get the field name from "flag" tags.
	Singleton.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("human")
		if name == "-" {
			return ""
		}

		return name
	})

	Singleton.RegisterAlias("mandatory", "required")

	eng := en.New()
	Translator, _ = (ut.New(eng, eng)).GetTranslator("en")
	if err := en_translations.RegisterDefaultTranslations(Singleton, Translator); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"mandatory",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("mandatory", "{0} is mandatory", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("mandatory", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}
}

func Validate(o interface{}) []error {
	if err := Singleton.Struct(o); err != nil {
		all := []error{}
		for _, e := range err.(ValidationErrors) {
			all = append(all, fmt.Errorf(e.Translate(Translator)))
		}

		return all
	}

	return nil
}
