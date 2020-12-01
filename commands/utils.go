package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"time"
)

const (
	httpContentType     = "Content-Type"
	httpContentTypeJSON = "application/json"
)

// Returns the Artifactory Details of the provided server-id, or the default one.
func getRtDetails(c *components.Context) (*config.ArtifactoryDetails, error) {
	details, err := commands.GetConfig(c.GetStringFlagValue(serverID), false)
	if err != nil {
		return nil, err
	}
	if details.Url == "" {
		return nil, errors.New("no server-id was found, or the server-id has no url")
	}
	details.Url = clientutils.AddTrailingSlashIfNeeded(details.Url)
	err = config.CreateInitialRefreshableTokensIfNeeded(details)
	if err != nil {
		return nil, err
	}
	return details, nil
}

func getTimeout(c *components.Context) time.Duration {
	defaultTimeout := 10 * time.Minute
	return getDurationOrDefault(c.GetStringFlagValue(downloadTimeout), defaultTimeout)
}

func getPromptOptions(c *components.Context) optionsProvider {
	var p optionsProvider = &defaultOptionsProvider{getDate: time.Now}
	if c.GetBoolFlagValue(promptOptions) {
		p = &promptOptionsProvider{getDate: time.Now}
	}
	return p
}

func getRetryInterval(c *components.Context) time.Duration {
	defaultRetryInterval := 5 * time.Second
	return getDurationOrDefault(c.GetStringFlagValue(retryInterval), defaultRetryInterval)
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

func parseJSON(bytes []byte) (JSONObject, error) {
	parsedResponse := make(JSONObject)
	err := json.Unmarshal(bytes, &parsedResponse)
	return parsedResponse, err
}

func (o JSONObject) getString(p string) (string, error) {
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
