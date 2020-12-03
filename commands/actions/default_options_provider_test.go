package actions

import (
	"github.com/google/go-cmp/cmp"
	"github.com/jfrog/jfrog-support-bundle-flunky/commands/http"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDefaultOptionsProvider_GetOptions(t *testing.T) {
	p := &DefaultOptionsProvider{getDate: func() time.Time {
		newYork, err := time.LoadLocation("America/New_York")
		require.NoError(t, err)
		return time.Date(2020, 12, 3, 22, 10, 0, 13, newYork)
	}}

	o, err := p.GetOptions("foo")
	require.NoError(t, err)
	require.Empty(t, cmp.Diff(o,
		http.SupportBundleCreationOptions{
			Name:        "JFrog Support Case number foo",
			Description: "Generated on 2020-12-04T03:10:00Z",
			Parameters:  nil,
		}))
}
