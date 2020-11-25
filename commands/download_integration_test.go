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

func Test_Download(t *testing.T) {
	tests := []IntegrationTest{
		{
			Name: "Success",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails) {
				bundleID := setUpSupportBundle(t, rtDetails)
				bundle, err := downloadSupportBundle(context.Background(), &HTTPClient{rtDetails: rtDetails},
					30*time.Second, 100*time.Millisecond, bundleID)
				require.NoError(t, err)
				assert.Contains(t, bundle, bundleID)
				assert.True(t, fileutils.IsZip(bundle))
				assertBundleIsAZipArchive(t, bundle)
			},
		},
		{
			Name: "Not found",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails) {
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

func setUpSupportBundle(t *testing.T, rtDetails *config.ArtifactoryDetails) string {
	t.Helper()
	conf := supportBundleCommandConfiguration{caseNumber: "foo"}
	r, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf, Now)
	require.NoError(t, err)
	require.NotEmpty(t, r.ID)
	return r.ID
}
