package analysisresult

import (
	"bytes"
	"encoding/json"
	"io"
)

type Severity string

const (
	SeverityLow    Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh   Severity = "high"
)

type Metadata map[string]interface{}

const (
	NPMPackageNameMetadataKey    = "npm_package_name"
	NPMPackageVersionMetadataKey = "npm_package_version"
)

type Verdict struct {
	Message  string                 `json:"message"`
	Severity Severity               `json:"severity"`
	Metadata map[string]interface{} `json:"metadata"`
}

type Verdicts []Verdict

func (v Verdicts) ToBuffer() (io.Reader, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	res := bytes.NewReader(buf)
	return res, nil

}
