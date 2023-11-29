package analysisrequest

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ AnalysisRequest = (*NOP)(nil)
var _ Publisher = (*NOP)(nil)
var _ Deliverer = (*NOP)(nil)
var _ Results = (*NOP)(nil)

type NOP struct {
	base
}

func NewNOP(snowflake string, priority uint8, force bool) AnalysisRequest {
	return &NOP{
		base: base{
			RequestType: Nop,
			Snowflake:   snowflake,
			Priority:    priority,
			Force:       force,
		},
	}
}

func (a *NOP) UnmarshalJSON(data []byte) error {
	return a.base.UnmarshalJSON(data)
}

func (a NOP) Validate() error {
	return a.base.Validate()
}

func (a NOP) String() string {
	return a.RequestType.String() + "-" + a.Snowflake
}

func (a NOP) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(&a)
}

func (a NOP) Publishing() (*amqp.Publishing, error) {
	return ComposeAMQPPublishing(&a)
}

func (a NOP) Delivery() (*amqp.Delivery, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	ret := &amqp.Delivery{
		ContentType: "application/json",
		Body:        body,
	}
	if a.Priority > 0 {
		ret.Priority = a.Priority
	}

	return ret, nil
}
