package npm

type Dist struct {
	Shasum string `json:"shasum"`
}

type DistTags struct {
	Latest string `json:"latest"`
	Next   string `json:"next"`
}
