package actions

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

type uploadClientStub struct {
	err                error
	statusCode         int
	receivedPath       string
	receivedRepo       string
	receivedCaseNumber string
	receivedFilename   string
}

func (ucs *uploadClientStub) UploadSupportBundle(sbFilePath string, repoKey string, caseNumber string,
	filename string) (status int, responseBytes []byte, err error) {
	ucs.receivedPath = sbFilePath
	ucs.receivedRepo = repoKey
	ucs.receivedCaseNumber = caseNumber
	ucs.receivedFilename = filename
	return ucs.statusCode, []byte("response"), ucs.err
}

func (ucs *uploadClientStub) GetURL() string {
	return "http://foo.bar/"
}

func Test_Upload(t *testing.T) {
	tests := []struct {
		name                 string
		clientStub           *uploadClientStub
		expectedErrorMessage string
	}{
		{
			name: "success",
			clientStub: &uploadClientStub{
				statusCode: http.StatusCreated,
			},
		},
		{
			name: "client error",
			clientStub: &uploadClientStub{
				statusCode: http.StatusInternalServerError,
				err:        errors.New("crash, bang, boom"),
			},
			expectedErrorMessage: "crash, bang, boom",
		},
		{
			name: "unexpected response code",
			clientStub: &uploadClientStub{
				statusCode: http.StatusBadRequest,
			},
			expectedErrorMessage: "http request failed with: 400 Bad Request",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			caseNumber := CaseNumber("1234")
			now := func() time.Time { return time.Unix(1, 1) }
			path, err := UploadSupportBundle(test.clientStub, caseNumber, "/some/file", "logsRepo", now)
			if test.expectedErrorMessage != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedErrorMessage)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, "/some/file", test.clientStub.receivedPath)
			assert.Equal(t, "logsRepo", test.clientStub.receivedRepo)
			assert.Equal(t, "1234", test.clientStub.receivedCaseNumber)
			assert.Equal(t, "SB-19700101-000001Z.zip", test.clientStub.receivedFilename)
			assert.Equal(t, "http://foo.bar/logsRepo/1234/SB-19700101-000001Z.zip", path)
		})
	}
}
