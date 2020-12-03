package actions

import "time"

// Clock is a provider of time
type Clock func() time.Time

func formattedString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
