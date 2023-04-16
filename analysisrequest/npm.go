package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ AnalysisRequest = (*NPM)(nil)
var _ Publisher = (*NPM)(nil)
var _ Results = (*NPM)(nil)

var (
	errNPMNameEmpty                              = errors.New("npm package name is empty")
	errNPMCouldNotRetrieveLastVersionTag         = errors.New("could not retrieve last npm version tag")
	errNPMCouldNotRetrieveLastVersion            = errors.New("could not retrieve last npm version")
	errNPMCouldNotRetrieveLastVersionShasum      = errors.New("could not retrieve last npm version shasum")
	errNPMCouldNotRetrieveSpecificPackageVersion = errors.New("could not retrieve specific npm package version")
	errNPMVersionDoesNotExistWithShasum          = errors.New("version does not exist on the npm registry with the given shasum")
	errNPMProvidedPackageVersionNotFound         = errors.New("provided npm package version not found")
	errNPMCouldNotRetrieveProvidedVersionShasum  = errors.New("could not retrieve provided npm version shasum")
)

type npmPackage struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Shasum  string `json:"shasum,omitempty"`
}

type NPM struct {
	base
	npmPackage
}

// NewNPM creates an AnalysisRequest for the NPM ecosystem.
func NewNPM(request Type, snowflake string, priority uint8, force bool, name, version, shasum string) (AnalysisRequest, error) {
	tc := request.Components()
	if !tc.HasEcosystem() {
		return nil, fmt.Errorf("couldn't instantiate an analysis request for NPM from a type without ecosystem at all")
	}
	if tc.Ecosystem == NPMEcosystem {
		return &NPM{
			base: base{
				RequestType: request,
				Snowflake:   snowflake,
				Priority:    priority,
				Force:       force,
			},
			npmPackage: npmPackage{
				Name:    name,
				Version: version,
				Shasum:  shasum,
			},
		}, nil
	}

	return nil, fmt.Errorf("couldn't instantiate an analysis request for NPM")
}

func (arn *NPM) UnmarshalJSON(data []byte) error {
	var baseResult base
	if err := json.Unmarshal(data, &baseResult); err != nil {
		return err
	}
	arn.base = baseResult

	var npmResult npmPackage
	if err := json.Unmarshal(data, &npmResult); err != nil {
		return err
	}
	arn.npmPackage = npmResult

	if err := arn.Validate(); err != nil {
		return err
	}

	return nil
}

func (arn NPM) Validate() error {
	if len(arn.Name) == 0 {
		return errNPMNameEmpty
	}

	return arn.base.Validate()
}

func (arn NPM) String() string {
	return arn.Name + "@" + arn.Version + "(" + arn.Type().String() + ")"
}

func (arn NPM) Publishing() (*amqp.Publishing, error) {
	body, err := json.Marshal(arn)
	if err != nil {
		return nil, err
	}

	ret := &amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	if arn.Priority > 0 {
		ret.Priority = arn.Priority
	}

	return ret, nil
}

func (arn *NPM) fillMissingData(parent context.Context, registryClient npm.RegistryClient) error {
	ctx, span := tracer.FromContext(parent).Start(parent, "analysisrequest[npm].fillMissingData")
	defer span.End()

	if len(arn.Version) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, arn.Name)
		if err != nil {
			return err
		}
		latestVersionTag := packageList.DistTags.Latest
		if len(latestVersionTag) == 0 {
			return errNPMCouldNotRetrieveLastVersionTag
		}
		if latestVersion, ok := packageList.Versions[latestVersionTag]; ok {
			arn.Version = latestVersion.Version
			arn.Shasum = latestVersion.Dist.Shasum
		}
		if len(arn.Version) == 0 {
			return errNPMCouldNotRetrieveLastVersion
		}
		if len(arn.Shasum) == 0 {
			return errNPMCouldNotRetrieveLastVersionShasum
		}

		return nil
	}

	if len(arn.Version) > 0 && len(arn.Shasum) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, arn.Name)
		if err != nil {
			return err
		}
		if version, ok := packageList.Versions[arn.Version]; ok {
			arn.Version = version.Version
			arn.Shasum = version.Dist.Shasum
		} else {
			return fmt.Errorf("%v: (version=%s)", errNPMProvidedPackageVersionNotFound, arn.Version)
		}
		if len(arn.Shasum) == 0 {
			return fmt.Errorf("%v: (version=%s)", errNPMCouldNotRetrieveProvidedVersionShasum, arn.Version)
		}

		return nil
	}

	return nil
}

func (arn NPM) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(arn)
}

func (arn NPM) Switch(t Type) (AnalysisRequest, error) {
	c := t.Components()
	if !c.HasEcosystem() {
		return nil, fmt.Errorf("couldn't switch the current NPM analysis request to an analysis request with a type without ecosystem")
	}
	if c.Ecosystem != NPMEcosystem {
		return nil, fmt.Errorf("couldn't switch the current NPM analysis request to a non NPM one")
	}
	arn.RequestType = t

	return &arn, nil
}
