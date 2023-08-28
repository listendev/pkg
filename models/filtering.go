package models

import (
	"bytes"
	"context"
	"encoding/json"
)

// Filter ...
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

	// FIXME: eval works on struct, check whether it also works on a structs slice
	raw, err := eval(c, iface)
	if err != nil {
		return nil, nil, err
	}

	// Convert result to JSON, to Verdicts
	jsonres, _ := json.Marshal(raw)
	verdicts, err := FromBuffer(bytes.NewReader(jsonres))

	return raw, verdicts, err
}
