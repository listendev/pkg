package pypi

import "context"

var _ Registry = (*NoOpRegistryClient)(nil)

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
