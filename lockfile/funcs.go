package lockfile

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

func (s Lockfile) String() string {
	switch s {
	case PackageLockJSON:
		fallthrough
	case PoetryLock:
		return string(s)
	}

	return string(None)
}

func FromString(input string) (Lockfile, error) {
	s := strings.ToLower(input)

	switch s {
	case PackageLockJSON.String():
		return PackageLockJSON, nil

	case PoetryLock.String():
		return PoetryLock, nil
	}

	return None, fmt.Errorf("the input %q is not a lockfile", input)
}

func FromPath(path string) (Lockfile, error) {
	_, file := filepath.Split(path)

	return FromString(file)
}

func Ecosystem(lockfile Lockfile) ecosystem.Ecosystem {
	for eco, files := range lockfiles {
		if goneric.SliceIn(files, lockfile) {
			return eco
		}
	}

	return ecosystem.None
}

// FromEcosystem returns the list of supported lock files for the given ecosystem.
func FromEcosystem(eco ecosystem.Ecosystem) []Lockfile {
	return lockfiles[eco]
}

// Map maps the given paths to the supported lockfiles.
func Map(paths []string) map[Lockfile][]string {
	ret := map[Lockfile][]string{}
	for _, f := range paths {
		l, e := FromPath(f)
		if e != nil {
			continue
		}

		if _, ok := ret[l]; !ok {
			ret[l] = []string{f}
		} else if !goneric.SliceIn(ret[l], f) {
			ret[l] = append(ret[l], f)
		}
	}

	return ret
}

// Existing splits the given maths into two maps:
// one containing the supported and existing lockfiles,
// the other containing the remainings.
func Existing(paths []string) (map[Lockfile][]string, map[Lockfile][]error) {
	errorsMap := map[Lockfile][]error{}
	lockfilesMap := Map(paths)

	err := validate.Singleton.Var(lockfilesMap, "filevalue")
	if err != nil {
		for _, e := range err.(validate.ValidationError) {
			field, idx := fromFieldError(e)
			// spew.Dump(field.String())
			messageErr := fmt.Errorf("%s", e.Translate(validate.Translator))
			if _, ok := errorsMap[field]; !ok {
				errorsMap[field] = []error{messageErr}
			} else {
				errorsMap[field] = append(errorsMap[field], messageErr)
			}
			if len(lockfilesMap[field]) < 2 {
				delete(lockfilesMap, field)
			} else {
				lockfilesMap[field] = typeutil.RemoveFromSliceAt(lockfilesMap[field], idx)
			}
		}
	}

	return lockfilesMap, errorsMap
}

func fromFieldError(e validator.FieldError) (Lockfile, int) {
	re := regexp.MustCompile(`\[(?P<field>.*?)\]\[(?P<index>\d+)\]`)
	m := re.FindStringSubmatch(e.Field())
	tmp := make(map[string]string)
	for i, group := range re.SubexpNames() {
		if i != 0 && group != "" {
			tmp[group] = m[i]
		}
	}
	idx, _ := strconv.Atoi(tmp["index"])

	return Lockfile(tmp["field"]), idx
}
