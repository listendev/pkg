package analysisrequest

import (
	"path"
)

type ResultUploadPath [4]string

func (r ResultUploadPath) ToS3Key() string {
	return path.Join(r[:]...)
}
