package manifest

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	type testCase struct {
		input []string
		want  map[Manifest]string
	}

	cases := []testCase{
		{
			want: map[Manifest]string{},
		},
		{
			input: []string{"unknown.json"},
			want:  map[Manifest]string{},
		},
		{
			input: []string{"package.json"},
			want:  map[Manifest]string{PackageJSON: "package.json"},
		},
		{
			input: []string{"package.JSON"},
			want:  map[Manifest]string{PackageJSON: "package.JSON"},
		},
		{
			input: []string{"working/dir/package.JSON"},
			want:  map[Manifest]string{PackageJSON: "working/dir/package.JSON"},
		},
		// TODO: uncomment when available
		// {
		// 	input: []string{"requirements.txt"},
		// 	want:  map[Manifest]string{RequirementsTxt: "requirements.txt"},
		// },
		// {
		// 	input: []string{"somedir/requirements.txt"},
		// 	want:  map[Manifest]string{RequirementsTxt: "somedir/requirements.txt"},
		// },
		// {
		// 	input: []string{"somedir/requirements.txt", "package.json"},
		// 	want:  map[Manifest]string{RequirementsTxt: "somedir/requirements.txt", PackageJSON: "package.json"},
		// },
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, Map(tc.input))
	}
}

func TestExisting(t *testing.T) {
	type testCase struct {
		input   []string
		want    map[Manifest]string
		wantErr map[Manifest]error
	}

	cases := []testCase{
		{
			want:    map[Manifest]string{},
			wantErr: map[Manifest]error{},
		},
		{
			input:   []string{"unknown.json"},
			want:    map[Manifest]string{},
			wantErr: map[Manifest]error{},
		},
		// FIXME: doesn't work in GitHub actions?!
		// {
		// 	input:   []string{"testdata/package.JSON"},
		// 	want:    map[Manifest]string{PackageJSON: "testdata/package.JSON"},
		// 	wantErr: map[Manifest]error{},
		// },
		{
			input:   []string{"package.json"},
			want:    map[Manifest]string{},
			wantErr: map[Manifest]error{PackageJSON: fmt.Errorf("package.json not found")},
		},
		// TODO: uncomment when available
		// {
		// 	input:   []string{"somedir/requirements.txt"},
		// 	want:    map[Manifest]string{},
		// 	wantErr: map[Manifest]error{RequirementsTxt: fmt.Errorf("somedir/requirements.txt not found")},
		// },
		// {
		// 	input:   []string{"testdata/poetry.lock", "testdata/package.json"},
		// 	want:    map[Manifest]string{RequirementsTxt: "testdata/poetry.lock", PackageJSON: "testdata/package.json"},
		// 	wantErr: map[Manifest]error{},
		// },
		// {
		// 	input:   []string{"unk/requirements.txt", "testdata/package.json"},
		// 	want:    map[Manifest]string{PackageJSON: "testdata/package.json"},
		// 	wantErr: map[Manifest]error{RequirementsTxt: fmt.Errorf("unk/requirements.txt not found")},
		// },
		// {
		// 	// Order matters: the last poetry.lock overrides the previous one
		// 	input:   []string{"testdata/package.json", "testdata/requirements.txt", "unk/requirements.txt"},
		// 	want:    map[Manifest]string{PackageJSON: "testdata/package.json"},
		// 	wantErr: map[Manifest]error{RequirementsTxt: fmt.Errorf("unk/requirements.txt not found")},
		// },
	}

	for _, tc := range cases {
		gotMap, gotErr := Existing(tc.input)
		require.Equal(t, tc.want, gotMap)
		require.Equal(t, tc.wantErr, gotErr)
	}
}
