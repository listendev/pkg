package analysisrequest

import (
	"path"
)

type ResultUploadPath []string

func (r ResultUploadPath) ToS3Key() string {
	return path.Join(r...)
}

func ComposeResultUploadPath(a AnalysisRequest) ResultUploadPath {
	t := a.Type()
	c := t.Components()
	filename := c.ResultFile()

	if c.Ecosystem == NPMEcosystem {
		arn := a.(NPM)

		return ResultUploadPath{string(c.Ecosystem), arn.Name, arn.Version, arn.Shasum, filename}
	}

	// Assuming there are no types - other than Nop - without ecosystem
	return ResultUploadPath{"nop", a.ID(), filename}
}

func GetResultFilesByEcosystem(eco Ecosystem) map[Type]string {
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
