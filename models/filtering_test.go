package models

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/models/category"
	"github.com/listendev/pkg/models/severity"
	"github.com/listendev/pkg/verdictcode"
	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	v1 := Verdict{
		CreatedAt: func() *time.Time {
			t := time.Now()

			return &t
		}(),
		Ecosystem:  ecosystem.Npm,
		Pkg:        "darcyclarke-manifest-pkg",
		Severity:   severity.High,
		Shasum:     "429eced1773fbc9ceea5cebda8338c0aaa21eeec",
		Version:    "2.1.15",
		File:       "metadata(mismatches).json",
		Message:    "This package has inconsistent name in the tarball's package.json",
		Code:       verdictcode.MDN05,
		Categories: []category.Category{category.Metadata},
	}
	v1b, _ := v1.MarshalJSON()
	var v1raw interface{}
	if err := json.NewDecoder(bytes.NewReader(v1b)).Decode(&v1raw); err != nil {
		t.Fatalf("failed to decode verdict #1: %v", err)
	}

	v2 := Verdict{
		CreatedAt: func() *time.Time {
			t := time.Now().Add(-time.Hour * 48)

			return &t
		}(),
		Ecosystem: ecosystem.Npm,
		Pkg:       "test",
		Version:   "0.0.1",
		Shasum:    "0123456789012345678901234567890123456789",
		File:      "dynamic!install!.json",
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
		Categories: []category.Category{category.AdjacentNetwork, category.CIS},
		Code:       verdictcode.FNI002,
	}

	v3 := Verdict{
		CreatedAt: func() *time.Time {
			t := time.Now().Add(-time.Hour * 72)

			return &t
		}(),
		Ecosystem: ecosystem.Npm,
		Org:       "@expires",
		Pkg:       "in",
		Version:   "5.0.0",
		Shasum:    "32b09b0e3ca7b757802f9a0b9ded8c2035ce7874",
		File:      "dynamic!install!.json",
		Message:   "spawn",
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
		ExpiresAt: func() *time.Time {
			t := time.Now()

			return &t
		}(),
	}

	v4 := Verdict{
		CreatedAt: func() *time.Time {
			t := time.Now().Add(-time.Hour * 72)

			return &t
		}(),
		Ecosystem: ecosystem.Npm,
		Org:       "@expires",
		Pkg:       "in",
		Version:   "4.0.0",
		Shasum:    "42b09b0e3ca7b757802f9a0b9ded8c2035ce7874",
		File:      "dynamic!install!.json",
		Message:   "some message",
		Metadata: map[string]interface{}{
			NPMPackageNameMetadataKey:    "electronx",
			NPMPackageVersionMetadataKey: "21.2.2",
			"commandline":                "sh -c node install.js",
			"parent_name":                "node",
			"executable_path":            "/bin/sh",
			"server_ip":                  "",
			"server_port":                float64(0),
			"file_descriptor":            "",
		},
		Severity:   severity.Low,
		Categories: []category.Category{category.Network, category.Process},
		Code:       verdictcode.FNI003,
		ExpiresAt: func() *time.Time {
			t := time.Now().Add(time.Hour * 72)

			return &t
		}(),
	}
	v4b, _ := v4.MarshalJSON()
	var v4raw interface{}
	if err := json.NewDecoder(bytes.NewReader(v4b)).Decode(&v4raw); err != nil {
		t.Fatalf("failed to decode verdict #4: %v", err)
	}

	v5 := func() Verdict {
		e, _ := NewEmptyVerdict(ecosystem.Npm, "", "test1", "0.0.2-alpha.1", "aaaaa12321321cssasaaaaaa12321321cssasa22", "typosquat.json")

		return *e
	}()
	v5b, _ := v5.MarshalJSON()
	var v5raw interface{}
	if err := json.NewDecoder(bytes.NewReader(v5b)).Decode(&v5raw); err != nil {
		t.Fatalf("failed to decode verdict #5: %v", err)
	}

	verdicts := Verdicts{
		v1,
		v2,
		v3,
		v4,
		v5,
	}

	type testCase struct {
		descr                string
		filter               string
		wantErr              bool
		wantFilterParsingErr bool
		wantVerdicts         Verdicts
		wantRaw              interface{}
	}

	cases := []testCase{
		{
			descr:                "wrong syntax",
			filter:               `$[]`,
			wantFilterParsingErr: true,
			wantErr:              true,
			wantVerdicts:         nil,
		},
		{
			descr:        "projection",
			filter:       `$[*].message`,
			wantErr:      true, // This also errors out because the output is not a verdicts slice
			wantVerdicts: nil,
			wantRaw: []interface{}{
				"This package has inconsistent name in the tarball's package.json",
				"@vue/devtools 6.5.0 1 B",
				"spawn",
				"some message",
			},
		},
		{
			descr:        "first verdict only",
			filter:       `$[0]`,
			wantErr:      true, // This also errors out because the output is not a slice (but a single verdict)
			wantVerdicts: nil,
			wantRaw:      v1raw,
		},
		{
			descr:        "select from the third verdict to the last one",
			filter:       `$[3:]`,
			wantErr:      false,
			wantVerdicts: Verdicts{v4, v5},
			wantRaw:      []interface{}{v4raw, v5raw},
		},
		{
			descr:        "root",
			filter:       `$`,
			wantErr:      false,
			wantVerdicts: Verdicts{v1, v2, v3, v4, v5},
		},
		{
			descr:        "invalid jsonpath expr",
			filter:       `$.created_at`,
			wantErr:      true,
			wantVerdicts: nil,
		},
		{
			descr:                "empty jsonpath expr",
			filter:               ``,
			wantErr:              true,
			wantFilterParsingErr: true,
			wantVerdicts:         nil,
		},
		{
			descr:        "select verdicts with severity gt low",
			filter:       `$[?(@.severity != "low")]`,
			wantErr:      false,
			wantVerdicts: Verdicts{v1, v2, v3},
		},
		{
			descr:        "select only FNI002 verdicts",
			filter:       `$[?(@.code == "FNI002")]`,
			wantErr:      false,
			wantVerdicts: Verdicts{v2},
		},
		{
			descr:   "select only NPM verdicts",
			filter:  `$[?(@.ecosystem == "npm")]`,
			wantErr: false,
			wantVerdicts: Verdicts{
				v1,
				v2,
				v3,
				v4,
				v5,
			},
		},
		{
			descr:        "select only verdicts with an expiration date",
			filter:       `$[?(@.expires_at ? true : false)]`, // ternary operator
			wantErr:      false,
			wantVerdicts: Verdicts{v3, v4},
		},
		{
			descr:        "select only verdicts without an expiration date",
			filter:       `$[?(@.expires_at ?? true)]`, // colaesce operator // NOTE: this works only if "expires_at" is not omitted during JSON serialization
			wantErr:      false,
			wantVerdicts: Verdicts{v1, v2, v5},
		},
		{
			descr: "select only verdicts without an expiration date or that expire before tomorrow",
			filter: fmt.Sprintf(`$[?(@.expires_at ? @.expires_at<%q : true)]`, func() string {
				t := time.Now().Add(time.Hour * 24)

				return t.Format(time.RFC3339)
			}()),
			wantErr:      false,
			wantVerdicts: Verdicts{v1, v2, v3, v5},
		},
		{
			descr:        "select manifest mismatch/confusion (metadata) verdicts with severity gt low",
			filter:       `$[?(@.severity!="low" && @.file=="metadata(mismatches).json")]`, // AND operator
			wantErr:      false,
			wantVerdicts: Verdicts{v1},
		},
		{
			descr: "select verdicts being created after yesterday",
			filter: fmt.Sprintf(`$[?(@.created_at > date(%q))]`, func() string {
				yesterday := time.Now().Add(-time.Hour * 24)

				return yesterday.Format(time.RFC3339)
			}()),
			wantErr:      false,
			wantVerdicts: Verdicts{v1, v5},
		},
		{
			descr:        "select only dynamic instrumentation verdicts",
			filter:       `$[?(@.file =~ "^dynamic")]`, // regex operator
			wantErr:      false,
			wantVerdicts: Verdicts{v2, v3, v4},
		},
		{
			descr:        "exclude metadata verdicts",
			filter:       `$[?(@.file !~ "^metadata")]`, // inverse regex operator
			wantErr:      false,
			wantVerdicts: Verdicts{v2, v3, v4, v5},
		},
		{
			descr:        "select only verdicts with CIS category",
			filter:       `$[?("cis" in @.categories)]`, // in array operator
			wantErr:      false,
			wantVerdicts: Verdicts{v2},
		},
		{
			descr:        "select only transitive dynamic verdicts (by checking the verdict's metadata)",
			filter:       `$[?(@.file =~ "^dynamic" && @.metadata.npm_package_name == "electronx")]`, // operators on nested fields
			wantErr:      false,
			wantVerdicts: Verdicts{v4},
		},
	}

	for _, tc := range cases {
		raw, got, err := verdicts.Filter(context.TODO(), tc.filter)
		if tc.wantErr {
			assert.Error(t, err)
			assert.Nil(t, got)
			if tc.wantFilterParsingErr {
				target := &FilterParsingError{}
				assert.ErrorAs(t, err, &target)
			}
		} else {
			assert.Nil(t, err)
			if !cmp.Equal(got, tc.wantVerdicts) {
				t.Fatalf("values are not the same:\n%s", cmp.Diff(got, tc.wantVerdicts))
			}
		}

		if tc.wantRaw != nil {
			if !cmp.Equal(raw, tc.wantRaw) {
				t.Fatalf("values are not the same:\n%s", cmp.Diff(raw, tc.wantRaw))
			}
		}
	}
}
