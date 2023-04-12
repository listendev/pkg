package analysisrequest

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability"
	"github.com/stretchr/testify/assert"
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

func (r *mockNpmregistryClient) GetPackageList(ctx context.Context, name string) (*npm.PackageList, error) {
	var packageList npm.PackageList
	err := json.Unmarshal(r.listContent, &packageList)
	if err != nil {
		return nil, err
	}

	return &packageList, nil
}

func (r *mockNpmregistryClient) GetPackageVersion(ctx context.Context, name, version string) (*npm.PackageVersion, error) {
	var packageVersion npm.PackageVersion
	err := json.Unmarshal(r.versionContent, &packageVersion)
	if err != nil {
		return nil, err
	}

	return &packageVersion, nil
}

func TestAnalysisRequestFromJSON(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name               string
		args               args
		want               AnalysisRequest
		wantErr            bool
		mockRegistryClient *mockNpmregistryClient
	}{
		{
			name: "valid full nop analysis request",
			args: args{
				body: []byte(`{"type": "urn:nop:nop", "snowflake_id": "1524854487523524609", "priority": 4}`),
			},
			want: &NOP{
				base: base{
					RequestType: Nop,
					Snowflake:   "1524854487523524609",
					Priority:    4,
				},
			},
			mockRegistryClient: nil,
			wantErr:            false,
		},
		{
			name: "valid full npm analysis request",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm.install", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileFalco,
					Snowflake:   "1524854487523524608",
					Priority:    5,
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.1.2",
					Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
				},
			},
			mockRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}
				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "request without shasum",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm.test", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMTestWhileFalco,
					Snowflake:   "1524854487523524608",
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.1.2",
					Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
				},
			},
			mockRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}
				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "request without version",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm.install", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileFalco,
					Snowflake:   "1524854487523524608",
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.2.0",
					Shasum:  "249623b7d66869c673699fb66d65723e54dfcfb3",
				},
			},
			mockRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_520.json")
				if err != nil {
					t.Fatal(err)
				}
				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "invalid analysis request",
			args: args{
				body: []byte(`{"type": "something"}`),
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder := NewBuilderWithContext(testCtx)
			arbuilder.WithNPMRegistryClient(tt.mockRegistryClient)
			got, err := arbuilder.FromJSON(tt.args.body)

			if (err != nil) != tt.wantErr {
				t.Errorf("arbuilder.FromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
