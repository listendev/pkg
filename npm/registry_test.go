package npm

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/listendev/pkg/observability"
	"github.com/stretchr/testify/require"
)

func TestRegistryClient_GetPackageVersion(t *testing.T) {
	tests := []struct {
		version     string
		name        string
		descr       string
		testFile    string
		wantName    string
		wantVersion string
		wantShasum  string
		wantErr     bool
	}{
		{
			descr:       "react package from upstream registry",
			name:        "react",
			version:     "15.4.0",
			testFile:    "package_version.json",
			wantName:    "react",
			wantVersion: "15.4.0",
			wantShasum:  "736c1c7c542e8088127106e1f450b010f86d172b",
		},
		{
			descr:       "@frontend-metrics/hotjar package from verdaccio registry",
			name:        "@frontend-metrics/hotjar",
			testFile:    "package_version_verdaccio.json",
			wantName:    "@frontend-metrics/hotjar",
			wantVersion: "951.512.2-garnet.0",
			wantShasum:  "4e43b7db05c8ba37b128058ba1659911e10ee971",
		},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				plist, err := os.ReadFile(path.Join("testdata/", tt.testFile))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(plist); err != nil {
					t.Fatal(err)
				}
			}))
			defer ts.Close()

			// Create a new client using the mock server as the base URL
			client, err := NewRegistryClient(RegistryClientConfig{
				BaseURL: ts.URL,
			})
			if err != nil {
				t.Fatal(err)
			}

			testCtx := observability.NewNopContext()
			packageVersion, err := client.GetPackageVersion(testCtx, tt.name, tt.version)
			if err != nil {
				t.Fatal(err)
			}

			if packageVersion.Name != tt.wantName {
				t.Errorf("Expected name to be '%s', got '%s'", tt.wantName, packageVersion.Name)
			}
			if packageVersion.Version != tt.wantVersion {
				t.Errorf("Expected version to be '%s', got '%s'", tt.wantVersion, packageVersion.Version)
			}
			if packageVersion.Dist.Shasum != tt.wantShasum {
				t.Errorf("Expected shasum to be '%s', got '%s'", tt.wantShasum, packageVersion.Dist.Shasum)
			}
		})
	}
}

func TestRegistryClient_GetPackageList(t *testing.T) {
	tests := []struct {
		descr                            string
		testFile                         string
		searchName                       string
		wantName                         string
		wantVersionsShasums              map[string]string
		wantLastVersionTag               string
		wantLastVersionShasum            string
		wantLastVersionTime              time.Time
		wantLastVersionMaintainersEmails []string
		wantErr                          bool
	}{
		{
			descr:                 "chalk package from upstream registry",
			testFile:              "package_list.json",
			searchName:            "chalk",
			wantName:              "chalk",
			wantLastVersionTag:    "5.2.0",
			wantLastVersionShasum: "249623b7d66869c673699fb66d65723e54dfcfb3",
			wantLastVersionTime: func() time.Time {
				ret, _ := time.Parse(time.RFC3339Nano, "2022-12-08T18:46:27.169Z")

				return ret
			}(),
			wantLastVersionMaintainersEmails: []string{"josh@junon.me", "sindresorhus@gmail.com"},
			wantVersionsShasums: map[string]string{
				"0.1.0": "69afbee2ffab5e0db239450767a6125cbea50fa2",
			},
		},
		{
			descr:                 "hotjar package from verdaccio registry",
			testFile:              "package_list_verdaccio.json",
			searchName:            "@frontend-metrics/hotjar",
			wantName:              "@frontend-metrics/hotjar",
			wantLastVersionTag:    "951.512.2-garnet.1",
			wantLastVersionShasum: "198dbaaef01cd7430b095ebc8928ee4e926fb04f",
			wantLastVersionTime: func() time.Time {
				ret, _ := time.Parse(time.RFC3339Nano, "2023-02-23T04:36:57.696Z")

				return ret
			}(),
			wantVersionsShasums: map[string]string{
				"0.0.1-security": "2f333c605d19e3be360cc541ad4521a750931968",
				"951.512.0":      "b46803072b62b7afb160b64a5df6c6bd74fb2f25",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				plist, err := os.ReadFile(path.Join("testdata/", tt.testFile))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(plist); err != nil {
					t.Fatal(err)
				}
			}))
			defer ts.Close()

			client, err := NewRegistryClient(RegistryClientConfig{
				BaseURL: ts.URL,
			})
			if err != nil {
				t.Fatal(err)
			}

			testCtx := observability.NewNopContext()
			packageList, err := client.GetPackageList(testCtx, tt.searchName)
			if err != nil {
				t.Fatal(err)
			}
			for version, shasum := range tt.wantVersionsShasums {
				ver := packageList.Versions[version]
				if ver.Dist.Shasum != shasum {
					t.Errorf("Expected version shasum to be '%s', got '%s'", shasum, ver.Dist.Shasum)
				}
			}

			lastVersionTag := packageList.DistTags.Latest
			if lastVersionTag != tt.wantLastVersionTag {
				t.Errorf("Expected last version tag to be '%s', got '%s'", tt.wantLastVersionTag, lastVersionTag)
			}
			lastVersion := packageList.Versions[lastVersionTag]
			if lastVersion.Dist.Shasum != tt.wantLastVersionShasum {
				t.Errorf("Expected last version shasum to be '%s', got '%s'", tt.wantLastVersionShasum, lastVersion.Dist.Shasum)
			}

			gotLatestVersionTime, gotLatestVersionTimeErr := packageList.LatestVersionTime()
			require.Nil(t, gotLatestVersionTimeErr)
			require.NotNil(t, gotLatestVersionTime)
			if !cmp.Equal(*gotLatestVersionTime, tt.wantLastVersionTime, cmpopts.EquateApproxTime(time.Millisecond*2)) {
				t.Fatal(cmp.Diff(tt.wantLastVersionTime, *gotLatestVersionTime))
			}

			if len(tt.wantLastVersionMaintainersEmails) > 0 {
				gotM, gotMaintainersErr := packageList.MaintainersByVersion("latest")
				require.Nil(t, gotMaintainersErr)
				require.NotNil(t, gotM)

				gotEmails := gotM.Emails()
				if !cmp.Equal(gotEmails, tt.wantLastVersionMaintainersEmails, cmpopts.SortSlices(func(x, y string) bool {
					return x < y
				})) {
					t.Fatal(cmp.Diff(tt.wantLastVersionMaintainersEmails, gotEmails))
				}
			}
		})
	}
}
