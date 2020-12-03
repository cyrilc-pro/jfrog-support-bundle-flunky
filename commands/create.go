package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"time"
)

// BundleID is the ID of a Support Bundle.
type BundleID string

type createSupportBundleHTTPClient interface {
	GetURL() string
	CreateSupportBundle(payload SupportBundleCreationOptions) (int, []byte, error)
}

// OptionsProvider provides options for the creation of a Support Bundle.
type OptionsProvider interface {
	GetOptions(caseNumber string) (SupportBundleCreationOptions, error)
}

// DefaultOptionsProvider provides default options for the creation of a Support Bundle.
type DefaultOptionsProvider struct {
	getDate Clock
}

// NewDefaultOptionsProvider creates a new DefaultOptionsProvider
func NewDefaultOptionsProvider() OptionsProvider {
	return &DefaultOptionsProvider{getDate: time.Now}
}

// GetOptions gets the default options.
func (p *DefaultOptionsProvider) GetOptions(caseNumber string) (SupportBundleCreationOptions, error) {
	return SupportBundleCreationOptions{
		Name:        fmt.Sprintf("JFrog Support Case number %s", caseNumber),
		Description: fmt.Sprintf("Generated on %s", toString(p.getDate())),
		Parameters:  nil,
	}, nil
}

// CreateSupportBundle creates a Support Bundle.
func CreateSupportBundle(httpClient createSupportBundleHTTPClient, conf *SupportBundleCommandConfiguration,
	optionsProvider OptionsProvider) (BundleID, error) {
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
