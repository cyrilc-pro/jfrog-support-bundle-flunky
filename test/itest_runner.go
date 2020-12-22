package test

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/sync/errgroup"
	"os"
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
	parentCtx := context.Background()

	eg, ctx := errgroup.WithContext(parentCtx)

	rt := &config.ArtifactoryDetails{}
	eg.Go(func() error {
		cleaner, err := startArtifactoryAndDeployLicense(ctx, t, rt, version, licenseKey)
		registerCleanup(parentCtx, t, cleaner)
		return err
	})

	targetRt := &config.ArtifactoryDetails{}
	eg.Go(func() error {
		cleaner, err := startArtifactoryAndDeployLicense(ctx, t, targetRt, version, licenseKey)
		registerCleanup(parentCtx, t, cleaner)
		if err != nil {
			return err
		}
		return setUpTargetRepositoryAndPermissions(ctx, t, targetRt)
	})

	err := eg.Wait()
	require.NoError(t, err)

	log.SetLogger(&testLogger{t: t})

	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			test.Function(t, rt, targetRt)
		})
	}
}

type cleaner func(ctx context.Context) error

func registerCleanup(ctx context.Context, t *testing.T, cleaner cleaner) {
	if cleaner != nil {
		t.Cleanup(func() {
			err := cleaner(ctx)
			if err != nil {
				t.Logf("Could not terminate container: %v", err)
			}
		})
	}
}

func startArtifactoryAndDeployLicense(ctx context.Context, l logger, rt *config.ArtifactoryDetails,
	version, licenseKey string) (cleaner, error) {
	cont, err := createContainer(ctx, version)
	if err != nil {
		return nil, err
	}
	err = buildRtDetails(ctx, cont, rt)
	if err != nil {
		return cont.Terminate, err
	}
	err = setUpLicense(ctx, l, rt, licenseKey)
	return cont.Terminate, err
}

func buildRtDetails(ctx context.Context, cont testcontainers.Container, rt *config.ArtifactoryDetails) error {
	ip, err := cont.Host(ctx)
	if err != nil {
		return err
	}
	port, err := cont.MappedPort(ctx, "8082")
	if err != nil {
		return err
	}
	rt.Url = fmt.Sprintf("http://%s:%d/artifactory/", ip, port.Int())
	rt.User = "admin"
	rt.Password = "password"
	return nil
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
