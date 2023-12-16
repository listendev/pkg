package pypi

import (
	"errors"
	"time"
)

var (
	ErrNoDist = errors.New("could not find a version with 'sdist' package type")
)

type PackageVersion struct {
	Name        string // FIXME: ...
	Version     string
	Digests     Digests   `json:"digests"`
	PackageType string    `json:"packagetype"`
	UploadTime  time.Time `json:"upload_time_iso_8601"`
	Filename    string    `json:"filename"`
}

type PackageVersions []PackageVersion

func (vs PackageVersions) GetSdist() (*PackageVersion, error) {
	for _, v := range vs {
		if v.PackageType == "sdist" {
			return &v, nil
		}
	}

	return nil, ErrNoDist
}
