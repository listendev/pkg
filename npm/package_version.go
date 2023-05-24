package npm

type PackageMaintainer struct {
	Name string `json:"name"`
	Mail string `json:"email"`
}

type PackageVersion struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Version     string              `json:"version"`
	Dist        Dist                `json:"dist"`
	Maintainers []PackageMaintainer `json:"maintainers"`
}
