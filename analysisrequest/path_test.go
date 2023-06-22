package analysisrequest

import (
	"testing"

	"github.com/garnet-org/pkg/ecosystem"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestGetResultFilesByEcosystem(t *testing.T) {
	wnt := map[Type]string{
		NPMInstallWhileFalco: "falco!install!.json",
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
	got := GetResultFilesByEcosystem(ecosystem.Npm)

	less := func(a, b string) bool { return a < b }
	if equal := cmp.Equal(wnt, got, cmpopts.SortMaps(less)); !equal {
		diff := cmp.Diff(wnt, got, cmpopts.SortMaps(less))
		t.Errorf("diff: (-got +want)\n%s", diff)
	}
}

func TestGetTypeForEcosystemFromResultFile(t *testing.T) {
	wnt := map[string]Type{
		"falco!install!.json": NPMInstallWhileFalco,
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
		got, err := GetTypeForEcosystemFromResultFile(ecosystem.Npm, f)

		if assert.Nil(t, err) {
			assert.Equal(t, typ, got)
		}
	}

	_, err := GetTypeForEcosystemFromResultFile(ecosystem.Npm, "unknown.json")
	if assert.Error(t, err) {
		assert.Equal(t, `couldn't find any type for ecosystem "npm" matching the results file "unknown.json"`, err.Error())
	}
}

func TestGetTypeFromResultFile(t *testing.T) {
	wnt := map[string]Type{
		"falco!install!.json": NPMInstallWhileFalco,
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
		got, err := GetTypeFromResultFile(f)

		if assert.Nil(t, err) {
			assert.Equal(t, typ, got)
		}
	}

	_, err := GetTypeFromResultFile("unknown.json")
	if assert.Error(t, err) {
		assert.Equal(t, `couldn't find any type in any ecosystem matching the results file "unknown.json"`, err.Error())
	}
}
