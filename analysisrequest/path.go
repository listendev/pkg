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

func GetResultFilesByEcosystem(eco Ecosystem) []string {
	res := []string{}
	for t := range typeURNs {
		c := t.Components()
		if c.Ecosystem == eco {
			res = append(res, c.ResultFile())
		}
	}

	return res
}
