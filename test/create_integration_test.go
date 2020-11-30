package test

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func (s *IntegrationTestSuite) Test_CreateSupportBundle() {
	s.T().Run("Success", func(t *testing.T) {
		conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
		id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: s.rtDetails}, &conf,
			&commands.DefaultOptionsProvider{GetDate: time.Now})
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})
	s.T().Run("Offline", func(t *testing.T) {
		conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
		_, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: &config.ArtifactoryDetails{
			Url: "http://unknown.invalid/",
		}}, &conf, &commands.DefaultOptionsProvider{GetDate: time.Now})
		require.Error(t, err)
		// exact message depends on OS
		require.Contains(t, err.Error(), "dial tcp:")
	})

	s.T().Run("Success with all options disabled", func(t *testing.T) {
		conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
		id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: s.rtDetails}, &conf,
			&commands.PromptOptionsProvider{GetDate: time.Now, Prompter: &commands.PrompterStub{
				IncludeLogs:          false,
				IncludeSystem:        false,
				IncludeConfiguration: false,
				IncludeThreadDump:    false,
			}})
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})

	s.T().Run("Success with all options enabled", func(t *testing.T) {
		conf := commands.SupportBundleCommandConfiguration{CaseNumber: "foo"}
		id, err := commands.CreateSupportBundle(&commands.HTTPClient{RtDetails: s.rtDetails}, &conf,
			&commands.PromptOptionsProvider{GetDate: time.Now, Prompter: &commands.PrompterStub{
				IncludeLogs:          true,
				IncludeSystem:        true,
				IncludeConfiguration: true,
				IncludeThreadDump:    true,
			}})
		require.NoError(t, err)
		require.NotEmpty(t, id)
	})
}
