package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/listendev/pkg/analysisrequest"
	"github.com/listendev/pkg/ecosystem"
	maputil "github.com/listendev/pkg/map/util"
	"github.com/listendev/pkg/models/category"
	"github.com/listendev/pkg/validate"
	"github.com/listendev/pkg/verdictcode"
	"golang.org/x/exp/maps"
)

const (
	NPMPackageNameMetadataKey    = "npm_package_name"
	NPMPackageVersionMetadataKey = "npm_package_version"
)

var (
	CompactMetadata = true
)

func NewEmptyVerdict(eco ecosystem.Ecosystem, org, pkg, version, digest, file string) (*Verdict, error) {
	now := time.Now()
	v := Verdict{
		Ecosystem:  eco,
		Org:        org,
		Pkg:        pkg,
		Version:    version,
		Digest:     digest,
		File:       file,
		CreatedAt:  &now,
		Categories: []category.Category{},    // Forcing empty slice instead of nil
		Metadata:   map[string]interface{}{}, // Forcing empty map instead of nil
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
		ft, _ := analysisrequest.GetTypeForEcosystemFromResultFile(o.Ecosystem, o.File)
		ct, err := o.Code.Type(false)
		// We assume o.Code has been already validated, thus we expect its Type() method to never error
		if err == nil {
			if ft != ct {
				all["Code"] = fmt.Errorf("verdict code is not coherent with the results file and its associated analysis type")
			}
			if !o.Code.UniquelyIdentifies() && o.Fingerprint == "" {
				all["CodeInstance"] = fmt.Errorf("a fingerprint is mandatory because the verdict code is not uniquely identifying it")
			}
		}
	}
	// Other contextual validations
	switch o.Ecosystem {
	case ecosystem.Npm:
		if o.Org != "" {
			if err := validate.Singleton.Var(o.Org, "npmorg"); err != nil {
				var orgErr error
				for _, e := range err.(validate.ValidationErrors) {
					orgErr = fmt.Errorf(e.Translate(validate.Translator))

					break
				}
				all["Org"] = orgErr
			}
			// TODO: use npm_package_name validator on org + pkg when ecosystem is NPM
		}
		if err := validate.Singleton.Var(o.Digest, "shasum"); err != nil {
			var digestErr error
			for _, e := range err.(validate.ValidationErrors) {
				digestErr = fmt.Errorf(e.Translate(validate.Translator))

				break
			}
			all["Digest"] = digestErr
		}
	case ecosystem.Pypi:
		if err := validate.Singleton.Var(o.Org, "pypiorg"); err != nil {
			var orgErr error
			for _, e := range err.(validate.ValidationErrors) {
				orgErr = fmt.Errorf(e.Translate(validate.Translator))

				break
			}
			all["Org"] = orgErr
		}
		if err := validate.Singleton.Var(o.Digest, "blake2b_256"); err != nil {
			var digestErr error
			for _, e := range err.(validate.ValidationErrors) {
				digestErr = fmt.Errorf(e.Translate(validate.Translator))

				break
			}
			all["Digest"] = digestErr
		}
	default:
	}

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
	if res.Categories == nil {
		res.Categories = []category.Category{}
	}
	if res.Metadata == nil {
		res.Metadata = map[string]interface{}{}
	}
	*o = Verdict(res)

	return o.Validate()
}

func (o Verdict) MarshalJSON() ([]byte, error) {
	if o.Categories == nil {
		o.Categories = []category.Category{}
	}
	if o.Metadata == nil {
		o.Metadata = map[string]interface{}{}
	}
	err := o.Validate()
	if err != nil {
		return nil, err
	}

	// Compact the metadata because we don't want a huge JSON with empty/zero values
	if CompactMetadata {
		var compactErr error
		o.Metadata, compactErr = maputil.Compact(o.Metadata)
		if compactErr != nil {
			return nil, compactErr
		}
	}

	type alias Verdict

	return json.Marshal(&struct {
		*alias
	}{
		alias: (*alias)(&o),
	})
}

func (o Verdict) Key() (string, error) {
	if err := o.Validate(); err != nil {
		return "", err
	}

	name := o.Pkg
	if o.Org != "" && o.Ecosystem == ecosystem.Npm {
		name = fmt.Sprintf("%s/%s", o.Org, o.Pkg)
	}

	return fmt.Sprintf("%s/%s/%s/%s/%s", o.Ecosystem.Case(), name, o.Version, o.Digest, o.File), nil
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
