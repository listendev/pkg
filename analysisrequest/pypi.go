package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/listendev/pkg/ecosystem"
	"github.com/listendev/pkg/observability/tracer"
	"github.com/listendev/pkg/pypi"
	amqp "github.com/rabbitmq/amqp091-go"
)

var _ AnalysisRequest = (*PyPi)(nil)
var _ Publisher = (*PyPi)(nil)
var _ Deliverer = (*PyPi)(nil)
var _ Results = (*PyPi)(nil)

var (
	errPyPiNameEmpty = errors.New("PyPi package name is empty")
)

type PyPiFillError struct {
	Err error
}

func (e PyPiFillError) Error() string {
	return e.Err.Error()
}

var (
	ErrMalfunctioningPyPiRegistryClient = errors.New("malfunctioning (no-op or similar) PyPi registry client")
	// PyPiFillError instances.
	ErrGivenVersionNotFoundOnPyPi        = PyPiFillError{errors.New("given PyPi package version not found on PyPi")}
	ErrGivenSha256DoesNotMatchOnPyPi     = PyPiFillError{errors.New("given PyPi version does not exist on PyPi with the given sha256 digest")}
	ErrGivenBlake2b256DoesNotMatchOnPyPi = PyPiFillError{errors.New("given PyPi version does not exist on PyPi with the given blake2b256 digest")}
)

type pypiPackage struct {
	Name       string `json:"name"`
	Version    string `json:"version,omitempty"`
	Sha256     string `json:"sha256,omitempty"`
	Blake2b256 string `json:"blake2b_256,omitempty"`
}

type PyPi struct {
	base
	pypiPackage
}

// NewPyPi creates an AnalysisRequest for the PyPi ecosystem.
func NewPyPi(request Type, snowflake string, priority uint8, force bool, name, version, digest string) (AnalysisRequest, error) {
	tc := request.Components()
	if !tc.HasEcosystem() {
		return nil, fmt.Errorf("couldn't instantiate an analysis request for PyPi from a type without ecosystem at all")
	}
	if tc.Ecosystem == ecosystem.Pypi {
		return &PyPi{
			base: base{
				RequestType: request,
				Snowflake:   snowflake,
				Priority:    priority,
				Force:       force,
			},
			pypiPackage: pypiPackage{
				Name:       name,
				Version:    version,
				Blake2b256: digest,
			},
		}, nil
	}

	return nil, fmt.Errorf("couldn't instantiate an analysis request for PyPi")
}

func (arp PyPi) PackageName() string {
	return arp.Name
}

func (arp PyPi) PackageVersion() string {
	return arp.Version
}

func (arp PyPi) PackageDigest() string {
	return arp.Blake2b256
}

func (arp PyPi) Publishing() (*amqp.Publishing, error) {
	return ComposeAMQPPublishing(&arp)
}

func (arp PyPi) ResultsPath() ResultUploadPath {
	return ComposeResultUploadPath(&arp)
}

func (arp PyPi) String() string {
	return arp.Name + "@" + arp.Version + "(" + arp.Type().String() + ")"
}

func (arp PyPi) Delivery() (*amqp.Delivery, error) {
	return ComposeAMQPDelivery(&arp)
}

func (arp PyPi) Validate() error {
	if len(arp.Name) == 0 {
		return errPyPiNameEmpty
	}

	return arp.base.Validate()
}

func (arp *PyPi) UnmarshalJSON(data []byte) error {
	var baseResult base
	if err := json.Unmarshal(data, &baseResult); err != nil {
		return err
	}
	arp.base = baseResult

	var pypiResult pypiPackage
	if err := json.Unmarshal(data, &pypiResult); err != nil {
		return err
	}
	arp.pypiPackage = pypiResult

	return arp.Validate()
}

func (arp *PyPi) fillMissingData(parent context.Context, client pypi.Registry) error {
	// Assuming the context contains a tracer...
	ctx, span := tracer.FromContext(parent).Start(parent, "analysisrequest[pypi].fillMissingData")
	defer span.End()

	if len(arp.Version) == 0 {
		pv, err := client.GetPackageLatestVersion(ctx, arp.Name)
		if err != nil {
			return err
		}
		if pv == nil {
			return ErrMalfunctioningPyPiRegistryClient
		}
		arp.Version = pv.Version
		arp.Sha256 = pv.Digests.SHA256
		arp.Blake2b256 = pv.Digests.Blake2bB256

		return nil
	}

	if len(arp.Blake2b256) == 0 || len(arp.Sha256) == 0 {
		pv, err := client.GetPackageVersion(ctx, arp.Name, arp.Version)
		if err != nil {
			if errors.Is(err, pypi.ErrVersionNotFound) {
				return ErrGivenVersionNotFoundOnPyPi
			}
			// all the other errors are considered as service unavailable or client errors

			return errors.Join(ErrMalfunctioningPyPiRegistryClient, err)
		}
		if pv == nil {
			return ErrMalfunctioningPyPiRegistryClient
		}
		if len(arp.Sha256) > 0 && pv.Digests.SHA256 != arp.Sha256 {
			return ErrGivenSha256DoesNotMatchOnPyPi
		}
		if len(arp.Blake2b256) > 0 && pv.Digests.Blake2bB256 != arp.Blake2b256 {
			return ErrGivenBlake2b256DoesNotMatchOnPyPi
		}
		arp.Sha256 = pv.Digests.SHA256
		arp.Blake2b256 = pv.Digests.Blake2bB256
	}

	return nil
}
