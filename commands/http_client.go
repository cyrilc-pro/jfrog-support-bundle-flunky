package commands

import (
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
)

type HTTPClient struct {
	rtDetails *config.ArtifactoryDetails
}

func (c *HTTPClient) GetURL() string {
	return c.rtDetails.Url
}

func (c *HTTPClient) CreateSupportBundle(requestPayload string) (status int, responseBytes []byte, err error) {
	servicesManager, err := utils.CreateServiceManager(c.rtDetails, false)
	if err != nil {
		return -1, nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	httpClientDetails.Headers[httpContentType] = httpContentTypeJSON
	response, bytes, err := servicesManager.Client().SendPost(getEndpoint(c.rtDetails, "api/system/support/bundle"),
		[]byte(requestPayload), &httpClientDetails)
	if err != nil {
		return -1, nil, err
	}
	return response.StatusCode, bytes, err
}
