package analysisrequest

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/listendev/pkg/pypi"
)

var _ pypi.Registry = (*mockPyPiRegistryClient)(nil)

type mockPyPiRegistryClient struct {
	listContent    []byte
	versionContent []byte
}

func newMockPyPiRegistryClient(listFilename, versionFilename string) (*mockPyPiRegistryClient, error) {
	prefix := path.Join("testdata", "pypi")
	plist, err := os.ReadFile(path.Join(prefix, listFilename))
	if err != nil {
		return nil, err
	}
	pversion, err := os.ReadFile(path.Join(prefix, versionFilename))
	if err != nil {
		return nil, err
	}

	return &mockPyPiRegistryClient{
		listContent:    plist,
		versionContent: pversion,
	}, nil
}

func (r *mockPyPiRegistryClient) GetPackageList(_ context.Context, name string) (*pypi.PackageList, error) {
	var packageList pypi.PackageList
	err := json.Unmarshal(r.listContent, &packageList)
	if err != nil {
		return nil, err
	}
	packageList.Fill()

	return &packageList, nil
}

func (r *mockPyPiRegistryClient) GetPackageVersion(_ context.Context, name, version string) (*pypi.PackageVersion, error) {
	var packageList pypi.PackageList
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

func (r *mockPyPiRegistryClient) GetPackageLatestVersion(_ context.Context, name string) (*pypi.PackageVersion, error) {
	var packageList pypi.PackageList
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
