package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/commands"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	clientutils "github.com/jfrog/jfrog-client-go/utils"
)

const (
	httpContentType     = "Content-Type"
	httpContentTypeJSON = "application/json"
)

// Returns the Artifactory Details of the provided server-id, or the default one.
func getRtDetails(c *components.Context) (*config.ArtifactoryDetails, error) {
	details, err := commands.GetConfig(c.GetStringFlagValue("server-id"), false)
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
