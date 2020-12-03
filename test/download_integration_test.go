package test

import (
	"archive/zip"
	"context"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_DownloadIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []integrationTest{
		{
			Name: "Success",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				supportBundle := setUpSupportBundle(t, rtDetails)
				bundle, err := actions.DownloadSupportBundle(context.Background(), &http.Client{RtDetails: rtDetails},
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
				bundle, err := actions.DownloadSupportBundle(context.Background(), &http.Client{RtDetails: rtDetails},
					1*time.Second, 100*time.Millisecond, "unknown")
				require.Empty(t, bundle)
				assert.EqualError(t, err, "http request failed with: 404 Not Found")
			},
		},
	}
	runIntegrationTests(t, tests)
}

func assertBundleIsAZipArchive(t *testing.T, bundle string) {
	r, err := zip.OpenReader(bundle)
	require.NoError(t, err)
	require.NoError(t, r.Close())
}

func setUpSupportBundle(t *testing.T, rtDetails *config.ArtifactoryDetails) actions.BundleID {
	t.Helper()
	supportBundle, err := actions.CreateSupportBundle(&http.Client{RtDetails: rtDetails}, "foo",
		actions.NewDefaultOptionsProvider())
	require.NoError(t, err)
	require.NotEmpty(t, supportBundle)
	return supportBundle
}
