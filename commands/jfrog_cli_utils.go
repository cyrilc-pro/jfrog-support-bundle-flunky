package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"time"
)

type flagValueProvider interface {
	GetStringFlagValue(flagName string) string
	GetBoolFlagValue(flagName string) bool
}

type argumentsProvider interface {
	GetArguments() []string
}

type artifactoryDetailsProvider interface {
	GetRtDetails() (*config.ArtifactoryDetails, error)
	GetTargetDetails() (*config.ArtifactoryDetails, error)
}

type serviceHelper interface {
	GetConfig(serverID string, excludeRefreshableTokens bool) (*config.ArtifactoryDetails, error)
	CreateInitialRefreshableTokensIfNeeded(artifactoryDetails *config.ArtifactoryDetails) error
}

// Returns the Artifactory Details of the provided server-id, or the default one.
func getRtDetails(flagProvider flagValueProvider, configHelper serviceHelper) (*config.ArtifactoryDetails, error) {
	serverID := flagProvider.GetStringFlagValue(serverIDFlag)
	return buildRtDetailsFromServerID(serverID, configHelper)
}

// Returns the Artifactory Details of the target-server-id, or JFrog support logs configured ArtifactoryDetails.
func getTargetDetails(flagProvider flagValueProvider, configProvider serviceHelper) (*config.ArtifactoryDetails, error) {
	serverID := flagProvider.GetStringFlagValue(targetServerIDFlag)
	if serverID == "" {
		// TODO change this when everything works correctly "https://supportlogs.jfrog.com/" (keep the trailing slash!)
		return &config.ArtifactoryDetails{Url: "https://supportlogs.jfrog.com.invalid/"}, nil
	}
	details, err := buildRtDetailsFromServerID(serverID, configProvider)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func buildRtDetailsFromServerID(serverID string, configHelper serviceHelper) (*config.ArtifactoryDetails, error) {
	details, err := configHelper.GetConfig(serverID, false)
	if err != nil {
		return nil, err
	}
	details.Url = clientutils.AddTrailingSlashIfNeeded(details.Url)
	err = configHelper.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func getTimeout(flagProvider flagValueProvider) time.Duration {
	defaultTimeout := 10 * time.Minute
	return getDurationOrDefault(flagProvider.GetStringFlagValue(downloadTimeoutFlag), defaultTimeout)
}

func shouldCleanup(flagProvider flagValueProvider) bool {
	return flagProvider.GetBoolFlagValue(cleanupFlag)
}

func getTargetRepo(flagProvider flagValueProvider) string {
	return flagProvider.GetStringFlagValue(targetRepoFlag)
}

func getPromptOptions(flagProvider flagValueProvider) actions.OptionsProvider {
	if flagProvider.GetBoolFlagValue(promptOptionsFlag) {
		return actions.NewPromptOptionsProvider()
	}
	return actions.NewDefaultOptionsProvider()
}

func getRetryInterval(flagProvider flagValueProvider) time.Duration {
	defaultRetryInterval := 5 * time.Second
	return getDurationOrDefault(flagProvider.GetStringFlagValue(retryIntervalFlag), defaultRetryInterval)
}

func getDurationOrDefault(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Debug(fmt.Sprintf("Error parsing duration: %+v", err))
		log.Warn(fmt.Sprintf("Error parsing duration %s, using default %s", value, defaultValue))
		return defaultValue
	}
	return duration
}
