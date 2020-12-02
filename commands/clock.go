package commands

import "time"

// Clock is a provider of time
type Clock func() time.Time

func toString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
