package manifest

import (
	"github.com/listendev/pkg/ecosystem"
)

var manifests = map[ecosystem.Ecosystem][]Manifest{
	ecosystem.Npm: {
		PackageJSON,
	},
	ecosystem.Pypi: {},
	ecosystem.None: {},
}
