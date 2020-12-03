package actions

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_FormattedString(t *testing.T) {
	newYork, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	tests := []struct {
		in     time.Time
		expect string
	}{
		{
			in:     time.Unix(0, 0),
			expect: "1970-01-01T00:00:00Z",
		},
		{
			in:     time.Date(2020, 12, 3, 22, 10, 0, 13, time.UTC),
			expect: "2020-12-03T22:10:00Z",
		},
		{
			in:     time.Date(2020, 12, 3, 22, 10, 0, 13, newYork),
			expect: "2020-12-04T03:10:00Z",
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.expect, func(t *testing.T) {
			assert.Equal(t, test.expect, formattedString(test.in))
		})
	}
}
