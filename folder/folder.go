package folder

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
)

type Folder struct {
	fs billy.Filesystem
}

func New(filenames ...string) *Folder {
	return &Folder{memfs.New()}
}

func (f *Folder) Write(filename string, body []byte) error {
	dir := filepath.Dir(filename)

	if err := f.fs.MkdirAll(dir, 0777); err != nil {
		return err
	}

	file, err := f.fs.Create(filename)
	if err != nil {
		return err
	}

	if _, err := file.Write(body); err != nil {
		return err
	}

	return nil
}

func (f *Folder) Exists(filename string) (bool, error) {
	_, err := f.fs.Stat(filename)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func (f *Folder) ForEachFile(dir string, cb func(path string, file fs.FileInfo)) error {
	exists, err := f.Exists(dir)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("dir %q doesn't exists", dir)
	}

	files, err := f.fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		subdir := filepath.Join(dir, file.Name())
		if file.IsDir() {
			f.ForEachFile(subdir, cb)
			continue
		}

		cb(subdir, file)
	}

	return nil
}

func (f *Folder) GroupLeaves() map[string][]fs.FileInfo {
	grouped := make(map[string][]fs.FileInfo)
	f.ForEachFile("/", func(path string, file fs.FileInfo) {
		dir, _ := filepath.Split(path)
		files, ok := grouped[dir]
		if !ok {
			files = []fs.FileInfo{}
		}
		grouped[dir] = append(files, file)
	})
	return grouped
}

func (f *Folder) ForEachMissing(targetFilenames []string, cb func(dir, missingFilename string)) {
	grouped := f.GroupLeaves()

	// There's no tree at all
	// Which means all the files are missing
	if len(grouped) == 0 {
		for _, target := range targetFilenames {
			cb("/", target)
		}

		return
	}

	for dir, files := range grouped {
		for _, target := range targetFilenames {
			found := false
			for _, file := range files {
				if file.Name() == target {
					found = true
					break
				}
			}
			if !found {
				cb(dir, target)
			}
		}
	}
}
