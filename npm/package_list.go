package npm

import (
	"fmt"
	"time"
)

type VersionTime map[string]time.Time

type PackageList struct {
	Name     string                    `json:"name"`
	Versions map[string]PackageVersion `json:"versions"`
	DistTags DistTags                  `json:"dist-tags"`
	Time     map[string]time.Time      `json:"time"`
}

func (l *PackageList) LatestVersionTime() (*time.Time, error) {
	t, ok := l.Time[l.DistTags.Latest]
	if !ok {
		return nil, fmt.Errorf("couldn't find the release time of the latest version")
	}

	return &t, nil
}
