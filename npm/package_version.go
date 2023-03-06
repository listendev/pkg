package npm

type PackageVersion struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Description     string            `json:"description"`
	License         string            `json:"license"`
	Repository      RepositoryUnion   `json:"repository"`
	Funding         string            `json:"funding"`
	Type            string            `json:"type"`
	Main            string            `json:"main"`
	Types           string            `json:"types"`
	Engines         Engines           `json:"engines"`
	Scripts         Scripts           `json:"scripts"`
	Keywords        []string          `json:"keywords"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string
	GitHead         string       `json:"gitHead"`
	Bugs            Bugs         `json:"bugs"`
	Homepage        string       `json:"homepage"`
	ID              string       `json:"_id"`
	NodeVersion     string       `json:"_nodeVersion"`
	NpmVersion      string       `json:"_npmVersion"`
	Dist            Dist         `json:"dist"`
	NpmUser         NPMUser      `json:"_npmUser"`
	Maintainers     []Maintainer `json:"maintainers"`
}
