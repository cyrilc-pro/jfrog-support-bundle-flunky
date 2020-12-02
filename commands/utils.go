package commands

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"time"
)

const (
	HTTPContentType     = "Content-Type"
	HTTPContentTypeJSON = "application/json"
	HTTPContentTypeXML  = "application/xml"
)

type flagValueProvider interface {
	GetStringFlagValue(flagName string) string
	GetBoolFlagValue(flagName string) bool
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
func getTargetDetails(flagProvider flagValueProvider, configProvider serviceHelper,
	conf *SupportBundleCommandConfiguration) (*config.ArtifactoryDetails, error) {
	serverID := flagProvider.GetStringFlagValue(targetServerID)
	if serverID == "" {
		return &config.ArtifactoryDetails{Url: conf.JfrogSupportLogsURL}, nil
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
	return getDurationOrDefault(flagProvider.GetStringFlagValue(downloadTimeout), defaultTimeout)
}

func getPromptOptions(flagProvider flagValueProvider) optionsProvider {
	if flagProvider.GetBoolFlagValue(promptOptions) {
		return newPromptOptionsProvider()
	}
	return newDefaultOptionsProvider()
}

func getRetryInterval(flagProvider flagValueProvider) time.Duration {
	defaultRetryInterval := 5 * time.Second
	return getDurationOrDefault(flagProvider.GetStringFlagValue(retryInterval), defaultRetryInterval)
}

func getDurationOrDefault(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Debug("Error parsing duration: %+v", err)
		log.Warn("Error parsing duration %s, using default %s", value, defaultValue)
		return defaultValue
	}
	return duration
}

type JSONObject map[string]interface{}

func ParseJSON(bytes []byte) (JSONObject, error) {
	parsedResponse := make(JSONObject)
	err := json.Unmarshal(bytes, &parsedResponse)
	return parsedResponse, err
}

func (o JSONObject) GetString(p string) (string, error) {
	v, ok := o[p]
	if !ok {
		return "", fmt.Errorf("property %s not found", p)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("property %s is not a string", p)
	}
	return s, nil
}

func getEndpoint(rtDetails *config.ArtifactoryDetails, endpoint string, args ...interface{}) string {
	return rtDetails.Url + fmt.Sprintf(endpoint, args...)
}
