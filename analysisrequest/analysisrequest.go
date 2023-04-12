package analysisrequest

import (
	"encoding/json"
	"fmt"
)

type BasicAnalysisRequest interface {
	// Type returns the type of the analysis request
	Type() Type
	// ID returns the snowflake ID of the analysis request
	ID() string
	// Validate tells whether the analysis request is ok or not
	Validate() error
}

type AnalysisRequest interface {
	BasicAnalysisRequest
	// ResultUploadPath returns the upload path of the analysis request result
	ResultUploadPath() ResultUploadPath
	//json.Unmarshaler // FIXME: ...
	json.Marshaler
	fmt.Stringer
}
