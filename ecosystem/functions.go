package ecosystem

import (
	"encoding/json"
	"fmt"
	"strings"
)

var all = []Ecosystem{}

func init() {
	index := _Ecosystem_index[1:] // NOTE: None (the default value, 0) is not a valid ecosystem
	for i := range index {
		if i == len(index)-1 {
			break
		}

		all = append(all, Ecosystem(i+1))
	}
}

func All() []Ecosystem {
	return all
}

type EcosystemsOutputOption uint8

const (
	ApplyCase EcosystemsOutputOption = iota
	SingleQuotes
	WithValue
)

func Ecosystems(opts ...EcosystemsOutputOption) []string {
	applyCase := false
	singleQuotes := false
	withValue := false
	for _, o := range opts {
		switch o {
		case ApplyCase:
			applyCase = true
		case SingleQuotes:
			singleQuotes = true
		case WithValue:
			withValue = true
		}
	}

	strs := []string{}
	for _, c := range all {
		s := c.String()
		if applyCase {
			s = c.Case()
		}
		if singleQuotes {
			s = fmt.Sprintf("'%s'", s)
		}
		if withValue {
			s = fmt.Sprintf("%s = %d", s, c)
		}
		strs = append(strs, s)
	}

	return strs
}

func FromUint64(input uint64) (Ecosystem, error) {
	for _, c := range all {
		if uint64(c) == input {
			return c, nil
		}
	}

	return 0, fmt.Errorf("couldn't find a valid ecosystem matching input (%d)", input)
}

func FromString(input string) (Ecosystem, error) {
	for _, c := range all {
		if c.String() == input || c.Case() == input || strings.EqualFold(c.String(), input) {
			return c, nil
		}
	}

	return 0, fmt.Errorf("couldn't find a valid ecosystem matching input %q", input)
}

func (e Ecosystem) Case() string {
	return strings.ToLower(e.String())
}

func (e Ecosystem) MarshalJSON() ([]byte, error) {
	s := e.String()
	if _, err := FromString(s); err != nil {
		return nil, err
	}
	lower := strings.ToLower(s)

	return []byte(fmt.Sprintf("%q", lower)), nil
}

func (e *Ecosystem) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("couldn't unmarshal ecosystem data %q", string(data))
	}

	o, err := FromString(str)
	if err != nil {
		return err
	}
	*e = o

	return nil
}

//go:generate stringer -type=Ecosystem
