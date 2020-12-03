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
	tests := []struct {
		name      string
		given     createSupportBundleHTTPClientStub
		expectErr string
		expectID  BundleID
	}{
		{
			name: "success",
			given: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{"id": "foo"}`,
				err:        nil,
			},
			expectErr: "",
			expectID:  "foo",
		},
		{
			name: "bad request",
			given: createSupportBundleHTTPClientStub{
				statusCode: 400,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "http request failed with: 400",
		},
		{
			name: "bad json",
			given: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `bad json`,
				err:        nil,
			},
			expectErr: "invalid character 'b' looking for beginning of value",
		},
		{
			name: "missing id",
			given: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "property id not found",
		},
		{
			name: "bad id",
			given: createSupportBundleHTTPClientStub{
				statusCode: 200,
				response:   `{"id":{}}`,
				err:        nil,
			},
			expectErr: "property id is not a string",
		},
		{
			name: "error",
			given: createSupportBundleHTTPClientStub{
				err: errors.New("oops"),
			},
			expectErr: "oops",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			caseNumber := CaseNumber("1234")
			clock := func() time.Time {
				timestamp, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
				require.NoError(t, err)
				return timestamp
			}
			id, err := CreateSupportBundle(&test.given, caseNumber, &DefaultOptionsProvider{getDate: clock})
			if test.expectErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, test.expectErr)
			} else {
				require.Equal(t, test.expectID, id)
			}
			assert.Empty(t, cmp.Diff(
				http.SupportBundleCreationOptions{
					Name:        "JFrog Support Case number 1234",
					Description: "Generated on 2012-11-01T22:08:41Z",
				},
				test.given.actualPayload))
		})
	}
}
