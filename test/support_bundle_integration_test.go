package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"strings"
	"testing"
)

func Test_SupportBundleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	tests := []integrationTest{
		{
			Name: "Success with temp file deletion",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				caseNumber := generateCaseNumber(t)
				r, err := runCmdWithFlags(
					[]string{caseNumber},
					rtDetails,
					targetRtDetails,
					"logs",
					true,
				)
				require.NoError(t, err)

				requireSupportBundleCreatedAndUploaded(t, rtDetails, targetRtDetails, r)

				_, err = os.Stat(r.LocalFilePath)
				require.Error(t, err)
				assert.True(t, os.IsNotExist(err))
			},
		},
		{
			Name: "Success without temp file deletion",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				caseNumber := generateCaseNumber(t)
				r, err := runCmdWithFlags(
					[]string{caseNumber},
					rtDetails,
					targetRtDetails,
					"logs",
					false,
				)
				require.NoError(t, err)

				requireSupportBundleCreatedAndUploaded(t, rtDetails, targetRtDetails, r)

				stat, err := os.Stat(r.LocalFilePath)
				require.NoError(t, err)
				assert.Greater(t, stat.Size(), int64(0))
			},
		},
		{
			Name: "Fail because no args",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := runCmd(
					nil,
					nil,
					nil,
				)
				assert.EqualError(t, err, "wrong number of arguments. Expected: 1, Received: 0")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to get RT details",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := runCmd(
					[]string{"1234"},
					nil,
					targetRtDetails,
				)
				assert.EqualError(t, err, "failed to get RT details")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to get target RT details",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := runCmd(
					[]string{"1234"},
					rtDetails,
					nil,
				)
				assert.EqualError(t, err, "failed to get Target RT details")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to create support bundle",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := runCmd(
					[]string{"1234"},
					&config.ArtifactoryDetails{
						Url: "http://rt.invalid",
					},
					targetRtDetails,
				)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "rt.invalid")
				require.NotNil(t, r)
				assert.Empty(t, r.BundleID)
			},
		},
		{
			Name: "Fail to upload support bundle",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := runCmd(
					[]string{"1234"},
					rtDetails,
					&config.ArtifactoryDetails{
						Url: "http://rt.invalid",
					},
				)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "rt.invalid")
				require.NotNil(t, r)
				assert.NotEmpty(t, r.BundleID)
				assert.NotEmpty(t, r.LocalFilePath)
				assert.NotEmpty(t, r.UploadURL)
			},
		},
	}
	runIntegrationTests(t, tests)
}

func generateCaseNumber(t *testing.T) string {
	return strings.ReplaceAll(t.Name(), "/", "_")
}

func runCmd(args []string, rtDetails, targetRtDetails *config.ArtifactoryDetails) (*commands.SupportBundleCmdResult, error) {
	return commands.SupportBundleCmd(
		context.Background(),
		&cliStub{
			arguments:       args,
			rtDetails:       rtDetails,
			targetRtDetails: targetRtDetails,
		})
}

func runCmdWithFlags(args []string, rtDetails, targetRtDetails *config.ArtifactoryDetails, targetRepo string,
	cleanup bool) (*commands.SupportBundleCmdResult, error) {
	return commands.SupportBundleCmd(
		context.Background(),
		&cliStub{
			arguments:       args,
			rtDetails:       rtDetails,
			targetRtDetails: targetRtDetails,
			stringFlags: map[string]string{
				"target-repo": targetRepo,
			},
			boolFlags: map[string]bool{
				"cleanup": cleanup,
			},
		})
}

func requireSupportBundleCreatedAndUploaded(t *testing.T, rtDetails, targetRtDetails *config.ArtifactoryDetails,
	result *commands.SupportBundleCmdResult) {
	require.NotNil(t, result)
	assert.NotEmpty(t, result.BundleID)
	assert.NotEmpty(t, result.LocalFilePath)
	assert.NotEmpty(t, result.UploadURL)

	exists, err := uploadedPathExists(targetRtDetails, result.UploadURL)
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = supportBundleExists(rtDetails, result.BundleID)
	require.NoError(t, err)
	assert.True(t, exists)
}

func uploadedPathExists(rtDetails *config.ArtifactoryDetails, uploadURL string) (bool, error) {
	endpoint := fmt.Sprintf("api/storage/%s", strings.ReplaceAll(uploadURL, rtDetails.Url, ""))
	return testExists(rtDetails, "HEAD", endpoint)
}

func supportBundleExists(rtDetails *config.ArtifactoryDetails, bundleID actions.BundleID) (bool, error) {
	endpoint := fmt.Sprintf("api/system/support/bundle/%s", bundleID)
	return testExists(rtDetails, "GET", endpoint)
}

func testExists(rtDetails *config.ArtifactoryDetails, method, endpoint string) (bool, error) {
	req, err := http.NewRequestWithContext(context.Background(), method, rtDetails.Url+endpoint, nil)
	if err != nil {
		return false, err
	}
	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer func() { _ = res.Body.Close() }()
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		return false, fmt.Errorf("%s %s -> %d", method, endpoint, res.StatusCode)
	}
}

type cliStub struct {
	arguments       []string
	stringFlags     map[string]string
	boolFlags       map[string]bool
	rtDetails       *config.ArtifactoryDetails
	targetRtDetails *config.ArtifactoryDetails
}

func (a *cliStub) GetRtDetails() (*config.ArtifactoryDetails, error) {
	if a.rtDetails == nil {
		return nil, errors.New("failed to get RT details")
	}
	return a.rtDetails, nil
}
func (a *cliStub) GetTargetDetails() (*config.ArtifactoryDetails, error) {
	if a.targetRtDetails == nil {
		return nil, errors.New("failed to get Target RT details")
	}
	return a.targetRtDetails, nil
}
func (a *cliStub) GetArguments() []string {
	return a.arguments
}
func (a *cliStub) GetStringFlagValue(flagName string) string {
	return a.stringFlags[flagName]
}
func (a *cliStub) GetBoolFlagValue(flagName string) bool {
	return a.boolFlags[flagName]
}
