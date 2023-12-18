package analysisrequest

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"testing"

	"github.com/hgsgtk/jsoncmp"
	"github.com/listendev/pkg/npm"
	"github.com/listendev/pkg/observability"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		name                   string
		args                   args
		want                   AnalysisRequest
		wantPublishing         *amqp.Publishing
		wantKey                string
		wantErr                bool
		mockNPMRegistryClient  *npm.MockRegistryClient
		mockPyPiRegistryClient *mockPyPiRegistryClient
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
		{
			name: "valid full pypi typosquat analysis request",
			args: args{
				body: []byte(`{"type": "urn:hoarding:typosquat!pypi.json", "snowflake_id": "1652803364692340737", "name": "cctx", "version": "1.0.0", "sha256": "1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd", "blake2b_256": "bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e", "priority": 5, "force": true}`),
			},
			want: &PyPi{
				base: base{
					RequestType: PypiTyposquat,
					Snowflake:   "1652803364692340737",
					Priority:    5,
					Force:       true,
				},
				pypiPackage: pypiPackage{
					Name:       "cctx",
					Version:    "1.0.0",
					Sha256:     "1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd",
					Blake2b256: "bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:hoarding:typosquat!pypi.json","snowflake_id":"1652803364692340737","name":"cctx","version":"1.0.0","priority":5,"force":true,"sha256":"1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd","blake2b_256":"bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e"}`),
			},
			wantKey: "pypi/cctx/1.0.0/bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e/typosquat.json",
			wantErr: false,
		},
		{
			name: "pypi typosquat analysis request without blake2b_256 digest",
			args: args{
				body: []byte(`{"type": "urn:hoarding:typosquat!pypi.json", "snowflake_id": "1652803364692340737", "name": "cctx", "version": "1.0.0", "sha256": "1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd", "priority": 5, "force": true}`),
			},
			want: &PyPi{
				base: base{
					RequestType: PypiTyposquat,
					Snowflake:   "1652803364692340737",
					Priority:    5,
					Force:       true,
				},
				pypiPackage: pypiPackage{
					Name:       "cctx",
					Version:    "1.0.0",
					Sha256:     "1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd",
					Blake2b256: "bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:hoarding:typosquat!pypi.json","snowflake_id":"1652803364692340737","name":"cctx","version":"1.0.0","priority":5,"force":true,"sha256":"1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd","blake2b_256":"bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e"}`),
			},
			wantKey: "pypi/cctx/1.0.0/bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e/typosquat.json",
			wantErr: false,
			mockPyPiRegistryClient: func() *mockPyPiRegistryClient {
				mockClient, err := newMockPyPiRegistryClient("cctx.json", "cctx_100.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
		},
		{
			name: "pypi typosquat analysis request without digests",
			args: args{
				body: []byte(`{"type": "urn:hoarding:typosquat!pypi.json", "snowflake_id": "1652803364692340737", "name": "boto3", "version": "1.33.8", "priority": 5, "force": true}`),
			},
			want: &PyPi{
				base: base{
					RequestType: PypiTyposquat,
					Snowflake:   "1652803364692340737",
					Priority:    5,
					Force:       true,
				},
				pypiPackage: pypiPackage{
					Name:       "boto3",
					Version:    "1.33.8",
					Sha256:     "d02a084b25aa8d46ef917b128e90877efab1ba45f9d1ba3a11f336930378e350",
					Blake2b256: "121f1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:hoarding:typosquat!pypi.json","snowflake_id":"1652803364692340737","name":"boto3","version":"1.33.8","priority":5,"force":true,"sha256":"d02a084b25aa8d46ef917b128e90877efab1ba45f9d1ba3a11f336930378e350","blake2b_256":"121f1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a"}`),
			},
			wantKey: "pypi/boto3/1.33.8/121f1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a/typosquat.json",
			wantErr: false,
			mockPyPiRegistryClient: func() *mockPyPiRegistryClient {
				mockClient, err := newMockPyPiRegistryClient("boto3.json", "boto3_1338.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
		},
		{
			name: "pypi typosquat analysis request with package name only",
			args: args{
				body: []byte(`{"type": "urn:hoarding:typosquat!pypi.json", "snowflake_id": "1652803364692340737", "name": "boto3", "priority": 5, "force": true}`),
			},
			want: &PyPi{
				base: base{
					RequestType: PypiTyposquat,
					Snowflake:   "1652803364692340737",
					Priority:    5,
					Force:       true,
				},
				pypiPackage: pypiPackage{
					Name:       "boto3",
					Version:    "1.34.2",
					Sha256:     "970fd9f9f522eb48f3cd5574e927b369279ebf5bcf0f2fae5ed9cc6306e58558",
					Blake2b256: "c86666f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f",
				},
			},
			wantPublishing: &amqp.Publishing{
				ContentType: "application/json",
				Priority:    5,
				Body:        []byte(`{"type":"urn:hoarding:typosquat!pypi.json","snowflake_id":"1652803364692340737","name":"boto3","version":"1.34.2","priority":5,"force":true,"sha256":"970fd9f9f522eb48f3cd5574e927b369279ebf5bcf0f2fae5ed9cc6306e58558","blake2b_256":"c86666f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f"}`),
			},
			wantKey: "pypi/boto3/1.34.2/c86666f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f/typosquat.json",
			wantErr: false,
			mockPyPiRegistryClient: func() *mockPyPiRegistryClient {
				mockClient, err := newMockPyPiRegistryClient("boto3.json", "boto3_1338.json")
				if err != nil {
					t.Fatal(err)
				}

				return mockClient
			}(),
		},
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_520.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_520.json")
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
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCtx := observability.NewNopContext()
			arbuilder, err := NewBuilder(testCtx)
			assert.Nil(t, err)
			assert.NotNil(t, arbuilder)
			arbuilder.WithNPMRegistryClient(tt.mockNPMRegistryClient)
			arbuilder.WithPyPiRegistryClient(tt.mockPyPiRegistryClient)
			got, err := arbuilder.FromJSON(tt.args.body)

			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, got)
				assert.Nil(t, tt.want)
				assert.Nil(t, tt.wantPublishing)
				assert.Empty(t, tt.wantKey)
			} else {
				require.Nil(t, err)
				require.NotNil(t, got)
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
		mockNPMRegistryClient *npm.MockRegistryClient
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
			mockNPMRegistryClient: func() *npm.MockRegistryClient {
				mockClient, err := npm.NewMockRegistryClient("chalk.json", "chalk_512.json")
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
