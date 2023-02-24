package npm

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
)

// activeFS provides a configurable entry point to the billy.Filesystem in active use
// for reading global and system level configuration files.
//
// Override this in tests to mock the filesystem
// (then reset to restore default behavior).
var activeFS = defaultFS()

// defaultFS provides a billy.Filesystem abstraction over the
// OS filesystem (via osfs.OS) scoped to the root directory
// in order to enable access to global and system configuration files
// via absolute paths.
func defaultFS() billy.Filesystem {
	return osfs.New("/")
}
