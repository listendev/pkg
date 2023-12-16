package analysisrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/listendev/pkg/npm"
)

var _ npm.Registry = (*mockNpmRegistryClient)(nil)

type mockNpmRegistryClient struct {
	listContent    []byte
	versionContent []byte
}

func newMockNpmRegistryClient(listFilename, versionFilename string) (*mockNpmRegistryClient, error) {
	prefix := path.Join("testdata", "npm")
	plist, err := os.ReadFile(path.Join(prefix, listFilename))
	if err != nil {
		return nil, err
	}
	pversion, err := os.ReadFile(path.Join(prefix, versionFilename))
	if err != nil {
		return nil, err
	}

	return &mockNpmRegistryClient{
		listContent:    plist,
		versionContent: pversion,
	}, nil
}

func (r *mockNpmRegistryClient) GetPackageList(_ context.Context, name string) (*npm.PackageList, error) {
	var packageList npm.PackageList
	err := json.Unmarshal(r.listContent, &packageList)
	if err != nil {
		return nil, err
	}
	if packageList.Name != name {
		return nil, fmt.Errorf("GetPackageList: name mismatch")
	}

	return &packageList, nil
}

func (r *mockNpmRegistryClient) GetPackageVersion(_ context.Context, name, _ string) (*npm.PackageVersion, error) {
	var packageVersion npm.PackageVersion
	err := json.Unmarshal(r.versionContent, &packageVersion)
	if err != nil {
		return nil, err
	}
	if packageVersion.Name != name {
		return nil, fmt.Errorf("GetPackageVersion: name mismatch")
	}

	return &packageVersion, nil
}

func (r *mockNpmRegistryClient) GetPackageLatestVersion(_ context.Context, name string) (*npm.PackageVersion, error) {
	var packageVersion npm.PackageVersion
	err := json.Unmarshal(r.versionContent, &packageVersion)
	if err != nil {
		return nil, err
	}
	if packageVersion.Name != name {
		return nil, fmt.Errorf("GetPackageVersion: name mismatch")
	}

	return &packageVersion, nil
}
