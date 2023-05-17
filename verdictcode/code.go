package verdictcode

import (
	"github.com/garnet-org/pkg/analysisrequest"
)

type Code uint64

// Constant codes.
//
// The naming convention is composed by a prefix using abbreviations of the following type components:
// <collector><ecosystem>[<ecosystem_action>,...]
// The prefix is followed by a progressive number.
// The number of digits (always with zero left-padding) of such a number depends on the range we wanna reserve.
// For example for the type "urn:scheduler:falco!npm,install.json" the prefix is FNI,
// followed by 3 digits because we wanna reserve 999 possible progressive numbers.
//
// Pay attention to:
// 1) do not change (their are meant to be constant) the underlying value of any Code when adding new ones
// 2) add new Code constants into their range
// 3) reserve space for new future Code constants
const (
	UNK Code = iota

	FNI001 // 1
	FNI002
	FNI003
	// 999 IDs for FNI* codes

	DDN01 Code = iota + 997 // 1001
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
}

//go:generate stringer -type=Code
