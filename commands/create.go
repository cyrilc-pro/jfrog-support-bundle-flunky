package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type createSupportBundleHttpClient interface {
	GetUrl() string
	CreateSupportBundle(payload string) (int, []byte, error)
}

func createSupportBundle(httpClient createSupportBundleHttpClient, conf *supportBundleCommandConfiguration, now Clock) (creationResponse, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.caseNumber, httpClient.GetUrl()))
	request := fmt.Sprintf(`{"name": "JFrog Support Case number %s","description": "Generated on %s","parameters":{}}`,
		conf.caseNumber,
		now())

	responseStatus, body, err := httpClient.CreateSupportBundle(request)
	if err != nil {
		return creationResponse{}, err
	}
	log.Debug(fmt.Sprintf("Got %d\n%s", responseStatus, string(body)))
	if responseStatus != http.StatusOK {
		return creationResponse{}, fmt.Errorf("http request failed with: %d", responseStatus)
	}
	json, err := parseJson(body)
	if err != nil {
		return creationResponse{}, err
	}
	id, err := json.getString("id")
	if err != nil {
		return creationResponse{}, err
	}
	return creationResponse{Id: id}, nil
}
