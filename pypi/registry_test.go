package pypi

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryClinet_GetPackageList(t *testing.T) {
	tests := []struct {
		descr                            string
		testFile                         string
		searchName                       string
		wantName                         string
		wantVersionsSha256               map[string]string
		wantLastVersionURL               string
		wantLastVersionTag               string
		wantLastVersionBlake2b256        string
		wantLastVersionSha256            string
		wantLastVersionTime              time.Time
		wantLastVersionMaintainersEmails []string
		wantErr                          bool
	}{
		{
			descr:                            "boto3 package from upstream registry",
			testFile:                         "package_list.json",
			searchName:                       "boto3",
			wantName:                         "boto3",
			wantLastVersionTag:               "1.34.2",
			wantLastVersionBlake2b256:        "c86666f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f",
			wantLastVersionSha256:            "970fd9f9f522eb48f3cd5574e927b369279ebf5bcf0f2fae5ed9cc6306e58558",
			wantLastVersionTime:              func() time.Time { ret, _ := time.Parse(time.RFC3339Nano, "2023-12-15T20:43:50.976124Z"); return ret }(),
			wantLastVersionURL:               "https://files.pythonhosted.org/packages/c8/66/66f4e87201f72a79c2bf600f2b7096988572447f4a3dae38e4b4873a346f/boto3-1.34.2.tar.gz",
			wantLastVersionMaintainersEmails: []string{},
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
			if lastSdistVersion.URL != tt.wantLastVersionURL {
				t.Errorf("Expected last version URL to be '%s', got '%s'", tt.wantLastVersionURL, lastSdistVersion.URL)
			}

			gotLastPackageVersion, err := client.GetPackageLatestVersion(testCtx, tt.searchName)
			assert.Nil(t, err)
			assert.Equal(t, lastSdistVersion, gotLastPackageVersion)

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

func TestRegistryClient_GetPackageVersion(t *testing.T) {
	tests := []struct {
		descr          string
		name           string
		version        string
		testFile       string
		wantName       string
		wantURL        string
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
			wantURL:        "https://files.pythonhosted.org/packages/12/1f/1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a/boto3-1.33.8.tar.gz",
			wantVersion:    "1.33.8",
			wantSha256:     "d02a084b25aa8d46ef917b128e90877efab1ba45f9d1ba3a11f336930378e350",
			wantBlake2b256: "121f1d4c5bbe89542b62ec6a6ba624ef0142e1d0c3267711b4f01f6258399a0a",
		},
		{
			descr:          "cctx 1.0.0 package from upstream registry",
			name:           "cctx",
			version:        "1.0.0",
			testFile:       "cctx_100.json",
			wantName:       "cctx",
			wantURL:        "https://files.pythonhosted.org/packages/bd/d2/35ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e/cctx-1.0.0.tar.gz",
			wantVersion:    "1.0.0",
			wantSha256:     "1d9ceb0603ed51a4f337cb8d53dd320339fd10814642d074b41a86d00be0bdbd",
			wantBlake2b256: "bdd235ad05b2669c50fc2756e35d0fe462bbd085a5b7afb571f443fd2ceb151e",
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
			if packageVersion.URL != tt.wantURL {
				t.Errorf("Expected URL to be '%s', got '%s'", tt.wantURL, packageVersion.URL)
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
