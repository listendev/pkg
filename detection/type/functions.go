package detectiontype

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

var events = []Event{}

func init() {
	// Fetch all valid event types
	evtIdx := _Event_index[1:] // NOTE: None (the default value, 0) is not a valid event type
	for i := range evtIdx {
		if i == len(evtIdx)-1 {
			break
		}

		events = append(events, Event(i+1))
	}
}

func EventTypes() []Event {
	return events
}

type EventOutputOption uint8

const (
	ApplyCase EventOutputOption = iota
	SingleQuotes
	WithValue
)

func EventTypesAsStrings(opts ...EventOutputOption) []string {
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
	for _, e := range events {
		s := e.String()
		if applyCase {
			s = e.Case()
		}
		if singleQuotes {
			s = fmt.Sprintf("'%s'", s)
		}
		if withValue {
			s = fmt.Sprintf("%s = %d", s, e)
		}
		strs = append(strs, s)
	}

	return strs
}

func FromUint64(input uint64) (Event, error) {
	for _, e := range events {
		if uint64(e) == input {
			return e, nil
		}
	}

	return 0, fmt.Errorf("couldn't find an event type matching the input %q", input)
}

func FromString(input string) (Event, error) {
	for _, e := range events {
		if e.String() == input || e.Case() == input || strings.EqualFold(e.String(), input) {
			return e, nil
		}
	}

	return 0, fmt.Errorf("couldn't find an event type matching the input %q", input)
}

func (e Event) MarshalJSON() ([]byte, error) {
	s := e.Case()
	if _, err := FromString(s); err != nil {
		return nil, err
	}
	lower := strings.ToLower(s)

	return []byte(fmt.Sprintf("%q", lower)), nil
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("couldn't unmarshal event type data %q", string(data))
	}

	o, err := FromString(str)
	if err != nil {
		return err
	}
	*e = o

	return nil
}

// Scan implements the sql.Scanner interface.
func (e *Event) Scan(source any) error {
	if x, ok := source.(string); ok {
		eco, err := FromString(x)
		if err != nil {
			return err
		}
		*e = eco

		return nil
	}
	if x, ok := source.(uint64); ok {
		eco, err := FromUint64(x)
		if err != nil {
			return err
		}
		*e = eco

		return nil
	}

	return fmt.Errorf("cannot scan %T into Event type", source)
}

func (e *Event) UnmarshalText(text []byte) error {
	eco, err := FromString(string(text))
	*e = eco

	return err
}

// Case returns the snake case string version of the receiving Event enum type.
func (e Event) Case() string {
	return strings.ToLower(strcase.ToSnake(e.String()))
}

//go:generate stringer -type=Event
