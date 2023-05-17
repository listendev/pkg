package models

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/garnet-org/pkg/models/category"
	"github.com/garnet-org/pkg/models/severity"
	"github.com/garnet-org/pkg/verdictcode"
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
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Severity = "ciao"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Code = verdictcode.DDN01
	v.Severity, _ = severity.New("HIGH")
	e = v.Validate()
	assert.Nil(t, e)
}

func TestVerdictMarshalOk(t *testing.T) {
	v := Verdict{
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
		Severity:   severity.Medium,
		Categories: []category.Category{category.Network, category.Process},
		Code:       verdictcode.FNI001,
	}

	want := heredoc.Doc(`{
        "message": "@vue/devtools 6.5.0 1 B",
        "severity": "medium",
        "categories": ["network", "process"],
        "metadata": {
            "npm_package_name": "electron",
            "npm_package_version": "21.4.2",
            "commandline": "sh -c node install.js",
            "parent_name": "node",
            "executable_path": "/bin/sh",
            "server_ip": "",
            "server_port": 0,
            "file_descriptor": ""
        },
		"code": "FNI001"
    }`)

	got, err := json.Marshal(v)
	assert.Nil(t, err)

	assert.JSONEq(t, want, string(got))
}

func TestVerdictWithoutCategoriesMarshal(t *testing.T) {
	v := Verdict{
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
		Severity: severity.Medium,
		Code:     verdictcode.FNI001,
	}

	want := heredoc.Doc(`{
        "message": "@vue/devtools 6.5.0 1 B",
        "severity": "medium",
        "categories": [],
        "metadata": {
            "npm_package_name": "electron",
            "npm_package_version": "21.4.2",
            "commandline": "sh -c node install.js",
            "parent_name": "node",
            "executable_path": "/bin/sh",
            "server_ip": "",
            "server_port": 0,
            "file_descriptor": ""
        },
		"code": "FNI001"
    }`)

	got, err := json.Marshal(v)
	assert.Nil(t, err)

	assert.JSONEq(t, want, string(got))
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
		Severity:   severity.Medium,
		Categories: []category.Category{category.Network, category.Process},
		Code:       verdictcode.FNI003,
	}
	input := heredoc.Doc(`{
        "message": "@vue/devtools 6.5.0 1 B",
        "severity": "medium",
        "categories": ["NETWORK", "process"],
        "metadata": {
            "npm_package_name": "electron",
            "npm_package_version": "21.4.2",
            "commandline": "sh -c node install.js",
            "parent_name": "node",
            "executable_path": "/bin/sh",
            "server_ip": "",
            "server_port": 0,
            "file_descriptor": ""
        },
		"code": "FNI003"
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
		assert.Equal(t, `the input "NONE" is not a severity`, err.Error())
	}
}

func TestBufferEmptyVerdicts(t *testing.T) {
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
			Severity:   "MEDIUM", // I'm testing it also work with upper case severity
			Categories: []category.Category{category.AdjacentNetwork, category.CIS},
			Code:       verdictcode.FNI002,
		},
	}

	want := data
	want[0].Severity = severity.Medium

	reader, err := data.Buffer()
	if !assert.Nil(t, err) {
		t.Fatalf("Buffer(): got error %q:", err)
	}
	assert.NotNil(t, reader)

	got, err := FromBuffer(reader)
	if !assert.Nil(t, err) {
		t.Fatalf("FromBuferr(): got error %q:", err)
	}
	assert.Equal(t, want, got)
}
