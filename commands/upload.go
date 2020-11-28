package commands

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
)

type uploadHTTPClient interface {
	UploadSupportBundle(sbFilePath string, repoKey string, caseNumber string,
		filename string) (status int, responseBytes []byte, err error)
}

func uploadSupportBundle(client uploadHTTPClient, conf *supportBundleCommandConfiguration, sbFilePath string,
	repoKey string, now Clock) error {
	filename := fmt.Sprintf("%s.zip", toString(now()))
	log.Debug(fmt.Sprintf("Uploading Support Bundle %s to repo %s with filename: %s",
		sbFilePath, repoKey, filename))

	statusCode, respBytes, err := client.UploadSupportBundle(sbFilePath, repoKey, conf.caseNumber, filename)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("Got HTTP response status: %d, body: %s", statusCode, respBytes))
	if statusCode != http.StatusCreated {
		return fmt.Errorf("http request failed with: %d %s", statusCode, http.StatusText(statusCode))
	}
	return nil
}
