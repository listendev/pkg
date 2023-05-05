package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/garnet-org/pkg/npm"
	"github.com/garnet-org/pkg/observability"
	"github.com/hgsgtk/jsoncmp"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func getFixture(filepath string) (string, error) {
	ret := path.Join("testdata", "fixtures", filepath)
	if _, err := os.ReadFile(ret); err != nil {
		return "", err
	}

	return ret, nil
}

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
	if packageList.Name != name {
		return nil, fmt.Errorf("GetPackageList: name mismatch")
	}

	return &packageList, nil
}

func (r *mockNpmregistryClient) GetPackageVersion(ctx context.Context, name, version string) (*npm.PackageVersion, error) {
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

func TestAnalysisRequestFromJSON(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name               string
		args               args
		want               AnalysisRequest
		wantPublishing     *amqp.Publishing
		wantS3Key          string
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
					Force:       false,
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    4,
				Body:        []byte(`{"type":"urn:nop:nop","snowflake_id":"1524854487523524609","priority":4,"force":false}`),
			},
			wantS3Key:          "nop/1524854487523524609/nop",
			mockRegistryClient: nil,
			wantErr:            false,
		},
		{
			name: "valid full npm deps dev analysis request",
			args: args{
				body: []byte(`{"type": "urn:hoarding:depsdev!npm.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5, "force": true}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMDepsDev,
					Snowflake:   "1524854487523524608",
					Priority:    5,
					Force:       true,
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.1.2",
					Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:hoarding:depsdev!npm.json","snowflake_id":"1524854487523524608","name":"chalk","version":"5.1.2","shasum":"d957f370038b75ac572471e83be4c5ca9f8e8c45","priority":5,"force":true}`),
			},
			wantS3Key: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/depsdev.json",
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
			name: "valid full npm falco install analysis request",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5, "force": true}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileFalco,
					Snowflake:   "1524854487523524608",
					Priority:    5,
					Force:       true,
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.1.2",
					Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:scheduler:falco!npm,install.json","snowflake_id":"1524854487523524608","name":"chalk","version":"5.1.2","shasum":"d957f370038b75ac572471e83be4c5ca9f8e8c45","priority":5,"force":true}`),
			},
			wantS3Key: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/falco[install].json",
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
			name: "npm falco test analysis request without shasum",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm,test.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2"}`),
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
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(`{"type": "urn:scheduler:falco!npm,test.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "force": false}`),
			},
			wantS3Key: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/falco[test].json",
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
			name: "npm falco install analysis request without version",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
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
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk","version": "5.2.0", "shasum": "249623b7d66869c673699fb66d65723e54dfcfb3", "force": false}`),
			},
			wantS3Key: "npm/chalk/5.2.0/249623b7d66869c673699fb66d65723e54dfcfb3/falco[install].json",
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
			name: "npm falco install analysis request enrichment with GPT without version",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMGPT4InstallWhileFalco,
					Snowflake:   "1524854487523524608",
				},
				npmPackage: npmPackage{
					Name:    "chalk",
					Version: "5.2.0",
					Shasum:  "249623b7d66869c673699fb66d65723e54dfcfb3",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(`{"type": "urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context", "snowflake_id": "1524854487523524608", "name": "chalk","version": "5.2.0", "shasum": "249623b7d66869c673699fb66d65723e54dfcfb3", "force": false}`),
			},
			wantS3Key: "npm/chalk/5.2.0/249623b7d66869c673699fb66d65723e54dfcfb3/falco[install].json",
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
			want:           nil,
			wantPublishing: nil,
			wantS3Key:      "",
			wantErr:        true,
		},
		{
			name: "valid npm falco install analysis request without version but no NPM registry set",
			args: args{
				body: []byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
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
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk","version": "5.2.0", "shasum": "249623b7d66869c673699fb66d65723e54dfcfb3", "force": false}`),
			},
			wantS3Key:          "npm/chalk/5.2.0/249623b7d66869c673699fb66d65723e54dfcfb3/falco[install].json",
			mockRegistryClient: nil,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder, err := NewBuilder(testCtx)
			assert.Nil(t, err)
			assert.NotNil(t, arbuilder)
			arbuilder.WithNPMRegistryClient(tt.mockRegistryClient)
			got, err := arbuilder.FromJSON(tt.args.body)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)

				// Test marshalling
				res, err := json.Marshal(got)
				assert.Nil(t, err)
				exp, err := json.Marshal(tt.want)
				assert.Nil(t, err)
				if diff := jsoncmp.Diff(string(res), string(exp)); diff != "" {
					t.Errorf("diff: (-got +want)\n%s", diff)
				}

				// Test amqp.Publishing
				pub, pubErr := got.Publishing()
				if assert.Nil(t, pubErr) {
					assert.Equal(t, tt.wantPublishing.Priority, pub.Priority)
					assert.Equal(t, tt.wantPublishing.ContentType, pub.ContentType)
					if diff := jsoncmp.Diff(string(tt.wantPublishing.Body), string(res)); diff != "" {
						t.Errorf("diff: (-got +want)\n%s", diff)
					}
				}

				// Test S3 key
				assert.Equal(t, tt.wantS3Key, got.ResultsPath().ToS3Key())
			}
		})
	}
}

func TestNewBuilderWithoutTracer(t *testing.T) {
	arbuilder, err := NewBuilder(context.TODO())
	assert.Nil(t, arbuilder)
	assert.NotNil(t, err)
}

func TestBuilderWithoutNPMRegistryClient(t *testing.T) {
	testCtx := observability.NewNopContext()
	arbuilder, err := NewBuilder(testCtx)
	assert.Nil(t, err)
	assert.NotNil(t, arbuilder)

	_, gotErr := arbuilder.FromJSON([]byte(`{"type": "urn:scheduler:falco!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`))

	if assert.Error(t, gotErr) {
		assert.True(t, errors.Is(gotErr, ErrMalfunctioningNPMRegistryClient))
		assert.Equal(t, ErrMalfunctioningNPMRegistryClient, gotErr)
	}
}

func TestAnalysisRequestFromFile(t *testing.T) {
	type args struct {
		fixture string
	}
	tests := []struct {
		name               string
		args               args
		want               []AnalysisRequest
		wantErr            bool
		mockRegistryClient *mockNpmregistryClient
	}{
		{
			name: "single",
			args: args{
				fixture: "single.json",
			},
			want: []AnalysisRequest{
				&NPM{
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
			},
			mockRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}
				return mockClient
			}(),
		},
		// TODO: list
		// {
		// 	name: "list",
		// 	args: args{
		// 		fixture: "list.json",
		// 	},
		// 	want: []AnalysisRequest{
		// 		// TODO: ...
		// 	},
		// 	mockRegistryClient: nil, // TODO: create mock registry client for all the packages in list.json
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder, err := NewBuilder(testCtx)
			assert.Nil(t, err)
			assert.NotNil(t, arbuilder)
			arbuilder.WithNPMRegistryClient(tt.mockRegistryClient)
			filepath, err := getFixture(tt.args.fixture)
			assert.Nil(t, err)

			got, err := arbuilder.FromFile(filepath)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
