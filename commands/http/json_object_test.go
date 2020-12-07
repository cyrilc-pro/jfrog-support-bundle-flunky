package http

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseJSON(t *testing.T) {
	tests := []struct {
		name          string
		payload       string
		expectedError string
	}{
		{
			name:    "valid json",
			payload: `{"key":"value"}`,
		},
		{
			name:          "invalid json",
			payload:       `{"key":"value"`,
			expectedError: "unexpected end of JSON input",
		},
		{
			name:          "not json",
			payload:       `"key"`,
			expectedError: "json: cannot unmarshal string into Go value of type http.JSONObject",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			result, err := ParseJSON([]byte(test.payload))
			if test.expectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestJSONObject_GetString(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		expectedValue string
		expectedError string
	}{
		{
			name:          "valid key",
			key:           "key",
			expectedValue: "value",
		},
		{
			name:          "unknown key",
			key:           "unknown",
			expectedError: "property unknown not found",
		},
		{
			name:          "not a string",
			key:           "object_key",
			expectedError: "property object_key is not a string",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			payload := []byte(`{"key":"value", "object_key":{}}`)
			result, err := ParseJSON(payload)
			require.NoError(t, err)
			got, err := result.GetString(test.key)
			if test.expectedError != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedValue, got)
			}
		})
	}
}
