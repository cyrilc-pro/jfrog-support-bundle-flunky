package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"time"
)

type BundleID string

type createSupportBundleHTTPClient interface {
	GetURL() string
	CreateSupportBundle(payload SupportBundleCreationOptions) (int, []byte, error)
}

type optionsProvider interface {
	GetOptions(caseNumber string) (SupportBundleCreationOptions, error)
}

type DefaultOptionsProvider struct {
	GetDate Clock
}

func newDefaultOptionsProvider() optionsProvider {
	return &DefaultOptionsProvider{GetDate: time.Now}
}

func (p *DefaultOptionsProvider) GetOptions(caseNumber string) (SupportBundleCreationOptions, error) {
	return SupportBundleCreationOptions{
		Name:        fmt.Sprintf("JFrog Support Case number %s", caseNumber),
		Description: fmt.Sprintf("Generated on %s", toString(p.GetDate())),
		Parameters:  nil,
	}, nil
}

func CreateSupportBundle(httpClient createSupportBundleHTTPClient, conf *SupportBundleCommandConfiguration,
	optionsProvider optionsProvider) (BundleID, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", conf.CaseNumber, httpClient.GetURL()))
	request, err := optionsProvider.GetOptions(conf.CaseNumber)
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
	json, err := ParseJSON(body)
	if err != nil {
		return "", err
	}
	id, err := json.GetString("id")
	if err != nil {
		return "", err
	}
	return BundleID(id), nil
}
