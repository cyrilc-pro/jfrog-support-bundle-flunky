package http

import (
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/artifactory/utils"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/utils/io/httputils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

const (
	// HTTPContentType is the HTTP header name for Content-Type
	HTTPContentType = "Content-Type"
	// HTTPContentTypeJSON is the header value for JSON Content-Type
	HTTPContentTypeJSON = "application/json"
	// HTTPContentTypeXML is the header value for XML Content-Type
	HTTPContentTypeXML  = "application/xml"
	undefinedStatusCode = -1
)

// Client is a facade for interacting with a JFrog Artifactory service through REST calls.
type Client struct {
	RtDetails *config.ArtifactoryDetails
}

// GetURL gives the URL of the JFrog Artifactory service
func (c *Client) GetURL() string {
	return c.RtDetails.Url
}

// CreateSupportBundle creates a Support Bundle.
// nolint: bodyclose // Body is closed by ArtifactoryHttpClient
func (c *Client) CreateSupportBundle(options SupportBundleCreationOptions) (status int, responseBytes []byte, err error) {
	servicesManager, httpClientDetails, err := c.createArtifactoryServicesManager()
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	httpClientDetails.Headers[HTTPContentType] = HTTPContentTypeJSON
	payload, err := json.Marshal(options)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	log.Debug(fmt.Sprintf("Sending %s", payload))
	response, bytes, err := servicesManager.Client().SendPost(fmt.Sprintf("%sapi/system/support/bundle", c.GetURL()),
		payload, &httpClientDetails)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	return response.StatusCode, bytes, nil
}

// DownloadSupportBundle downloads a Support Bundle. This returns the support bundle in the response.Body.
// Closing the body is the caller's responsibility.
func (c *Client) DownloadSupportBundle(bundleID string) (*http.Response, error) {
	servicesManager, httpClientDetails, err := c.createArtifactoryServicesManager()
	if err != nil {
		return nil, err
	}
	downloadSbURL := fmt.Sprintf("%sapi/system/support/bundle/%s/archive", c.GetURL(), bundleID)
	resp, _, _, err := servicesManager.Client().Send("GET", downloadSbURL, nil, true, false, &httpClientDetails)
	return resp, err
}

// GetSupportBundleStatus gets the status of a Support Bundle creation process.
// nolint: bodyclose // Body is closed by ArtifactoryHttpClient
func (c *Client) GetSupportBundleStatus(bundleID string) (status int, responseBytes []byte, err error) {
	servicesManager, httpClientDetails, err := c.createArtifactoryServicesManager()
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	sbStatusURL := fmt.Sprintf("%sapi/system/support/bundle/%s", c.GetURL(), bundleID)
	resp, responseBytes, _, err := servicesManager.Client().SendGet(sbStatusURL, true, &httpClientDetails)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	return resp.StatusCode, responseBytes, nil
}

// UploadSupportBundle uploads a Support Bundle.
// nolint: bodyclose // Body is closed by ArtifactoryHttpClient
func (c *Client) UploadSupportBundle(sbFilePath string, repoKey string, supportCaseDirectory string,
	filename string) (status int, responseBytes []byte, err error) {
	// TODO add flag for number of retries
	const retries = 5
	servicesManager, httpClientDetails, err := c.createArtifactoryServicesManager()
	if err != nil {
		return undefinedStatusCode, nil, err
	}

	url := fmt.Sprintf("%s%s/%s/%s;uploadedBy=support-bundle-flunky", c.RtDetails.Url, repoKey, supportCaseDirectory,
		filename)
	resp, body, err := servicesManager.Client().UploadFile(sbFilePath, url, "",
		&httpClientDetails, retries, nil)
	if err != nil {
		return undefinedStatusCode, nil, err
	}
	return resp.StatusCode, body, err
}

func (c *Client) createArtifactoryServicesManager() (artifactory.ArtifactoryServicesManager,
	httputils.HttpClientDetails, error) {
	servicesManager, err := utils.CreateServiceManager(c.RtDetails, false)
	if err != nil {
		return nil, httputils.HttpClientDetails{}, err
	}
	httpClientDetails := servicesManager.GetConfig().GetServiceDetails().CreateHttpClientDetails()
	return servicesManager, httpClientDetails, nil
}
