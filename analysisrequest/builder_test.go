package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/hgsgtk/jsoncmp"
	"github.com/listendev/pkg/observability"
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

// TODO: test also output AMQP delivery message
func TestAnalysisRequestFromJSON(t *testing.T) {
	type args struct {
		body []byte
	}
	tests := []struct {
		name                  string
		args                  args
		want                  AnalysisRequest
		wantPublishing        *amqp.Publishing
		wantKey               string
		wantErr               bool
		mockNPMRegistryClient *mockNpmregistryClient
		// mockPyPiRegistryClient *mockPyPiRegistryClient // TODO: impl
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
			wantKey:               "nop/1524854487523524609/nop",
			mockNPMRegistryClient: nil,
			wantErr:               false,
		},

		// FIXME: mock pypi registry response
		// {
		// 	name: "valid full pypi typosquat analysis request",
		// 	args: args{
		// 		body: []byte(`{"type": "urn:hoarding:typosquat!pypi.json", "snowflake_id": "1652803364692340737", "name": "cctx", "version": "1.0.0", "priority": 5, "force": true}`),
		// 	},
		// 	want: &PyPi{
		// 		base: base{
		// 			RequestType: PypiTyposquat,
		// 			Snowflake:   "1652803364692340737",
		// 			Priority:    5,
		// 			Force:       true,
		// 		},
		// 		pypiPackage: pypiPackage{
		// 			Name:    "cctx",
		// 			Version: "1.0.0",
		// 		},
		// 	},
		// 	wantPublishing: &amqp.Publishing{
		// 		ContentType: "application/json",
		// 		Priority:    5,
		// 		Body:        []byte(`{"type":"urn:hoarding:typosquat!pypi.json","snowflake_id":"1652803364692340736","name":"cctx","version":"1.0.0","priority":5,"force":true}`),
		// 	},
		// 	wantKey: "pypi/cctx/1.0.0/typosquat.json",
		// 	// mockNPMRegistryClient: , // TODO: implement
		// 	wantErr: false,
		// },

		{
			name: "valid full npm advisory analysis request",
			args: args{
				body: []byte(`{"type": "urn:hoarding:advisory!npm.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5, "force": true}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMAdvisory,
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
				Body:        []byte(`{"type":"urn:hoarding:advisory!npm.json","snowflake_id":"1524854487523524608","name":"chalk","version":"5.1.2","shasum":"d957f370038b75ac572471e83be4c5ca9f8e8c45","priority":5,"force":true}`),
			},
			wantKey: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/advisory.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "valid full npm typosquat analysis request",
			args: args{
				body: []byte(`{"type": "urn:hoarding:typosquat!npm.json", "snowflake_id": "1652803364692340736", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5, "force": true}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMTyposquat,
					Snowflake:   "1652803364692340736",
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
				Body:        []byte(`{"type":"urn:hoarding:typosquat!npm.json","snowflake_id":"1652803364692340736","name":"chalk","version":"5.1.2","shasum":"d957f370038b75ac572471e83be4c5ca9f8e8c45","priority":5,"force":true}`),
			},
			wantKey: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/typosquat.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "valid full npm install analysis request",
			args: args{
				body: []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "priority": 5, "force": true}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileDynamicInstrumentation,
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
				Body:        []byte(`{"type":"urn:scheduler:dynamic!npm,install.json","snowflake_id":"1524854487523524608","name":"chalk","version":"5.1.2","shasum":"d957f370038b75ac572471e83be4c5ca9f8e8c45","priority":5,"force":true}`),
			},
			wantKey: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/dynamic!install!.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "npm install dynamic instrumentation analysis request without shasum",
			args: args{
				body: []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileDynamicInstrumentation,
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
				Body:        []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk", "version": "5.1.2", "shasum": "d957f370038b75ac572471e83be4c5ca9f8e8c45", "force": false}`),
			},
			wantKey: "npm/chalk/5.1.2/d957f370038b75ac572471e83be4c5ca9f8e8c45/dynamic!install!.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "npm install dynamic instrumentation analysis request without version",
			args: args{
				body: []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileDynamicInstrumentation,
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
				Body:        []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk","version": "5.2.0", "shasum": "249623b7d66869c673699fb66d65723e54dfcfb3", "force": false}`),
			},
			wantKey: "npm/chalk/5.2.0/249623b7d66869c673699fb66d65723e54dfcfb3/dynamic!install!.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_520.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
			wantErr: false,
		},
		{
			name: "npm install dynamic instrumentation analysis request enrichment with AI without version",
			args: args{
				body: []byte(`{"type": "urn:scheduler:dynamic!npm,install.json+urn:hoarding:ai,context", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			want: &NPM{
				base: base{
					RequestType: NPMInstallWhileDynamicInstrumentationAIEnriched,
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
				Body:        []byte(`{"type": "urn:scheduler:dynamic!npm,install.json+urn:hoarding:ai,context", "snowflake_id": "1524854487523524608", "name": "chalk","version": "5.2.0", "shasum": "249623b7d66869c673699fb66d65723e54dfcfb3", "force": false}`),
			},
			wantKey: "npm/chalk/5.2.0/249623b7d66869c673699fb66d65723e54dfcfb3/dynamic!install!.json",
			mockNPMRegistryClient: func() *mockNpmregistryClient {
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
			wantKey:        "",
			wantErr:        true,
		},
		{
			name: "valid npm install dynamic instrumentation analysis request without version but no NPM registry set",
			args: args{
				body: []byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`),
			},
			wantErr:               true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder, err := NewBuilder(testCtx)
			assert.Nil(t, err)
			assert.NotNil(t, arbuilder)
			arbuilder.WithNPMRegistryClient(tt.mockNPMRegistryClient)
			got, err := arbuilder.FromJSON(tt.args.body)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
				assert.Nil(t, tt.want)
				assert.Nil(t, tt.wantPublishing)
				assert.Empty(t, tt.wantKey)
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
				assert.Equal(t, tt.wantKey, got.ResultsPath().Key())
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

	_, gotErr := arbuilder.FromJSON([]byte(`{"type": "urn:scheduler:dynamic!npm,install.json", "snowflake_id": "1524854487523524608", "name": "chalk"}`))

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
		name                  string
		args                  args
		want                  []AnalysisRequest
		wantErr               bool
		mockNPMRegistryClient *mockNpmregistryClient
	}{
		{
			name: "single",
			args: args{
				fixture: "single.json",
			},
			want: []AnalysisRequest{
				&NPM{
					base: base{
						RequestType: NPMInstallWhileDynamicInstrumentation,
						Snowflake:   "1524854487523524608",
					},
					npmPackage: npmPackage{
						Name:    "chalk",
						Version: "5.1.2",
						Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
					},
				},
			},
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
		},
		{
			name: "list",
			args: args{
				fixture: "list.json",
			},
			want: []AnalysisRequest{
				&NPM{
					base: base{
						RequestType: NPMInstallWhileDynamicInstrumentation,
						Snowflake:   "1524854487523524608",
					},
					npmPackage: npmPackage{
						Name:    "chalk",
						Version: "5.1.2",
						Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
					},
				},
				&NPM{
					base: base{
						RequestType: NPMInstallWhileDynamicInstrumentation,
						Snowflake:   "1524854487523524608",
					},
					npmPackage: npmPackage{
						Name:    "chalk",
						Version: "5.1.2",
						Shasum:  "d957f370038b75ac572471e83be4c5ca9f8e8c45",
					},
				},
			},
			mockNPMRegistryClient: func() *mockNpmregistryClient {
				mockClient, err := newMockNpmregistryClient("testdata/chalk.json", "testdata/chalk_512.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder, err := NewBuilder(testCtx)
			assert.Nil(t, err)
			assert.NotNil(t, arbuilder)
			arbuilder.WithNPMRegistryClient(tt.mockNPMRegistryClient)
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
