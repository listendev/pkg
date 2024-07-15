package validate

import (
	"strings"
	"testing"

	detectiontype "github.com/listendev/pkg/detection/type"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/maps"
)

type test struct {
	AsEnum      detectiontype.Event `validate:"is_detection_event_type=enum"`
	AsSnakeCase string              `validate:"is_detection_event_type=case"`
	AsCamelCase string              `validate:"is_detection_event_type"`
}

type badTypeTest struct {
	Wrong []uint8 `validate:"is_detection_event_type"`
}

func TestBadFieldType(t *testing.T) {
	input := badTypeTest{
		Wrong: []uint8("pam_config_modification"),
	}
	require.PanicsWithValue(t, "bad field type: []uint8", func() { Singleton.Struct(input) })
}

func TestIsDetectionTypeValidator(t *testing.T) {
	allValid := test{
		AsEnum:      detectiontype.CapabilitiesModification,
		AsSnakeCase: "capabilities_modification",
		AsCamelCase: "CapabilitiesModification",
	}
	require.NoError(t, Singleton.Struct(allValid))

	invalid := test{
		AsEnum:      detectiontype.ContainerEscapeAttempt,
		AsSnakeCase: "xyz_escape",
		AsCamelCase: "XYZEscape",
	}
	fieldsWithError := []string{
		"AsSnakeCase",
		"AsCamelCase",
	}
	errors := Singleton.Struct(invalid).(ValidationErrors)
	fields := map[string]string{}
	for _, err := range errors {
		errMsg := err.Translate(Translator)
		assert.True(t, strings.HasSuffix(errMsg, "a valid detection event type"))
		fields[err.Field()] = errMsg
	}
	require.Equal(t, fieldsWithError, maps.Keys(fields))
}
