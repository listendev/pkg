package models

import (
	"bytes"
	"context"
	"encoding/json"
)

type FilterParsingError struct {
	message string
}

func (e *FilterParsingError) Error() string {
	return e.message
}

// Filter applies a JSONPath expression to the receiving Verdicts and returns the result.
//
// See the JSONPath expressions documentation for more information about the syntax:
// https://goessner.net/articles/JsonPath/index.html#e2.
//
// The result is returned as raw (interface{} or interface{} slice) and as a Verdicts instance.
// Notice that, depending on the JSONPath input, the result may have a different structure than the Verdicts one:
// in this case this function returns the raw result, a nil Verdicts, and an error due to the conversion failure.
func (v *Verdicts) Filter(c context.Context, jsonpath string) (interface{}, Verdicts, error) {
	r, err := v.Buffer()
	if err != nil {
		return nil, nil, err
	}

	var iface interface{}
	if decodeErr := json.NewDecoder(r).Decode(&iface); decodeErr != nil {
		return nil, nil, decodeErr
	}

	eval, err := lang.NewEvaluableWithContext(c, jsonpath)
	if err != nil {
		e := &FilterParsingError{message: err.Error()}

		return nil, nil, e
	}

	raw, err := eval(c, iface)
	if err != nil {
		return nil, nil, err
	}

	// Convert result to JSON, to Verdicts
	jsonres, _ := json.Marshal(raw)
	verdicts, err := FromBuffer(bytes.NewReader(jsonres))

	return raw, verdicts, err
}
