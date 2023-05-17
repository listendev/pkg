package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/garnet-org/pkg/validate"
)

const (
	NPMPackageNameMetadataKey    = "npm_package_name"
	NPMPackageVersionMetadataKey = "npm_package_version"
)

func (o *Verdict) Validate() error {
	errors := validate.Validate(o)
	if errors != nil {
		ret := "validation error"
		if len(errors) > 1 {
			ret += "s"
		}
		ret += ": "
		for i, e := range errors {
			if i > 0 {
				ret += "; "
			}
			ret += e.Error()
		}

		return fmt.Errorf(ret)
	}

	return nil
}

func (o *Verdict) UnmarshalJSON(data []byte) error {
	type alias Verdict
	var res alias
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	*o = Verdict(res)

	if err := o.Validate(); err != nil {
		return err
	}

	return nil
}

func (o *Verdict) MarshalJSON() ([]byte, error) {
	err := o.Validate()
	if err != nil {
		return nil, err
	}
	type alias Verdict

	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(o),
	})
}

type Verdicts []Verdict

func FromBuffer(stream io.Reader) (Verdicts, error) {
	b := new(bytes.Buffer)
	b.ReadFrom(stream)

	res := []Verdict{}
	if err := json.Unmarshal(b.Bytes(), &res); err != nil {
		return nil, err
	}

	return res, nil
}

func (v *Verdicts) Buffer() (io.Reader, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	res := bytes.NewReader(buf)
	return res, nil
}
