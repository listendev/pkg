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
		NPMDepsDev:                                "depsdev.json",
		NPMTyposquat:                              "typosquat.json",
		NPMMetadataEmptyDescription:               "metadata(empty_descr).json",
		NPMMetadataVersion:                        "metadata(version).json",
		NPMMetadataMaintainersEmailCheck:          "metadata(email_check).json",
		NPMStaticAnalysisEnvExfiltration:          "static(exfiltrate_env).json",
		NPMStaticAnalysisEvalBase64:               "static(base64_eval).json",
		NPMStaticAnalysisDetachedProcessExecution: "static(detached_process_exec).json",
		NPMStaticAnalysisShadyLinks:               "static(shady_links).json",
		NPMStaticAnalysisInstallScript:            "static(install_script).json",
		NPMStaticNonRegistryDependency:            "static(non_registry_dependency).json",
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
		"depsdev.json":                         NPMDepsDev,
		"typosquat.json":                       NPMTyposquat,
		"metadata(empty_descr).json":           NPMMetadataEmptyDescription,
		"metadata(version).json":               NPMMetadataVersion,
		"metadata(email_check).json":           NPMMetadataMaintainersEmailCheck,
		"static(exfiltrate_env).json":          NPMStaticAnalysisEnvExfiltration,
		"static(shady_links).json":             NPMStaticAnalysisShadyLinks,
		"static(detached_process_exec).json":   NPMStaticAnalysisDetachedProcessExecution,
		"static(base64_eval).json":             NPMStaticAnalysisEvalBase64,
		"static(install_script).json":          NPMStaticAnalysisInstallScript,
		"static(non_registry_dependency).json": NPMStaticNonRegistryDependency,
	}
	for f, typ := range wnt {
		got, err := GetTypeFromResultFile(NPMEcosystem, f)

		if assert.Nil(t, err) {
			assert.Equal(t, typ, got)
		}
	}
}
