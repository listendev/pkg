package analysisrequest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/listendev/pkg/ecosystem"
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
			input:   NPMGPT4InstallWhileDynamicInstrumentation,
			want:    NPMInstallWhileDynamicInstrumentation,
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
	got := NPMGPT4InstallWhileDynamicInstrumentation.Components().ResultFile()
	assert.Equal(t, NPMInstallWhileDynamicInstrumentation.Components().ResultFile(), got)
}

func TestEnrichersEquality(t *testing.T) {
	tt, err := ToType("urn:scheduler:dynamic!npm,install.json+urn:hoarding:gpt4,context")
	assert.Nil(t, err)
	assert.Equal(t, NPMGPT4InstallWhileDynamicInstrumentation, tt)
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
			input: NPMGPT4InstallWhileDynamicInstrumentation,
			want: want{
				urn:  "urn:scheduler:dynamic!npm,install.json+urn:hoarding:gpt4,context",
				json: []byte(`"urn:scheduler:dynamic!npm,install.json+urn:hoarding:gpt4,context"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       GPT4Collector,
					CollectorAction: "context",
					Ecosystem:       ecosystem.Npm, // From parent
					EcosystemAction: "install",     // From parent
					Format:          "json",        // From parent
					Parent: &TypeComponents{
						Framework:       Scheduler,
						Collector:       DynamicInstrumentationCollector,
						CollectorAction: "",
						Ecosystem:       ecosystem.Npm,
						EcosystemAction: "install",
						Format:          "json",
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
					Framework:       None,
					Collector:       NoCollector,
					CollectorAction: "",
					EcosystemAction: "",
					Format:          "",
				},
			},
		},
		{
			input: NPMInstallWhileDynamicInstrumentation,
			want: want{
				urn:  "urn:scheduler:dynamic!npm,install.json",
				json: []byte(`"urn:scheduler:dynamic!npm,install.json"`),
				TypeComponents: TypeComponents{
					Framework:       Scheduler,
					Collector:       DynamicInstrumentationCollector,
					CollectorAction: "",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "install",
					Format:          "json",
				},
			},
		},
		// {
		// 	input: NPMTestWhileDynamicInstrumentation,
		// 	want: want{
		// 		urn:  "urn:scheduler:dynamic!npm,test.json",
		// 		json: []byte(`"urn:scheduler:dynamic!npm,test.json"`),
		// 		TypeComponents: TypeComponents{
		// 			Framework:        Scheduler,
		// 			Collector:        "dynamic",
		// 			CollectorAction: "",
		// 			Ecosystem:        ecosystem.Npm,
		// 			EcosystemAction: []string{"test"},
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
					Framework:       Hoarding,
					Collector:       DepsDevCollector,
					CollectorAction: "",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMTyposquat,
			want: want{
				urn:  "urn:hoarding:typosquat!npm.json",
				json: []byte(`"urn:hoarding:typosquat!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       TyposquatCollector,
					CollectorAction: "",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMMetadataEmptyDescription,
			want: want{
				urn:  "urn:hoarding:metadata,empty_descr!npm.json",
				json: []byte(`"urn:hoarding:metadata,empty_descr!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       MetadataCollector,
					CollectorAction: "empty_descr",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMMetadataMaintainersEmailCheck,
			want: want{
				urn:  "urn:hoarding:metadata,email_check!npm.json",
				json: []byte(`"urn:hoarding:metadata,email_check!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       MetadataCollector,
					CollectorAction: "email_check",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMMetadataVersion,
			want: want{
				urn:  "urn:hoarding:metadata,version!npm.json",
				json: []byte(`"urn:hoarding:metadata,version!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       MetadataCollector,
					CollectorAction: "version",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMMetadataMismatches,
			want: want{
				urn:  "urn:hoarding:metadata,mismatches!npm.json",
				json: []byte(`"urn:hoarding:metadata,mismatches!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       MetadataCollector,
					CollectorAction: "mismatches",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticAnalysisEnvExfiltration,
			want: want{
				urn:  "urn:hoarding:static,exfiltrate_env!npm.json",
				json: []byte(`"urn:hoarding:static,exfiltrate_env!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "exfiltrate_env",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticAnalysisDetachedProcessExecution,
			want: want{
				urn:  "urn:hoarding:static,detached_process_exec!npm.json",
				json: []byte(`"urn:hoarding:static,detached_process_exec!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "detached_process_exec",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticAnalysisEvalBase64,
			want: want{
				urn:  "urn:hoarding:static,base64_eval!npm.json",
				json: []byte(`"urn:hoarding:static,base64_eval!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "base64_eval",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticAnalysisShadyLinks,
			want: want{
				urn:  "urn:hoarding:static,shady_links!npm.json",
				json: []byte(`"urn:hoarding:static,shady_links!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "shady_links",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticAnalysisInstallScript,
			want: want{
				urn:  "urn:hoarding:static,install_script!npm.json",
				json: []byte(`"urn:hoarding:static,install_script!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "install_script",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
				},
			},
		},
		{
			input: NPMStaticNonRegistryDependency,
			want: want{
				urn:  "urn:hoarding:static,non_registry_dependency!npm.json",
				json: []byte(`"urn:hoarding:static,non_registry_dependency!npm.json"`),
				TypeComponents: TypeComponents{
					Framework:       Hoarding,
					Collector:       StaticAnalysisCollector,
					CollectorAction: "non_registry_dependency",
					Ecosystem:       ecosystem.Npm,
					EcosystemAction: "",
					Format:          "json",
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

	assert.Equal(t, NPMStaticNonRegistryDependency, got)
}
