package analysisrequest

import (
	"fmt"
	"path"
	"strings"
)

type ResultUploadPath []string

func (r ResultUploadPath) ToS3Key() string {
	return path.Join(r...)
}

func ComposeResultUploadPath(a AnalysisRequest) ResultUploadPath {
	t := a.Type()
	c := t.Components()
	filename := c.Collector
	suffix := strings.Join(c.Actions, ",")
	if len(suffix) > 0 {
		filename = fmt.Sprintf("%s(%s)", filename, suffix)
	}
	filename += "." + strings.TrimPrefix(c.Format, ".")

	if c.Ecosystem == EcosystemNPM {
		arn := a.(NPM)

		return ResultUploadPath{c.Ecosystem, arn.Name, arn.Version, arn.Shasum, filename}
	}

	// Assuming there are no types - other than Nop - without ecosystem
	return ResultUploadPath{"nop", a.ID(), filename}
}
