package pypi

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/listendev/pkg/observability"
	"github.com/stretchr/testify/assert"
)

func TestRegistryClinet_GetPackageList(t *testing.T) {
	tests := []struct {
		descr                     string
		testFile                  string
		searchName                string
		wantName                  string
		wantVersionsSha256        map[string]string
		wantLastVersionTag        string
		wantLastVersionBlake2b256 string
		wantLastVersionSha256     string
		wantErr                   bool
	}{
		{
			descr:                     "boto3 package from upstream registry",
			testFile:                  "package_list.json",
			searchName:                "boto3",
			wantName:                  "boto3",
			wantLastVersionTag:        "1.34.2",
			wantLastVersionBlake2b256: "c86666f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f",
			wantLastVersionSha256:     "970fd9f9f522eb48f3cd5574e927b369279ebf5bcf0f2fae5ed9cc6306e58558",
			wantVersionsSha256: map[string]string{
				"0.0.1":   "bc018a3aedc5cf7329dcdeb435ece8a296b605c19fb09842c1821935f1b14cfd",
				"1.24.36": "b1855ede59e725b968d6336908ffc864b65985ca441d730625b09c43ccd6413b",
				"1.9.99":  "817b6f5e5277a9e370702314adbfcaa6957e138540e50d6b557a717846c6c999",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			for version, sha256 := range tt.wantVersionsSha256 {
				pkgVersions, ok := packageList.Versions[version]
				if !ok {
					t.Errorf("Expected version '%s' to be present in package list", version)
				}
				ver, sdistErr := pkgVersions.GetSdist()
				assert.Nil(t, sdistErr)
				assert.NotNil(t, ver)
				assert.Equal(t, version, ver.Version)
				assert.Equal(t, tt.searchName, ver.Name)
				if ver.Digests.SHA256 != sha256 {
					t.Errorf("Expected version sha256 digest to be '%s', got '%s'", sha256, ver.Digests.SHA256)
				}
			}

			lastVersionTag := packageList.Info.Version
			if lastVersionTag != tt.wantLastVersionTag {
				t.Errorf("Expected last version tag to be '%s', got '%s'", tt.wantLastVersionTag, lastVersionTag)
			}
			lastVersions := packageList.Versions[lastVersionTag]
			lastSdistVersion, err := lastVersions.GetSdist()
			assert.Nil(t, err)
			assert.NotNil(t, lastSdistVersion)
			if lastSdistVersion.Digests.Blake2bB256 != tt.wantLastVersionBlake2b256 {
				t.Errorf("Expected last version blake2b_256 digest to be '%s', got '%s'", tt.wantLastVersionBlake2b256, lastSdistVersion.Digests.Blake2bB256)
			}
			if lastSdistVersion.Digests.SHA256 != tt.wantLastVersionSha256 {
				t.Errorf("Expected last version sha256 digest to be '%s', got '%s'", tt.wantLastVersionSha256, lastSdistVersion.Digests.SHA256)
			}

			gotLastPackageVersion, err := client.GetPackageLatestVersion(testCtx, tt.searchName)
			assert.Nil(t, err)
			assert.Equal(t, lastSdistVersion, gotLastPackageVersion)
		})
	}
}

func TestRegistryClient_GetPackageVersion(t *testing.T) {
	tests := []struct {
		descr          string
		name           string
		version        string
		testFile       string
		wantName       string
		wantVersion    string
		wantSha256     string
		wantBlake2b256 string
		wantErr        bool
	}{
		{
			descr:          "boto3 1.33.8 package from upstream registry",
			name:           "boto3",
			version:        "1.33.8",
			testFile:       "package_version.json",
			wantName:       "boto3",
			wantVersion:    "1.33.8",
			wantSha256:     "d02a084b25aa8d46ef917b128e90877efab1ba45f9d1ba3a11f336930378e350",
			wantBlake2b256: "121f1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.descr, func(t *testing.T) {
			// Set up a mock HTTP server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			if packageVersion.Digests.Blake2bB256 != tt.wantBlake2b256 {
				t.Errorf("Expected blake2b_256 digest to be '%s', got '%s'", tt.wantBlake2b256, packageVersion.Digests.Blake2bB256)
			}
			if packageVersion.Digests.SHA256 != tt.wantSha256 {
				t.Errorf("Expected sha256 digest to be '%s', got '%s'", tt.wantSha256, packageVersion.Digests.SHA256)
			}
		})
	}
}
