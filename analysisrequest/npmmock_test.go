package analysisrequest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/listendev/pkg/npm"
)

type mockNpmregistryClient struct {
	listContent    []byte
	versionContent []byte
}

func newMockNpmregistryClient(listFilePath, versionFilePath string) (*mockNpmregistryClient, error) {
	plist, err := os.ReadFile(listFilePath)
	if err != nil {
		return nil, err
	}
	pversion, err := os.ReadFile(versionFilePath)
	if err != nil {
		return nil, err
	}

	return &mockNpmregistryClient{
		listContent:    plist,
		versionContent: pversion,
	}, nil
}

func (r *mockNpmregistryClient) GetPackageList(_ context.Context, name string) (*npm.PackageList, error) {
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

func (r *mockNpmregistryClient) GetPackageVersion(_ context.Context, name, _ string) (*npm.PackageVersion, error) {
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

func (r *mockNpmregistryClient) GetPackageLatestVersion(_ context.Context, name string) (*npm.PackageVersion, error) {
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
