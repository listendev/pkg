package folder

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type FolderSuite struct {
	suite.Suite
}

func TestFolderSuite(t *testing.T) {
	suite.Run(t, new(FolderSuite))
}

func (s *FolderSuite) TestWriteAndExists() {
	tcs := []struct {
		name      string
		filenames []string
	}{
		{
			"single file",
			[]string{"/test"},
		},
		{
			"single file in dir",
			[]string{"/dir/test"},
		},
		{
			"multiple files",
			[]string{"/foo.txt", "/bar.go"},
		},
		{
			"multiple files in dirs",
			[]string{"/dir1/foo.txt", "/dir1/bar.go", "/dir2/baz.js"},
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			folder := New()
			for _, filename := range tc.filenames {
				assert.Nil(t, folder.Write(filename, []byte(nil)))
				exists, err := folder.Exists(filename)
				assert.Nil(t, err)
				assert.True(t, exists)
			}
		})
	}
}

func (s *FolderSuite) TestForEachFile() {
	folder := New()
	filenames := []string{"/foo", "/dir1/foo.go", "/dir1/bar", "/dir1/dir2/baz.txt", "/dir3/zam"}
	for _, filename := range filenames {
		assert.Nil(s.T(), folder.Write(filename, []byte(nil)))
	}

	tcs := []struct {
		name  string
		dir   string
		paths []string
		err   error
	}{
		{
			"root",
			"/",
			append([]string{}, filenames...),
			nil,
		},
		{
			"not existent dir",
			"/nope",
			[]string{},
			fmt.Errorf("dir %q doesn't exists", "/nope"),
		},
		{
			"subdir",
			"/dir1",
			[]string{"/dir1/foo.go", "/dir1/bar", "/dir1/dir2/baz.txt"},
			nil,
		},
		{
			"subsubdir",
			"/dir1/dir2",
			[]string{"/dir1/dir2/baz.txt"},
			nil,
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			var paths []string
			assert.Equal(t, tc.err, folder.ForEachFile(tc.dir, func(path string, file fs.FileInfo) {
				paths = append(paths, path)
			}), tc.err)
			assert.ElementsMatch(t, paths, tc.paths)
		})
	}
}

func (s *FolderSuite) TestGroupLeaves() {
	t := s.T()

	folder := New()
	filenames := []string{"/foo", "/dir1/foo.go", "/dir1/bar", "/dir1/dir2/baz.txt", "/dir3/zam"}
	for _, filename := range filenames {
		assert.Nil(s.T(), folder.Write(filename, []byte(nil)))
	}

	actual := folder.GroupLeaves()
	expected := map[string][]string{
		"/":           {"foo"},
		"/dir1/":      {"bar", "foo.go"},
		"/dir1/dir2/": {"baz.txt"},
		"/dir3/":      {"zam"},
	}

	// Assert the map has the expected keys
	ekeys := make([]string, 0, len(expected))
	for key := range expected {
		ekeys = append(ekeys, key)
	}
	lkeys := make([]string, 0, len(actual))
	for key := range actual {
		lkeys = append(lkeys, key)
	}
	assert.ElementsMatch(t, ekeys, lkeys)

	for dir, filenames := range expected {
		files, ok := actual[dir]
		assert.True(t, ok)

		var read []string
		for _, file := range files {
			read = append(read, file.Name())
		}

		assert.Equal(t, filenames, read)
	}
}

func (s *FolderSuite) TestForEachMissing() {
	tcs := []struct {
		name               string
		factory            func(f *Folder)
		collectorFilenames []string
		missingFilenames   []string
	}{
		{
			"empty folder",
			func(f *Folder) {},
			[]string{"foo", "bar"},
			[]string{"/foo", "/bar"},
		},
		{
			"both present",
			func(f *Folder) {
				assert.Nil(s.T(), f.Write("/foo", []byte(nil)))
				assert.Nil(s.T(), f.Write("/bar", []byte(nil)))
			},
			[]string{"foo", "bar"},
			[]string{},
		},
		{
			"one missing",
			func(f *Folder) {
				assert.Nil(s.T(), f.Write("/foo", []byte(nil)))
			},
			[]string{"foo", "bar"},
			[]string{"/bar"},
		},
		{
			"with nesting",
			func(f *Folder) {
				assert.Nil(s.T(), f.Write("/a/foo", []byte(nil)))
				assert.Nil(s.T(), f.Write("/a/bar", []byte(nil)))
				assert.Nil(s.T(), f.Write("/b/foo", []byte(nil)))
				assert.Nil(s.T(), f.Write("/b/zam", []byte(nil)))
			},
			[]string{"foo", "bar"},
			[]string{"/b/bar"},
		},
		{
			"filled folders but not collector file",
			func(f *Folder) {
				assert.Nil(s.T(), f.Write("/a/foo", []byte(nil)))
				assert.Nil(s.T(), f.Write("/a/bar", []byte(nil)))
				assert.Nil(s.T(), f.Write("/b/foo", []byte(nil)))
				assert.Nil(s.T(), f.Write("/b/bar", []byte(nil)))
			},
			[]string{"zam"},
			[]string{"/a/zam", "/b/zam"},
		},
	}

	for _, tc := range tcs {
		s.T().Run(tc.name, func(t *testing.T) {
			var read []string
			folder := New()
			tc.factory(folder)

			folder.ForEachMissing(tc.collectorFilenames, func(dir, missingFilename string) {
				read = append(read, filepath.Join(dir, missingFilename))
			})
			assert.ElementsMatch(t, tc.missingFilenames, read)
		})
	}
}
