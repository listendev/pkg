package analysisrequest

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetResultFilesByEcosystem(t *testing.T) {
	wnt := []string{
		"falco(install).json",
		"falco(test).json",
		"depsdev.json",
	}
	got := GetResultFilesByEcosystem(NPMEcosystem)

	less := func(a, b string) bool { return a < b }
	if equal := cmp.Equal(wnt, got, cmpopts.SortSlices(less)); !equal {
		diff := cmp.Diff(wnt, got, cmpopts.SortSlices(less))
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}
