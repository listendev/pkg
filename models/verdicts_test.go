package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestVerdictValidation(t *testing.T) {
	v := &Verdict{}
	e := v.Validate()

	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Message = "some message"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
	}

	v.Severity = "ciao"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.Equal(t, "validation error: the verdict severity must be low, medium, or high", e.Error())
	}

	v.Severity = VerdictSeverity(strings.ToUpper(string(High)))
	e = v.Validate()
	assert.Nil(t, e)
}

func TestVerdictUnmarshalOk(t *testing.T) {
	want := Verdict{
		Message: "@vue/devtools 6.5.0 1 B",
		Metadata: map[string]interface{}{
			NPMPackageNameMetadataKey:    "electron",
			NPMPackageVersionMetadataKey: "21.4.2",
			"commandline":                "sh -c node install.js",
			"parent_name":                "node",
			"executable_path":            "/bin/sh",
			"server_ip":                  "",
			"server_port":                float64(0),
			"file_descriptor":            "",
		},
		Severity: Medium,
	}
	input := heredoc.Doc(`{
        "message": "@vue/devtools 6.5.0 1 B",
        "severity": "medium",
        "metadata": {
            "npm_package_name": "electron",
            "npm_package_version": "21.4.2",
            "commandline": "sh -c node install.js",
            "parent_name": "node",
            "executable_path": "/bin/sh",
            "server_ip": "",
            "server_port": 0,
            "file_descriptor": ""
        }
    }`)

	var v Verdict
	err := json.Unmarshal([]byte(input), &v)
	assert.Nil(t, err)
	assert.Equal(t, want, v)
}

func TestVerdictUnmarshalErrorSeverity(t *testing.T) {
	input := heredoc.Doc(`{
        "message": "some message",
        "severity": "NONE"
    }`)

	var v Verdict
	err := json.Unmarshal([]byte(input), &v)
	if assert.Error(t, err) {
		assert.Equal(t, "validation error: the verdict severity must be low, medium, or high", err.Error())
	}
}

func TestBufferInvalidVerdict(t *testing.T) {
	data := Verdicts{
		{},
	}

	reader, err := data.Buffer()
	assert.Nil(t, reader)
	assert.NotNil(t, err)
}

func TestFromBuffer(t *testing.T) {
	data := Verdicts{
		Verdict{
			Message: "@vue/devtools 6.5.0 1 B",
			Metadata: map[string]interface{}{
				NPMPackageNameMetadataKey:    "electron",
				NPMPackageVersionMetadataKey: "21.4.2",
				"commandline":                "sh -c node install.js",
				"parent_name":                "node",
				"executable_path":            "/bin/sh",
				"server_ip":                  "",
				"server_port":                float64(0),
				"file_descriptor":            "",
			},
			Severity: Medium,
		},
	}

	reader, err := data.Buffer()
	assert.NotNil(t, reader)
	assert.Nil(t, err)

	got, err := FromBuffer(reader)
	assert.Nil(t, err)
	assert.Equal(t, data, got)
}
