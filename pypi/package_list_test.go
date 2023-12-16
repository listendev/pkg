package pypi

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		descr          string
		fromList       bool
		wantLatest     string
		wantName       string
		wantVersion    string
		wantVersionErr error
	}{
		{
			descr:       "instance from list endpoint",
			fromList:    true,
			wantLatest:  "1.34.2",
			wantName:    "boto3",
			wantVersion: "1.33.8",
		},
		{
			descr:       "instance from version (1.33.8) endpoint",
			fromList:    false,
			wantLatest:  "",
			wantName:    "boto3",
			wantVersion: "1.33.8",
		},
		{
			descr:          "asking wrong existing version to instance from version (1.33.8) endpoint",
			fromList:       false,
			wantLatest:     "",
			wantName:       "boto3",
			wantVersion:    "0.0.1",
			wantVersionErr: ErrVersionMismatch,
		},
		{
			descr:       "asking existing version to instance from list endpoint",
			fromList:    true,
			wantLatest:  "1.34.2",
			wantName:    "boto3",
			wantVersion: "0.0.1",
		},
		{
			descr:          "asking not existing version to instance from list endpoint",
			fromList:       true,
			wantLatest:     "1.34.2",
			wantName:       "boto3",
			wantVersion:    "x.y.z",
			wantVersionErr: ErrVersionNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			filename := "package_version.json"
			if tt.fromList {
				filename = "package_list.json"
			}
			plistBytes, err := os.ReadFile(path.Join("testdata/", filename))
			if err != nil {
				t.Fatal(err)
			}
			var plist PackageList
			err = json.NewDecoder(bytes.NewReader(plistBytes)).Decode(&plist)
			if err != nil {
				t.Fatal(err)
			}

			latestV, latestErr := plist.GetVersion("latest")
			if !tt.fromList {
				require.Error(t, latestErr)
				require.Nil(t, latestV)
				require.Empty(t, tt.wantLatest)
			} else {
				require.NoError(t, latestErr)
				require.NotNil(t, latestV)
				require.Equal(t, tt.wantLatest, latestV.Version)
				require.Equal(t, tt.wantName, latestV.Name)
			}

			if tt.wantVersion != "" {
				pv, err := plist.GetVersion(tt.wantVersion)
				require.Equal(t, tt.wantVersionErr, err)
				if err != nil {
					require.Nil(t, pv)
				} else {
					assert.Equal(t, tt.wantVersion, pv.Version)
					assert.Equal(t, tt.wantName, pv.Name)
				}
			}
		})
	}
}
