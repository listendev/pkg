package validate

import (
	"strings"
	"testing"

	informationaltype "github.com/listendev/pkg/informational/type"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type isInformationalTypetest struct {
	AsEnum      informationaltype.Event `validate:"is_informational_event_type=enum"`
	AsSnakeCase string                  `validate:"is_informational_event_type=case"`
	AsCamelCase string                  `validate:"is_informational_event_type"`
}

type badIsInformationalTypeTest struct {
	Wrong []uint8 `validate:"is_informational_event_type"`
}

func TestIsInformationalTypeValidatorBadFieldType(t *testing.T) {
	input := badIsInformationalTypeTest{
		Wrong: []uint8("change_summary"),
	}
	require.PanicsWithValue(t, "bad field type: []uint8", func() {
		//nolint:errcheck // we are checking it panics
		Singleton.Struct(input)
	})
}

func TestIsInformationalTypeValidator(t *testing.T) {
	allValid := isInformationalTypetest{
		AsEnum:      informationaltype.DetectionsSummary,
		AsSnakeCase: "detections_summary",
		AsCamelCase: "DetectionsSummary",
	}
	require.NoError(t, Singleton.Struct(allValid))

	invalid := isInformationalTypetest{
		AsEnum:      informationaltype.ChangeSummary,
		AsSnakeCase: "xyz_escape",
		AsCamelCase: "XYZEscape",
	}
	fieldsWithError := map[string]bool{
		"AsSnakeCase": true,
		"AsCamelCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid informational event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestNoneIsNotAValidInformationalType(t *testing.T) {
	invalid := isInformationalTypetest{
		AsEnum:      informationaltype.None,
		AsSnakeCase: "none",
		AsCamelCase: "None",
	}
	fieldsWithError := map[string]bool{
		"AsEnum":      true,
		"AsSnakeCase": true,
		"AsCamelCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid informational event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestIsInformationalTypeValidatorDetectingWrongStringRepresentation(t *testing.T) {
	invalid := isInformationalTypetest{
		AsEnum:      informationaltype.PullSummary,
		AsSnakeCase: "pull_summary",
		AsCamelCase: "AaaBbb", // String representation
	}
	fieldsWithError := map[string]bool{
		"AsCamelCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid informational event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestIsInformationalTypeValidatorDetectingWrongCaseRepresentation(t *testing.T) {
	invalid := isInformationalTypetest{
		AsEnum:      informationaltype.FlowsSummary,
		AsSnakeCase: "flowssss_summary",
		AsCamelCase: "FlowsSummary",
	}
	fieldsWithError := map[string]bool{
		"AsSnakeCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid informational event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}
