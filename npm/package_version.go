package npm

import (
	stringutil "github.com/listendev/pkg/string/util"
	"golang.org/x/exp/maps"
)

type PackageMaintainer struct {
	Name string `json:"name"`
	Mail string `json:"email"`
}

type PackageMaintainers []PackageMaintainer

// PackageVersion represents the NPM registry response for the route <package_name>/<version>.
type PackageVersion struct {
	Name            string             `json:"name"`
	Description     string             `json:"description"`
	Version         string             `json:"version"`
	Dist            Dist               `json:"dist"`
	Maintainers     PackageMaintainers `json:"maintainers"`
	Scripts         map[string]string  `json:"scripts"`
	Dependencies    map[string]string  `json:"dependencies"`
	DevDependencies map[string]string  `json:"devDependencies"`
}

func (pm PackageMaintainers) Emails() []string {
	ret := map[string]bool{}
	for _, pm := range pm {
		emails := stringutil.MatchEmails(pm.Mail)
		for _, e := range emails {
			ret[e] = true
		}
	}

	return maps.Keys(ret)
}
