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
		// TODO:
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

func TestMaxType(t *testing.T) {
	got := MaxType()

	assert.Equal(t, NPMTyposquat, got)
}
