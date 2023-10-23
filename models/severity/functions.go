package severity

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (s Severity) String() string {
	switch s {
	case Low:
		fallthrough
	case Medium:
		fallthrough
	case High:
		return string(s)
	}

	return string(Unknown)
}

// Scan implements the sql.Scanner interface.
func (s *Severity) Scan(source any) error {
	if x, ok := source.(string); ok {
		cat, err := New(x)
		if cat != Unknown && err != nil {
			return err
		}
		*s = cat

		return nil
	}

	return fmt.Errorf("cannot scan %T into Category", source)
}

func New(input string) (Severity, error) {
	s := strings.ToLower(input)

	switch s {
	case Low.String():
		return Low, nil
	case Medium.String():
		return Medium, nil
	case High.String():
		return High, nil
	}

	return Unknown, fmt.Errorf("the input %q is not a severity", input)
}

func (s *Severity) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("couldn't unmarshal severity data %q", string(data))
	}

	o, e := New(str)
	if e != nil {
		return e
	}
	*s = o

	return nil
}

func (s Severity) MarshalJSON() ([]byte, error) {
	str := s.String()
	if _, err := New(str); err != nil {
		return nil, err
	}
	lower := strings.ToLower(str)

	return []byte(fmt.Sprintf("%q", lower)), nil
}
