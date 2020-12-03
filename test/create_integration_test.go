package test

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
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
				id, err := createSupportBundle(rtDetails, actions.NewDefaultOptionsProvider())
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options disabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				id, err := createSupportBundle(rtDetails, newPromptOptionsProviderStub(false))
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Success with all options enabled",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				id, err := createSupportBundle(rtDetails, newPromptOptionsProviderStub(true))
				require.NoError(t, err)
				require.NotEmpty(t, id)
			},
		},
		{
			Name: "Offline",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails,
				targetRtDetails *config.ArtifactoryDetails) {
				_, err := createSupportBundle(&config.ArtifactoryDetails{Url: "http://unknown.invalid/"},
					actions.NewDefaultOptionsProvider())
				require.Error(t, err)
				// exact message depends on OS
				require.Contains(t, err.Error(), "dial tcp:")
			},
		},
	}
	runIntegrationTests(t, tests)
}

func newPromptOptionsProviderStub(includeAll bool) *actions.PromptOptionsProvider {
	return &actions.PromptOptionsProvider{GetDate: time.Now, Prompter: &actions.PrompterStub{
		IncludeLogs:          includeAll,
		IncludeSystem:        includeAll,
		IncludeConfiguration: includeAll,
		IncludeThreadDump:    includeAll,
	}}
}

func createSupportBundle(rtDetails *config.ArtifactoryDetails, optionsProvider actions.OptionsProvider) (
	actions.BundleID, error) {
	return actions.CreateSupportBundle(&http.Client{RtDetails: rtDetails}, "foo", optionsProvider)
}
