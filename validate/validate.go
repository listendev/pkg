package validate

import (
	"fmt"
	"reflect"

	"github.com/garnet-org/pkg/analysisrequest"
	"github.com/garnet-org/pkg/verdictcode"
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
	Singleton.RegisterAlias("severity", "eq_ignore_case=low|eq_ignore_case=medium|eq_ignore_case=high")
	Singleton.RegisterAlias("amqp", "startswith=amqp://|startswith=amqps://")
	Singleton.RegisterAlias("store", "startswith=file:///|startswith=s3://")
	Singleton.RegisterAlias("shasum", "len=40")
	if err := Singleton.RegisterValidation("is_analysisrequest_type", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.String {
			_, err := analysisrequest.ToType(f.String())

			return err == nil
		}

		panic(fmt.Sprintf("bad field type: %T", f.Interface()))
	}); err != nil {
		panic(err)
	}
	if err := Singleton.RegisterValidation("is_verdictcode", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.Uint64 {
			_, err := verdictcode.FromUint64(f.Uint(), false)

			return err == nil
		}

		panic(fmt.Sprintf("bad field type: %T", f.Interface()))
	}); err != nil {
		panic(err)
	}
	if err := Singleton.RegisterValidation("npm_package_name", isNpmPackageName); err != nil {
		panic(err)
	}

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

	if err := Singleton.RegisterTranslation(
		"severity",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("severity", "{0} must be low, medium, or high", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("severity", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"amqp",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("amqp", "{0} must start with amqp:// or with amqps://", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("amqp", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"store",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("store", "{0} must start with s3://... or file:///...", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("store", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"shasum",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("shasum", "{0} must be a valid SHA1 (40 characters long)", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("shasum", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"is_analysisrequest_type",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_analysisrequest_type", "{0} is not a valid analysis request type", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_analysisrequest_type", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"is_verdictcode",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_verdictcode", "{0} is not a valid verdict code", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_verdictcode", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"npm_package_name",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("npm_package_name", "{0} is not a valid npm package name", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("npm_package_name", fe.Field())

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
