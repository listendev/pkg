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
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/models/category"
	"github.com/listendev/pkg/models/severity"
	"github.com/listendev/pkg/verdictcode"
	"github.com/stretchr/testify/assert"
)

func TestNPMVerdictValidation(t *testing.T) {
	v := &Verdict{}
	e := v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
	}

	v.Ecosystem = 0
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "must be one of [npm"))
	}
	v.Ecosystem = ecosystem.Npm

	v.Org = "org"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "organization name must start with @"))
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

	v.Digest = "1"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "SHA1 (40 characters long)"))
	}
	v.Digest = "aaaaaaaaaa1aaaaaaaaaa1aaaaaaaaaa12345678"

	v.CreatedAt = &time.Time{}

	v.File = "INVALID"
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:")) // last mandatory field
		assert.True(t, strings.Contains(e.Error(), "valid results file"))
	}
	v.File = "dynamic!install!.json"

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
		assert.True(t, strings.HasPrefix(e.Error(), "validation errors:"))
		assert.True(t, strings.Contains(e.Error(), "verdict code is not coherent"))
		assert.True(t, strings.Contains(e.Error(), "a fingerprint is mandatory because the verdict code is not uniquely identifying it"))
	}

	v.Code = verdictcode.Code(uint64(math.MaxUint64))
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "the code identifying the verdict type is not a valid verdict code"))
	}

	v.Code = verdictcode.FNI003
	e = v.Validate()
	if assert.Error(t, e) {
		assert.True(t, strings.HasPrefix(e.Error(), "validation error:"))
		assert.True(t, strings.Contains(e.Error(), "a fingerprint is mandatory because the verdict code is not uniquely identifying it"))
	}

	k, ke := v.Key()
	assert.NotNil(t, ke)
	assert.Empty(t, k)
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

func TestKey(t *testing.T) {
	v1, err1 := NewEmptyVerdict(ecosystem.Npm, "@phantom", "synpress", "4.0.0-alpha.19", "879524fd7166d1a8659cd0f5d81800afb268d8c2", "metadata(mismatches).json")
	assert.Nil(t, err1)
	assert.NotNil(t, v1)

	k1, k1Err := v1.Key()
	assert.Nil(t, k1Err)
	assert.NotNil(t, k1)
	assert.Equal(t, "npm/@phantom/synpress/4.0.0-alpha.19/879524fd7166d1a8659cd0f5d81800afb268d8c2/metadata(mismatches).json", k1)

	v2, err2 := NewEmptyVerdict(ecosystem.Pypi, "", "boto3", "1.33.8", "879524fd7166d1a8659cd0f5d81800afb268d8c21231312213aasadsda213321", "typosquat.json")
	assert.Nil(t, err2)
	assert.NotNil(t, v2)

	k2, k2Err := v2.Key()
	assert.Nil(t, k2Err)
	assert.NotNil(t, k2)
	assert.Equal(t, "pypi/boto3/1.33.8/879524fd7166d1a8659cd0f5d81800afb268d8c21231312213aasadsda213321/typosquat.json", k2)
}

func TestMarshalNPMOkVerdict(t *testing.T) {
	now := time.Now()
	v := Verdict{
		CreatedAt:   &now,
		Ecosystem:   ecosystem.Npm,
		Pkg:         "test",
		Version:     "0.0.1",
		Digest:      "0123456789012345678901234567890123456789",
		File:        "dynamic!install!.json",
		Fingerprint: "something",
		Message:     "@vue/devtools 6.5.0 1 B",
		Metadata: map[string]interface{}{
			// Expecting this to not be present into the resulting JSON
			"end": map[string]int{
				"col":    0,
				"line":   0,
				"offset": 0,
			},
			NPMPackageNameMetadataKey:    "electron",
			NPMPackageVersionMetadataKey: "21.4.2",
			"commandline":                "sh -c node install.js",
			"parent_name":                "node",
			"executable_path":            "/bin/sh",
			// Expecting this to not be present into the resulting JSON
			"server_ip": "",
			// Expecting this to not be present into the resulting JSON
			"server_port": float64(0),
			// Expecting this to not be present into the resulting JSON
			"file_descriptor": "",
			// Expecting this to not be present into the resulting JSON
			"emptyarray": []string{},
			// Expecting this to not be present into the resulting JSON
			"emptymap": map[string]string{},
			// Expecting this to not be present into the resulting JSON
			"start": map[string]int{
				"col":    0,
				"line":   0,
				"offset": 0,
			},
		},
		Severity:   severity.Medium,
		Categories: []category.Category{category.Network, category.Process},
		Code:       verdictcode.FNI001,
	}
	v.ExpiresIn(time.Second * 5)

	want := heredoc.Docf(`{
		"digest": "0123456789012345678901234567890123456789",
		"ecosystem": "npm",
		"pkg": "test",
		"version": "0.0.1",
		"file": "dynamic!install!.json",
		"fingerprint": "something",
		"created_at": %q,
		"message": "@vue/devtools 6.5.0 1 B",
		"severity": "medium",
		"categories": ["network", "process"],
		"metadata": {
			"npm_package_name": "electron",
			"npm_package_version": "21.4.2",
			"commandline": "sh -c node install.js",
			"parent_name": "node",
			"executable_path": "/bin/sh"
		},
		"code": "FNI001",
		"expires_at": %q
	}`, now.Format(time.RFC3339Nano), v.ExpiresAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestMarshalNPMEmptyVerdict(t *testing.T) {
	v, e := NewEmptyVerdict(ecosystem.Npm, "", "test", "0.0.1", "0123456789012345678901234567890123456789", "dynamic!install!.json")
	assert.Nil(t, e)
	assert.NotNil(t, v)
	k, ke := v.Key()
	assert.Nil(t, ke)
	assert.NotNil(t, k)
	assert.Equal(t, "npm/test/0.0.1/0123456789012345678901234567890123456789/dynamic!install!.json", k)

	want := heredoc.Docf(`{
		"digest": "0123456789012345678901234567890123456789",
		"ecosystem": "npm",
		"pkg": "test",
		"version": "0.0.1",
		"expires_at": null,
		"file": "dynamic!install!.json",
		"created_at": %q
	}`, v.CreatedAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestMarshalPyPiEmptyVerdict(t *testing.T) {
	v, e := NewEmptyVerdict(ecosystem.Pypi, "", "boto3", "1.33.8", "0123456789012345678901234567890123456789sasdadadadsadasdaq233332", "typosquat.json")
	assert.Nil(t, e)
	assert.NotNil(t, v)
	k, ke := v.Key()
	assert.Nil(t, ke)
	assert.NotNil(t, k)
	assert.Equal(t, "pypi/boto3/1.33.8/0123456789012345678901234567890123456789sasdadadadsadasdaq233332/typosquat.json", k)

	want := heredoc.Docf(`{
		"digest": "0123456789012345678901234567890123456789sasdadadadsadasdaq233332",
		"ecosystem": "pypi",
		"pkg": "boto3",
		"version": "1.33.8",
		"expires_at": null,
		"file": "typosquat.json",
		"created_at": %q
	}`, v.CreatedAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestMarshalNPMEmptyVerdictWithExplicitUnknownCodeAndSeverity(t *testing.T) {
	v, e := NewEmptyVerdict(ecosystem.Npm, "", "test", "0.0.1", "0123456789012345678901234567890123456789", "dynamic!install!.json")
	v.Code = verdictcode.UNK
	v.Severity = severity.Empty
	assert.Nil(t, e)
	assert.NotNil(t, v)
	k, ke := v.Key()
	assert.Nil(t, ke)
	assert.NotNil(t, k)
	assert.Equal(t, "npm/test/0.0.1/0123456789012345678901234567890123456789/dynamic!install!.json", k)

	want := heredoc.Docf(`{
		"digest": "0123456789012345678901234567890123456789",
		"ecosystem": "npm",
		"pkg": "test",
		"version": "0.0.1",
		"expires_at": null,
		"file": "dynamic!install!.json",
		"created_at": %q
	}`, v.CreatedAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestMarshalPyPiEmptyVerdictWithExplicitUnknownCodeAndSeverity(t *testing.T) {
	v, e := NewEmptyVerdict(ecosystem.Pypi, "", "boto3", "1.33.8", "0123456789012345678901234567890123456789sasdadadadsadasdaq233332", "typosquat.json")
	v.Code = verdictcode.UNK
	v.Severity = severity.Empty
	assert.Nil(t, e)
	assert.NotNil(t, v)
	k, ke := v.Key()
	assert.Nil(t, ke)
	assert.NotNil(t, k)
	assert.Equal(t, "pypi/boto3/1.33.8/0123456789012345678901234567890123456789sasdadadadsadasdaq233332/typosquat.json", k)

	want := heredoc.Docf(`{
		"digest": "0123456789012345678901234567890123456789sasdadadadsadasdaq233332",
		"ecosystem": "pypi",
		"pkg": "boto3",
		"version": "1.33.8",
		"expires_at": null,
		"file": "typosquat.json",
		"created_at": %q
	}`, v.CreatedAt.Format(time.RFC3339Nano))

	got, err := json.Marshal(v)
	if assert.Nil(t, err, "Marshalling error") {
		assert.JSONEq(t, want, string(got))
	}
}

func TestUnmarshalNPMEmptyVerdict(t *testing.T) {
	want, e := NewEmptyVerdict(ecosystem.Npm, "", "test1", "0.0.2-alpha.1", "aaaaa12321321cssasaaaaaa12321321cssasa22", "typosquat.json")
	assert.Nil(t, e)
	assert.NotNil(t, want)

	input := heredoc.Docf(`{
		"ecosystem": "npm",
		"digest": "aaaaa12321321cssasaaaaaa12321321cssasa22",
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

func TestUnmarshalPyPiEmptyVerdict(t *testing.T) {
	want, e := NewEmptyVerdict(ecosystem.Pypi, "", "cctx", "1.0.0", "6435a4aaaa12321321cssasaaaaaa12321321cssasa221234567889012132131", "typosquat.json")
	assert.Nil(t, e)
	assert.NotNil(t, want)

	input := heredoc.Docf(`{
		"ecosystem": "pypi",
		"digest": "6435a4aaaa12321321cssasaaaaaa12321321cssasa221234567889012132131",
		"pkg": "cctx",
		"version": "1.0.0",
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

func TestPyPiVerdictValidations(t *testing.T) {
	_, err1 := NewEmptyVerdict(ecosystem.Pypi, "ORG", "cctx", "1.0.0", "123", "typosquat.json")
	assert.NotNil(t, err1)
	if assert.Error(t, err1) {
		assert.True(t, strings.HasPrefix(err1.Error(), "validation errors:"))
		assert.True(t, strings.Contains(err1.Error(), "organization name must be empty"))
		assert.True(t, strings.Contains(err1.Error(), "digest must be a valid blake2b digest (64 characters long)"))
	}
}

func TestUnmarshalNPMOkVerdict(t *testing.T) {
	now := time.Now()
	want := Verdict{
		Ecosystem:   ecosystem.Npm,
		CreatedAt:   &now,
		Pkg:         "test",
		Version:     "0.0.1",
		Digest:      "0123456789012345678901234567890123456789",
		File:        "dynamic!install!.json",
		Message:     "@vue/devtools 6.5.0 1 B",
		Fingerprint: "something",
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
		"digest": "0123456789012345678901234567890123456789",
		"pkg": "test",
		"version": "0.0.1",
		"file": "dynamic!install!.json",
		"created_at": %q,
		"message": "@vue/devtools 6.5.0 1 B",
		"severity": "medium",
		"categories": ["NETWORK", "process"],
		"fingerprint": "something",
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
	vvvv1 := Verdict{
		Ecosystem:   ecosystem.Npm,
		CreatedAt:   &now,
		Pkg:         "test",
		Version:     "0.0.1",
		Digest:      "0123456789012345678901234567890123456789",
		File:        "dynamic!install!.json",
		Message:     "@vue/devtools 6.5.0 1 B",
		Fingerprint: "something",
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
	vvvv1.ExpiresIn(time.Second * 5)
	vvvv2 := Verdict{
		CreatedAt:  &now,
		Ecosystem:  ecosystem.Npm,
		Pkg:        "darcyclarke-manifest-pkg",
		Severity:   severity.Medium,
		Digest:     "429eced1773fbc9ceea5cebda8338c0aaa21eeec",
		Version:    "2.1.15",
		File:       "metadata(mismatches).json",
		Message:    "This package has inconsistent name in the tarball's package.json",
		Code:       verdictcode.MDN05,
		Categories: []category.Category{category.Metadata},
	}
	empty, _ := NewEmptyVerdict(ecosystem.Npm, "", "testone", "0.0.2", "a123456789012345678901234567890123456789", "typosquat.json")
	data := Verdicts{
		*empty,
		vvvv1,
		vvvv2,
	}

	reader, err := data.Buffer()
	if assert.Nil(t, err) {
		emptyJSON, _ := json.Marshal(empty)
		vvvv1JSON, _ := json.Marshal(vvvv1)
		vvvv2JSON, _ := json.Marshal(vvvv2)
		got, _ := io.ReadAll(reader)
		assert.JSONEq(t, fmt.Sprintf("[%s,%s,%s]", emptyJSON, vvvv1JSON, vvvv2JSON), string(got))
	}
}

func TestFromBuffer(t *testing.T) {
	now := time.Now()
	noResults, _ := NewEmptyVerdict(ecosystem.Npm, "", "test1", "0.0.2-alpha.1", "aaaaa12321321cssasaaaaaa12321321cssasa22", "typosquat.json")

	want := Verdicts{
		Verdict{
			CreatedAt:  &now,
			Ecosystem:  ecosystem.Npm,
			Pkg:        "darcyclarke-manifest-pkg",
			Severity:   severity.Medium,
			Digest:     "429eced1773fbc9ceea5cebda8338c0aaa21eeec",
			Version:    "2.1.15",
			File:       "metadata(mismatches).json",
			Message:    "This package has inconsistent name in the tarball's package.json",
			Code:       verdictcode.MDN05,
			Categories: []category.Category{category.Metadata},
			Metadata:   map[string]interface{}{},
		},
		Verdict{
			CreatedAt:   &now,
			Ecosystem:   ecosystem.Npm,
			Pkg:         "test",
			Version:     "0.0.1",
			Digest:      "0123456789012345678901234567890123456789",
			File:        "dynamic!install!.json",
			Fingerprint: "something",
			Message:     "@vue/devtools 6.5.0 1 B",
			Metadata: map[string]interface{}{
				NPMPackageNameMetadataKey:    "electron",
				NPMPackageVersionMetadataKey: "21.4.2",
				"commandline":                "sh -c node install.js",
				"parent_name":                "node",
				"executable_path":            "/bin/sh",
			},
			Severity:   severity.Medium,
			Categories: []category.Category{category.AdjacentNetwork, category.CIS},
			Code:       verdictcode.FNI002,
		},
		*noResults,
	}

	reader, err := want.Buffer()
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
