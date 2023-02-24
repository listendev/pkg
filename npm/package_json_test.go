package npm

import (
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"
)

func TestNewPackageJSONFromReader(t *testing.T) {
	tests := []struct {
		desc    string
		input   string
		output  *PackageJSON
		wantErr string
	}{
		{
			desc:    "empty",
			input:   "",
			output:  nil,
			wantErr: "couldn't instantiate from the input package.json contents",
		},
		{
			desc: "dep-devdep-peer-peermeta-bundle",
			input: heredoc.Doc(`{
	"name": "xxx",
	"dependencies": {
		"@isaacs/import-jsx": "^4.0.1",
		"@types/react": "^17.0.52",
		"chokidar": "^3.3.0",
		"findit": "^2.0.0",
		"foreground-child": "^2.0.0",
		"fs-exists-cached": "^1.0.0",
		"glob": "^7.2.3",
		"ink": "^3.2.0",
		"isexe": "^2.0.0",
		"istanbul-lib-processinfo": "^2.0.3",
		"jackspeak": "^1.4.2",
		"libtap": "^1.4.0",
		"minipass": "^3.3.4",
		"mkdirp": "^1.0.4",
		"nyc": "^15.1.0",
		"opener": "^1.5.1",
		"react": "^17.0.2",
		"rimraf": "^3.0.0",
		"signal-exit": "^3.0.6",
		"source-map-support": "^0.5.16",
		"tap-mocha-reporter": "^5.0.3",
		"tap-parser": "^11.0.2",
		"tap-yaml": "^1.0.2",
		"tcompare": "^5.0.7",
		"treport": "^3.0.4",
		"which": "^2.0.2"
	},
	"devDependencies": {
		"coveralls": "^3.1.1",
		"eslint": "^7.32.0",
		"flow-remove-types": "^2.193.0",
		"node-preload": "^0.2.1",
		"process-on-spawn": "^1.0.0",
		"ts-node": "^8.5.2",
		"typescript": "^3.7.2"
	},
	"peerDependencies": {
		"coveralls": "^3.1.1",
		"flow-remove-types": ">=2.112.0",
		"ts-node": ">=8.5.2",
		"typescript": ">=3.7.2"
	},
	"peerDependenciesMeta": {
		"coveralls": {
		"optional": true
		},
		"flow-remove-types": {
		"optional": true
		},
		"ts-node": {
		"optional": true
		},
		"typescript": {
		"optional": true
		}
	},
	"bundleDependencies": [
		"ink",
		"treport",
		"@types/react",
		"@isaacs/import-jsx",
		"react"
	]
}`),
			output: &PackageJSON{
				Name: "xxx",
				Dependencies: map[string]string{
					"@isaacs/import-jsx":       "^4.0.1",
					"@types/react":             "^17.0.52",
					"chokidar":                 "^3.3.0",
					"findit":                   "^2.0.0",
					"foreground-child":         "^2.0.0",
					"fs-exists-cached":         "^1.0.0",
					"glob":                     "^7.2.3",
					"ink":                      "^3.2.0",
					"isexe":                    "^2.0.0",
					"istanbul-lib-processinfo": "^2.0.3",
					"jackspeak":                "^1.4.2",
					"libtap":                   "^1.4.0",
					"minipass":                 "^3.3.4",
					"mkdirp":                   "^1.0.4",
					"nyc":                      "^15.1.0",
					"opener":                   "^1.5.1",
					"react":                    "^17.0.2",
					"rimraf":                   "^3.0.0",
					"signal-exit":              "^3.0.6",
					"source-map-support":       "^0.5.16",
					"tap-mocha-reporter":       "^5.0.3",
					"tap-parser":               "^11.0.2",
					"tap-yaml":                 "^1.0.2",
					"tcompare":                 "^5.0.7",
					"treport":                  "^3.0.4",
					"which":                    "^2.0.2",
				},
				DevDependencies: map[string]string{
					"coveralls":         "^3.1.1",
					"eslint":            "^7.32.0",
					"flow-remove-types": "^2.193.0",
					"node-preload":      "^0.2.1",
					"process-on-spawn":  "^1.0.0",
					"ts-node":           "^8.5.2",
					"typescript":        "^3.7.2",
				},
				PeerDependencies: map[string]string{
					"coveralls":         "^3.1.1",
					"flow-remove-types": ">=2.112.0",
					"ts-node":           ">=8.5.2",
					"typescript":        ">=3.7.2",
				},
				BundleDependencies: []string{
					"ink",
					"treport",
					"@types/react",
					"@isaacs/import-jsx",
					"react",
				},
			},
			wantErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			res, err := NewPackageJSONFromReader(strings.NewReader(tc.input))
			if err != nil {
				assert.Nil(t, res)
				if assert.Error(t, err) {
					assert.Equal(t, tc.wantErr, err.Error())
				}
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.output, res)
			}
		})
	}
}
