package analysisrequest

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BasicAnalysisRequest interface {
	// Type returns the type of the analysis request
	Type() Type
	// ID returns the snowflake ID of the analysis request
	ID() string
	// Priority returns the priority of the analysis request
	Prio() uint8
	// SetPrio lets the change the priority of the analysis request
	SetPrio(uint8)
	// MustProcess returns whether the analysis request must be forcibly processed or not
	MustProcess() bool
	// Validate tells whether the analysis request is ok or not
	Validate() error
}

type AnalysisRequest interface {
	BasicAnalysisRequest
	fmt.Stringer
	Publisher
	Results
}

type Publisher interface {
	Publishing() (*amqp.Publishing, error)
}

// Results
type Results interface {
	// ResultsPath returns the upload path of the analysis request result
	ResultsPath() ResultUploadPath
}

type Builder interface {
	FromJSON([]byte) AnalysisRequest
}
