package validate

import (
	"strings"
	"testing"

	detectiontype "github.com/listendev/pkg/detection/type"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type isDetectionTypetest struct {
	AsEnum      detectiontype.Event `validate:"is_detection_event_type=enum"`
	AsSnakeCase string              `validate:"is_detection_event_type=case"`
	AsCamelCase string              `validate:"is_detection_event_type"`
}

type badIsDetectionTypeTest struct {
	Wrong []uint8 `validate:"is_detection_event_type"`
}

func TestIsDetectionTypeValidatorBadFieldType(t *testing.T) {
	input := badIsDetectionTypeTest{
		Wrong: []uint8("pam_config_modification"),
	}
	require.PanicsWithValue(t, "bad field type: []uint8", func() {
		//nolint:errcheck // we are checking it panics
		Singleton.Struct(input)
	})
}

func TestIsDetectionTypeValidator(t *testing.T) {
	allValid := isDetectionTypetest{
		AsEnum:      detectiontype.CapabilitiesModification,
		AsSnakeCase: "capabilities_modification",
		AsCamelCase: "CapabilitiesModification",
	}
	require.NoError(t, Singleton.Struct(allValid))

	invalid := isDetectionTypetest{
		AsEnum:      detectiontype.ContainerEscapeAttempt,
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
		assert.True(t, strings.HasSuffix(errMsg, "a valid detection event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestNoneIsNotAValidDetectionType(t *testing.T) {
	invalid := isDetectionTypetest{
		AsEnum:      detectiontype.None,
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
		assert.True(t, strings.HasSuffix(errMsg, "a valid detection event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestIsDetectionTypeValidatorDetectingWrongStringRepresentation(t *testing.T) {
	invalid := isDetectionTypetest{
		AsEnum:      detectiontype.ProcessFingerprint,
		AsSnakeCase: "process_fingerprint",
		AsCamelCase: "AaaBbb", // String representation
	}
	fieldsWithError := map[string]bool{
		"AsCamelCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid detection event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}

func TestIsDetectionTypeValidatorDetectingWrongCaseRepresentation(t *testing.T) {
	invalid := isDetectionTypetest{
		AsEnum:      detectiontype.OsNetworkFingerprint,
		AsSnakeCase: "osnetwork_fingerprint",
		AsCamelCase: "OsNetworkFingerprint",
	}
	fieldsWithError := map[string]bool{
		"AsSnakeCase": true,
	}
	errors := Singleton.Struct(invalid).(ValidationError)
	fields := map[string]bool{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid detection event type"))
		fields[err.Field()] = true
	}
	require.Equal(t, fieldsWithError, fields)
}
