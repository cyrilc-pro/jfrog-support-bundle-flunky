package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"sync"
	"testing"
	"time"
)

type integrationTestFunction func(*testing.T, *config.ArtifactoryDetails, *config.ArtifactoryDetails)

type integrationTest struct {
	Name     string
	Function integrationTestFunction
}

type logger interface {
	Logf(format string, args ...interface{})
}

type rtInit struct {
	cont    testcontainers.Container
	details *config.ArtifactoryDetails
	err     error
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

	wg := sync.WaitGroup{}
	rt := &rtInit{}
	targetRt := &rtInit{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		initRt(ctx, t, rt, licenseKey, version)
	}()
	go func() {
		defer wg.Done()
		initTargetRt(ctx, t, targetRt, licenseKey, version)
	}()
	wg.Wait()

	t.Cleanup(func() {
		terminate(ctx, t, rt)
	})
	t.Cleanup(func() {
		terminate(ctx, t, targetRt)
	})

	require.NoError(t, rt.err)
	require.NoError(t, targetRt.err)

	require.NotNil(t, rt.details)
	require.NotNil(t, targetRt.details)

	log.SetLogger(&testLogger{t: t})

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			test.Function(t, rt.details, targetRt.details)
		})
	}
}

func terminate(ctx context.Context, l logger, rt *rtInit) {
	if rt.cont != nil {
		err := rt.cont.Terminate(ctx)
		if err != nil {
			l.Logf("Could not terminate container: %v", err)
		}
	}
}

func initTargetRt(ctx context.Context, l logger, r *rtInit, licenseKey, version string) {
	r.cont, r.err = createContainer(ctx, version)
	if r.err != nil {
		return
	}
	var ip string
	ip, r.err = r.cont.Host(ctx)
	if r.err != nil {
		return
	}
	var port nat.Port
	port, r.err = r.cont.MappedPort(ctx, "8082")
	if r.err != nil {
		return
	}
	r.details = buildRtDetails(ip, port)
	r.err = setUpLicense(ctx, l, licenseKey, r.details)
	if r.err != nil {
		return
	}
	r.err = setUpTargetRepositoryAndPermissions(ctx, l, r.details)
}

func initRt(ctx context.Context, l logger, r *rtInit, licenseKey, version string) {
	r.cont, r.err = createContainer(ctx, version)
	if r.err != nil {
		return
	}
	var ip string
	ip, r.err = r.cont.Host(ctx)
	if r.err != nil {
		return
	}
	var port nat.Port
	port, r.err = r.cont.MappedPort(ctx, "8082")
	if r.err != nil {
		return
	}
	r.details = buildRtDetails(ip, port)
	r.err = setUpLicense(ctx, l, licenseKey, r.details)
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
