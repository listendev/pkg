package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability/tracer"
)

var (
	errBuilderInvalidAnalysisRequest = errors.New("invalid analysis request")
)

type Builder struct {
	ctx                       context.Context
	npmRegistryRegistryClient npm.RegistryClient
}

func NewBuilder() *Builder {
	return &Builder{
		ctx:                       context.Background(),
		npmRegistryRegistryClient: npm.NewNoOpRegistryClient(),
	}
}

func NewBuilderWithContext(ctx context.Context) *Builder {
	return &Builder{
		ctx:                       ctx,
		npmRegistryRegistryClient: npm.NewNoOpRegistryClient(),
	}
}

func (a *Builder) WithNPMRegistryClient(npmRegistry npm.RegistryClient) {
	a.npmRegistryRegistryClient = npmRegistry
}

func (b *Builder) FromJSON(body []byte) (AnalysisRequest, error) {
	t := tracer.FromContext(b.ctx)
	if t == nil {
		return nil, fmt.Errorf("couldn't retrieve the tracer from context")
	}
	ctx, span := t.Start(b.ctx, "analysysrequest.Builder.UnmarshalJSON")
	defer span.End()

	var arb base
	if err := json.Unmarshal(body, &arb); err != nil {
		return nil, err
	}

	// TODO: adjust condition while evolving
	switch arb.RequestType {
	// NPM
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

		// TODO: this validation below is likely not needed (also in wrong place)
		if len(arn.Version) > 0 && len(arn.Shasum) > 0 {
			packageVersion, err := b.npmRegistryRegistryClient.GetPackageVersion(ctx, arn.Name, arn.Version)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", errNPMCouldNotRetrieveSpecificPackageVersion, err)
			}
			if packageVersion.Dist.Shasum != arn.Shasum {
				return nil, errNPMVersionDoesNotExistWithShasum
			}
		}

		return &arn, nil

	// NOP
	case Nop:
		return &NOP{arb}, nil

	}

	return nil, errBuilderInvalidAnalysisRequest
}
