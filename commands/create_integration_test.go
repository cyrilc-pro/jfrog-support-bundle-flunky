package commands

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CreateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []IntegrationTest{
		{
			Name: "Success with default options",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				id, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf, &defaultOptionsProvider{getDate: time.Now})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options disabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				id, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf,
					&promptOptionsProvider{getDate: time.Now, prompter: &prompterStub{
						includeLogs:          false,
						includeSystem:        false,
						includeConfiguration: false,
						includeThreadDump:    false,
					}})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options enabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				id, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf,
					&promptOptionsProvider{getDate: time.Now, prompter: &prompterStub{
						includeLogs:          true,
						includeSystem:        true,
						includeConfiguration: true,
						includeThreadDump:    true,
					}})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Offline",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				_, err := createSupportBundle(&HTTPClient{rtDetails: &config.ArtifactoryDetails{
					Url: "http://unknown.invalid/",
				}}, &conf, &defaultOptionsProvider{getDate: time.Now})
				require.Error(t, err)
				// exact message depends on OS
				require.Contains(t, err.Error(), "dial tcp:")
			},
		},
	}
	RunIntegrationTests(t, tests)
}
