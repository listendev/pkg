package analysisrequest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestGetResultFilesByEcosystem(t *testing.T) {
	wnt := map[Type]string{
		NPMInstallWhileFalco: "falco[install].json",
		// NPMTestWhileFalco:    "falco[test].json",
		NPMDepsDev:   "depsdev.json",
		NPMTyposquat: "typosquat.json",
	}
	got := GetResultFilesByEcosystem(NPMEcosystem)

	less := func(a, b string) bool { return a < b }
	if equal := cmp.Equal(wnt, got, cmpopts.SortMaps(less)); !equal {
		diff := cmp.Diff(wnt, got, cmpopts.SortMaps(less))
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}

func TestGetTypeFromResultFile(t *testing.T) {
	wnt := map[string]Type{
		"falco[install].json": NPMInstallWhileFalco,
		// "falco[test].json":    NPMTestWhileFalco,
		"depsdev.json":   NPMDepsDev,
		"typosquat.json": NPMTyposquat,
	}
	for f, typ := range wnt {
		got, err := GetTypeFromResultFile(NPMEcosystem, f)

		if assert.Nil(t, err) {
			assert.Equal(t, typ, got)
		}
	}
}
