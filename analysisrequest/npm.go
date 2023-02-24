package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/npm"
	"github.com/listendev/pkg/observability/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ AnalysisRequest = (*NPM)(nil)
var _ Publisher = (*NPM)(nil)
var _ Deliverer = (*NPM)(nil)
var _ Results = (*NPM)(nil)

var (
	errNPMNameEmpty = errors.New("npm package name is empty")
)

type NPMFillError struct {
	Err error
}

func (e NPMFillError) Error() string {
	return e.Err.Error()
}

var (
	ErrMalfunctioningNPMRegistryClient = errors.New("malfunctioning (no-op or similar) NPM registry client")
	// NPMFillError instances
	ErrCouldNotRetrieveLastVersionTagFromNPM        = NPMFillError{errors.New("could not retrieve last npm version tag")}
	ErrCouldNotRetrieveLastVersionFromNPM           = NPMFillError{errors.New("could not retrieve last npm version")}
	ErrCouldNotRetrieveLastShasumFromNPM            = NPMFillError{errors.New("could not retrieve last npm version shasum")}
	ErrCouldNotRetrieveShasumForGivenVersionFromNPM = NPMFillError{errors.New("could not retrieve the shasum for the given npm version")}
	ErrGivenVersionNotFoundOnNPM                    = NPMFillError{errors.New("given npm package version not found on npm")}
	ErrGivenShasumDoesntMatchGivenVersionOnNPM      = NPMFillError{errors.New("given npm version does not exist on npm with the given shasum")}
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
	if tc.Ecosystem == ecosystem.Npm {
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

func (arn NPM) Delivery() (*amqp.Delivery, error) {
	body, err := json.Marshal(arn)
	if err != nil {
		return nil, err
	}

	ret := &amqp.Delivery{
		ContentType: "application/json",
		Body:        body,
	}
	if arn.Priority > 0 {
		ret.Priority = arn.Priority
	}

	return ret, nil
}

func (arn *NPM) fillMissingData(parent context.Context, registryClient npm.RegistryClient) error {
	// Assuming the context contains a tracer...
	ctx, span := tracer.FromContext(parent).Start(parent, "analysisrequest[npm].fillMissingData")
	defer span.End()

	if len(arn.Version) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, arn.Name)
		if err != nil {
			return err
		}
		if packageList == nil {
			return ErrMalfunctioningNPMRegistryClient
		}
		latestVersionTag := packageList.DistTags.Latest
		if len(latestVersionTag) == 0 {
			return ErrCouldNotRetrieveLastVersionTagFromNPM
		}
		if latestVersion, ok := packageList.Versions[latestVersionTag]; ok {
			arn.Version = latestVersion.Version
			arn.Shasum = latestVersion.Dist.Shasum
		}
		if len(arn.Version) == 0 {
			return ErrCouldNotRetrieveLastVersionFromNPM
		}
		if len(arn.Shasum) == 0 {
			return ErrCouldNotRetrieveLastShasumFromNPM
		}

		return nil
	}

	if len(arn.Version) > 0 && len(arn.Shasum) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, arn.Name)
		if err != nil {
			return err
		}
		if packageList == nil {
			return ErrMalfunctioningNPMRegistryClient
		}
		if version, ok := packageList.Versions[arn.Version]; ok {
			arn.Version = version.Version
			arn.Shasum = version.Dist.Shasum
		} else {
			return ErrGivenVersionNotFoundOnNPM
		}
		if len(arn.Shasum) == 0 {
			return ErrCouldNotRetrieveShasumForGivenVersionFromNPM
		}

		return nil
	}

	if len(arn.Version) > 0 && len(arn.Shasum) > 0 {
		packageVersion, err := registryClient.GetPackageVersion(ctx, arn.Name, arn.Version)
		if err != nil {
			return ErrGivenVersionNotFoundOnNPM
		}
		if packageVersion == nil {
			return ErrMalfunctioningNPMRegistryClient
		}
		if packageVersion.Dist.Shasum != arn.Shasum {
			return ErrGivenShasumDoesntMatchGivenVersionOnNPM
		}
	}

	return nil
}

func (arn NPM) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(&arn)
}

func (arn NPM) Switch(t Type) (AnalysisRequest, error) {
	c := t.Components()
	if !c.HasEcosystem() {
		return nil, fmt.Errorf("couldn't switch the current NPM analysis request to an analysis request with a type without ecosystem")
	}
	if c.Ecosystem != ecosystem.Npm {
		return nil, fmt.Errorf("couldn't switch the current NPM analysis request to a non NPM one")
	}
	arn.RequestType = t

	return &arn, nil
}
