package test

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"testing"
	"time"
)

// Main entry point for integration tests
func Test_RunAllTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests")
	}
	suite.Run(t, &IntegrationTestSuite{})
}

type IntegrationTestSuite struct {
	suite.Suite
	ctx         context.Context
	rtContainer testcontainers.Container
	rtDetails   *config.ArtifactoryDetails
}

func (s *IntegrationTestSuite) SetupSuite() {
	licenseKey, exists := os.LookupEnv("TEST_LICENSE")
	if !exists || licenseKey == "" {
		s.T().Skip("Environment variable TEST_LICENSE does not contain a license key")
		return
	}
	version, exists := os.LookupEnv("ARTIFACTORY_VERSION")
	if !exists {
		version = "latest"
	}
	s.ctx = context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "docker.bintray.io/jfrog/artifactory-pro:" + version,
		ExposedPorts: []string{"8082"},
		WaitingFor: wait.ForHTTP("/artifactory/api/system/ping").
			WithPort("8082").
			WithStartupTimeout(3 * time.Minute).
			WithPollInterval(10 * time.Second),
	}
	var err error
	s.rtContainer, err = testcontainers.GenericContainer(s.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(s.T(), err)
	s.T().Cleanup(func() {})
	ip, err := s.rtContainer.Host(s.ctx)
	require.NoError(s.T(), err)
	port, err := s.rtContainer.MappedPort(s.ctx, "8082")
	require.NoError(s.T(), err)

	s.rtDetails = &config.ArtifactoryDetails{
		Url:      fmt.Sprintf("http://%s:%d/artifactory/", ip, port.Int()),
		User:     "admin",
		Password: "password",
	}

	setUpLicense(s.ctx, s.T(), licenseKey, s.rtDetails)

	log.SetLogger(&testLog{t: s.T()})
}

func (s *IntegrationTestSuite) TearDownSuite() {
	err := s.rtContainer.Terminate(s.ctx)
	require.NoError(s.T(), err)
}
