package npm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type PackageJSON struct {
	Name                 string            `json:"name"`
	Scripts              map[string]string `json:"scripts"`
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	BundleDependencies   []string          `json:"bundleDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func NewPackageJSONFromDir(dir string) (*PackageJSON, error) {
	reader, err := readPackageJSON(dir)
	if err != nil {
		return nil, err
	}

	return NewPackageJSONFromReader(reader)
}

func NewPackageJSONFromReader(reader io.Reader) (*PackageJSON, error) {
	ret := &PackageJSON{}
	if err := json.NewDecoder(reader).Decode(ret); err != nil {
		return nil, errors.New("couldn't instantiate from the input package.json contents")
	}

	return ret, nil
}

func readPackageJSON(dir string) (io.Reader, error) {
	name := filepath.Join(dir, "package.json")

	f, err := activeFS.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory %s does not contain a package.json file", dir)
		}

		return nil, errors.New("couldn't read the package.json file")
	}

	return f, nil
}
