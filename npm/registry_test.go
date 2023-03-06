package npm

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/garnet-org/pkg/observability"
)

func TestGetPackageList(t *testing.T) {
	// Set up a mock HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		plist, err := ioutil.ReadFile("testdata/package_list.json")
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
	packageList, err := client.GetPackageList(testCtx, "chalk")
	if err != nil {
		t.Fatal(err)
	}

	if packageList.Name != "chalk" {
		t.Errorf("Expected name to be 'chalk', got '%s'", packageList.Name)
	}
	if len(packageList.Versions) != 37 {
		t.Errorf("Expected 37 versions, got %d", len(packageList.Versions))
	}

	firstVersion := packageList.Versions["0.1.0"]
	if firstVersion.Dist.Shasum != "69afbee2ffab5e0db239450767a6125cbea50fa2" {
		t.Errorf("Expected shasum to be '69afbee2ffab5e0db239450767a6125cbea50fa2', got '%s'", firstVersion.Dist.Shasum)
	}

	lastVersionTag := packageList.DistTags.Latest
	lastVersion := packageList.Versions[lastVersionTag]
	if lastVersion.Dist.Shasum != "249623b7d66869c673699fb66d65723e54dfcfb3" {
		t.Errorf("Expected shasum to be '249623b7d66869c673699fb66d65723e54dfcfb3', got '%s'", lastVersion.Dist.Shasum)
	}

}

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
