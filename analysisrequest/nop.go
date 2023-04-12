package analysisrequest

import "encoding/json"

type NOP struct {
	base
}

func (a *NOP) UnmarshalJSON(data []byte) error {
	return a.base.UnmarshalJSON(data)
}

func (a NOP) MarshalJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a NOP) Validate() error {
	return a.base.Validate()
}

func (a NOP) ResultUploadPath() ResultUploadPath {
	return ResultUploadPath{
		"nop",
		a.Snowflake,
	}
}

func (a NOP) String() string {
	return a.RequestType.String() + "-" + a.Snowflake
}
