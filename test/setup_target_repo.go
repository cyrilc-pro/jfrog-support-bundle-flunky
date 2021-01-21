package test

import (
	"context"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
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
	body, err := doPutJSON(ctx, rtDetails, "api/v2/security/permissions/logsPerm", payload)
	l.Logf("Create anonymous permission on logs repository response: %s", body)
	return err
}

func setAnonAccess(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	endpoint := "api/system/configuration"
	getRequest, err := newHTTPGETRequest(ctx, rtDetails, endpoint)
	if err != nil {
		return err
	}

	getRespBody, err := do(getRequest)
	if err != nil {
		return err
	}

	payload := strings.Replace(string(getRespBody), "<anonAccessEnabled>false</anonAccessEnabled>",
		"<anonAccessEnabled>true</anonAccessEnabled>", 1)
	postRequest, err := newHTTPRequestWithBody(ctx, rtDetails, "POST", endpoint, flunkyhttp.HTTPContentTypeXML, payload)
	if err != nil {
		return err
	}

	postResponseBody, err := do(postRequest)
	l.Logf("Set anonymous access response: %s", postResponseBody)
	return err
}

func createLogRepository(ctx context.Context, l logger, targetRtDetails *config.ArtifactoryDetails) error {
	payload := `{"key": "logs","rclass": "local","packageType": "generic"}`
	respBody, err := doPutJSON(ctx, targetRtDetails, "api/repositories/logs", payload)
	l.Logf("Create logs repository response: %s", respBody)
	return err
}

func doPutJSON(ctx context.Context, targetRtDetails *config.ArtifactoryDetails, endpoint, payload string) ([]byte, error) {
	req, err := newHTTPRequestWithBody(ctx, targetRtDetails, "PUT", endpoint,
		flunkyhttp.HTTPContentTypeJSON, payload)
	if err != nil {
		return nil, err
	}
	return do(req)
}
