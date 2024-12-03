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

var _ Registry = (*RegistryClient)(nil)

const (
	defaultRegistryBaseURL = "https://registry.npmjs.org"
	defaultUserAgent       = "listendev/pkg/npm"
)

var (
	ErrPackageNotFound        = errors.New("package not found")
	ErrVersionNotFound        = errors.New("version not found")
	ErrLatestVersionNotFound  = errors.New("latest version not found")
	ErrCouldNotDecodeResponse = errors.New("could not decode registry response")
	ErrCouldNotDoRequest      = errors.New("could not start request to the registry")
	ErrCouldNotCreateRequest  = errors.New("could not create request to the registry")
)

type ServiceError struct {
	StatusCode int
	Message    string
}

func (e *ServiceError) Error() string {
	return e.Message
}

type Registry interface {
	GetPackageList(ctx context.Context, name string) (*PackageList, error)
	GetPackageVersion(ctx context.Context, name, version string) (*PackageVersion, error)
	GetPackageLatestVersion(ctx context.Context, name string) (*PackageVersion, error)
}

type RegistryClient struct {
	client    *http.Client
	baseURL   *url.URL
	userAgent string
}

type RegistryClientConfig struct {
	Timeout   time.Duration
	BaseURL   string
	UserAgent string
}

func NewRegistryClient(config RegistryClientConfig) (Registry, error) {
	timeout := time.Second * 10
	if config.Timeout != 0 {
		timeout = config.Timeout
	}
	ua := defaultUserAgent
	if len(config.UserAgent) > 0 {
		ua = config.UserAgent
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
		client:    c,
		baseURL:   url,
		userAgent: ua,
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
	req.Header.Set("User-Agent", c.userAgent)

	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrPackageNotFound
		}

		return nil, &ServiceError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
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
	req.Header.Set("User-Agent", c.userAgent)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrVersionNotFound
		}

		return nil, &ServiceError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
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
	req.Header.Set("User-Agent", c.userAgent)
	response, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDoRequest, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusNotFound {
			return nil, ErrLatestVersionNotFound
		}

		return nil, &ServiceError{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}

	var packageVersion PackageVersion
	err = json.NewDecoder(response.Body).Decode(&packageVersion)
	if err != nil {
		return nil, errors.Join(ErrCouldNotDecodeResponse, err)
	}

	return &packageVersion, nil
}
