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

func setUpTargetRepositoryAndPermissions(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	err := createLogRepository(ctx, l, rtDetails)
	if err != nil {
		return err
	}
	err = setAnonAccess(ctx, l, rtDetails)
	if err != nil {
		return err
	}
	return createAnonymousPermission(ctx, l, rtDetails)
}

func createAnonymousPermission(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	payload := `{"name":"logsPerm","repo":{"repositories":["logs"],"actions":{"users":{"anonymous":["write"]}}}}`
	url := fmt.Sprintf("%sapi/v2/security/permissions/logsPerm", rtDetails.Url)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	req.Header[flunkyhttp.HTTPContentType] = []string{flunkyhttp.HTTPContentTypeJSON}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()
	l.Logf("Create anonymous permission on logs repository: status %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add anonymous permission failed: %d", resp.StatusCode)
	}
	return nil
}

func setAnonAccess(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	url := fmt.Sprintf("%sapi/system/configuration", rtDetails.Url)
	getRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	getRequest.SetBasicAuth(rtDetails.User, rtDetails.Password)
	getRequest.Header[flunkyhttp.HTTPContentType] = []string{flunkyhttp.HTTPContentTypeXML}
	getResp, err := http.DefaultClient.Do(getRequest)
	if err != nil {
		return err
	}
	defer func() { _ = getResp.Body.Close() }()

	getRespBody, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		return err
	}

	payload := strings.Replace(string(getRespBody), "<anonAccessEnabled>false</anonAccessEnabled>",
		"<anonAccessEnabled>true</anonAccessEnabled>", 1)
	postRequest, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	postRequest.SetBasicAuth(rtDetails.User, rtDetails.Password)
	postRequest.Header[flunkyhttp.HTTPContentType] = []string{flunkyhttp.HTTPContentTypeXML}
	postResponse, err := http.DefaultClient.Do(postRequest)
	if err != nil {
		return err
	}
	defer func() { _ = postResponse.Body.Close() }()

	postResponseBody, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		return err
	}

	l.Logf("Set anonymous access: status %s, response: %s", getResp.Status, postResponseBody)
	if getResp.StatusCode != http.StatusOK {
		return fmt.Errorf("set anonymous access failed: %d", getResp.StatusCode)
	}
	return nil
}

func createLogRepository(ctx context.Context, l logger, targetRtDetails *config.ArtifactoryDetails) error {
	payload := `{"key": "logs","rclass": "local","packageType": "generic"}`
	url := fmt.Sprintf("%sapi/repositories/logs", targetRtDetails.Url)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.SetBasicAuth(targetRtDetails.User, targetRtDetails.Password)
	req.Header[flunkyhttp.HTTPContentType] = []string{flunkyhttp.HTTPContentTypeJSON}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	l.Logf("Create logs repository: status %s, response: %s", resp.Status, respBody)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("create logs repo failed: %d", resp.StatusCode)
	}
	return nil
}
