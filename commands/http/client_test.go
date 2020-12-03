package http

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"github.com/jfrog/jfrog-cli-core/utils/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const authorizationHeader = "Authorization"

type request struct {
	Method        string
	ContentType   []string
	Body          string
	RequestURI    string
	Authorization []string
}

func TestClient_CreateSupportBundle_Success(t *testing.T) {
	ts, c := startedServer(t)
	defer ts.Close()

	s, res, err := createSupportBundle(c)

	require.NoError(t, err)
	require.Equal(t, s, http.StatusOK)
	var req request
	err = json.Unmarshal(res, &req)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(req, request{
		Method:      "POST",
		ContentType: []string{"application/json"},
		Body: `{"name":"foo","description":"desc","parameters":{"configuration":true,` +
			`"logs":{"include":false,"start_date":"2020-12-01","end_date":"2020-12-20"},"system":false,` +
			`"thread_dump":{"count":1,"interval":2}}}`,
		RequestURI:    "/api/system/support/bundle",
		Authorization: []string{"Basic YWRtaW46cGFzc3dvcmQ="},
	}))
}

func TestClient_DownloadSupportBundle_Success(t *testing.T) {
	ts, c := startedServer(t)
	defer ts.Close()

	res, err := c.DownloadSupportBundle("foo")
	require.NoError(t, err)
	defer func() { _ = res.Body.Close() }()

	require.Equal(t, res.StatusCode, http.StatusOK)

	bytes, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)

	var req request
	err = json.Unmarshal(bytes, &req)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(req, request{
		Method:        "GET",
		ContentType:   nil,
		Body:          "",
		RequestURI:    "/api/system/support/bundle/foo/archive",
		Authorization: []string{"Basic YWRtaW46cGFzc3dvcmQ="},
	}))
}

func TestClient_GetSupportBundleStatus_Success(t *testing.T) {
	ts, c := startedServer(t)
	defer ts.Close()

	status, bytes, err := c.GetSupportBundleStatus("foo")

	require.NoError(t, err)
	require.Equal(t, status, http.StatusOK)
	var req request
	err = json.Unmarshal(bytes, &req)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(req, request{
		Method:        "GET",
		ContentType:   nil,
		Body:          "",
		RequestURI:    "/api/system/support/bundle/foo",
		Authorization: []string{"Basic YWRtaW46cGFzc3dvcmQ="},
	}))
}

func TestClient_UploadSupportBundleStatus_Success(t *testing.T) {
	ts, c := startedServer(t)
	defer ts.Close()

	file, err := createTempFile()
	require.NoError(t, err)
	defer func() { _ = os.Remove(file.Name()) }()
	status, bytes, err := c.UploadSupportBundle(file.Name(), "r", "c", "f")

	require.NoError(t, err)
	require.Equal(t, status, http.StatusOK)
	var req request
	err = json.Unmarshal(bytes, &req)
	require.NoError(t, err)

	assert.Empty(t, cmp.Diff(req, request{
		Method:        "PUT",
		ContentType:   nil,
		Body:          "hello world",
		RequestURI:    "/r/c/f;uploadedBy=support-bundle-flunky",
		Authorization: []string{"Basic YWRtaW46cGFzc3dvcmQ="},
	}))
}

func TestClient_Offline(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T, c *Client) error
	}{
		{
			name: "Create",
			run: func(t *testing.T, c *Client) error {
				_, _, err := createSupportBundle(c)
				return err
			},
		},
		{
			name: "Download",
			run: func(t *testing.T, c *Client) error {
				res, err := c.DownloadSupportBundle("foo")
				if err == nil {
					_ = res.Body.Close()
				}
				return err
			},
		},
		{
			name: "Get Status",
			run: func(t *testing.T, c *Client) error {
				_, _, err := c.GetSupportBundleStatus("foo")
				return err
			},
		},
		{
			name: "Upload",
			run: func(t *testing.T, c *Client) error {
				file, err := createTempFile()
				require.NoError(t, err)
				defer func() { _ = os.Remove(file.Name()) }()
				_, _, err = c.UploadSupportBundle(file.Name(), "r", "c", "f")
				return err
			},
		},
	}

	ts, c := startedServer(t)
	ts.Close()

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			err := test.run(t, c)
			require.Error(t, err)
			require.Contains(t, err.Error(), "dial tcp")
		})
	}
}

func createTempFile() (*os.File, error) {
	file, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(file.Name(), []byte("hello world"), 0600)
	return file, err
}

func createSupportBundle(c *Client) (status int, responseBytes []byte, err error) {
	status, responseBytes, err = c.CreateSupportBundle(SupportBundleCreationOptions{
		Name:        "foo",
		Description: "desc",
		Parameters: &SupportBundleParameters{
			Configuration: true,
			Logs: &SupportBundleParametersLogs{
				Include:   false,
				StartDate: "2020-12-01",
				EndDate:   "2020-12-20",
			},
			System: false,
			ThreadDump: &SupportBundleParametersThreadDump{
				Count:    1,
				Interval: 2,
			},
		},
	})
	return
}

func newHTTPClient(ts *httptest.Server) *Client {
	c := Client{
		RtDetails: &config.ArtifactoryDetails{
			Url:      ts.URL + "/",
			User:     "admin",
			Password: "password",
		}}
	return &c
}

func startedServer(t *testing.T) (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		defer func() { _ = r.Body.Close() }()
		require.NoError(t, err)
		req := &request{
			Method:        r.Method,
			RequestURI:    r.RequestURI,
			ContentType:   r.Header[HTTPContentType],
			Body:          string(bytes),
			Authorization: r.Header[authorizationHeader],
		}
		res, err := json.Marshal(req)
		require.NoError(t, err)
		_, err = w.Write(res)
		require.NoError(t, err)
		w.WriteHeader(http.StatusOK)
	}))
	return ts, newHTTPClient(ts)
}
