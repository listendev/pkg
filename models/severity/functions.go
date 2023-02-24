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

	return "unknown"
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

	return "unknown", fmt.Errorf("the input %q is not a severity", input)
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
