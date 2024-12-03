package verdictcode

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/XANi/goneric"
	"github.com/listendev/pkg/analysisrequest"
	"golang.org/x/exp/maps"
)

// Scan implements the sql.Scanner interface.
func (c *Code) Scan(source any) error {
	if x, ok := source.(string); ok {
		cat, err := FromString(x, true) // consider also deprecated values
		if cat != UNK && err != nil {
			return err
		}
		*c = cat

		return nil
	}
	if x, ok := source.(uint64); ok {
		cat, err := FromUint64(x, true) // consider also deprecated values
		if cat != UNK && err != nil {
			return err
		}
		*c = cat
	}

	return fmt.Errorf("cannot scan %T into Code", source)
}

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

// FromString converts the input string to a Code.
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

	return analysisrequest.Nop, errors.New("not found")
}

// UniquelyIdentifies tells whether the receiving Code uniquely identifies a verdict.
//
// A code uniquely identifies a verdict when a collector can only generate one instance of a verdict with such a code
// for the tuple (ecosystem, package, version, collector itself).
//
// This means that when a collector, for the tuple (ecosystem, package, version, collector itself), can generate more instance of verdicts with the same code,
// than such a code is not uniquely identifying the verdict.
func (c Code) UniquelyIdentifies() bool {
	return !goneric.SliceIn(nonUniquelyIdentifying, c)
}
