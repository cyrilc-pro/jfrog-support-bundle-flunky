package commands

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type createSupportBundleHttpClientStub struct {
	statusCode    int
	response      string
	err           error
	actualPayload string
}

func (c *createSupportBundleHttpClientStub) GetUrl() string {
	return "stub"
}

func (c *createSupportBundleHttpClientStub) CreateSupportBundle(payload string) (int, []byte, error) {
	c.actualPayload = payload
	return c.statusCode, []byte(c.response), c.err
}

func Test_CreateSupportBundle(t *testing.T) {
	tests := []struct {
		name       string
		given      createSupportBundleHttpClientStub
		expectErr  string
		expectResp creationResponse
	}{
		{
			name: "success",
			given: createSupportBundleHttpClientStub{
				statusCode: 200,
				response:   `{"id": "foo"}`,
				err:        nil,
			},
			expectErr:  "",
			expectResp: creationResponse{Id: "foo"},
		},
		{
			name: "bad request",
			given: createSupportBundleHttpClientStub{
				statusCode: 400,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "http request failed with: 400",
		},
		{
			name: "bad json",
			given: createSupportBundleHttpClientStub{
				statusCode: 200,
				response:   `bad json`,
				err:        nil,
			},
			expectErr: "invalid character 'b' looking for beginning of value",
		},
		{
			name: "missing id",
			given: createSupportBundleHttpClientStub{
				statusCode: 200,
				response:   `{}`,
				err:        nil,
			},
			expectErr: "property id not found",
		},
		{
			name: "bad id",
			given: createSupportBundleHttpClientStub{
				statusCode: 200,
				response:   `{"id":{}}`,
				err:        nil,
			},
			expectErr: "property id is not a string",
		},
		{
			name: "error",
			given: createSupportBundleHttpClientStub{
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
			resp, err := createSupportBundle(&test.given, conf, now)
			if test.expectErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, test.expectErr)
			} else {
				require.Equal(t, test.expectResp, resp)
			}
			require.Equal(t, `{"name": "JFrog Support Case number 1234","description": "Generated on now","parameters":{}}`, test.given.actualPayload)
		})
	}
}
