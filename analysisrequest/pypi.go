package analysisrequest

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type pypiPackage struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	// FIXME: shasum?
	// Shasum  string `json:"shasum,omitempty"`
}

type PyPi struct {
	base
	pypiPackage
}

func (arn PyPi) Publishing() (*amqp.Publishing, error) {
	return ComposeAMQPPublishing(&arn)
}

func (arn PyPi) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(&arn)
}

func (arn PyPi) String() string {
	return arn.Name + "@" + arn.Version + "(" + arn.Type().String() + ")"
}

func (arn PyPi) Delivery() (*amqp.Delivery, error) {
	return ComposeAMQPDelivery(&arn)
}

// FIXME: implement missing methods
