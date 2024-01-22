package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/listendev/pkg/npm"
	"github.com/listendev/pkg/observability/tracer"
	"github.com/listendev/pkg/pypi"
)

var (
	errBuilderInvalidAnalysisRequest = errors.New("invalid analysis request")
)

type builder struct {
	ctx                context.Context
	npmRegistryClient  npm.Registry
	pypiRegistryClient pypi.Registry
}

//nolint:revive // we are doing this on purpose (for now)
func NewBuilder(ctx context.Context) (*builder, error) {
	t := tracer.FromContext(ctx)
	if t == nil {
		return nil, fmt.Errorf("couldn't retrieve the tracer from context")
	}

	return &builder{
		ctx:               ctx,
		npmRegistryClient: npm.NewNoOpRegistryClient(),
	}, nil
}

func (b *builder) WithNPMRegistryClient(client npm.Registry) {
	if client == nil || reflect.ValueOf(client).IsNil() {
		b.npmRegistryClient = npm.NewNoOpRegistryClient()

		return
	}
	b.npmRegistryClient = client
}

func (b *builder) WithPyPiRegistryClient(client pypi.Registry) {
	if client == nil || reflect.ValueOf(client).IsNil() {
		b.pypiRegistryClient = pypi.NewNoOpRegistryClient()

		return
	}
	b.pypiRegistryClient = client
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
			res, errJSON := b.FromJSON(msg)
			if errJSON != nil {
				return nil, errJSON
			}
			results = append(results, res)
		}

		return results, nil
	}
	// Try single element
	content := json.RawMessage{}
	if errJSON := json.Unmarshal(data, &content); errJSON != nil {
		return nil, errJSON
	}

	res, err := b.FromJSON(content)
	if err != nil {
		return nil, err
	}

	return []AnalysisRequest{res}, nil
}

func (b *builder) getNPMAnalysisRequest(body []byte) (AnalysisRequest, error) {
	var arn NPM
	if err := json.Unmarshal(body, &arn); err != nil {
		return nil, err
	}

	if err := arn.fillMissingData(b.ctx, b.npmRegistryClient); err != nil {
		return nil, err
	}

	return &arn, nil
}

func (b *builder) getPyPiAnalysisRequest(body []byte) (AnalysisRequest, error) {
	var arp PyPi
	if err := json.Unmarshal(body, &arp); err != nil {
		return nil, err
	}

	if err := arp.fillMissingData(b.ctx, b.pypiRegistryClient); err != nil {
		return nil, err
	}

	return &arp, nil
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
	// TODO: uncomment when ready
	// case NPMTestWhileDynamicInstrumentation:
	// 	fallthrough
	case NPMInstallWhileDynamicInstrumentationAIEnriched:
		fallthrough
	case NPMAdvisory:
		fallthrough
	case NPMTyposquat:
		fallthrough
	case NPMMetadataEmptyDescription:
		fallthrough
	case NPMMetadataMaintainersEmailCheck:
		fallthrough
	case NPMMetadataVersion:
		fallthrough
	case NPMMetadataMismatches:
		fallthrough
	case NPMStaticAnalysisEnvExfiltration:
		fallthrough
	case NPMStaticAnalysisDetachedProcessExecution:
		fallthrough
	case NPMStaticAnalysisEvalBase64:
		fallthrough
	case NPMStaticAnalysisShadyLinks:
		fallthrough
	case NPMStaticAnalysisInstallScript:
		fallthrough
	case NPMStaticNonRegistryDependency:
		fallthrough
	case NPMInstallWhileDynamicInstrumentation:
		return b.getNPMAnalysisRequest(body)

	// PyPi
	case PypiTyposquat:
		fallthrough
	case PypiMetadataMaintainersEmailCheck:
		return b.getPyPiAnalysisRequest(body)

	// NOP
	case Nop:
		return &NOP{arb}, nil
	}

	return nil, errBuilderInvalidAnalysisRequest
}

type noOpBuilder struct {
}

//nolint:revive // we are doing this on purpose (for now)
func NewNoOpBuilder() *noOpBuilder {
	return &noOpBuilder{}
}

func (b *noOpBuilder) FromJSON(_ []byte) (AnalysisRequest, error) {
	return &NOP{}, nil
}
