package commands

import (
	"archive/zip"
	"context"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_DownloadIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []IntegrationTest{
		{
			Name: "Success",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				supportBundle := setUpSupportBundle(t, rtDetails)
				bundle, err := downloadSupportBundle(context.Background(), &HTTPClient{rtDetails: rtDetails},
					30*time.Second, 100*time.Millisecond, supportBundle)
				require.NoError(t, err)
				assert.Contains(t, bundle, supportBundle)
				assert.True(t, fileutils.IsZip(bundle))
				assertBundleIsAZipArchive(t, bundle)
			},
		},
		{
			Name: "Not found",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				bundle, err := downloadSupportBundle(context.Background(), &HTTPClient{rtDetails: rtDetails},
					1*time.Second, 100*time.Millisecond, "unknown")
				require.Empty(t, bundle)
				assert.EqualError(t, err, "http request failed with: 404 Not Found")
			},
		},
	}
	RunIntegrationTests(t, tests)
}

func assertBundleIsAZipArchive(t *testing.T, bundle string) {
	r, err := zip.OpenReader(bundle)
	require.NoError(t, err)
	require.NoError(t, r.Close())
}

func setUpSupportBundle(t *testing.T, rtDetails *config.ArtifactoryDetails) bundleID {
	t.Helper()
	conf := supportBundleCommandConfiguration{caseNumber: "foo"}
	supportBundle, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf, &defaultOptionsProvider{getDate: time.Now})
	require.NoError(t, err)
	require.NotEmpty(t, supportBundle)
	return supportBundle
}
