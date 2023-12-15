package analysisrequest

import (
	"fmt"
	"path"

	"github.com/listendev/pkg/ecosystem"
)

type ResultUploadPath []string

// Key returns a path-like key.
func (r ResultUploadPath) Key() string {
	return path.Join(r...)
}

func ComposeResultUploadPath(a AnalysisRequest) ResultUploadPath {
	t := a.Type()
	c := t.Components()
	filename := c.ResultFile()

	switch c.Ecosystem {
	case ecosystem.Npm:
		arn := a.(*NPM)

		return ResultUploadPath{c.Ecosystem.Case(), arn.Name, arn.Version, arn.Shasum, filename}

	case ecosystem.Pypi:
		arp := a.(*PyPi)

		return ResultUploadPath{c.Ecosystem.Case(), arp.Name, arp.Version, filename} // FIXME: shasum or similar?
	}

	// Assuming there are no types - other than Nop - without ecosystem
	return ResultUploadPath{"nop", a.ID(), filename}
}

func GetResultFilesByEcosystem(eco ecosystem.Ecosystem) map[Type]string {
	tmp := map[string]Type{}
	for t := range typeURNs {
		c := t.Components()
		_, notEnricherErr := t.Parent()
		if c.Ecosystem == eco && notEnricherErr != nil {
			f := c.ResultFile()
			if _, ok := tmp[f]; !ok {
				tmp[f] = t
			}
		}
	}

	res := make(map[Type]string, len(tmp))
	for k, v := range tmp {
		res[v] = k
	}

	return res
}

func GetTypeForEcosystemFromResultFile(eco ecosystem.Ecosystem, filename string) (Type, error) {
	all := GetResultFilesByEcosystem(eco)
	for t, f := range all {
		if f == filename {
			return t, nil
		}
	}

	return Nop, fmt.Errorf("couldn't find any type for ecosystem %q matching the results file %q", eco.Case(), filename)
}

func GetTypesFromResultFile(filename string) ([]Type, error) {
	res := []Type{}
	for _, e := range ecosystem.All() {
		t, err := GetTypeForEcosystemFromResultFile(e, filename)
		if err == nil {
			res = append(res, t)
		}
	}
	if len(res) > 0 {
		return res, nil
	}

	return nil, fmt.Errorf("couldn't find any type in any ecosystem matching the results file %q", filename)
}
