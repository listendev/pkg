package npm

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/listendev/pkg/observability/tracer"
)

const defaultRegistryBaseURL = "https://registry.npmjs.org"

var (
	ErrVersionNotFound        = errors.New("version not found")
	ErrServiceUnavail         = errors.New("service unavailable")
	ErrLatestVersionNotFound  = errors.New("latest version not found")
	ErrCouldNotDecodeResponse = errors.New("could not decode registry response")
	ErrCouldNotDoRequest      = errors.New("could not start request to the registry")
	ErrCouldNotCreateRequest  = errors.New("could not create request to the registry")
)

type Registry interface {
	GetPackageList(ctx context.Context, name string) (*PackageList, error)
	GetPackageVersion(ctx context.Context, name, version string) (*PackageVersion, error)
	GetPackageLatestVersion(ctx context.Context, name string) (*PackageVersion, error)
}

type RegistryClient struct {
	client  *http.Client
	baseURL *url.URL
}

type RegistryClientConfig struct {
	Timeout time.Duration
	BaseURL string
}

func NewRegistryClient(config RegistryClientConfig) (Registry, error) {
	timeout := time.Second * 10
	if config.Timeout != 0 {
		timeout = config.Timeout
	}
	c := &http.Client{Timeout: timeout}

	registryURL := defaultRegistryBaseURL
	if config.BaseURL != "" {
		registryURL = config.BaseURL
	}
	url, err := url.Parse(registryURL)
	if err != nil {
		return nil, err
	}

	return &RegistryClient{
		client:  c,
		baseURL: url,
	}, nil
}

func (c *RegistryClient) GetPackageList(parent context.Context, name string) (*PackageList, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "RegistryClient.GetPackageList")
	defer span.End()
	endpoint := c.baseURL.ResolveReference(&url.URL{Path: name})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, errors.Join(ErrCouldNotCreateRequest, err)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrVersionNotFound
		}
		return nil, ErrServiceUnavail
	}

	var packageList PackageList
	err = json.NewDecoder(response.Body).Decode(&packageList)
	if err != nil {
		return nil, ErrCouldNotDecodeResponse
	}

	return &packageList, nil
}

func (c *RegistryClient) GetPackageVersion(parent context.Context, name, version string) (*PackageVersion, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "RegistryClient.GetPackageVersion")
	defer span.End()
	endpoint := c.baseURL.ResolveReference(&url.URL{Path: path.Join(name, version)})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, errors.Join(ErrCouldNotCreateRequest, err)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrVersionNotFound
		}
		return nil, ErrServiceUnavail
	}
	var packageVersion PackageVersion
	err = json.NewDecoder(response.Body).Decode(&packageVersion)
	if err != nil {
		return nil, ErrCouldNotDecodeResponse
	}

	return &packageVersion, nil
}

func (c *RegistryClient) GetPackageLatestVersion(parent context.Context, name string) (*PackageVersion, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "RegistryClient.GetPackageLatestVersion")
	defer span.End()
	endpoint := c.baseURL.ResolveReference(&url.URL{Path: path.Join(name, "latest")})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, errors.Join(ErrCouldNotCreateRequest, err)
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrLatestVersionNotFound
		}
		return nil, ErrServiceUnavail
	}

	var packageVersion PackageVersion
	err = json.NewDecoder(response.Body).Decode(&packageVersion)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDecodeResponse, err)
	}

	return &packageVersion, nil
}

type NoOpRegistryClient struct{}

func NewNoOpRegistryClient() Registry {
	return &NoOpRegistryClient{}
}

func (c *NoOpRegistryClient) GetPackageList(_ context.Context, _ string) (*PackageList, error) {
	//nolint:nilnil // this is a mock
	return nil, nil
}

func (c *NoOpRegistryClient) GetPackageVersion(_ context.Context, _, _ string) (*PackageVersion, error) {
	//nolint:nilnil // this is a mock
	return nil, nil
}

func (c *NoOpRegistryClient) GetPackageLatestVersion(_ context.Context, _ string) (*PackageVersion, error) {
	//nolint:nilnil // this is a mock
	return nil, nil
}
