package test

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_CreateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []integrationTest{
		{
			Name: "Success with default options",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
				id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: rtDetails}, &conf,
					&commands.DefaultOptionsProvider{GetDate: time.Now})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options disabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
				id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: rtDetails}, &conf,
					&commands.PromptOptionsProvider{GetDate: time.Now, Prompter: &commands.PrompterStub{
						IncludeLogs:          false,
						IncludeSystem:        false,
						IncludeConfiguration: false,
						IncludeThreadDump:    false,
					}})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options enabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
				id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: rtDetails}, &conf,
					&commands.PromptOptionsProvider{GetDate: time.Now, Prompter: &commands.PrompterStub{
						IncludeLogs:          true,
						IncludeSystem:        true,
						IncludeConfiguration: true,
						IncludeThreadDump:    true,
					}})
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Offline",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
				_, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: &config.ArtifactoryDetails{
					Url: "http://unknown.invalid/",
				}}, &conf, &commands.DefaultOptionsProvider{GetDate: time.Now})
				require.Error(t, err)
				// exact message depends on OS
				require.Contains(t, err.Error(), "dial tcp:")
			},
		},
	}
	runIntegrationTests(t, tests)
}
