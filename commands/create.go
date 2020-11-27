package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type bundleID string

type createSupportBundleHTTPClient interface {
	GetURL() string
	CreateSupportBundle(payload string) (int, []byte, error)
}

func createSupportBundle(httpClient createSupportBundleHTTPClient, conf *supportBundleCommandConfiguration,
	now Clock) (bundleID, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.caseNumber, httpClient.GetURL()))
	request := fmt.Sprintf(`{"name": "JFrog Support Case number %s","description": "Generated on %s","parameters":{}}`,
		conf.caseNumber,
		now())

	responseStatus, body, err := httpClient.CreateSupportBundle(request)
	if err != nil {
		return "", err
	}
	log.Debug(fmt.Sprintf("Got %d\n%s", responseStatus, string(body)))
	if responseStatus != http.StatusOK {
		return "", fmt.Errorf("http request failed with: %d", responseStatus)
	}
	json, err := parseJSON(body)
	if err != nil {
		return "", err
	}
	id, err := json.getString("id")
	if err != nil {
		return "", err
	}
	return bundleID(id), nil
}
