package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type downloadSupportBundleHTTPClient interface {
	GetURL() string
	DownloadSupportBundle(bundleID bundleID) (*http.Response, error)
	GetSupportBundleStatus(bundleID bundleID) (int, []byte, error)
}

func downloadSupportBundle(ctx context.Context, client downloadSupportBundleHTTPClient, timeout time.Duration,
	retryInterval time.Duration, bundleID bundleID) (string, error) {
	log.Debug(fmt.Sprintf("Download Support Bundle %s from %s", bundleID, client.GetURL()))

	err := waitUntilSupportBundleIsReady(ctx, client, retryInterval, timeout, bundleID)
	if err != nil {
		return "", err
	}

	dirPath, err := fileutils.CreateTempDir()
	if err != nil {
		return "", err
	}
	tmpFilePath := filepath.Join(dirPath, fmt.Sprintf("%s.zip", bundleID))
	tmpZipFile, err := os.Create(tmpFilePath)
	if err != nil {
		return "", err
	}
	defer handleClose(tmpZipFile)

	err = downloadSupportBundleAndWriteToFile(client, tmpZipFile, bundleID)
	if err != nil {
		return "", err
	}

	log.Debug(fmt.Sprintf("Downloaded Support Bundle to %s", tmpFilePath))
	return tmpFilePath, nil
}

func downloadSupportBundleAndWriteToFile(client downloadSupportBundleHTTPClient, tmpZipFile *os.File, bundleID bundleID) error {
	resp, err := client.DownloadSupportBundle(bundleID)
	if err != nil {
		return err
	}
	defer handleClose(resp.Body)
	log.Debug(fmt.Sprintf("Got %d", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed with: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	_, err = io.Copy(tmpZipFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func waitUntilSupportBundleIsReady(ctx context.Context, client downloadSupportBundleHTTPClient,
	retryInterval time.Duration, timeout time.Duration, bundleID bundleID) error {
	ctxWithTimeout, cancelCtx := context.WithTimeout(ctx, timeout)
	defer cancelCtx()
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctxWithTimeout.Done():
			return errors.New("timeout waiting for support bundle to be ready")
		case <-ticker.C:
			log.Debug(fmt.Sprintf("Attempting to get status for support bundle %s", bundleID))
			statusCode, body, err := client.GetSupportBundleStatus(bundleID)
			if err != nil {
				return err
			}

			log.Debug(fmt.Sprintf("Got HTTP response status: %d", statusCode))
			if statusCode != http.StatusOK {
				return fmt.Errorf("http request failed with: %d %s", statusCode, http.StatusText(statusCode))
			}

			parsedBody, err := parseJSON(body)
			if err != nil {
				return err
			}

			sbStatus, err := parsedBody.getString("status")
			if err != nil {
				return err
			}

			log.Debug(fmt.Sprintf("Support bundle status: %s", sbStatus))
			if sbStatus != "in progress" {
				return nil
			}
		}
	}
}

func handleClose(closer io.Closer) {
	if closer != nil {
		err := closer.Close()
		if err != nil {
			log.Warn("error occurred while closing: %+v", err)
		}
	}
}
