package analysisrequest

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ComposeAMQPPublishing(a AnalysisRequest) (*amqp.Publishing, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	ret := &amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	if a.Prio() > 0 {
		ret.Priority = a.Prio()
	}

	return ret, nil
}
