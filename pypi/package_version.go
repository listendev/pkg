package pypi

import (
	"errors"
	"time"

	"golang.org/x/exp/maps"
)

var (
	ErrNoDist = errors.New("could not find a version with 'sdist' package type")
)

type MaintainerType string

const (
	PackageAuthorType     MaintainerType = "author"
	PackageMaintainerType MaintainerType = "maintainer"
)

type PackageMaintainer struct {
	Name string
	Mail string
	Type MaintainerType
}

type PackageMaintainers []PackageMaintainer

func (pm PackageMaintainers) Emails() []string {
	ret := map[string]bool{}
	for _, pm := range pm {
		ret[pm.Mail] = true
	}

	return maps.Keys(ret)
}

type PackageVersion struct {
	// These without json tags get filled in
	Name        string
	Version     string
	Authors     PackageMaintainers
	Maintainers PackageMaintainers

	URL         string    `json:"url"`
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
