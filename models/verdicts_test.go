package models

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/garnet-org/pkg/ecosystem"
	"github.com/garnet-org/pkg/models/category"
	"github.com/garnet-org/pkg/models/severity"
	"github.com/garnet-org/pkg/verdictcode"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestVerdictValidation(t *testing.T) {
	v := &Verdict{}
	e := v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Org = "org"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "valid NPM organization"))
	}
	v.Org = ""

	v.Pkg = "test"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Version = "INVALID"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "valid semantic version"))
	}
	v.Version = "0.0.1-beta.1+b1234567"

	v.Shasum = "1"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "40 characters long"))
	}
	v.Shasum = "aaaaaaaaaa1aaaaaaaaaa1aaaaaaaaaa12345678"

	v.Ecosystem = 0
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "must be on of [npm"))
	}
	v.Ecosystem = ecosystem.Npm

	v.CreatedAt = &time.Time{}

	v.File = "INVALID"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:")) // last mandatory field
		assert.True(t, strings.Contains(e.Error(), "valid results file"))
	}
	v.File = "falco!install!.json"

	// From now on, all fields are optional until the verdict gets a message...
	e = v.Validate()
	assert.Nil(t, e)

	v.Message = "1"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "severity is required when the message field has a value"))
		assert.True(t, strings.Contains(e.Error(), "code identifying the verdict type is required when the message field has a value"))
		assert.True(t, strings.Contains(e.Error(), "one or more verdict category is required when the message field has a value"))
		assert.True(t, strings.Contains(e.Error(), "the verdict message must be greater than 1 character in length"))
	}

	// So, from this point on, we expect again errors if we don't set the severity, the code, and the category/categories
	v.Message = "spawn something"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "severity is required when the message field has a value"))
		assert.True(t, strings.Contains(e.Error(), "code identifying the verdict type is required when the message field has a value"))
		assert.True(t, strings.Contains(e.Error(), "one or more verdict category is required when the message field has a value"))
	}

	v.Severity = "INVALID"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "severity must be low, medium, or high"))
	}
	v.Severity, _ = severity.New("HIGH")

	v.Categories = []category.Category{(category.Category)(0)}
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.Contains(e.Error(), "not a valid verdict category"))
	}
	v.Categories = []category.Category{category.Network, category.Process}

	// No code means invalid code at this point
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "code identifying the verdict type is required when the message field has a value"))
	}

	// Unknown code is invalid
	v.Code = verdictcode.UNK
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "code identifying the verdict type is required when the message field has a value"))
	}

	v.Code = verdictcode.DDN01
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "verdict code is not coherent"))
	}

	v.Code = verdictcode.Code(uint64(math.MaxUint64))
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "the code identifying the verdict type is not a valid verdict code"))
	}

	v.Code = verdictcode.FNI003
	e = v.Validate()
	assert.Nil(t, e)
}

func TestExpiration(t *testing.T) {
	d := time.Microsecond * 20
	v := &Verdict{}
	assert.False(t, v.HasExpired())

	v.ExpiresIn(d)

	assert.False(t, v.HasExpired())

	time.Sleep(d)

	assert.True(t, v.HasExpired())
}

func TestMarshalOkVerdict(t *testing.T) {
	now := time.Now()
	v := Verdict{
		CreatedAt: &now,
		Ecosystem: ecosystem.Npm,
		Pkg:       "test",
		Version:   "0.0.1",
		Shasum:    "0123456789012345678901234567890123456789",
		File:      "falco!install!.json",
		Message:   "@vue/devtools 6.5.0 1 B",
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
	v.ExpiresIn(time.Second * 5)

	want := heredoc.Docf(`{
		"shasum": "0123456789012345678901234567890123456789",
		"ecosystem": "npm",
		"pkg": "test",
		"version": "0.0.1",
		"file": "falco!install!.json",
		"created_at": %q,
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
		"code": "FNI001",
		"expires_at": %q
	}`, now.Format(time.RFC3339Nano), v.ExpiresAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestMarshalEmptyVerdict(t *testing.T) {
	v, e := NewEmptyVerdict(ecosystem.Npm, "", "test", "0.0.1", "0123456789012345678901234567890123456789", "falco!install!.json")
	assert.Nil(t, e)
	assert.NotNil(t, v)

	want := heredoc.Docf(`{
		"shasum": "0123456789012345678901234567890123456789",
		"ecosystem": "npm",
		"pkg": "test",
		"version": "0.0.1",
		"file": "falco!install!.json",
		"created_at": %q
	}`, v.CreatedAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestUnmarshalEmptyVerdict(t *testing.T) {
	want, e := NewEmptyVerdict(ecosystem.Npm, "", "test1", "0.0.2-alpha.1", "aaaaa12321321cssasaaaaaa12321321cssasa22", "typosquat.json")
	assert.Nil(t, e)
	assert.NotNil(t, want)

	input := heredoc.Docf(`{
		"ecosystem": "npm",
		"shasum": "aaaaa12321321cssasaaaaaa12321321cssasa22",
		"pkg": "test1",
		"version": "0.0.2-alpha.1",
		"file": "typosquat.json",
		"created_at": %q
	}`, want.CreatedAt.Format(time.RFC3339Nano))

	var v Verdict
	err := json.Unmarshal([]byte(input), &v)
	if assert.Nil(t, err) {
		opt := cmpopts.IgnoreFields(Verdict{}, "CreatedAt")
		if !cmp.Equal(v, *want, opt) {
			t.Fatalf("values are not the same:\n%s", cmp.Diff(v, *want))
		}
	}
}

func TestUnmarshalOkVerdict(t *testing.T) {
	now := time.Now()
	want := Verdict{
		Ecosystem: ecosystem.Npm,
		CreatedAt: &now,
		Pkg:       "test",
		Version:   "0.0.1",
		Shasum:    "0123456789012345678901234567890123456789",
		File:      "falco!install!.json",
		Message:   "@vue/devtools 6.5.0 1 B",
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
	input := heredoc.Docf(`{
		"ecosystem": "npm",
		"shasum": "0123456789012345678901234567890123456789",
		"pkg": "test",
		"version": "0.0.1",
		"file": "falco!install!.json",
		"created_at": %q,
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
	}`, want.CreatedAt.Format(time.RFC3339Nano))

	var v Verdict
	err := json.Unmarshal([]byte(input), &v)
	if assert.Nil(t, err) {
		if !cmp.Equal(v, want) {
			t.Fatalf("values are not the same:\n%s", cmp.Diff(v, want))
		}
	}
}

func TestUnmarshalErrorSeverityVerdict(t *testing.T) {
	// Notice we don't need a valid verdict here because severity has its own custom unmarshaller that runs before the verdict's one
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

func TestBuffer(t *testing.T) {
	now := time.Now()
	vvvvv := Verdict{
		Ecosystem: ecosystem.Npm,
		CreatedAt: &now,
		Pkg:       "test",
		Version:   "0.0.1",
		Shasum:    "0123456789012345678901234567890123456789",
		File:      "falco!install!.json",
		Message:   "@vue/devtools 6.5.0 1 B",
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
	vvvvv.ExpiresIn(time.Second * 5)
	empty, _ := NewEmptyVerdict(ecosystem.Npm, "", "testone", "0.0.2", "a123456789012345678901234567890123456789", "typosquat.json")
	data := Verdicts{
		*empty,
		vvvvv,
	}

	reader, err := data.Buffer()
	if assert.Nil(t, err) {
		emptyJSON, _ := json.Marshal(empty)
		vvvvvJSON, _ := json.Marshal(vvvvv)
		got, _ := io.ReadAll(reader)
		assert.JSONEq(t, fmt.Sprintf("[%s,%s]", emptyJSON, vvvvvJSON), string(got))
	}
}

func TestFromBuffer(t *testing.T) {
	now := time.Now()
	data := Verdicts{
		Verdict{
			CreatedAt: &now,
			Ecosystem: ecosystem.Npm,
			Pkg:       "test",
			Version:   "0.0.1",
			Shasum:    "0123456789012345678901234567890123456789",
			File:      "falco!install!.json",
			Message:   "@vue/devtools 6.5.0 1 B",
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
	if !cmp.Equal(got, want) {
		t.Fatalf("values are not the same:\n%s", cmp.Diff(got, want))
	}
}
