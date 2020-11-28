package commands

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type createSupportBundleHTTPClientStub struct {
	statusCode    int
	response      string
	err           error
	actualPayload string
}

func (c *createSupportBundleHTTPClientStub) GetURL() string {
	return "stub"
}

func (c *createSupportBundleHTTPClientStub) CreateSupportBundle(payload string) (status int, responseBytes []byte, err error) {
	c.actualPayload = payload
	return c.statusCode, []byte(c.response), c.err
}

func Test_CreateSupportBundle(t *testing.T) {
	tests := []struct {
		name      string
		given     createSupportBundleHTTPClientStub
		expectErr string
		expectID  bundleID
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
			conf := &supportBundleCommandConfiguration{
				caseNumber: "1234",
			}
			now := func() string { return "now" }
			id, err := createSupportBundle(&test.given, conf, now)
			if test.expectErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, test.expectErr)
			} else {
				require.Equal(t, test.expectID, id)
			}
			require.Equal(t, `{"name": "JFrog Support Case number 1234","description": "Generated on now","parameters":{}}`,
				test.given.actualPayload)
		})
	}
}
