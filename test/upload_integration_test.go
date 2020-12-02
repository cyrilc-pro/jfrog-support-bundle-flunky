package test

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-cli-core/utils/ioutils"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func Test_UploadIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []integrationTest{
		{
			Name: "Upload to specified target using target credentials",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				testBundle := getSupportBundle(t)
				err := commands.UploadSupportBundle(&commands.HTTPClient{RtDetails: targetRtDetails},
					&commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}, testBundle, "logs",
					func() time.Time { return time.Unix(1, 1) })
				assert.NoError(t, err)
			},
		},
		{
			Name: "Upload to default target without credentials",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				testBundle := getSupportBundle(t)
				targetDetailsWithoutCreds := &config.ArtifactoryDetails{Url: targetRtDetails.Url}
				err := commands.UploadSupportBundle(&commands.HTTPClient{RtDetails: targetDetailsWithoutCreds},
					&commands.SupportBundleCommandConfiguration{
						CaseNumber:          "foo",
						JfrogSupportLogsURL: targetRtDetails.GetUrl(),
					}, testBundle, "logs", func() time.Time { return time.Unix(2, 2) })
				assert.NoError(t, err)
			},
		},
		{
			Name: "Upload when target is offline",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				testBundle := getSupportBundle(t)
				invalidTarget := &config.ArtifactoryDetails{Url: "http://invalid"}
				err := commands.UploadSupportBundle(&commands.HTTPClient{RtDetails: invalidTarget},
					&commands.SupportBundleCommandConfiguration{
						CaseNumber:          "foo",
						JfrogSupportLogsURL: targetRtDetails.GetUrl(),
					}, testBundle, "logs", func() time.Time { return time.Unix(3, 3) })
				require.Error(t, err)
				assert.Contains(t, err.Error(), "dial tcp:")
			},
		},
	}
	runIntegrationTests(t, tests)
}

func getSupportBundle(t *testing.T) string {
	dir := os.TempDir()
	testBundle := filepath.Join(dir, "foo")
	// nolint: gocritic // octalLiteral
	err := ioutils.CopyFile("testdata/sb.zip", testBundle, 0644)
	require.NoError(t, err)
	return testBundle
}
