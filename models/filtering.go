package models

import (
	"bytes"
	"context"
	"encoding/json"
)

// Filter applies a JSONPath expression to the receiving Verdicts and returns the result.
//
// See the JSONPath documentation for more information about the syntax:
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
	if err := json.NewDecoder(r).Decode(&iface); err != nil {
		return nil, nil, err
	}

	eval, err := lang.NewEvaluableWithContext(c, jsonpath)
	if err != nil {
		return nil, nil, err
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
