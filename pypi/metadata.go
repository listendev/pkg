package pypi

type Digests struct {
	MD5         string `json:"md5"`
	SHA256      string `json:"sha256"`
	Blake2bB256 string `json:"blake2b_256"`
}

type Info struct {
	Name            string `json:"name"`
	Author          string `json:"author,omitempty"`
	AuthorEmail     string `json:"author_email,omitempty"`
	Maintainer      string `json:"maintainer,omitempty"`
	MaintainerEmail string `json:"maintainer_email,omitempty"`
	Description     string `json:"description"`
	// Version is the last version (list endpoint) or the current version (version endpoint)
	Version string `json:"version"`
}
