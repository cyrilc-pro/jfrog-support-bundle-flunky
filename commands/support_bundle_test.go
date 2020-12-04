package commands

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func Test_GetSupportBundleCommand(t *testing.T) {
	expectedFlags := []components.Flag{
		components.StringFlag{
			Name: "server-id",
			Description: "Artifactory server ID configured using the config command. " +
				"If not provided the default configuration will be used.",
		},
		components.StringFlag{
			Name: "target-server-id",
			Description: "Artifactory server ID configured using the config command to be used as the target for " +
				"uploading the generated Support Bundle. If not provided JFrog support logs will be used.",
		},
		components.StringFlag{
			Name:         "download-timeout",
			Description:  "The timeout for download.",
			DefaultValue: "10m",
		},
		components.StringFlag{
			Name:         "retry-interval",
			Description:  "The duration to wait between retries.",
			DefaultValue: "5s",
		},
		components.BoolFlag{
			Name:        "prompt-options",
			Description: "Ask for support bundle options or use Artifactory default options.",
		},
		components.BoolFlag{
			Name:         "cleanup",
			Description:  "Delete the support bundle local temp file after upload.",
			DefaultValue: true,
		},
		components.StringFlag{
			Name:         "target-repo",
			Description:  "The target repository key where the support bundle will be uploaded to.",
			DefaultValue: "logs",
		},
	}

	expectedArgs := []components.Argument{
		{
			Name:        "case",
			Description: "JFrog Support case number.",
		},
	}

	expected := components.Command{
		Name:        "support-case",
		Description: `Creates a Support Bundle and uploads it to JFrog Support "dropbox" service`,
		Aliases:     []string{"c", "case"},
		Arguments:   expectedArgs,
		Flags:       expectedFlags,
		EnvVars:     nil,
	}
	assert.Empty(t, cmp.Diff(expected, GetSupportBundleCommand(),
		cmpopts.IgnoreFields(components.Command{}, "Action")))
}

func Test_parseArguments(t *testing.T) {
	tests := []struct {
		name          string
		ctx           *components.Context
		expected      actions.CaseNumber
		expectedError string
	}{
		{
			name:     "parse valid argument",
			ctx:      &components.Context{Arguments: []string{"1234"}},
			expected: "1234",
		},
		{
			name:     "parse valid argument with whitespace",
			ctx:      &components.Context{Arguments: []string{"   1234  "}},
			expected: "1234",
		},
		{
			name:          "parse too many arguments",
			ctx:           &components.Context{Arguments: []string{"1234", "5678"}},
			expectedError: "Wrong number of arguments. Expected: 1, Received: 2",
		},
		{
			name:          "not enough arguments",
			ctx:           &components.Context{},
			expectedError: "Wrong number of arguments. Expected: 1, Received: 0",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			caseNumber, err := parseArguments(test.ctx)
			if test.expectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expected, caseNumber)
			}
		})
	}
}

func Test_deleteFile(t *testing.T) {
	dir := os.TempDir()
	path := filepath.Join(dir, "testfile")
	f, err := os.Create(path)
	require.NoError(t, err)
	defer assert.NoError(t, f.Close())

	deleteSupportBundleArchive(path)

	assert.True(t, !exists(path))
}

func Test_deleteNonExistentFile(t *testing.T) {
	path := "file/does/not/exist"
	deleteSupportBundleArchive(path)
	assert.False(t, exists(path))
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
