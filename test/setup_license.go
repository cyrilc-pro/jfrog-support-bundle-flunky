package test

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func setUpLicense(ctx context.Context, t *testing.T, licenseKey string, rtDetails *config.ArtifactoryDetails) {
	t.Helper()
	deployTestLicense(ctx, t, licenseKey, rtDetails)
	waitForLicenseDeployed(ctx, t, rtDetails)
}

func waitForLicenseDeployed(ctx context.Context, t *testing.T, rtDetails *config.ArtifactoryDetails) {
	req, err := http.NewRequestWithContext(ctx, "GET", getLicensesEndpointURL(rtDetails), nil)
	require.NoError(t, err)
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	retry := true
	for retry {
		select {
		case <-ctx.Done():
			require.Fail(t, "Timed out waiting for license to be applied")
		case <-ticker.C:
			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode, "License check failed")
			bytes, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			_ = resp.Body.Close()
			t.Logf("Get license: %s %s", resp.Status, string(bytes))
			json, err := commands.ParseJSON(bytes)
			require.NoError(t, err)
			licenseType, err := json.GetString("type")
			require.NoError(t, err)
			if licenseType != "N/A" {
				t.Logf("License %v applied", licenseType)
				retry = false
			}
		}
	}
}

func deployTestLicense(ctx context.Context, t *testing.T, licenseKey string, rtDetails *config.ArtifactoryDetails) {
	licensePayload := strings.NewReader(fmt.Sprintf(`{"licenseKey":"%s"}`, licenseKey))
	req, err := http.NewRequestWithContext(ctx, "POST", getLicensesEndpointURL(rtDetails), licensePayload)
	require.NoError(t, err)
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	req.Header[commands.HTTPContentType] = []string{commands.HTTPContentTypeJSON}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	_, err = ioutil.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)
	// DO NOT PRINT RESPONSE BODY: it may contain the license key in clear-text
	t.Logf("Deploy license: %s", resp.Status)
	require.Equal(t, http.StatusOK, resp.StatusCode, "License deploy failed")
}

func getLicensesEndpointURL(rtDetails *config.ArtifactoryDetails) string {
	return fmt.Sprintf("%sapi/system/licenses", rtDetails.Url)
}
