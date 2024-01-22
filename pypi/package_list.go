package pypi

import (
	"time"

	"github.com/XANi/goneric"
)

type PackageList struct {
	Info     Info                       `json:"info"`
	Versions map[string]PackageVersions `json:"releases"`
	URLs     PackageVersions            `json:"urls"`
}

// Fill the version and name field for all the package versions.
func (p *PackageList) Fill() {
	for v := range p.Versions {
		for i := range p.Versions[v] {
			p.Versions[v][i].Version = v
			p.Versions[v][i].Name = p.Info.Name
			p.Versions[v][i].Authors = []PackageMaintainer{
				{
					Name: p.Info.Author,
					Mail: p.Info.AuthorEmail,
					Type: PackageAuthorType,
				},
			}
			p.Versions[v][i].Maintainers = []PackageMaintainer{
				{
					Name: p.Info.Maintainer,
					Mail: p.Info.MaintainerEmail,
					Type: PackageMaintainerType,
				},
			}
		}
	}
}

func (p *PackageList) GetVersion(version string) (*PackageVersion, error) {
	// Detect if the receiving PackageList instance was created from the version endpoint response.
	if len(p.Versions) == 0 {
		if version == "latest" {
			return nil, ErrLatestVersionNotFound
		}
		if p.Info.Version != version {
			return nil, ErrVersionMismatch
		}
		pv, err := p.URLs.GetSdist()
		if err != nil {
			return nil, ErrMissingSdistPackageVersion
		}
		// We store version and name manually because the response doesn't contain them at this level
		pv.Version = version
		pv.Name = p.Info.Name

		return pv, nil
	}

	// Otherwise, the receiving PackageList instance was created from the list endpoint response.
	latest := false
	if version == "latest" {
		// In this case the version in the info part is the latest version
		version = p.Info.Version
		latest = true
	}
	pvs, ok := p.Versions[version]
	if !ok {
		if latest {
			return nil, ErrLatestVersionNotFound
		}

		return nil, ErrVersionNotFound
	}
	pv, err := pvs.GetSdist()
	if err != nil {
		return nil, ErrMissingSdistPackageVersion
	}
	// We store version and name manually because the response doesn't contain them at this level
	pv.Version = version
	pv.Name = p.Info.Name

	return pv, nil
}

func (p *PackageList) LatestVersionTime() (*time.Time, error) {
	// TODO: keep in mind we are only considering the "sdist" version (not the wheel distribution one) for now
	latest, e := p.GetVersion("latest")
	if e != nil {
		return nil, e
	}

	return &latest.UploadTime, nil
}

func (p *PackageList) MaintainersByVersion(version string) (PackageMaintainers, error) {
	v, e := p.GetVersion(version)
	if e != nil {
		return nil, e
	}

	return goneric.SliceDedupe(append(v.Authors, v.Maintainers...)), nil
}
