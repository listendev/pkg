package npm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var _ Registry = (*MockRegistryClient)(nil)

type MockRegistryClient struct {
	listContent    []byte
	versionContent []byte
}

func NewMockRegistryClient(listFilename, versionFilename string) (*MockRegistryClient, error) {
	prefix := path.Join("testdata", "npm")
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
	if packageList.Name != name {
		return nil, fmt.Errorf("GetPackageList: name mismatch")
	}

	return &packageList, nil
}

func (r *MockRegistryClient) GetPackageVersion(_ context.Context, name, _ string) (*PackageVersion, error) {
	var packageVersion PackageVersion
	err := json.Unmarshal(r.versionContent, &packageVersion)
	if err != nil {
		return nil, err
	}
	if packageVersion.Name != name {
		return nil, fmt.Errorf("GetPackageVersion: name mismatch")
	}

	return &packageVersion, nil
}

func (r *MockRegistryClient) GetPackageLatestVersion(_ context.Context, name string) (*PackageVersion, error) {
	var packageVersion PackageVersion
	err := json.Unmarshal(r.versionContent, &packageVersion)
	if err != nil {
		return nil, err
	}
	if packageVersion.Name != name {
		return nil, fmt.Errorf("GetPackageVersion: name mismatch")
	}

	return &packageVersion, nil
}
