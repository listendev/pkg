package lockfile

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	type testCase struct {
		input []string
		want  map[Lockfile]string
	}

	cases := []testCase{
		{
			want: map[Lockfile]string{},
		},
		{
			input: []string{"unknown.json"},
			want:  map[Lockfile]string{},
		},
		{
			input: []string{"package-lock.json"},
			want:  map[Lockfile]string{PackageLockJSON: "package-lock.json"},
		},
		{
			input: []string{"package-lock.JSON"},
			want:  map[Lockfile]string{PackageLockJSON: "package-lock.JSON"},
		},
		{
			input: []string{"working/dir/package-lock.JSON"},
			want:  map[Lockfile]string{PackageLockJSON: "working/dir/package-lock.JSON"},
		},
		{
			input: []string{"poetry.lock"},
			want:  map[Lockfile]string{PoetryLock: "poetry.lock"},
		},
		{
			input: []string{"somedir/poetry.lock"},
			want:  map[Lockfile]string{PoetryLock: "somedir/poetry.lock"},
		},
		{
			input: []string{"somedir/poetry.lock", "package-lock.json"},
			want:  map[Lockfile]string{PoetryLock: "somedir/poetry.lock", PackageLockJSON: "package-lock.json"},
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, Map(tc.input))
	}
}

func TestExisting(t *testing.T) {
	type testCase struct {
		input   []string
		want    map[Lockfile]string
		wantErr map[Lockfile]error
	}

	cases := []testCase{
		{
			want:    map[Lockfile]string{},
			wantErr: map[Lockfile]error{},
		},
		{
			input:   []string{"unknown.json"},
			want:    map[Lockfile]string{},
			wantErr: map[Lockfile]error{},
		},
		// FIXME: doesn't work in GitHub actions?!
		// {
		// 	input:   []string{"testdata/package-lock.JSON"},
		// 	want:    map[Lockfile]string{PackageLockJSON: "testdata/package-lock.JSON"},
		// 	wantErr: map[Lockfile]error{},
		// },
		{
			input:   []string{"package-lock.json"},
			want:    map[Lockfile]string{},
			wantErr: map[Lockfile]error{PackageLockJSON: fmt.Errorf("package-lock.json not found")},
		},
		{
			input:   []string{"somedir/poetry.lock"},
			want:    map[Lockfile]string{},
			wantErr: map[Lockfile]error{PoetryLock: fmt.Errorf("somedir/poetry.lock not found")},
		},
		{
			input:   []string{"testdata/poetry.lock", "testdata/package-lock.json"},
			want:    map[Lockfile]string{PoetryLock: "testdata/poetry.lock", PackageLockJSON: "testdata/package-lock.json"},
			wantErr: map[Lockfile]error{},
		},
		{
			input:   []string{"unk/poetry.lock", "testdata/package-lock.json"},
			want:    map[Lockfile]string{PackageLockJSON: "testdata/package-lock.json"},
			wantErr: map[Lockfile]error{PoetryLock: fmt.Errorf("unk/poetry.lock not found")},
		},
		{
			// Order matters: the last poetry.lock overrides the previous one
			input:   []string{"testdata/package-lock.json", "testdata/poetry.lock", "unk/poetry.lock"},
			want:    map[Lockfile]string{PackageLockJSON: "testdata/package-lock.json"},
			wantErr: map[Lockfile]error{PoetryLock: fmt.Errorf("unk/poetry.lock not found")},
		},
	}

	for _, tc := range cases {
		gotMap, gotErr := Existing(tc.input)
		require.Equal(t, tc.want, gotMap)
		require.Equal(t, tc.wantErr, gotErr)
	}
}
