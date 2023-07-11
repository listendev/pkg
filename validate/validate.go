package validate

import (
	"fmt"
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/listendev/pkg/analysisrequest"
	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/models/category"
	"github.com/listendev/pkg/models/severity"
	"github.com/listendev/pkg/verdictcode"
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
	Singleton.RegisterAlias("amqp", "startswith=amqp://|startswith=amqps://")
	Singleton.RegisterAlias("store", "startswith=file:///|startswith=s3://")
	Singleton.RegisterAlias("shasum", "len=40")
	Singleton.RegisterAlias("npmorg", "startswith=@")

	if err := Singleton.RegisterValidation("is_severity", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.String {
			_, err := severity.New(f.String())

			return err == nil
		}

		panic(fmt.Sprintf("bad field type: %T", f.Interface()))
	}); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterValidation("is_category", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.Uint64 {
			_, err := category.FromUint64(f.Uint())

			return err == nil
		}

		panic(fmt.Sprintf("bad field type: %T", f.Interface()))
	}); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterValidation("is_ecosystem", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.Uint64 {
			_, err := ecosystem.FromUint64(f.Uint())

			return err == nil
		}

		panic(fmt.Sprintf("bad field type: %T", f.Interface()))
	}); err != nil {
		panic(err)
	}

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

	if err := Singleton.RegisterValidation("is_resultsfile", func(fl validator.FieldLevel) bool {
		f := fl.Field()

		if f.Kind() == reflect.String {
			_, err := analysisrequest.GetTypeFromResultFile(f.String())

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
		"is_severity",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_severity", "{0} must be low, medium, or high", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_severity", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"is_ecosystem",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_ecosystem", fmt.Sprintf("{0} must be on of [%s]", strings.Join(ecosystem.Ecosystems(), ", ")), true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_ecosystem", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"isdefault|is_severity",
		Translator,

		func(ut ut.Translator) error {
			return ut.Add("isdefault|is_severity", "{0} must be low, medium, or high", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isdefault|is_severity", fe.Field())

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
		"npmorg",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("npmorg", "{0} must be a valid NPM organization (starting with @)", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("npmorg", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"semver",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("semver", "{0} must be a valid semantic version (https://semver.org)", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("semver", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"is_resultsfile",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_resultsfile", "{0} must be a valid results file", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("is_resultsfile", fe.Field())

			return t
		},
	); err != nil {
		panic(err)
	}

	if err := Singleton.RegisterTranslation(
		"is_category",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("is_category", "{0} is not a valid verdict category", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			// FIXME: handle cardinals translation
			t, _ := ut.T("is_category", fe.Field())

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
		"isdefault|is_verdictcode",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("isdefault|is_verdictcode", "{0} is not a valid verdict code", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("isdefault|is_verdictcode", fe.Field())

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

	if err := Singleton.RegisterTranslation(
		"required_with",
		Translator,
		func(ut ut.Translator) error {
			return ut.Add("required_with", "{0} is required when the {1} field has a value", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required_with", fe.Field(), strings.ToLower(fe.Param()))

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
