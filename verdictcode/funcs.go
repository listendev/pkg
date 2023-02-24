package verdictcode

import (
	"encoding/json"
	"fmt"

	"github.com/listendev/pkg/analysisrequest"
	"golang.org/x/exp/maps"
)

// FromUint64 converts the input number to a Code.
//
// It returns only Code instanced that are associated to an analysis request type.
func FromUint64(input uint64, deprecatedToo bool) (Code, error) {
	for _, codemap := range mapping {
		for k, v := range codemap {
			// Skip when the current code is deprecated
			if deprecatedToo && !v {
				continue
			}
			if uint64(k) == input {
				return k, nil
			}
		}
	}

	return UNK, fmt.Errorf("couldn't find a code matching input (%d)", input)
}

// FromString converst the input string to a Code.
//
// It returns only Code instanced that are associated to an analysis request type.
func FromString(input string, deprecatedToo bool) (Code, error) {
	for _, codemap := range mapping {
		for k, v := range codemap {
			// Skip deprecated
			if deprecatedToo && !v {
				continue
			}
			if k.String() == input {
				return k, nil
			}
		}
	}

	return UNK, fmt.Errorf("couldn't find a code matching input %q", input)
}

// GetBy gives the codes for the input analysis request type.
//
// It returns the deprecated codes too.
func GetBy(t analysisrequest.Type) ([]Code, error) {
	submap, ok := mapping[t]
	if !ok {
		return nil, fmt.Errorf("couldn't find codes for the analysis request type %q", t.String())
	}

	return maps.Keys(submap), nil
}

func (c Code) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c *Code) UnmarshalJSON(data []byte) error {
	var cStr string
	if err := json.Unmarshal(data, &cStr); err != nil {
		return err
	}

	res, err := FromString(cStr, true)
	if err != nil {
		return err
	}
	*c = res

	return nil
}

func (c Code) Type(deprecatedToo bool) (analysisrequest.Type, error) {
	for t, codemap := range mapping {
		for k, v := range codemap {
			// Skip deprecated
			if deprecatedToo && !v {
				continue
			}
			if k == c {
				return t, nil
			}
		}
	}

	return analysisrequest.Nop, fmt.Errorf("not found")
}
