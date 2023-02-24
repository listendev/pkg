package analysisrequest

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability"
	"github.com/google/go-cmp/cmp"
)

type mockNpmregistryClient struct {
	listContent    []byte
	versionContent []byte
}

func newMockNpmregistryClient(listFilePath, versionFilePath string) (*mockNpmregistryClient, error) {
	plist, err := ioutil.ReadFile(listFilePath)
	if err != nil {
		return nil, err
	}
	pversion, err := ioutil.ReadFile(versionFilePath)
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
			name: "valid full npm analysis request",
			args: args{
				body: []byte(`{"type": "npm", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45"}`),
			},
			want: NPMAnalysisRequest{
				AnalysisRequestBase: AnalysisRequestBase{
					RequestType: AnalysisRequestTypeNPM,
					Snowflake:   "1524854487523524608",
				},
				Name:    "chalk",
				Version: "5.1.2",
				Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
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
				body: []byte(`{"type": "npm", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2"}`),
			},
			want: NPMAnalysisRequest{
				AnalysisRequestBase: AnalysisRequestBase{
					RequestType: AnalysisRequestTypeNPM,
					Snowflake:   "1524854487523524608",
				},
				Name:    "chalk",
				Version: "5.1.2",
				Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
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
				body: []byte(`{"type": "npm", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			want: NPMAnalysisRequest{
				AnalysisRequestBase: AnalysisRequestBase{
					RequestType: AnalysisRequestTypeNPM,
					Snowflake:   "1524854487523524608",
				},
				Name:    "chalk",
				Version: "5.2.0",
				Shasum:  "249623b7d66869c673699fb66d65723e54dfcfb3",
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
			arbuilder := NewAnalysisRequestBuilder()
			arbuilder.WithNPMRegistryClient(tt.mockRegistryClient)
			got, err := arbuilder.FromJSON(testCtx, tt.args.body)

			if (err != nil) != tt.wantErr {
				t.Errorf(" arbuilder.FromJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !cmp.Equal(got, tt.want) {
				t.Errorf(" arbuilder.FromJSON(): %s", cmp.Diff(got, tt.want))

			}
		})
	}
}
