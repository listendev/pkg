package npm

import "time"

type VersionTime map[string]time.Time

type PackageList struct {
	Name     string                    `json:"name"`
	Versions map[string]PackageVersion `json:"versions"`
	DistTags DistTags                  `json:"dist-tags"`
}
