package actions

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"net/http"
	"strings"
)

type uploadHTTPClient interface {
	UploadSupportBundle(sbFilePath string, repoKey string, supportCaseDirectory string,
		filename string) (status int, responseBytes []byte, err error)
}

// UploadSupportBundle uploads a Support Bundle.
func UploadSupportBundle(client uploadHTTPClient, caseNumber CaseNumber, sbFilePath string,
	repoKey string, now Clock) error {
	filename := fmt.Sprintf("%s.zip", strings.ReplaceAll(formattedString(now()), ":", "_"))
	log.Debug(fmt.Sprintf("Uploading Support Bundle %s to repo %s with filename: %s",
		sbFilePath, repoKey, filename))

	statusCode, respBytes, err := client.UploadSupportBundle(sbFilePath, repoKey, string(caseNumber), filename)
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("Got HTTP response status: %d, body: %s", statusCode, respBytes))
	if statusCode != http.StatusCreated {
		return fmt.Errorf("http request failed with: %d %s", statusCode, http.StatusText(statusCode))
	}
	return nil
}
