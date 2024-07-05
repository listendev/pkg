package detectiontype

import (
	"encoding/json"
	"fmt"
	"strings"
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

func EventTypesAsStrings() []string {
	strs := []string{}
	for _, e := range events {
		strs = append(strs, e.String())
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
		if e.String() == input || strings.EqualFold(e.String(), input) {
			return e, nil
		}
	}

	return 0, fmt.Errorf("couldn't find an event type matching the input %q", input)
}

func (e Event) MarshalJSON() ([]byte, error) {
	s := e.String()
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

//go:generate stringer -type=Event
