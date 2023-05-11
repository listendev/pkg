package models

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	NPMPackageNameMetadataKey    = "npm_package_name"
	NPMPackageVersionMetadataKey = "npm_package_version"
)

type Verdicts []Verdict

func (v Verdicts) ToBuffer() (io.Reader, error) {
	buf, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	res := bytes.NewReader(buf)
	return res, nil
}
