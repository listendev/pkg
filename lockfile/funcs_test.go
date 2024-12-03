package lockfile

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	type testCase struct {
		input []string
		want  map[Lockfile][]string
	}

	cases := []testCase{
		{
			want: map[Lockfile][]string{},
		},
		{
			input: []string{"unknown.json"},
			want:  map[Lockfile][]string{},
		},
		{
			input: []string{"package-lock.json"},
			want:  map[Lockfile][]string{PackageLockJSON: {"package-lock.json"}},
		},
		{
			input: []string{"package-lock.JSON"},
			want:  map[Lockfile][]string{PackageLockJSON: {"package-lock.JSON"}},
		},
		{
			input: []string{"working/dir/package-lock.JSON"},
			want:  map[Lockfile][]string{PackageLockJSON: {"working/dir/package-lock.JSON"}},
		},
		{
			input: []string{"poetry.lock"},
			want:  map[Lockfile][]string{PoetryLock: {"poetry.lock"}},
		},
		{
			input: []string{"somedir/poetry.lock"},
			want:  map[Lockfile][]string{PoetryLock: {"somedir/poetry.lock"}},
		},
		{
			input: []string{"somedir/poetry.lock", "package-lock.json"},
			want:  map[Lockfile][]string{PoetryLock: {"somedir/poetry.lock"}, PackageLockJSON: {"package-lock.json"}},
		},
		{
			input: []string{"somedir/poetry.lock", "somedir/poetry.lock", "package-lock.json"},
			want:  map[Lockfile][]string{PoetryLock: {"somedir/poetry.lock"}, PackageLockJSON: {"package-lock.json"}},
		},
		{
			input: []string{"somedir/poetry.lock", "package-lock.json", "otherdir/poetry.lock"},
			want:  map[Lockfile][]string{PoetryLock: {"somedir/poetry.lock", "otherdir/poetry.lock"}, PackageLockJSON: {"package-lock.json"}},
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want, Map(tc.input))
	}
}

func TestExisting(t *testing.T) {
	type testCase struct {
		input   []string
		want    map[Lockfile][]string
		wantErr map[Lockfile][]error
	}

	cases := []testCase{
		{
			want:    map[Lockfile][]string{},
			wantErr: map[Lockfile][]error{},
		},
		{
			input:   []string{"unknown.json"},
			want:    map[Lockfile][]string{},
			wantErr: map[Lockfile][]error{},
		},
		{
			input:   []string{"package-lock.json"},
			want:    map[Lockfile][]string{},
			wantErr: map[Lockfile][]error{PackageLockJSON: {errors.New("package-lock.json not found")}},
		},
		{
			input:   []string{"somedir/poetry.lock"},
			want:    map[Lockfile][]string{},
			wantErr: map[Lockfile][]error{PoetryLock: {errors.New("somedir/poetry.lock not found")}},
		},
		{
			input:   []string{"testdata/poetry.lock", "testdata/package-lock.json"},
			want:    map[Lockfile][]string{PoetryLock: {"testdata/poetry.lock"}, PackageLockJSON: {"testdata/package-lock.json"}},
			wantErr: map[Lockfile][]error{},
		},
		{
			input:   []string{"unk/poetry.lock", "testdata/package-lock.json"},
			want:    map[Lockfile][]string{PackageLockJSON: {"testdata/package-lock.json"}},
			wantErr: map[Lockfile][]error{PoetryLock: {errors.New("unk/poetry.lock not found")}},
		},
		{
			input:   []string{"testdata/package-lock.json", "testdata/poetry.lock", "unk/poetry.lock", "testdata/1/poetry.lock", "boh/package-lock.json"},
			want:    map[Lockfile][]string{PackageLockJSON: {"testdata/package-lock.json"}, PoetryLock: {"testdata/poetry.lock", "testdata/1/poetry.lock"}},
			wantErr: map[Lockfile][]error{PoetryLock: {errors.New("unk/poetry.lock not found")}, PackageLockJSON: {errors.New("boh/package-lock.json not found")}},
		},
	}

	for _, tc := range cases {
		gotMap, gotErr := Existing(tc.input)
		require.Equal(t, tc.want, gotMap)
		require.Equal(t, tc.wantErr, gotErr)
	}
}
