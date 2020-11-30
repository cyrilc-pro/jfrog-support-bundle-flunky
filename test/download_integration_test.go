package test

import (
	"archive/zip"
	"context"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func (s *IntegrationTestSuite) Test_DownloadIntegration() {
	s.T().Run("Success", func(t *testing.T) {
		supportBundle := setUpSupportBundle(t, s.rtDetails)
		bundle, err := commands.DownloadSupportBundle(context.Background(), &commands.HTTPClient{RtDetails: s.rtDetails},
			30*time.Second, 100*time.Millisecond, supportBundle)
		require.NoError(t, err)
		assert.Contains(t, bundle, supportBundle)
		assert.True(t, fileutils.IsZip(bundle))
		assertBundleIsAZipArchive(t, bundle)
	})

	s.T().Run("NotFound", func(t *testing.T) {
		bundle, err := commands.DownloadSupportBundle(context.Background(), &commands.HTTPClient{RtDetails: s.rtDetails},
			1*time.Second, 100*time.Millisecond, "unknown")
		require.Empty(t, bundle)
		assert.EqualError(t, err, "http request failed with: 404 Not Found")
	})
}

func assertBundleIsAZipArchive(t *testing.T, bundle string) {
	r, err := zip.OpenReader(bundle)
	require.NoError(t, err)
	require.NoError(t, r.Close())
}

func setUpSupportBundle(t *testing.T, rtDetails *config.ArtifactoryDetails) commands.BundleID {
	t.Helper()
	conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
	supportBundle, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: rtDetails}, &conf,
		&commands.DefaultOptionsProvider{GetDate: time.Now})
	require.NoError(t, err)
	require.NotEmpty(t, supportBundle)
	return supportBundle
}
