package test

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"io/ioutil"
	"net/http"
	"strings"
)

func newHTTPRequestWithBody(ctx context.Context, rtDetails *config.ArtifactoryDetails,
	method, endpoint, contentType, body string) (*http.Request, error) {
	url := rtDetails.Url + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	if err != nil {
		return req, err
	}
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	req.Header[flunkyhttp.HTTPContentType] = []string{contentType}
	return req, nil
}

func newHTTPGETRequest(ctx context.Context, rtDetails *config.ArtifactoryDetails, endpoint string) (*http.Request, error) {
	url := rtDetails.Url + endpoint
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return req, err
	}
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	return req, nil
}

func do(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s %s: %w", req.Method, req.URL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return bytes, fmt.Errorf("%s %s: %w", req.Method, req.URL, err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("%s %s: %d", req.Method, req.URL, resp.StatusCode)
	}

	return bytes, err
}
