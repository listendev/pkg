package analysisrequest

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/listendev/pkg/ecosystem"
	"github.com/stretchr/testify/assert"
)

func TestGetResultFilesByEcosystem(t *testing.T) {
	wnt := map[Type]string{}
	for _, e := range ecosystem.All() {
		switch e {
		case ecosystem.Npm:
			wnt = map[Type]string{
				NPMInstallWhileDynamicInstrumentation: "dynamic!install!.json",
				// NPMTestWhileDynamicInstrumentation:    "dynamic[test].json",
				NPMAdvisory:                               "advisory.json",
				NPMTyposquat:                              "typosquat.json",
				NPMMetadataEmptyDescription:               "metadata(empty_descr).json",
				NPMMetadataVersion:                        "metadata(version).json",
				NPMMetadataMaintainersEmailCheck:          "metadata(email_check).json",
				NPMMetadataMismatches:                     "metadata(mismatches).json",
				NPMStaticAnalysisEnvExfiltration:          "static(exfiltrate_env).json",
				NPMStaticAnalysisEvalBase64:               "static(base64_eval).json",
				NPMStaticAnalysisDetachedProcessExecution: "static(detached_process_exec).json",
				NPMStaticAnalysisShadyLinks:               "static(shady_links).json",
				NPMStaticAnalysisInstallScript:            "static(install_script).json",
				NPMStaticNonRegistryDependency:            "static(non_registry_dependency).json",
			}
		case ecosystem.Pypi:
			wnt = map[Type]string{
				PypiTyposquat:                              "typosquat.json",
				PypiMetadataMaintainersEmailCheck:          "metadata(email_check).json",
				PypiStaticAnalysisEnvExfiltration:          "static(exfiltrate_env).json",
				PypiStaticAnalysisEvalBase64:               "static(base64_eval).json",
				PypiStaticAnalysisDetachedProcessExecution: "static(detached_process_exec).json",
				PypiStaticAnalysisShadyLinks:               "static(shady_links).json",
				PypiStaticNonRegistryDependency:            "static(non_registry_dependency).json",
				PypiStaticAnalysisCodeExecutionAtSetup:     "static(code_exec_at_setup).json",
			}
		}
		got := GetResultFilesByEcosystem(e)

		less := func(a, b string) bool { return a < b }
		if equal := cmp.Equal(wnt, got, cmpopts.SortMaps(less)); !equal {
			diff := cmp.Diff(wnt, got, cmpopts.SortMaps(less))
			t.Errorf("diff: (-got +want)\n%s", diff)
		}
	}
}

func TestGetTypeForEcosystemFromResultFile(t *testing.T) {
	wnt := map[string]Type{}
	for _, e := range ecosystem.All() {
		switch e {
		case ecosystem.Npm:
			wnt = map[string]Type{
				"dynamic!install!.json": NPMInstallWhileDynamicInstrumentation,
				// "dynamic[test].json":    NPMTestWhileDynamicInstrumentation,
				"advisory.json":                        NPMAdvisory,
				"typosquat.json":                       NPMTyposquat,
				"metadata(empty_descr).json":           NPMMetadataEmptyDescription,
				"metadata(version).json":               NPMMetadataVersion,
				"metadata(email_check).json":           NPMMetadataMaintainersEmailCheck,
				"metadata(mismatches).json":            NPMMetadataMismatches,
				"static(exfiltrate_env).json":          NPMStaticAnalysisEnvExfiltration,
				"static(shady_links).json":             NPMStaticAnalysisShadyLinks,
				"static(detached_process_exec).json":   NPMStaticAnalysisDetachedProcessExecution,
				"static(base64_eval).json":             NPMStaticAnalysisEvalBase64,
				"static(install_script).json":          NPMStaticAnalysisInstallScript,
				"static(non_registry_dependency).json": NPMStaticNonRegistryDependency,
			}
		case ecosystem.Pypi:
			wnt = map[string]Type{
				"typosquat.json":                       PypiTyposquat,
				"metadata(email_check).json":           PypiMetadataMaintainersEmailCheck,
				"static(exfiltrate_env).json":          PypiStaticAnalysisEnvExfiltration,
				"static(base64_eval).json":             PypiStaticAnalysisEvalBase64,
				"static(detached_process_exec).json":   PypiStaticAnalysisDetachedProcessExecution,
				"static(shady_links).json":             PypiStaticAnalysisShadyLinks,
				"static(non_registry_dependency).json": PypiStaticNonRegistryDependency,
				"static(code_exec_at_setup).json":      PypiStaticAnalysisCodeExecutionAtSetup,
			}
		}

		for f, typ := range wnt {
			got, err := GetTypeForEcosystemFromResultFile(e, f)

			if assert.Nil(t, err) {
				assert.Equal(t, typ, got)
			}
		}

		_, err := GetTypeForEcosystemFromResultFile(e, "unknown.json")
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Sprintf(`couldn't find any type for ecosystem %q matching the results file "unknown.json"`, e.Case()), err.Error())
		}
	}
}

func TestGetTypesFromResultFile(t *testing.T) {
	wnt := map[string][]Type{
		"dynamic!install!.json": {NPMInstallWhileDynamicInstrumentation},
		// "dynamic[test].json":    {NPMTestWhileDynamicInstrumentation},
		"advisory.json":                        {NPMAdvisory},
		"typosquat.json":                       {NPMTyposquat, PypiTyposquat},
		"metadata(empty_descr).json":           {NPMMetadataEmptyDescription},
		"metadata(version).json":               {NPMMetadataVersion},
		"metadata(email_check).json":           {NPMMetadataMaintainersEmailCheck, PypiMetadataMaintainersEmailCheck},
		"metadata(mismatches).json":            {NPMMetadataMismatches},
		"static(exfiltrate_env).json":          {NPMStaticAnalysisEnvExfiltration, PypiStaticAnalysisEnvExfiltration},
		"static(shady_links).json":             {NPMStaticAnalysisShadyLinks, PypiStaticAnalysisShadyLinks},
		"static(detached_process_exec).json":   {NPMStaticAnalysisDetachedProcessExecution, PypiStaticAnalysisDetachedProcessExecution},
		"static(base64_eval).json":             {NPMStaticAnalysisEvalBase64, PypiStaticAnalysisEvalBase64},
		"static(install_script).json":          {NPMStaticAnalysisInstallScript},
		"static(non_registry_dependency).json": {NPMStaticNonRegistryDependency, PypiStaticNonRegistryDependency},
	}
	for f, typ := range wnt {
		got, err := GetTypesFromResultFile(f)

		if assert.Nil(t, err) {
			assert.Equal(t, typ, got)
		}
	}

	_, err := GetTypesFromResultFile("unknown.json")
	if assert.Error(t, err) {
		assert.Equal(t, `couldn't find any type in any ecosystem matching the results file "unknown.json"`, err.Error())
	}
}
