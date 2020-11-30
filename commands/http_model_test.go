package commands

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
	"testing"
)

func Test_MarshallJSON(t *testing.T) {
	tests := []struct {
		name   string
		input  SupportBundleCreationOptions
		expect string
	}{
		{
			name: "Nil parameters",
			input: SupportBundleCreationOptions{
				Name:        "n",
				Description: "d",
			},
			expect: `{"name":"n","description":"d","parameters":{}}`,
		},
		{
			name: "With parameters",
			input: SupportBundleCreationOptions{
				Name:        "n",
				Description: "d",
				Parameters: &SupportBundleParameters{
					Configuration: true,
					Logs: &SupportBundleParametersLogs{
						Include:   false,
						StartDate: "s",
						EndDate:   "e",
					},
					System: false,
					ThreadDump: &SupportBundleParametersThreadDump{
						Count:    1,
						Interval: 2,
					},
				},
			},
			expect: `{"name":"n","description":"d","parameters":{"configuration":true,"logs":` +
				`{"include":false,"start_date":"s","end_date":"e"},"system":false,"thread_dump":{"count":1,"interval":2}}}`,
		},
		{
			name: "Nil logs and threaddump",
			input: SupportBundleCreationOptions{
				Name:        "n",
				Description: "d",
				Parameters: &SupportBundleParameters{
					Configuration: true,
					System:        false,
				},
			},
			expect: `{"name":"n","description":"d","parameters":{"configuration":true,"logs":null,"system":false,"thread_dump":null}}`,
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			bytes, err := json.Marshal(test.input)
			require.NoError(t, err)
			assert.Equal(t, string(bytes), test.expect)
		})
	}
}
