package commands

import (
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Create(t *testing.T) {
	tests := []IntegrationTest{
		{
			Name: "Success",
			Function: func(t *testing.T, rtDetails *config.ArtifactoryDetails) {
				conf := supportBundleCommandConfiguration{caseNumber: "foo"}
				r, err := createSupportBundle(&HTTPClient{rtDetails: rtDetails}, &conf, Now)
				require.NoError(t, err)
				require.NotEmpty(t, r.ID)
			},
		},
	}
	RunIntegrationTests(t, tests)
}
