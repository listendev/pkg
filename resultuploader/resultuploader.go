package resultuploader

import (
	"context"
	"io"

	"github.com/garnet-org/pkg/analysisrequest"
)

type ResultUploader interface {
	String() string
	Upload(ctx context.Context, src io.Reader, uploadPath analysisrequest.ResultUploadPath) error
	AlreadyExists(ctx context.Context, uploadPath analysisrequest.ResultUploadPath) bool
}
