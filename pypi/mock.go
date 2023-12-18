package pypi

import (
	"context"
	"encoding/json"
	"os"
	"path"
)

var _ Registry = (*MockRegistryClient)(nil)

type MockRegistryClient struct {
	listContent    []byte
	versionContent []byte
}

func NewMockRegistryClient(listFilename, versionFilename string) (*MockRegistryClient, error) {
	prefix := path.Join("testdata", "pypi")
	plist, err := os.ReadFile(path.Join(prefix, listFilename))
	if err != nil {
		return nil, err
	}
	pversion, err := os.ReadFile(path.Join(prefix, versionFilename))
	if err != nil {
		return nil, err
	}

	return &MockRegistryClient{
		listContent:    plist,
		versionContent: pversion,
	}, nil
}

func (r *MockRegistryClient) GetPackageList(_ context.Context, name string) (*PackageList, error) {
	var packageList PackageList
	err := json.Unmarshal(r.listContent, &packageList)
	if err != nil {
		return nil, err
	}
	packageList.Fill()

	return &packageList, nil
}

func (r *MockRegistryClient) GetPackageVersion(_ context.Context, name, version string) (*PackageVersion, error) {
	var packageList PackageList
	err := json.Unmarshal(r.versionContent, &packageList)
	if err != nil {
		return nil, err
	}
	packageVersion, err := packageList.GetVersion(version)
	if err != nil {
		return nil, err
	}

	return packageVersion, nil
}

func (r *MockRegistryClient) GetPackageLatestVersion(_ context.Context, name string) (*PackageVersion, error) {
	var packageList PackageList
	err := json.Unmarshal(r.listContent, &packageList)
	if err != nil {
		return nil, err
	}
	packageVersion, err := packageList.GetVersion("latest")
	if err != nil {
		return nil, err
	}

	return packageVersion, nil
}
