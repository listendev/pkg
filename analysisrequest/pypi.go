package analysisrequest

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ AnalysisRequest = (*PyPi)(nil)
var _ Publisher = (*PyPi)(nil)
var _ Deliverer = (*PyPi)(nil)
var _ Results = (*PyPi)(nil)

var (
	errPyPiNameEmpty = errors.New("PyPi package name is empty")
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

func (arp PyPi) Publishing() (*amqp.Publishing, error) {
	return ComposeAMQPPublishing(&arp)
}

func (arp PyPi) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(&arp)
}

func (arp PyPi) String() string {
	return arp.Name + "@" + arp.Version + "(" + arp.Type().String() + ")"
}

func (arp PyPi) Delivery() (*amqp.Delivery, error) {
	return ComposeAMQPDelivery(&arp)
}

func (arp PyPi) Validate() error {
	if len(arp.Name) == 0 {
		return errPyPiNameEmpty
	}

	return arp.base.Validate()
}

// FIXME: implement missing methods
