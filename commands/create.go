package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"time"
)

type bundleID string

type createSupportBundleHTTPClient interface {
	GetURL() string
	CreateSupportBundle(payload SupportBundleCreationOptions) (int, []byte, error)
}

type optionsProvider interface {
	GetOptions(caseNumber string) (SupportBundleCreationOptions, error)
}

type defaultOptionsProvider struct {
	getDate Clock
}

func newDefaultOptionsProvider() optionsProvider {
	return &defaultOptionsProvider{getDate: time.Now}
}

func (p *defaultOptionsProvider) GetOptions(caseNumber string) (SupportBundleCreationOptions, error) {
	return SupportBundleCreationOptions{
		Name:        fmt.Sprintf("JFrog Support Case number %s", caseNumber),
		Description: fmt.Sprintf("Generated on %s", toString(p.getDate())),
		Parameters:  nil,
	}, nil
}

func createSupportBundle(httpClient createSupportBundleHTTPClient, conf *supportBundleCommandConfiguration,
	optionsProvider optionsProvider) (bundleID, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.caseNumber, httpClient.GetURL()))
	request, err := optionsProvider.GetOptions(conf.caseNumber)
	if err != nil {
		return "", err
	}
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
