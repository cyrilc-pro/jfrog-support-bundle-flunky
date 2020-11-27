package commands

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_CreateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []IntegrationTest{
		{
			Name: "Success",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				r, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf, Now)
				require.NoError(t, err)
				require.NotEmpty(t, r.ID)
			},
		},
		{
			Name: "Offline",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				_, err := createSupportBundle(&HTTPClient{rtDetails: &config.ArtifactoryDetails{
					Url: "http://unknown.invalid/",
				}}, &conf, Now)
				require.Error(t, err)
				// exact message depends on OS
				require.Contains(t, err.Error(), "dial tcp:")
			},
		},
	}
	RunIntegrationTests(t, tests)
}
