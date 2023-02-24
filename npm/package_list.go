package npm

import "time"

type VersionTime map[string]time.Time

type PackageList struct {
	ID             string                    `json:"_id"`
	Rev            string                    `json:"_rev"`
	Name           string                    `json:"name"`
	Description    string                    `json:"description"`
	DistTags       DistTags                  `json:"dist-tags"`
	Versions       map[string]PackageVersion `json:"versions"`
	Readme         string                    `json:"readme"`
	Maintainers    []Maintainer              `json:"maintainers"`
	Time           VersionTime               `json:"time"`
	Repository     Repository                `json:"repository"`
	Users          map[string]bool           `json:"users"`
	Homepage       string                    `json:"homepage"`
	Keywords       []string                  `json:"keywords"`
	Bugs           Bugs                      `json:"bugs"`
	License        string                    `json:"license"`
	ReadmeFilename string                    `json:"readmeFilename"`
}
