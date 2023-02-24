package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/listendev/pkg/observability/tracer"
)

const defaultRegistryBaseURL = "https://registry.npmjs.org"

type RegistryClient interface {
	GetPackageList(ctx context.Context, name string) (*PackageList, error)
	GetPackageVersion(ctx context.Context, name, version string) (*PackageVersion, error)
}

type NPMRegistryClient struct {
	client  *http.Client
	baseURL *url.URL
}

type NPMRegistryClientConfig struct {
	Timeout time.Duration
	BaseURL string
}

func NewNPMRegistryClient(config NPMRegistryClientConfig) (RegistryClient, error) {
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

	return &NPMRegistryClient{
		client:  c,
		baseURL: url,
	}, nil
}

func (c *NPMRegistryClient) GetPackageList(parent context.Context, name string) (*PackageList, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "NPMRegistryClient.GetPackageList")
	defer span.End()
	endpoint := c.baseURL.ResolveReference(&url.URL{Path: name})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response not ok")
	}

	var packageList PackageList
	err = json.NewDecoder(response.Body).Decode(&packageList)
	if err != nil {
		return nil, err
	}

	return &packageList, nil
}

func (c *NPMRegistryClient) GetPackageVersion(parent context.Context, name, version string) (*PackageVersion, error) {
	ctx, span := tracer.FromContext(parent).Start(parent, "NPMRegistryClient.GetPackageVersion")
	defer span.End()
	endpoint := c.baseURL.ResolveReference(&url.URL{Path: path.Join(name, version)})
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response not ok")
	}

	var packageVersion PackageVersion
	err = json.NewDecoder(response.Body).Decode(&packageVersion)
	if err != nil {
		return nil, err
	}

	return &packageVersion, nil
}

type NoOpRegistryClient struct{}

func NewNoOpRegistryClient() RegistryClient {
	return &NoOpRegistryClient{}
}

func (c *NoOpRegistryClient) GetPackageList(ctx context.Context, name string) (*PackageList, error) {
	return nil, nil
}

func (c *NoOpRegistryClient) GetPackageVersion(ctx context.Context, name, version string) (*PackageVersion, error) {
	return nil, nil
}
