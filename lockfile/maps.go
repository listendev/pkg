package lockfile

import (
	"github.com/listendev/pkg/ecosystem"
)

var lockfiles = map[ecosystem.Ecosystem][]Lockfile{
	ecosystem.Npm: {
		PackageLockJSON,
	},
	ecosystem.Pypi: {
		PoetryLock,
	},
	ecosystem.None: {},
}
