package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

const (
	body = `{"status":"%s"}`
	url  = "url"
)

type checkStatusClientStub struct {
	count            int
	statusCode       int
	payloads         []string
	err              error
	receivedBundleID bundleID
}

func (cs *checkStatusClientStub) GetSupportBundleStatus(bundleID bundleID) (status int, responseBytes []byte, err error) {
	responseBytes = []byte(cs.payloads[cs.count])
	cs.receivedBundleID = bundleID
	cs.count++
	return cs.statusCode, responseBytes, cs.err
}

func (cs *checkStatusClientStub) DownloadSupportBundle(_ bundleID) (*http.Response, error) {
	return nil, nil
}

func (cs *checkStatusClientStub) GetURL() string {
	return url
}

func Test_WaitUntilReady(t *testing.T) {
	tests := []struct {
		name                 string
		timeout              time.Duration
		retryInterval        time.Duration
		clientStub           *checkStatusClientStub
		expectedErrorMessage string
	}{
		{
			name:          "first retry successful",
			timeout:       100 * time.Millisecond,
			retryInterval: 10 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{fmt.Sprintf(body, "success")},
			},
		},
		{
			name:          "second retry successful",
			timeout:       100 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{fmt.Sprintf(body, "in progress"), fmt.Sprintf(body, "success")},
			},
		},
		{
			name:          "support bundle not found",
			timeout:       100 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusNotFound,
				payloads:   []string{`{}`},
			},
			expectedErrorMessage: "http request failed with: 404 Not Found",
		},
		{
			name:          "retry fails with 500",
			timeout:       100 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusInternalServerError,
				payloads:   []string{fmt.Sprintf(body, "in progress")},
			},
			expectedErrorMessage: "http request failed with: 500 Internal Server Error",
		},
		{
			name:          "timeout",
			timeout:       10 * time.Millisecond,
			retryInterval: 6 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{fmt.Sprintf(body, "in progress"), fmt.Sprintf(body, "in progress")},
			},
			expectedErrorMessage: "timeout waiting for support bundle to be ready",
		},
		{
			name:          "client returns error",
			timeout:       100 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: -1,
				payloads:   []string{""},
				err:        errors.New("some error"),
			},
			expectedErrorMessage: "some error",
		},
		{
			name:          "client returns invalid json",
			timeout:       100 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{"{"},
			},
			expectedErrorMessage: "unexpected end of JSON input",
		},
		{
			name:          "client returns unexpected JSON response",
			timeout:       10 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{`{"some":"unexpected json"}`},
			},
			expectedErrorMessage: "property status not found",
		},
		{
			name:          "client returns unexpected type in response",
			timeout:       10 * time.Millisecond,
			retryInterval: 5 * time.Millisecond,
			clientStub: &checkStatusClientStub{
				statusCode: http.StatusOK,
				payloads:   []string{`{"status":{}}`},
			},
			expectedErrorMessage: "property status is not a string",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			err := waitUntilSupportBundleIsReady(ctx, test.clientStub, test.retryInterval, test.timeout,
				"bundleID")
			if test.expectedErrorMessage != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedErrorMessage)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, bundleID("bundleID"), test.clientStub.receivedBundleID)
		})
	}
}

type downloadClientStub struct {
	response         *http.Response
	err              error
	receivedBundleID bundleID
}

func (dc *downloadClientStub) GetSupportBundleStatus(bundleID) (status int, responseBytes []byte, err error) {
	return http.StatusOK, []byte(fmt.Sprintf(body, "success")), nil
}

func (dc *downloadClientStub) DownloadSupportBundle(bundleID bundleID) (*http.Response, error) {
	dc.receivedBundleID = bundleID
	return dc.response, dc.err
}

func (dc *downloadClientStub) GetURL() string {
	return "url"
}

func Test_DownloadSupportBundle(t *testing.T) {
	tests := []struct {
		name                 string
		clientStub           *downloadClientStub
		expectedErrorMessage string
	}{
		{
			name: "successful download",
			clientStub: &downloadClientStub{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       ioutil.NopCloser(strings.NewReader("file-contents")),
				}},
		},
		{
			name: "client returns error",
			clientStub: &downloadClientStub{
				response: &http.Response{
					StatusCode: -1,
				},
				err: errors.New("yikes, something really bad happened"),
			},
			expectedErrorMessage: "yikes, something really bad happened",
		},
		{
			name: "support bundle not found",
			clientStub: &downloadClientStub{
				response: &http.Response{
					StatusCode: http.StatusNotFound,
				},
			},
			expectedErrorMessage: "http request failed with: 404 Not Found",
		},
		{
			name: "client returns Internal Server Error",
			clientStub: &downloadClientStub{
				response: &http.Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			expectedErrorMessage: "http request failed with: 500 Internal Server Error",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			timeout := 10 * time.Millisecond
			retryInterval := 5 * time.Millisecond
			filePath, err := downloadSupportBundle(ctx, test.clientStub, timeout, retryInterval, "bundleID")
			if test.expectedErrorMessage != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedErrorMessage)
			} else {
				require.NoError(t, err)
				assert.Contains(t, filePath, "bundleID.zip")
			}
			assert.Equal(t, bundleID("bundleID"), test.clientStub.receivedBundleID)
		})
	}
}
