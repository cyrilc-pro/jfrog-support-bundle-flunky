package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func setUpLicense(ctx context.Context, l logger, licenseKey string, rtDetails *config.ArtifactoryDetails) error {
	err := deployTestLicense(ctx, l, licenseKey, rtDetails)
	if err != nil {
		return err
	}
	return waitForLicenseDeployed(ctx, l, rtDetails)
}

func waitForLicenseDeployed(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	req, err := http.NewRequestWithContext(ctx, "GET", getLicensesEndpointURL(rtDetails), nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	retry := true
	for retry {
		select {
		case <-ctx.Done():
			return errors.New("timed out waiting for license to be applied")
		case <-ticker.C:
			licenseType, err2 := getLicenseType(req, l)
			if err2 != nil {
				return err2
			}

			if licenseType != "N/A" {
				l.Logf("License %v applied", licenseType)
				retry = false
			}
		}
	}
	return nil
}

func getLicenseType(req *http.Request, l logger) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	l.Logf("Get license: %s %s", resp.Status, string(bytes))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("license check failed: %d", resp.StatusCode)
	}

	json, err := flunkyhttp.ParseJSON(bytes)
	if err != nil {
		return "", err
	}

	licenseType, err := json.GetString("type")
	if err != nil {
		return "", err
	}
	return licenseType, nil
}

func deployTestLicense(ctx context.Context, l logger, licenseKey string, rtDetails *config.ArtifactoryDetails) error {
	licensePayload := strings.NewReader(fmt.Sprintf(`{"licenseKey":"%s"}`, licenseKey))
	req, err := http.NewRequestWithContext(ctx, "POST", getLicensesEndpointURL(rtDetails), licensePayload)
	if err != nil {
		return err
	}

	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	req.Header[flunkyhttp.HTTPContentType] = []string{flunkyhttp.HTTPContentTypeJSON}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	if err != nil {
		return err
	}

	// DO NOT PRINT RESPONSE BODY: it may contain the license key in clear-text
	l.Logf("Deploy license: %s", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("license deploy failed: %d", resp.StatusCode)
	}

	return nil
}

func getLicensesEndpointURL(rtDetails *config.ArtifactoryDetails) string {
	return fmt.Sprintf("%sapi/system/licenses", rtDetails.Url)
}
