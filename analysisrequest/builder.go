package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability/tracer"
)

var (
	errBuilderInvalidAnalysisRequest = errors.New("invalid analysis request")
)

type builder struct {
	ctx                       context.Context
	npmRegistryRegistryClient npm.RegistryClient // FIXME: make this come from context too?
}

func NewBuilder(ctx context.Context) (*builder, error) {
	t := tracer.FromContext(ctx)
	if t == nil {
		return nil, fmt.Errorf("couldn't retrieve the tracer from context")
	}

	return &builder{
		ctx:                       ctx,
		npmRegistryRegistryClient: npm.NewNoOpRegistryClient(),
	}, nil
}

func (a *builder) WithNPMRegistryClient(npmRegistry npm.RegistryClient) {
	if npmRegistry == nil || reflect.ValueOf(npmRegistry).IsNil() {
		a.npmRegistryRegistryClient = npm.NewNoOpRegistryClient()

		return
	}
	a.npmRegistryRegistryClient = npmRegistry
}

func (b *builder) FromFile(path string) ([]AnalysisRequest, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could not find the file at path %s: %v", path, err)
	}
	if !fileInfo.Mode().IsRegular() || filepath.Ext(path) != ".json" {
		return nil, fmt.Errorf("the file at path %q is not a valid json file: %v", path, err)
	}

	// Read file content
	data, readErr := os.ReadFile(path)
	if readErr != nil {
		return nil, fmt.Errorf("could not open the file at path %q: %v", path, readErr)
	}

	contents := []json.RawMessage{}
	decodeErr := json.Unmarshal(data, &contents)
	if decodeErr == nil {
		// Try list of elements
		results := []AnalysisRequest{}
		for _, msg := range contents {
			res, err := b.FromJSON(msg)
			if err != nil {
				return nil, err
			}
			results = append(results, res)
		}

		return results, nil
	} else {
		// Try single element
		content := json.RawMessage{}
		if err := json.Unmarshal(data, &content); err != nil {
			return nil, err
		}

		res, err := b.FromJSON(content)
		if err != nil {
			return nil, err
		}

		return []AnalysisRequest{res}, nil
	}
}

func (b *builder) FromJSON(body []byte) (AnalysisRequest, error) {
	t := tracer.FromContext(b.ctx)
	_, span := t.Start(b.ctx, "analysysrequest.Builder.UnmarshalJSON")
	defer span.End()

	var arb base
	if err := json.Unmarshal(body, &arb); err != nil {
		return nil, err
	}

	// TODO: adjust condition while evolving
	switch arb.RequestType {
	// NPM
	case NPMGPT4InstallWhileFalco:
		// Not sure we wanna create enrichers from JSON too, but it shouldn't do any harm
		fallthrough

	case NPMDepsDev:
		fallthrough

	// TODO: uncomment when ready
	// case NPMTestWhileFalco:
	// 	fallthrough

	case NPMTyposquat:
		fallthrough

	case NPMMetadataEmptyDescription:
		fallthrough

	// TODO: uncomment when ready
	// case NPMMetadataMaintainersEmailCheck:
	// 	fallthrough

	case NPMMetadataVersion:
		fallthrough

	// TODO: uncomment when ready
	// case NPMSemgrepEnvExfiltration:
	// 	fallthrough

	// TODO: uncomment when ready
	// case NPMSemgrepProcessExecution:
	// 	fallthrough

	// TODO: uncomment when ready
	// case NPMSemgrepEvalBase64:
	// 	fallthrough

	// TODO: uncomment when ready
	// case NPMSemgrepShadyLinks:
	// 	fallthrough

	case NPMInstallWhileFalco:
		var arn NPM
		if err := json.Unmarshal(body, &arn); err != nil {
			return nil, err
		}

		if err := arn.fillMissingData(b.ctx, b.npmRegistryRegistryClient); err != nil {
			return nil, err
		}

		return &arn, nil

	// NOP
	case Nop:
		return &NOP{arb}, nil

	}

	return nil, errBuilderInvalidAnalysisRequest
}

type noOpBuilder struct {
}

func NewNoOpBuilder() *noOpBuilder {
	return &noOpBuilder{}
}

func (b *noOpBuilder) FromJSON(body []byte) (AnalysisRequest, error) {
	return &NOP{}, nil
}
