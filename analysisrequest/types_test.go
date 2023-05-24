package analysisrequest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParent(t *testing.T) {
	type testCase struct {
		input   Type
		want    Type
		wantErr bool
	}

	cases := []testCase{
		{
			input:   NPMGPT4InstallWhileFalco,
			want:    NPMInstallWhileFalco,
			wantErr: false,
		},
	}

	for _, tc := range cases {
		got, err := tc.input.Parent()
		if tc.wantErr {
			assert.Error(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		}
	}
}

func TestEnricherResultFileIsTheParentOne(t *testing.T) {
	got := NPMGPT4InstallWhileFalco.Components().ResultFile()
	assert.Equal(t, NPMInstallWhileFalco.Components().ResultFile(), got)
}

func TestEnrichersEquality(t *testing.T) {
	tt, err := ToType("urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context")
	assert.Nil(t, err)
	assert.Equal(t, NPMGPT4InstallWhileFalco, tt)
}

func TestTypes(t *testing.T) {
	type want struct {
		urn  string
		json []byte
		TypeComponents
	}
	type testCase struct {
		input Type
		want  want
	}

	cases := []testCase{
		{
			input: NPMGPT4InstallWhileFalco,
			want: want{
				urn:  "urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context",
				json: []byte(`"urn:scheduler:falco!npm,install.json+urn:hoarding:gpt4,context"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        GPT4Collector,
					CollectorActions: []string{"context"},
					Ecosystem:        NPMEcosystem,        // From parent
					EcosystemActions: []string{"install"}, // From parent
					Format:           "json",              // From parent
					Parent: &TypeComponents{
						Framework:        Scheduler,
						Collector:        FalcoCollector,
						CollectorActions: []string{},
						Ecosystem:        NPMEcosystem,
						EcosystemActions: []string{"install"},
						Format:           "json",
					},
				},
			},
		},
		{
			input: Nop,
			want: want{
				urn:  "urn:nop:nop",
				json: []byte(`"urn:nop:nop"`),
				TypeComponents: TypeComponents{
					Framework:        None,
					Collector:        NoCollector,
					CollectorActions: []string{},
					EcosystemActions: []string{},
					Format:           "",
				},
			},
		},
		{
			input: NPMInstallWhileFalco,
			want: want{
				urn:  "urn:scheduler:falco!npm,install.json",
				json: []byte(`"urn:scheduler:falco!npm,install.json"`),
				TypeComponents: TypeComponents{
					Framework:        Scheduler,
					Collector:        FalcoCollector,
					CollectorActions: []string{},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{"install"},
					Format:           "json",
				},
			},
		},
		// {
		// 	input: NPMTestWhileFalco,
		// 	want: want{
		// 		urn:  "urn:scheduler:falco!npm,test.json",
		// 		json: []byte(`"urn:scheduler:falco!npm,test.json"`),
		// 		TypeComponents: TypeComponents{
		// 			Framework:        Scheduler,
		// 			Collector:        "falco",
		// 			CollectorActions: []string{},
		// 			Ecosystem:        NPMEcosystem,
		// 			EcosystemActions: []string{"test"},
		// 			Format:           "json",
		// 		},
		// 	},
		// },
		{
			input: NPMDepsDev,
			want: want{
				urn:  "urn:hoarding:depsdev!npm.json",
				json: []byte(`"urn:hoarding:depsdev!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        DepsDevCollector,
					CollectorActions: []string{},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMTyposquat,
			want: want{
				urn:  "urn:hoarding:typosquat!npm.json",
				json: []byte(`"urn:hoarding:typosquat!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        TyposquatCollector,
					CollectorActions: []string{},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMMetadataEmptyDescription,
			want: want{
				urn:  "urn:hoarding:metadata,empty_descr!npm.json",
				json: []byte(`"urn:hoarding:metadata,empty_descr!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        MetadataCollector,
					CollectorActions: []string{"empty_descr"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMMetadataMaintainersEmailCheck,
			want: want{
				urn:  "urn:hoarding:metadata,email_check!npm.json",
				json: []byte(`"urn:hoarding:metadata,email_check!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        MetadataCollector,
					CollectorActions: []string{"email_check"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMMetadataZeroVersion,
			want: want{
				urn:  "urn:hoarding:metadata,zero_version!npm.json",
				json: []byte(`"urn:hoarding:metadata,zero_version!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        MetadataCollector,
					CollectorActions: []string{"zero_version"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMSemgrepEnvExfiltration,
			want: want{
				urn:  "urn:hoarding:semgrep,exfiltrate_env!npm.json",
				json: []byte(`"urn:hoarding:semgrep,exfiltrate_env!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        SemgrepCollector,
					CollectorActions: []string{"exfiltrate_env"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMSemgrepProcessExecution,
			want: want{
				urn:  "urn:hoarding:semgrep,process_exec!npm.json",
				json: []byte(`"urn:hoarding:semgrep,process_exec!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        SemgrepCollector,
					CollectorActions: []string{"process_exec"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMSemgrepEvalBase64,
			want: want{
				urn:  "urn:hoarding:semgrep,base64_eval!npm.json",
				json: []byte(`"urn:hoarding:semgrep,base64_eval!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        SemgrepCollector,
					CollectorActions: []string{"base64_eval"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
		{
			input: NPMSemgrepShadyLinks,
			want: want{
				urn:  "urn:hoarding:semgrep,shady_links!npm.json",
				json: []byte(`"urn:hoarding:semgrep,shady_links!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:        Hoarding,
					Collector:        SemgrepCollector,
					CollectorActions: []string{"shady_links"},
					Ecosystem:        NPMEcosystem,
					EcosystemActions: []string{},
					Format:           "json",
				},
			},
		},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.want.urn, tc.input.String())
		t.Run(tc.want.urn, func(t *testing.T) {
			assert.Equal(t, tc.want.TypeComponents, tc.input.Components())
			got, err := json.Marshal(tc.input)
			assert.Nil(t, err)
			assert.Equal(t, tc.want.json, got)

			var t1 Type
			assert.Nil(t, json.Unmarshal(tc.want.json, &t1))
			assert.Equal(t, tc.input, t1)

			var t2 Type
			assert.Nil(t, json.Unmarshal([]byte(fmt.Sprintf("%q", tc.want.urn)), &t2))
			assert.Equal(t, tc.input, t2)

			urnObj := tc.want.TypeComponents.ToURN()
			assert.NotNil(t, urnObj)
			typeObj, err := ToType(urnObj.String())
			assert.Nil(t, err)
			assert.Equal(t, tc.input, typeObj)
		})
	}
}

func TestLastType(t *testing.T) {
	got := LastType()

	assert.Equal(t, NPMSemgrepEvalBase64, got)
}
