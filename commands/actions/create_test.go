package actions

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type createSupportBundleHTTPClientStub struct {
	statusCode    int
	response      string
	err           error
	actualPayload http.SupportBundleCreationOptions
}

func (c *createSupportBundleHTTPClientStub) GetURL() string {
	return "stub"
}

func (c *createSupportBundleHTTPClientStub) CreateSupportBundle(payload http.SupportBundleCreationOptions) (status int,
	responseBytes []byte, err error) {
	c.actualPayload = payload
	return c.statusCode, []byte(c.response), c.err
}

func Test_CreateSupportBundle(t *testing.T) {
	clock := func() time.Time {
		return time.Unix(1351807721, 0)
	}

	tests := []struct {
		name                 string
		givenHTTPClient      createSupportBundleHTTPClientStub
		givenOptionsProvider OptionsProvider
		expectErr            string
		expectClientSkipped  bool
		expectID             BundleID
	}{
		{
			name: "success",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{"id": "foo"}`,
				err:        nil,
			},
			expectErr: "",
			expectID:  "foo",
		},
		{
			name: "bad request",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				statusCode: 400,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "http request failed with: 400",
		},
		{
			name: "bad json",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `bad json`,
				err:        nil,
			},
			expectErr: "invalid character 'b' looking for beginning of value",
		},
		{
			name: "missing id",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "property id not found",
		},
		{
			name: "bad id",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{"id":{}}`,
				err:        nil,
			},
			expectErr: "property id is not a string",
		},
		{
			name: "error in HTTPClient",
			givenHTTPClient: createSupportBundleHTTPClientStub{
				err: errors.New("oops"),
			},
			expectErr: "oops",
		},
		{
			name: "error in OptionsProvider",
			givenOptionsProvider: &PromptOptionsProvider{
				GetDate: clock,
				Prompter: &PrompterStub{
					IncludeLogsErr: errors.New("oops"),
				},
			},
			expectClientSkipped: true,
			expectErr:           "oops",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			caseNumber := CaseNumber("1234")
			optionsProvider := test.givenOptionsProvider
			if optionsProvider == nil {
				optionsProvider = &DefaultOptionsProvider{getDate: clock}
			}
			id, err := CreateSupportBundle(&test.givenHTTPClient, caseNumber, optionsProvider)
			if test.expectErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, test.expectErr)
			} else {
				require.Equal(t, test.expectID, id)
			}
			if !test.expectClientSkipped {
				assert.Empty(t, cmp.Diff(
					http.SupportBundleCreationOptions{
						Name:        "JFrog Support Case number 1234",
						Description: "Generated on 2012-11-01T22:08:41Z",
					},
					test.givenHTTPClient.actualPayload))
			}
		})
	}
}
