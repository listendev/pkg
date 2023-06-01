package verdictcode

import (
	"github.com/garnet-org/pkg/analysisrequest"
)

// mapping maps the codes to the analysis request type that can generate them.
var mapping = map[analysisrequest.Type]map[Code]bool{
	analysisrequest.NPMInstallWhileFalco: map[Code]bool{
		FNI001: true,
		FNI002: true,
		FNI003: true,
	},
	analysisrequest.NPMDepsDev: map[Code]bool{
		DDN01: true,
	},
	analysisrequest.NPMTyposquat: map[Code]bool{
		TSN01: true,
	},
	analysisrequest.NPMMetadataEmptyDescription: map[Code]bool{
		MDN01: true,
	},
	analysisrequest.NPMMetadataVersion: map[Code]bool{
		MDN02: true,
		MDN03: true,
	},
	analysisrequest.NPMMetadataMaintainersEmailCheck: map[Code]bool{
		MDN04: true,
	},
}
