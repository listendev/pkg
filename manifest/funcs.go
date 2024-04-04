package manifest

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/XANi/goneric"
	"github.com/go-playground/validator/v10"
	"github.com/listendev/pkg/ecosystem"
	typeutil "github.com/listendev/pkg/type/util"
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
func Map(paths []string) map[Manifest][]string {
	ret := map[Manifest][]string{}
	for _, f := range paths {
		l, e := FromPath(f)
		if e != nil {
			continue
		}

		if _, ok := ret[l]; !ok {
			ret[l] = []string{f}
		} else {
			if !goneric.SliceIn(ret[l], f) {
				ret[l] = append(ret[l], f)
			}
		}
	}

	return ret
}

// Existing splits the given maths into two maps:
// one containing the supported and existing manifests,
// the other containing the remainings.
func Existing(paths []string) (map[Manifest][]string, map[Manifest][]error) {
	errorsMap := map[Manifest][]error{}
	manifestsMap := Map(paths)

	err := validate.Singleton.Var(manifestsMap, "filevalue")
	if err != nil {
		for _, e := range err.(validate.ValidationErrors) {
			field, idx := fromFieldError(e)
			// spew.Dump(field.String())
			messageErr := fmt.Errorf(e.Translate(validate.Translator))
			if _, ok := errorsMap[field]; !ok {
				errorsMap[field] = []error{messageErr}
			} else {
				errorsMap[field] = append(errorsMap[field], messageErr)
			}
			if len(manifestsMap[field]) < 2 {
				delete(manifestsMap, field)
			} else {
				manifestsMap[field] = typeutil.RemoveFromSliceAt(manifestsMap[field], idx)
			}
		}
	}

	return manifestsMap, errorsMap
}

func fromFieldError(e validator.FieldError) (Manifest, int) {
	re := regexp.MustCompile(`\[(?P<field>.*?)\]\[(?P<index>\d+)\]`)
	m := re.FindStringSubmatch(e.Field())
	tmp := make(map[string]string)
	for i, group := range re.SubexpNames() {
		if i != 0 && group != "" {
			tmp[group] = m[i]
		}
	}
	idx, _ := strconv.Atoi(tmp["index"])

	return Manifest(tmp["field"]), idx
}
