package verdictcode

import (
	"github.com/listendev/pkg/analysisrequest"
)

// mapping maps the codes to the analysis request type that can generate them.
var mapping = map[analysisrequest.Type]map[Code]bool{
	analysisrequest.NPMInstallWhileDynamicInstrumentation: {
		FNI001: true,
		FNI002: true,
		FNI003: true,
	},
	analysisrequest.NPMDepsDev: {
		DDN01: true,
	},
	analysisrequest.NPMTyposquat: {
		TSN01: true,
	},
	analysisrequest.NPMMetadataEmptyDescription: {
		MDN01: true,
	},
	analysisrequest.NPMMetadataVersion: {
		MDN02: true,
		MDN03: true,
	},
	analysisrequest.NPMMetadataMaintainersEmailCheck: {
		MDN04: true,
	},
	analysisrequest.NPMMetadataMismatches: {
		MDN05: true,
		MDN06: true,
		MDN07: true,
		MDN08: true,
	},
	analysisrequest.NPMStaticAnalysisEnvExfiltration: {
		STN001: true,
	},
	analysisrequest.NPMStaticAnalysisDetachedProcessExecution: {
		STN002: true,
	},
	analysisrequest.NPMStaticAnalysisShadyLinks: {
		STN003: true,
	},
	analysisrequest.NPMStaticAnalysisEvalBase64: {
		STN004: true,
	},
	analysisrequest.NPMStaticAnalysisInstallScript: {
		STN005: true,
	},
	analysisrequest.NPMStaticNonRegistryDependency: {
		STN006: true,
		STN007: true,
		STN008: true,
		STN009: true,
	},
}
