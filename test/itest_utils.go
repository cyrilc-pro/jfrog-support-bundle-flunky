package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

type integrationTestFunction func(*testing.T, *config.ArtifactoryDetails, *config.ArtifactoryDetails)

type integrationTest struct {
	Name     string
	Function integrationTestFunction
}

func runIntegrationTests(t *testing.T, tests []integrationTest) {
	t.Helper()
	licenseKey, exists := os.LookupEnv("TEST_LICENSE")
	if !exists || licenseKey == "" {
		t.Skip("Environment variable TEST_LICENSE does not contain a license key")
		return
	}
	version, exists := os.LookupEnv("ARTIFACTORY_VERSION")
	if !exists {
		version = "latest"
	}
	ctx := context.Background()
	rtContainer, err := createContainer(ctx, version)
	require.NoError(t, err)

	targetContainer, err := createContainer(ctx, version)
	require.NoError(t, err)

	t.Cleanup(func() {
		terminate(ctx, t, []testcontainers.Container{rtContainer, targetContainer})
	})

	ip, err := rtContainer.Host(ctx)
	require.NoError(t, err)
	port, err := rtContainer.MappedPort(ctx, "8082")
	require.NoError(t, err)
	rtDetails := buildRtDetails(ip, port)

	targetIP, err := targetContainer.Host(ctx)
	require.NoError(t, err)
	targetPort, err := targetContainer.MappedPort(ctx, "8082")
	require.NoError(t, err)
	targetRtDetails := buildRtDetails(targetIP, targetPort)

	setUpLicense(ctx, t, licenseKey, rtDetails)
	setUpLicense(ctx, t, licenseKey, targetRtDetails)
	setUpTargetRepositoryAndPermissions(ctx, t, targetRtDetails)

	log.SetLogger(&testLog{t: t})

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			test.Function(t, rtDetails, targetRtDetails)
		})
	}
}

func buildRtDetails(ip string, port nat.Port) *config.ArtifactoryDetails {
	return &config.ArtifactoryDetails{
		Url:      fmt.Sprintf("http://%s:%d/artifactory/", ip, port.Int()),
		User:     "admin",
		Password: "password",
	}
}

func createContainer(ctx context.Context, version string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "docker.bintray.io/jfrog/artifactory-pro:" + version,
		ExposedPorts: []string{"8082"},
		WaitingFor: wait.ForHTTP("/artifactory/api/system/ping").
			WithPort("8082").
			WithStartupTimeout(10 * time.Minute).
			WithPollInterval(10 * time.Second),
	}
	rtContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return rtContainer, err
}

func terminate(ctx context.Context, t *testing.T, containers []testcontainers.Container) {
	for _, c := range containers {
		err := c.Terminate(ctx)
		assert.NoError(t, err)
	}
}

func setUpLicense(ctx context.Context, t *testing.T, licenseKey string, rtDetails *config.ArtifactoryDetails) {
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

func setUpTargetRepositoryAndPermissions(ctx context.Context, t *testing.T, rtDetails *config.ArtifactoryDetails) {
	createLogRepository(ctx, t, rtDetails)
	setAnonAccess(ctx, t, rtDetails)
	createAnonymousPermission(ctx, t, rtDetails)
}

func createAnonymousPermission(ctx context.Context, t *testing.T, rtDetails *config.ArtifactoryDetails) {
	payload := `{"name":"logsPerm","repo":{"repositories":["logs"],"actions":{"users":{"anonymous":["write"]}}}}`
	url := fmt.Sprintf("%sapi/v2/security/permissions/logsPerm", rtDetails.Url)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(payload))
	require.NoError(t, err)

	req.SetBasicAuth(rtDetails.User, rtDetails.Password)
	req.Header[commands.HTTPContentType] = []string{commands.HTTPContentTypeJSON}
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	t.Logf("Create anonymous permission on logs repository: status %s", resp.Status)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Add anonymous permission failed.")
}

func setAnonAccess(ctx context.Context, t *testing.T, rtDetails *config.ArtifactoryDetails) {
	url := fmt.Sprintf("%sapi/system/configuration", rtDetails.Url)
	getRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	require.NoError(t, err)
	getRequest.SetBasicAuth(rtDetails.User, rtDetails.Password)
	getRequest.Header[commands.HTTPContentType] = []string{commands.HTTPContentTypeXML}
	getResp, err := http.DefaultClient.Do(getRequest)
	require.NoError(t, err)
	getRespBody, err := ioutil.ReadAll(getResp.Body)
	require.NoError(t, err)
	defer func() { _ = getResp.Body.Close() }()

	payload := strings.Replace(string(getRespBody), "<anonAccessEnabled>false</anonAccessEnabled>",
		"<anonAccessEnabled>true</anonAccessEnabled>", 1)
	postRequest, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
	require.NoError(t, err)
	postRequest.SetBasicAuth(rtDetails.User, rtDetails.Password)
	postRequest.Header[commands.HTTPContentType] = []string{commands.HTTPContentTypeXML}
	postResponse, err := http.DefaultClient.Do(postRequest)
	require.NoError(t, err)
	postResponseBody, err := ioutil.ReadAll(postResponse.Body)
	require.NoError(t, err)
	defer func() { _ = postResponse.Body.Close() }()
	t.Logf("Set anonymous access: status %s, response: %s", getResp.Status, postResponseBody)
	require.Equal(t, http.StatusOK, getResp.StatusCode, "Set anonymous access failed.")
}

func createLogRepository(ctx context.Context, t *testing.T, targetRtDetails *config.ArtifactoryDetails) {
	payload := `{"key": "logs","rclass": "local","packageType": "generic"}`
	url := fmt.Sprintf("%sapi/repositories/logs", targetRtDetails.Url)
	req, err := http.NewRequestWithContext(ctx, "PUT", url, strings.NewReader(payload))
	require.NoError(t, err)
	req.SetBasicAuth(targetRtDetails.User, targetRtDetails.Password)
	req.Header[commands.HTTPContentType] = []string{commands.HTTPContentTypeJSON}

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	defer func() { _ = resp.Body.Close() }()
	require.NoError(t, err)
	t.Logf("Create logs repository: status %s, response: %s", resp.Status, respBody)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Create logs repo failed")
}

type testLog struct {
	t *testing.T
}

func (l *testLog) GetLogLevel() log.LevelType {
	return log.DEBUG
}
func (l *testLog) SetLogLevel(_ log.LevelType) {}
func (l *testLog) SetOutputWriter(_ io.Writer) {}
func (l *testLog) SetLogsWriter(_ io.Writer)   {}
func (l *testLog) Debug(a ...interface{}) {
	l.print("DEBUG", a)
}
func (l *testLog) Info(a ...interface{}) {
	l.print("INFO ", a)
}
func (l *testLog) Warn(a ...interface{}) {
	l.print("WARN ", a)
}
func (l *testLog) Error(a ...interface{}) {
	l.print("ERROR", a)
}
func (l *testLog) Output(a ...interface{}) {
	l.print("OUT  ", a)
}

func (l *testLog) print(level string, msgParts ...interface{}) {
	msg := level
	for i := range msgParts {
		msg += " " + fmt.Sprintf("%v", msgParts[i])
	}
	l.t.Log(msg)
}
