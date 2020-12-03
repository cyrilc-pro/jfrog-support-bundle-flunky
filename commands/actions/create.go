package actions

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"net/http"
)

type createSupportBundleHTTPClient interface {
	GetURL() string
	CreateSupportBundle(payload flunkyhttp.SupportBundleCreationOptions) (int, []byte, error)
}

// OptionsProvider provides options for the creation of a Support Bundle.
type OptionsProvider interface {
	GetOptions(caseNumber CaseNumber) (flunkyhttp.SupportBundleCreationOptions, error)
}

// CreateSupportBundle creates a Support Bundle.
func CreateSupportBundle(httpClient createSupportBundleHTTPClient, caseNumber CaseNumber,
	optionsProvider OptionsProvider) (BundleID, error) {
	log.Debug(fmt.Sprintf("Create Support Bundle %s on %s", caseNumber, httpClient.GetURL()))
	request, err := optionsProvider.GetOptions(caseNumber)
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
	json, err := flunkyhttp.ParseJSON(body)
	if err != nil {
		return "", err
	}
	id, err := json.GetString("id")
	if err != nil {
		return "", err
	}
	return BundleID(id), nil
}
