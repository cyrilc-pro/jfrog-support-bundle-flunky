package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
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

type IntegrationTestFunction func(*testing.T, *config.ArtifactoryDetails)

type IntegrationTest struct {
	Name     string
	Function IntegrationTestFunction
}

func RunIntegrationTests(t *testing.T, tests []IntegrationTest) {
	t.Helper()
	licenseKey, exists := os.LookupEnv("TEST_LICENSE")
	if !exists {
		t.Skip("No license found in env")
		return
	}
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.bintray.io/jfrog/artifactory-pro:7.9.0",
		ExposedPorts: []string{"8082"},
		WaitingFor: wait.ForHTTP("/artifactory/api/system/ping").
			WithPort("8082").
			WithStartupTimeout(3 * time.Minute).
			WithPollInterval(10 * time.Second),
	}
	rtContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	t.Cleanup(func() { _ = rtContainer.Terminate(ctx) })
	ip, err := rtContainer.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	port, err := rtContainer.MappedPort(ctx, "8082")
	if err != nil {
		t.Error(err)
	}

	rtDetails := config.ArtifactoryDetails{
		Url:      fmt.Sprintf("http://%s:%d/artifactory/", ip, port.Int()),
		User:     "admin",
		Password: "password",
	}

	setUpLicense(t, ctx, licenseKey, err, rtDetails)

	log.SetLogger(&testLog{t: t})

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			test.Function(t, &rtDetails)
		})
	}
}

func setUpLicense(t *testing.T, ctx context.Context, licenseKey string, err error, rtDetails config.ArtifactoryDetails) {
	licensePayload := strings.NewReader(fmt.Sprintf(`{"licenseKey":"%s"}`, licenseKey))
	licensesEndpointUrl := fmt.Sprintf("%sapi/system/licenses", rtDetails.Url)
	postLicenseReq, err := http.NewRequestWithContext(ctx, "POST", licensesEndpointUrl, licensePayload)
	require.NoError(t, err)
	postLicenseReq.SetBasicAuth(rtDetails.User, rtDetails.Password)
	postLicenseReq.Header["Content-Type"] = []string{"application/json"}
	resp, err := http.DefaultClient.Do(postLicenseReq)
	require.NoError(t, err)
	bytes, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	t.Logf("Deploy license: %s %s", resp.Status, string(bytes))
	require.Equal(t, http.StatusOK, resp.StatusCode, "License deploy failed")
	getLicenseReq, err := http.NewRequestWithContext(ctx, "GET", licensesEndpointUrl, licensePayload)
	require.NoError(t, err)
	getLicenseReq.SetBasicAuth(rtDetails.User, rtDetails.Password)
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
			resp, err := http.DefaultClient.Do(getLicenseReq)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode, "License check failed")
			bytes, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			t.Logf("Get license: %s %s", resp.Status, string(bytes))
			payload := make(map[string]interface{})
			err = json.Unmarshal(bytes, &payload)
			require.NoError(t, err)
			licenseType := payload["type"]
			if licenseType != "N/A" {
				t.Logf("License %v applied", licenseType)
				retry = false
			}
		}
	}
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
