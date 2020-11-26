package commands

import (
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
)

type HttpClient struct {
	rtDetails *config.ArtifactoryDetails
}

func (c *HttpClient) GetUrl() string {
	return c.rtDetails.Url
}

func (c *HttpClient) CreateSupportBundle(requestPayload string) (int, []byte, error) {
	servicesManager, err := utils.CreateServiceManager(c.rtDetails, false)
	if err != nil {
		return -1, nil, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	httpClientDetails.Headers[httpContentType] = httpJsonContentTypeJson
	response, bytes, err := servicesManager.Client().SendPost(getEndpoint(c.rtDetails, "api/system/support/bundle"), []byte(requestPayload), &httpClientDetails)
	if err != nil {
		return -1, nil, err
	}
	return response.StatusCode, bytes, err
}
