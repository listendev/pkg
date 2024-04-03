package lockfile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/XANi/goneric"
	"github.com/listendev/pkg/ecosystem"
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

func GetEcosystem(lockfile Lockfile) ecosystem.Ecosystem {
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
func Map(paths []string) map[Lockfile]string {
	lockfiles := goneric.SliceMapSkip(func(f string) (Lockfile, bool) {
		l, e := FromPath(f)
		if e != nil {
			return None, true
		}

		return l, false
	}, paths)

	return lockfiles
}

// Existing splits the given maths into two maps:
// one containing the supported and existing lockfiles,
// the other containing the remainings.
func Existing(paths []string) (map[Lockfile]string, map[Lockfile]error) {
	errorsMap := map[Lockfile]error{}
	lockfilesMap := Map(paths)

	err := validate.Singleton.Var(lockfilesMap, "filevalue")
	if err != nil {
		for _, e := range err.(validate.ValidationErrors) {
			field := Lockfile(strings.TrimRight(strings.TrimLeft(e.Field(), "["), "]"))
			errorsMap[field] = fmt.Errorf(e.Translate(validate.Translator))
			delete(lockfilesMap, field)
		}
	}

	return lockfilesMap, errorsMap
}
