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

func Test_Prompt(t *testing.T) {
	tests := []struct {
		name             string
		stub             PrompterStub
		expectErr        string
		expectParameters http.SupportBundleParameters
	}{
		{
			name: "Include all",
			stub: PrompterStub{
				IncludeLogs:          true,
				IncludeSystem:        true,
				IncludeConfiguration: true,
				IncludeThreadDump:    true,
			},
			expectParameters: http.SupportBundleParameters{
				Configuration: true,
				Logs: &http.SupportBundleParametersLogs{
					Include:   true,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: true,
				ThreadDump: &http.SupportBundleParametersThreadDump{
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
			expectParameters: http.SupportBundleParameters{
				Configuration: false,
				Logs: &http.SupportBundleParametersLogs{
					Include:   true,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &http.SupportBundleParametersThreadDump{
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
			expectParameters: http.SupportBundleParameters{
				Configuration: false,
				Logs: &http.SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: true,
				ThreadDump: &http.SupportBundleParametersThreadDump{
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
			expectParameters: http.SupportBundleParameters{
				Configuration: true,
				Logs: &http.SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &http.SupportBundleParametersThreadDump{
					Count:    0,
					Interval: 0,
				},
			},
		},
		{
			name: "Include nothing",
			stub: PrompterStub{},
			expectParameters: http.SupportBundleParameters{
				Configuration: false,
				Logs: &http.SupportBundleParametersLogs{
					Include:   false,
					StartDate: "2012-10-31",
					EndDate:   "2012-11-01",
				},
				System: false,
				ThreadDump: &http.SupportBundleParametersThreadDump{
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
			provider := PromptOptionsProvider{
				GetDate: func() time.Time {
					return time.Unix(1351807721, 0)
				},
				Prompter: &test.stub,
			}
			options, err := provider.GetOptions("foo")
			if test.expectErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, test.expectErr)
			} else {
				require.NoError(t, err)
				assert.Empty(t, cmp.Diff(
					http.SupportBundleCreationOptions{
						Name:        "JFrog Support Case number foo",
						Description: "Generated on 2012-11-01T22:08:41Z",
						Parameters:  &test.expectParameters,
					},
					options))
			}
		})
	}
}
