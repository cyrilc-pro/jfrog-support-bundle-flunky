package commands

import (
	"context"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
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
		Url:      fmt.Sprintf("http://%s:%s/artifactory/", ip, port),
		User:     "admin",
		Password: "password",
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.Name, func(t *testing.T) {
			test.Function(t, &rtDetails)
		})
	}
}
