package npm

type PackageVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Dist    Dist   `json:"dist"`
}
