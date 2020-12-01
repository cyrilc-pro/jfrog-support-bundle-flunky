package commands

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const undefinedStatusCode = -1

type HTTPClient struct {
	rtDetails *config.ArtifactoryDetails
}

func (c *HTTPClient) GetURL() string {
	return c.rtDetails.Url
}

// nolint: bodyclose
func (c *HTTPClient) CreateSupportBundle(options SupportBundleCreationOptions) (status int, responseBytes []byte, err error) {
	servicesManager, err := utils.CreateServiceManager(c.rtDetails, false)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	httpClientDetails.Headers[httpContentType] = httpContentTypeJSON
	payload, err := json.Marshal(options)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	log.Debug(fmt.Sprintf("Sending %s", payload))
	response, bytes, err := servicesManager.Client().SendPost(getEndpoint(c.rtDetails, "api/system/support/bundle"),
		payload, &httpClientDetails)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	return response.StatusCode, bytes, nil
}

// This returns the support bundle in the response.Body. Closing the body is the caller's responsibility.
func (c *HTTPClient) DownloadSupportBundle(bundleID bundleID) (*http.Response, error) {
	servicesManager, err := utils.CreateServiceManager(c.rtDetails, false)
	if err != nil {
		return nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	downloadSbURL := fmt.Sprintf("%sapi/system/support/bundle/%s/archive", c.GetURL(), bundleID)
	resp, _, _, err := servicesManager.Client().Send("GET", downloadSbURL, nil, true, false, &httpClientDetails)
	return resp, err
}

// nolint: bodyclose
func (c *HTTPClient) GetSupportBundleStatus(bundleID bundleID) (status int, responseBytes []byte, err error) {
	servicesManager, err := utils.CreateServiceManager(c.rtDetails, false)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	sbStatusURL := fmt.Sprintf("%sapi/system/support/bundle/%s", c.GetURL(), bundleID)
	resp, responseBytes, _, err := servicesManager.Client().SendGet(sbStatusURL, true, &httpClientDetails)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	return resp.StatusCode, responseBytes, nil
}
