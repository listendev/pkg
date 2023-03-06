package npm

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/garnet-org/pkg/observability"
)

func TestNPMRegistryClient_GetPackageVersion(t *testing.T) {
	tests := []struct {
		name        string
		testFile    string
		wantName    string
		wantVersion string
		wantShasum  string
		wantErr     bool
	}{
		{
			name:        "react package from upstream registry",
			testFile:    "package_version.json",
			wantName:    "react",
			wantVersion: "15.4.0",
			wantShasum:  "736c1c7c542e8088127106e1f450b010f86d172b",
		},
		{
			name:        "package from verdaccio registry",
			testFile:    "package_version_verdaccio.json",
			wantName:    "@frontend-metrics/hotjar",
			wantVersion: "951.512.2-garnet.0",
			wantShasum:  "4e43b7db05c8ba37b128058ba1659911e10ee971",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				plist, err := ioutil.ReadFile(path.Join("testdata/", tt.testFile))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(plist); err != nil {
					t.Fatal(err)
				}
			}))
			defer ts.Close()

			// Create a new client using the mock server as the base URL
			client, err := NewNPMRegistryClient(NPMRegistryClientConfig{
				BaseURL: ts.URL,
			})
			if err != nil {
				t.Fatal(err)
			}

			testCtx := observability.NewNopContext()
			packageVersion, err := client.GetPackageVersion(testCtx, "chalk", "5.1.2")
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

func TestNPMRegistryClient_GetPackageList(t *testing.T) {
	tests := []struct {
		name                  string
		testFile              string
		searchName            string
		wantName              string
		wantVersionsShasums   map[string]string
		wantLastVersionTag    string
		wantLastVersionShasum string
		wantErr               bool
	}{
		{
			name:                  "chalk package from upstream registry",
			testFile:              "package_list.json",
			searchName:            "chalk",
			wantName:              "chalk",
			wantLastVersionTag:    "5.2.0",
			wantLastVersionShasum: "249623b7d66869c673699fb66d65723e54dfcfb3",
			wantVersionsShasums: map[string]string{
				"0.1.0": "69afbee2ffab5e0db239450767a6125cbea50fa2",
			},
		},
		{
			name:                  "hotjar package from verdaccio registry",
			testFile:              "package_list_verdaccio.json",
			searchName:            "@frontend-metrics/hotjar",
			wantName:              "@frontend-metrics/hotjar",
			wantLastVersionTag:    "951.512.2-garnet.1",
			wantLastVersionShasum: "198dbaaef01cd7430b095ebc8928ee4e926fb04f",
			wantVersionsShasums: map[string]string{
				"0.0.1-security": "2f333c605d19e3be360cc541ad4521a750931968",
				"951.512.0":      "b46803072b62b7afb160b64a5df6c6bd74fb2f25",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				plist, err := ioutil.ReadFile(path.Join("testdata/", tt.testFile))
				if err != nil {
					t.Fatal(err)
				}
				if _, err := w.Write(plist); err != nil {
					t.Fatal(err)
				}
			}))
			defer ts.Close()

			client, err := NewNPMRegistryClient(NPMRegistryClientConfig{
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
		})
	}
}
