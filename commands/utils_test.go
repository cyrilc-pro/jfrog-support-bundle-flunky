package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_getDurationFromContext(t *testing.T) {
	defaultDuration := 10 * time.Hour
	tests := []struct {
		name     string
		value    string
		expected time.Duration
	}{
		{
			name:     "empty string uses default",
			value:    "",
			expected: defaultDuration,
		},
		{
			name:     "unparsable duration uses default",
			value:    "30 seconds",
			expected: defaultDuration,
		},
		{
			name:     "valid duration",
			value:    "25s",
			expected: 25 * time.Second,
		},
	}

	for i := range tests {
		test := tests[i]
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, getDurationOrDefault(test.value, defaultDuration))
		})
	}
}
