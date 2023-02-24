package npm

type Dist struct {
	Shasum       string `json:"shasum"`
	TarballURL   string `json:"tarball"`
	NumFiles     int    `json:"fileCount"`
	UnpackedSize int    `json:"unpackedSize"`
}

type DistTags struct {
	Latest string `json:"latest"`
	Next   string `json:"next"`
}
