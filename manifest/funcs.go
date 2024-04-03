package manifest

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/XANi/goneric"
	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/validate"
)

func (s Manifest) String() string {
	//nolint:gocritic // temporary single case switch
	switch s {
	case PackageJSON:
		return string(s)
	}

	return string(None)
}

func FromString(input string) (Manifest, error) {
	s := strings.ToLower(input)

	//nolint:gocritic // temporary single case switch
	switch s {
	case PackageJSON.String():
		return PackageJSON, nil
	}

	return None, fmt.Errorf("the input %q is not a manifest", input)
}

func FromPath(path string) (Manifest, error) {
	_, file := filepath.Split(path)

	return FromString(file)
}

func Ecosystem(manifest Manifest) ecosystem.Ecosystem {
	for eco, files := range manifests {
		if goneric.SliceIn(files, manifest) {
			return eco
		}
	}

	return ecosystem.None
}

// Map maps the given paths to the supported manifests.
func Map(paths []string) map[Manifest]string {
	m := goneric.SliceMapSkip(func(f string) (Manifest, bool) {
		l, e := FromPath(f)
		if e != nil {
			return None, true
		}

		return l, false
	}, paths)

	return m
}

// Existing splits the given maths into two maps:
// one containing the supported and existing manifests,
// the other containing the remainings.
func Existing(paths []string) (map[Manifest]string, map[Manifest]error) {
	errorsMap := map[Manifest]error{}
	manifestsMap := Map(paths)

	err := validate.Singleton.Var(manifestsMap, "filevalue")
	if err != nil {
		for _, e := range err.(validate.ValidationErrors) {
			field := Manifest(strings.TrimRight(strings.TrimLeft(e.Field(), "["), "]"))
			errorsMap[field] = fmt.Errorf(e.Translate(validate.Translator))
			delete(manifestsMap, field)
		}
	}

	return manifestsMap, errorsMap
}
