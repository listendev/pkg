package npm

type PackageMaintainer struct {
	Name string `json:"name"`
	Mail string `json:"email"`
}

// PackageVersion represents the NPM registry response for the route <package_name>/<version>.
type PackageVersion struct {
	Name            string              `json:"name"`
	Description     string              `json:"description"`
	Version         string              `json:"version"`
	Dist            Dist                `json:"dist"`
	Maintainers     []PackageMaintainer `json:"maintainers"`
	Scripts         map[string]string   `json:"scripts"`
	Dependencies    map[string]string   `json:"dependencies"`
	DevDependencies map[string]string   `json:"devDependencies"`
}
