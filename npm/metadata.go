package npm

type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

type Bugs struct {
	URL string `json:"url"`
}

type Engines struct {
	Node string `json:"node"`
}

type Scripts struct {
	Test  string `json:"test"`
	Bench string `json:"bench"`
}

type Signature struct {
	Keyid string `json:"keyid"`
	Sig   string `json:"sig"`
}

type Dist struct {
	Integrity    string      `json:"integrity"`
	Shasum       string      `json:"shasum"`
	Tarball      string      `json:"tarball"`
	FileCount    int         `json:"fileCount"`
	UnpackedSize int         `json:"unpackedSize"`
	Signatures   []Signature `json:"signatures"`
	NpmSignature string      `json:"npm-signature"`
}

type NPMUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DistTags struct {
	Latest string `json:"latest"`
	Next   string `json:"next"`
}
