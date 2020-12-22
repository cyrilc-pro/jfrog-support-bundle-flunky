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
				caseNumber := strings.ReplaceAll(t.Name(), "/", "_")
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments: []string{caseNumber},
						stringFlags: map[string]string{
							"target-repo": "logs",
						},
						boolFlags: map[string]bool{
							"cleanup": true,
						},
						rtDetails:       rtDetails,
						targetRtDetails: targetRtDetails,
					})
				require.NoError(t, err)

				require.NotNil(t, r)
				assert.NotEmpty(t, r.BundleID)
				assert.NotEmpty(t, r.LocalFilePath)
				assert.NotEmpty(t, r.UploadURL)

				exists, err := uploadedPathExists(targetRtDetails, r.UploadURL)
				require.NoError(t, err)
				assert.True(t, exists)

				exists, err = supportBundleExists(rtDetails, r.BundleID)
				require.NoError(t, err)
				assert.True(t, exists)

				_, err = os.Stat(r.LocalFilePath)
				require.Error(t, err)
				assert.True(t, os.IsNotExist(err))
			},
		},
		{
			Name: "Success without temp file deletion",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				caseNumber := strings.ReplaceAll(t.Name(), "/", "_")
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments: []string{caseNumber},
						stringFlags: map[string]string{
							"target-repo": "logs",
						},
						boolFlags: map[string]bool{
							"cleanup": false,
						},
						rtDetails:       rtDetails,
						targetRtDetails: targetRtDetails,
					})
				require.NoError(t, err)

				require.NotNil(t, r)
				assert.NotEmpty(t, r.BundleID)
				assert.NotEmpty(t, r.LocalFilePath)
				assert.NotEmpty(t, r.UploadURL)

				exists, err := uploadedPathExists(targetRtDetails, r.UploadURL)
				require.NoError(t, err)
				assert.True(t, exists)

				exists, err = supportBundleExists(rtDetails, r.BundleID)
				require.NoError(t, err)
				assert.True(t, exists)

				stat, err := os.Stat(r.LocalFilePath)
				require.NoError(t, err)
				assert.Greater(t, stat.Size(), int64(0))
			},
		},
		{
			Name: "Fail because no args",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := commands.SupportBundleCmd(context.Background(), &cliStub{})
				assert.EqualError(t, err, "wrong number of arguments. Expected: 1, Received: 0")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to get RT details",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments:       []string{"1234"},
						rtDetails:       nil,
						targetRtDetails: targetRtDetails,
					})
				assert.EqualError(t, err, "failed to get RT details")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to get target RT details",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments:       []string{"1234"},
						rtDetails:       rtDetails,
						targetRtDetails: nil,
					})
				assert.EqualError(t, err, "failed to get Target RT details")
				assert.Nil(t, r)
			},
		},
		{
			Name: "Fail to create support bundle",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments: []string{"1234"},
						rtDetails: &config.ArtifactoryDetails{
							Url: "http://rt.invalid",
						},
						targetRtDetails: targetRtDetails,
					})
				require.Error(t, err)
				assert.Contains(t, err.Error(), "rt.invalid")
				require.NotNil(t, r)
				assert.Empty(t, r.BundleID)
			},
		},
		{
			Name: "Fail to upload support bundle",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails, targetRtDetails *config.ArtifactoryDetails) {
				r, err := commands.SupportBundleCmd(
					context.Background(),
					&cliStub{
						arguments: []string{"1234"},
						rtDetails: rtDetails,
						targetRtDetails: &config.ArtifactoryDetails{
							Url: "http://rt.invalid",
						},
					})
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

func uploadedPathExists(rtDetails *config.ArtifactoryDetails, path string) (bool, error) {
	endpoint := fmt.Sprintf("%sapi/storage/%s", rtDetails.Url, strings.ReplaceAll(path, rtDetails.Url, ""))
	req, err := http.NewRequestWithContext(context.Background(), "HEAD", endpoint, nil)
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
		return false, fmt.Errorf("GET %s -> %d", endpoint, res.StatusCode)
	}
}

func supportBundleExists(rtDetails *config.ArtifactoryDetails, bundleID actions.BundleID) (bool, error) {
	endpoint := fmt.Sprintf("%sapi/system/support/bundle/%s", rtDetails.Url, bundleID)
	req, err := http.NewRequestWithContext(context.Background(), "GET", endpoint, nil)
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
		return false, fmt.Errorf("GET %s -> %d", endpoint, res.StatusCode)
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
