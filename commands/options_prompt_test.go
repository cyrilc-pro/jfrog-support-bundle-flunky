package commands

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Prompt(t *testing.T) {
	clock := func() time.Time {
		timestamp, err := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")
		require.NoError(t, err)
		return timestamp
	}

	tests := []struct {
		name             string
		stub             PrompterStub
		expectErr        string
		expectParameters SupportBundleParameters
	}{
		{
			name: "Include all",
			stub: PrompterStub{
				IncludeLogs:          true,
				IncludeSystem:        true,
				IncludeConfiguration: true,
				IncludeThreadDump:    true,
			},
			expectParameters: SupportBundleParameters{
				Configuration: true,
				Logs: &SupportBundleParametersLogs{
					Include:   true,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: true,
				ThreadDump: &SupportBundleParametersThreadDump{
					Count:    1,
					Interval: 0,
				},
			},
		},
		{
			name: "Include logs only",
			stub: PrompterStub{
				IncludeLogs: true,
			},
			expectParameters: SupportBundleParameters{
				Configuration: false,
				Logs: &SupportBundleParametersLogs{
					Include:   true,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &SupportBundleParametersThreadDump{
					Count:    0,
					Interval: 0,
				},
			},
		},
		{
			name: "Include system only",
			stub: PrompterStub{
				IncludeSystem: true,
			},
			expectParameters: SupportBundleParameters{
				Configuration: false,
				Logs: &SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: true,
				ThreadDump: &SupportBundleParametersThreadDump{
					Count:    0,
					Interval: 0,
				},
			},
		},
		{
			name: "Include configuration only",
			stub: PrompterStub{
				IncludeConfiguration: true,
			},
			expectParameters: SupportBundleParameters{
				Configuration: true,
				Logs: &SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &SupportBundleParametersThreadDump{
					Count:    0,
					Interval: 0,
				},
			},
		},
		{
			name: "Include nothing",
			stub: PrompterStub{},
			expectParameters: SupportBundleParameters{
				Configuration: false,
				Logs: &SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &SupportBundleParametersThreadDump{
					Count:    0,
					Interval: 0,
				},
			},
		},
		{
			name: "Error",
			stub: PrompterStub{
				err: errors.New("oops"),
			},
			expectErr: "oops",
		},
	}
	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			provider := PromptOptionsProvider{GetDate: clock, Prompter: &test.stub}
			options, err := provider.GetOptions("foo")
			if test.expectErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectErr)
			} else {
				require.NoError(t, err)
				assert.Empty(t, cmp.Diff(
					SupportBundleCreationOptions{
						Name:        "JFrog Support Case number foo",
						Description: "Generated on 2012-11-01T22:08:41Z",
						Parameters:  &test.expectParameters,
					},
					options))
			}
		})
	}
}
