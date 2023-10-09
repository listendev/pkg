package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/listendev/pkg/analysisrequest"
	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/models/category"
	"github.com/listendev/pkg/validate"
	"github.com/listendev/pkg/verdictcode"
	"golang.org/x/exp/maps"
)

const (
	NPMPackageNameMetadataKey    = "npm_package_name"
	NPMPackageVersionMetadataKey = "npm_package_version"
)

func NewEmptyVerdict(eco ecosystem.Ecosystem, org, pkg, version, shasum, file string) (*Verdict, error) {
	now := time.Now()
	v := Verdict{
		Ecosystem:  eco,
		Org:        org,
		Pkg:        pkg,
		Version:    version,
		Shasum:     shasum,
		File:       file,
		CreatedAt:  &now,
		Categories: []category.Category{}, // Forcing empty slice instead of nil
	}
	err := v.Validate()
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (o *Verdict) ExpiresIn(duration time.Duration) {
	t := time.Now().Add(duration)
	o.ExpiresAt = &t
}

func (o *Verdict) HasExpired() bool {
	if o.ExpiresAt == nil {
		return false
	}

	return time.Now().After(*o.ExpiresAt)
}

func (o *Verdict) Validate() error {
	all := map[string]error{}
	if err := validate.Singleton.Struct(o); err != nil {
		for _, e := range err.(validate.ValidationErrors) {
			all[e.StructField()] = fmt.Errorf(e.Translate(validate.Translator))
		}
	}
	// When there aren't already errors on the File field and on the Code field...
	_, fileError := all["File"]
	_, codeError := all["Code"]
	if !fileError && !codeError && o.Code != verdictcode.UNK {
		// We assume o.File has been already validated, thus we don't check for the error
		ft, _ := analysisrequest.GetTypeFromResultFile(o.File)
		ct, err := o.Code.Type(false)
		// We assum o.Code has been already validated, thus we expect its Type() method to never error
		if err == nil && ft != ct {
			all["Code"] = fmt.Errorf("verdict code is not coherent with the results file and its associated analysis type")
		}
	}
	// TODO: use npm_package_name validator on org + pkg
	errors := maps.Values(all)
	if len(errors) > 0 {
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

	return o.Validate()
}

func (o Verdict) MarshalJSON() ([]byte, error) {
	err := o.Validate()
	if err != nil {
		return nil, err
	}
	type alias Verdict
	if o.Categories == nil {
		o.Categories = []category.Category{}
	}

	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(&o),
	})
}

type Verdicts []Verdict

func FromBuffer(stream io.Reader) (Verdicts, error) {
	b := new(bytes.Buffer)
	_, err := b.ReadFrom(stream)
	if err != nil {
		return nil, err
	}

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

	return bytes.NewReader(buf), nil
}
