package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability/tracer"
	amqp "github.com/rabbitmq/amqp091-go"
)

type AnalysisRequestType = string

const (
	AnalysisRequestTypeNPM AnalysisRequestType = "npm"
)

var (
	// generic errors
	errInvalidAnalysisRequest        = errors.New("invalid analysis request")
	errAnalysisRequestSnowflakeEmpty = errors.New("analysis request snowflake_id is empty")

	// npm errors
	errNPMAnalysisRequestNameEmpty              = errors.New("npm analysis request package name is empty")
	errNPMCouldNotRetrieveLastVersionTagb       = errors.New("could not retrieve last version tag")
	errNPMCouldNotRetrieveLastVersion           = errors.New("could not retrieve last version")
	errNPMCouldNotRetrieveLastVersionShasum     = errors.New("could not retrieve last version shasum")
	errCouldNotretireveSpecificPackageVersion   = errors.New("could not retrieve specific package version")
	errNPMVersionDoesNotExistWithShasum         = errors.New("version does not exist on the NPM registry with the provided shasum")
	errNPMProvidedPackageVersionNotFound        = errors.New("provided npm package version not found")
	errNPMCouldNotRetrieveProvidedVersionShasum = errors.New("could not retrieve provided version shasum")
)

type AnalysisRequest interface {
	// String returns the string representation of the analysis request
	String() string
	// ResultUploaderPath returns the path to the result uploader
	ResultUploaderPath() ResultUploaderPath
	// Type returns the type of the analysis request
	Type() AnalysisRequestType
	// SnowflakeID returns the snowflake ID of the analysis request
	SnowflakeID() string
	// ToJSON returns the JSON representation of the analysis request
	ToJSON() ([]byte, error)
}

type ResultUploaderPath []string

func (r ResultUploaderPath) TOS3Key() string {
	return path.Join(r...)
}

type NPMAnalysisRequest struct {
	AnalysisRequestBase
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
	Shasum  string `json:"shasum,omitempty"`
}

func NewNPMAnalysisRequest(snowflake, name, version, shasum string) NPMAnalysisRequest {
	return NPMAnalysisRequest{
		AnalysisRequestBase: AnalysisRequestBase{
			RequestType: AnalysisRequestTypeNPM,
			Snowflake:   snowflake,
		},
		Name:    name,
		Version: version,
		Shasum:  shasum,
	}
}

type AnalysisRequestBase struct {
	RequestType string `json:"type"`
	Snowflake   string `json:"snowflake_id"`
	Priority    uint8  `json:"priority"`
}

type AnalysisRequestBuilder struct {
	npmRegistryRegistryClient npm.RegistryClient
}

func NewAnalysisRequestBuilder() *AnalysisRequestBuilder {
	return &AnalysisRequestBuilder{
		npmRegistryRegistryClient: npm.NewNoOpRegistryClient(),
	}
}

func (a *AnalysisRequestBuilder) WithNPMRegistryClient(npmRegistry npm.RegistryClient) {
	a.npmRegistryRegistryClient = npmRegistry
}

func (a *AnalysisRequestBuilder) FromJSON(parent context.Context, body []byte) (AnalysisRequest, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "AnalysisRequestBuilder.FromJSON")
	defer span.End()
	var err error

	var art AnalysisRequestBase

	err = json.Unmarshal(body, &art)
	if err != nil {
		return nil, err
	}

	if err := validateAnalysisRequestBase(art); err != nil {
		return nil, err
	}

	switch art.RequestType {
	case AnalysisRequestTypeNPM:
		ar, err := npmAnalysisRequestFromJSON(body)
		if err != nil {
			return nil, err
		}
		if err := validateNPManalysisRequest(ar); err != nil {
			return nil, err
		}

		ar, err = populateNPMAnalysisRequestWithMissingData(ctx, a.npmRegistryRegistryClient, ar)
		if err != nil {
			return nil, err
		}

		if len(ar.Version) > 0 && len(ar.Shasum) > 0 {
			packageVersion, err := a.npmRegistryRegistryClient.GetPackageVersion(ctx, ar.Name, ar.Version)
			if err != nil {
				return nil, fmt.Errorf("%v: %w", errCouldNotretireveSpecificPackageVersion, err)
			}
			if packageVersion.Dist.Shasum != ar.Shasum {
				return nil, errNPMVersionDoesNotExistWithShasum
			}
		}

		return ar, nil
	}
	return nil, errInvalidAnalysisRequest
}

func npmAnalysisRequestFromJSON(body []byte) (NPMAnalysisRequest, error) {
	var ar NPMAnalysisRequest
	err := json.Unmarshal(body, &ar)

	if err != nil {
		return ar, err
	}

	return ar, nil
}

func populateNPMAnalysisRequestWithMissingData(parent context.Context, registryClient npm.RegistryClient, ar NPMAnalysisRequest) (NPMAnalysisRequest, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "analysisrequest.populateNPMAnalysisRequestWithMissingData")
	defer span.End()
	if len(ar.Version) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, ar.Name)
		if err != nil {
			return ar, err
		}
		latestVersionTag := packageList.DistTags.Latest
		if len(latestVersionTag) == 0 {
			return ar, errNPMCouldNotRetrieveLastVersionTagb
		}
		if latestVersion, ok := packageList.Versions[latestVersionTag]; ok {
			ar.Version = latestVersion.Version
			ar.Shasum = latestVersion.Dist.Shasum
		}
		if len(ar.Version) == 0 {
			return ar, errNPMCouldNotRetrieveLastVersion
		}
		if len(ar.Shasum) == 0 {
			return ar, errNPMCouldNotRetrieveLastVersionShasum
		}
		return ar, nil
	}

	if len(ar.Version) > 0 && len(ar.Shasum) == 0 {
		packageList, err := registryClient.GetPackageList(ctx, ar.Name)
		if err != nil {
			return ar, err
		}
		if version, ok := packageList.Versions[ar.Version]; ok {
			ar.Version = version.Version
			ar.Shasum = version.Dist.Shasum
		} else {
			return ar, fmt.Errorf("%v: (version=%s)", errNPMProvidedPackageVersionNotFound, ar.Version)
		}
		if len(ar.Shasum) == 0 {
			return ar, fmt.Errorf("%v: (version=%s)", errNPMCouldNotRetrieveProvidedVersionShasum, ar.Version)
		}
		return ar, nil
	}
	return ar, nil
}

func validateAnalysisRequestBase(ar AnalysisRequestBase) error {
	if len(ar.Snowflake) == 0 {
		return errAnalysisRequestSnowflakeEmpty
	}
	return nil
}

func validateNPManalysisRequest(ar NPMAnalysisRequest) error {
	if len(ar.Name) == 0 {
		return errNPMAnalysisRequestNameEmpty
	}
	return nil
}

func (a NPMAnalysisRequest) SnowflakeID() string {
	return a.Snowflake
}

func (a NPMAnalysisRequest) String() string {
	return a.Name + "@" + a.Version
}

func (a NPMAnalysisRequest) ResultUploaderPath() ResultUploaderPath {
	return ResultUploaderPath{
		"npm",
		a.Name,
		a.Version,
		a.Shasum,
	}
}

func (a NPMAnalysisRequest) Type() AnalysisRequestType {
	return AnalysisRequestTypeNPM
}

func (a NPMAnalysisRequest) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func (a NPMAnalysisRequest) ToPublishing() (*amqp.Publishing, error) {
	body, err := a.ToJSON()
	if err != nil {
		return nil, err
	}

	ret := &amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}
	if a.Priority > 0 {
		ret.Priority = a.Priority
	}

	return ret, nil
}
