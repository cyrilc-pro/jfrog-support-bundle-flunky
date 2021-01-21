package commands

import (
	"errors"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/actions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type flagProviderStub struct {
	value            string
	receivedFlagName string
	boolVal          bool
}

func (fps *flagProviderStub) GetStringFlagValue(flagName string) string {
	fps.receivedFlagName = flagName
	return fps.value
}

func (fps *flagProviderStub) GetBoolFlagValue(flagName string) bool {
	fps.receivedFlagName = flagName
	return fps.boolVal
}

type serviceHelperStub struct {
	details   *config.ArtifactoryDetails
	configErr error
	initErr   error
}

func (chs *serviceHelperStub) GetConfig(string, bool) (*config.ArtifactoryDetails, error) {
	return chs.details, chs.configErr
}

func (chs *serviceHelperStub) CreateInitialRefreshableTokensIfNeeded(*config.ArtifactoryDetails) error {
	return chs.initErr
}

func Test_getTimeoutAndRetryInterval(t *testing.T) {
	defaultTimeout := 10 * time.Minute
	defaultRetry := 5 * time.Second
	tests := []struct {
		name                  string
		flagProvider          *flagProviderStub
		expectedTimeout       time.Duration
		expectedRetryInterval time.Duration
	}{
		{
			name: "empty string uses default",
			flagProvider: &flagProviderStub{
				value: "",
			},
			expectedTimeout:       defaultTimeout,
			expectedRetryInterval: defaultRetry,
		},
		{
			name: "parse error uses default",
			flagProvider: &flagProviderStub{
				value: "30 seconds",
			},
			expectedTimeout:       defaultTimeout,
			expectedRetryInterval: defaultRetry,
		},
		{
			name: "valid duration",
			flagProvider: &flagProviderStub{
				value: "25s",
			},
			expectedTimeout:       25 * time.Second,
			expectedRetryInterval: 25 * time.Second,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expectedTimeout, getTimeout(test.flagProvider))
			assert.Equal(t, "download-timeout", test.flagProvider.receivedFlagName)
			assert.Equal(t, test.expectedRetryInterval, getRetryInterval(test.flagProvider))
			assert.Equal(t, "retry-interval", test.flagProvider.receivedFlagName)
		})
	}
}

func Test_getPromptOptions(t *testing.T) {
	tests := []struct {
		name          string
		flagProvider  *flagProviderStub
		expectDefault bool
	}{
		{
			name:          "no prompt options specified uses default",
			flagProvider:  &flagProviderStub{},
			expectDefault: true,
		},
		{
			name: "false prompt options specified uses default",
			flagProvider: &flagProviderStub{
				boolVal: false,
			},
			expectDefault: true,
		},
		{
			name: "true prompt options specified uses custom",
			flagProvider: &flagProviderStub{
				boolVal: true,
			},
			expectDefault: false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			provider := getPromptOptions(test.flagProvider)
			_, isDefault := provider.(*actions.DefaultOptionsProvider)
			assert.Equal(t, test.expectDefault, isDefault)
			assert.Equal(t, "prompt-options", test.flagProvider.receivedFlagName)
		})
	}
}

func Test_shouldCleanup(t *testing.T) {
	tests := []struct {
		name         string
		flagProvider *flagProviderStub
		expected     bool
	}{
		{
			name: "cleanup",
			flagProvider: &flagProviderStub{
				boolVal: true,
			},
			expected: true,
		},
		{
			name: "do not cleanup",
			flagProvider: &flagProviderStub{
				boolVal: false,
			},
			expected: false,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			result := shouldCleanup(test.flagProvider)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, "cleanup", test.flagProvider.receivedFlagName)
		})
	}
}

func Test_getTargetRepo(t *testing.T) {
	tests := []struct {
		name         string
		flagProvider *flagProviderStub
		expected     string
	}{
		{
			name: "get repo",
			flagProvider: &flagProviderStub{
				value: "repo-local",
			},
			expected: "repo-local",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			result := getTargetRepo(test.flagProvider)
			assert.Equal(t, test.expected, result)
			assert.Equal(t, "target-repo", test.flagProvider.receivedFlagName)
		})
	}
}

type getRtDetailsTest struct {
	name            string
	flagProvider    *flagProviderStub
	configHelper    serviceHelper
	expectedDetails *config.ArtifactoryDetails
	expectedError   string
}

func runGetRtDetailsTests(t *testing.T, tests []getRtDetailsTest, expectedFlagName string,
	getter func(flagProvider flagValueProvider, configProvider serviceHelper) (*config.ArtifactoryDetails, error)) {
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			result, err := getter(test.flagProvider, test.configHelper)
			if test.expectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedDetails, result)
			}
			assert.Equal(t, expectedFlagName, test.flagProvider.receivedFlagName)
		})
	}
}

func Test_getRtDetails(t *testing.T) {
	expectedDetailsWithCreds := &config.ArtifactoryDetails{
		Url:      "http://myurl.test/",
		User:     "me",
		Password: "top-secret",
	}
	tests := []getRtDetailsTest{
		{
			name:         "adds trailing slash to URL",
			flagProvider: &flagProviderStub{value: "my-favorite-artifactory"},
			configHelper: &serviceHelperStub{
				details: &config.ArtifactoryDetails{
					Url:      "http://myurl.test",
					User:     "me",
					Password: "top-secret",
				},
			},
			expectedDetails: expectedDetailsWithCreds,
		},
		{
			name:         "URL already has trailing slash",
			flagProvider: &flagProviderStub{value: "my-favorite-artifactory"},
			configHelper: &serviceHelperStub{
				details: expectedDetailsWithCreds,
			},
			expectedDetails: expectedDetailsWithCreds,
		},
		{
			name:         "get config returns error",
			flagProvider: &flagProviderStub{value: "my-favorite-artifactory"},
			configHelper: &serviceHelperStub{
				configErr: errors.New("oops"),
			},
			expectedError: "oops",
		},
		{
			name:         "init token returns error",
			flagProvider: &flagProviderStub{value: "my-favorite-artifactory"},
			configHelper: &serviceHelperStub{
				details: expectedDetailsWithCreds,
				initErr: errors.New("oops"),
			},
			expectedError: "oops",
		},
	}

	runGetRtDetailsTests(t, tests, "server-id", getRtDetails)
}

func Test_getRtTargetDetails(t *testing.T) {
	tests := []getRtDetailsTest{
		{
			name:         "default",
			flagProvider: &flagProviderStub{},
			configHelper: &serviceHelperStub{},
			expectedDetails: &config.ArtifactoryDetails{
				Url: "https://supportlogs.jfrog.com/",
			},
		},
		{
			name:         "specific target service",
			flagProvider: &flagProviderStub{value: "my-artifactory"},
			configHelper: &serviceHelperStub{
				details: &config.ArtifactoryDetails{
					Url: "my-artifactory-url",
				},
			},
			expectedDetails: &config.ArtifactoryDetails{
				Url: "my-artifactory-url/",
			},
		},
		{
			name:         "get config returns error",
			flagProvider: &flagProviderStub{value: "my-favorite-artifactory"},
			configHelper: &serviceHelperStub{
				configErr: errors.New("oops"),
			},
			expectedError: "oops",
		},
	}

	runGetRtDetailsTests(t, tests, "target-server-id", getTargetDetails)
}
