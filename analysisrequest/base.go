package analysisrequest

import (
	"encoding/json"
	"errors"
)

var _ BasicAnalysisRequest = (*base)(nil)

var (
	errBaseSnowflakeEmpty = errors.New("missing snowflake ID")
)

// base is a struct containing the fields common to all the analysis requests.
//
// Notice it doesn't implement `AnalysisRequest` interface,
// but only the `BasicAnalysisRequest` interface.
type base struct {
	RequestType Type   `json:"type"`
	Snowflake   string `json:"snowflake_id"`
	Priority    uint8  `json:"priority,omitempty"`
	Force       bool   `json:"force"`
}

func (arb *base) UnmarshalJSON(data []byte) error {
	type alias base
	var res alias
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*arb = base(res)

	if err := arb.Validate(); err != nil {
		return err
	}

	return nil
}

func (arb base) Validate() error {
	if len(arb.Snowflake) == 0 {
		return errBaseSnowflakeEmpty
	}

	return nil
}

func (arb base) HasEcosystem() bool {
	return arb.RequestType.HasEcosystem()
}

func (arb base) Type() Type {
	return arb.RequestType
}

func (arb base) ID() string {
	return arb.Snowflake
}

func (arb base) MustProcess() bool {
	return arb.Force
}

func (arb base) Prio() uint8 {
	return arb.Priority
}
