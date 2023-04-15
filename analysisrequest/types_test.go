package analysisrequest

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			input: Nop,
			want: want{
				urn:  "urn:NOP:nop",
				json: []byte(`"urn:nop:nop"`),
				TypeComponents: TypeComponents{
					Framework: "nop",
					Collector: "nop",
					Actions:   []string{},
					Format:    "json",
				},
			},
		},
		{
			input: NPMInstallWhileFalco,
			want: want{
				urn:  "urn:scheduler:falco!npm.install",
				json: []byte(`"urn:scheduler:falco!npm.install"`),
				TypeComponents: TypeComponents{
					Framework: "scheduler",
					Collector: "falco",
					Ecosystem: NPMEcosystem,
					Actions:   []string{"install"},
					Format:    "json",
				},
			},
		},
		{
			input: NPMTestWhileFalco,
			want: want{
				urn:  "urn:scheduler:falco!npm.test",
				json: []byte(`"urn:scheduler:falco!npm.test"`),
				TypeComponents: TypeComponents{
					Framework: "scheduler",
					Collector: "falco",
					Ecosystem: NPMEcosystem,
					Actions:   []string{"test"},
					Format:    "json",
				},
			},
		},
		{
			input: NPMDepsDev,
			want: want{
				urn:  "urn:hoarding:depsdev!npm",
				json: []byte(`"urn:hoarding:depsdev!npm"`),
				TypeComponents: TypeComponents{
					Framework: "hoarding",
					Collector: "depsdev",
					Ecosystem: NPMEcosystem,
					Actions:   []string{},
					Format:    "json",
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

	assert.Equal(t, NPMDepsDev, got)
}
