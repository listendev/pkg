package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	case NPMTestWhileFalco:
		fallthrough

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
