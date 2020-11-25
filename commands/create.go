package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"time"
)

func createSupportBundle(rtDetails *config.ArtifactoryDetails, conf *supportBundleCommandConfiguration) (creationResponse, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.caseNumber, rtDetails.Url))
	response, body, err := sendHttpRequest(rtDetails, conf)
	if err != nil {
		return creationResponse{}, err
	}
	log.Debug(fmt.Sprintf("Got %s\n%s", response.Status, string(body)))
	if response.StatusCode != http.StatusOK {
		return creationResponse{}, fmt.Errorf("http request failed with: %s", response.Status)
	}
	parsedResponse := make(map[string]interface{})
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		return creationResponse{}, err
	}
	untypedId, found := parsedResponse["id"]
	if !found {
		return creationResponse{}, errors.New(`unexpected JSON response (missing "id" property)`)
	}
	id, ok := untypedId.(string)
	if !ok {
		return creationResponse{}, errors.New(`unexpected JSON response ("id" property is not a string)`)
	}
	return creationResponse{Id: id}, nil
}

func sendHttpRequest(rtDetails *config.ArtifactoryDetails, conf *supportBundleCommandConfiguration) (*http.Response, []byte, error) {
	servicesManager, err := utils.CreateServiceManager(rtDetails, false)
	if err != nil {
		return nil, nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	httpClientDetails.Headers["Content-Type"] = "application/json"
	request := fmt.Sprintf(`{"name": "JFrog Support Case number %s","description": "Generated on %s","parameters":{}}`,
		conf.caseNumber,
		time.Now().Format(time.RFC3339))
	return servicesManager.Client().SendPost(fmt.Sprintf("%sapi/system/support/bundle", rtDetails.Url), []byte(request), &httpClientDetails)
}
