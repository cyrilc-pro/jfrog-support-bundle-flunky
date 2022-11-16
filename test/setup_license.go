package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	flunkyhttp "github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"net/http"
	"time"
)

func setUpLicense(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails, licenseKey string) error {
	err := deployTestLicense(ctx, licenseKey, rtDetails)
	if err != nil {
		return err
	}
	return waitForLicenseDeployed(ctx, l, rtDetails)
}

func waitForLicenseDeployed(ctx context.Context, l logger, rtDetails *config.ArtifactoryDetails) error {
	req, err := newHTTPGETRequest(ctx, rtDetails, "api/system/licenses")
	if err != nil {
		return err
	}
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
	bytes, err := do(req)
	if err != nil {
		return "", fmt.Errorf("license check failed: %w", err)
	}

	l.Logf("Get license: %s", string(bytes))

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

func deployTestLicense(ctx context.Context, licenseKey string, rtDetails *config.ArtifactoryDetails) error {
	licensePayload := fmt.Sprintf(`{"licenseKey":"%s"}`, licenseKey)
	req, err := newHTTPRequestWithBody(ctx, rtDetails, "POST", "api/system/licenses",
		flunkyhttp.HTTPContentTypeJSON, licensePayload)
	if err != nil {
		return err
	}

	_, err = do(req)
	if err != nil {
		return fmt.Errorf("license deploy failed: %w", err)
	}

	// DO NOT PRINT RESPONSE BODY: it may contain the license key in clear-text

	return nil
}
