package commands

import "time"

type Clock func() string

func Now() string {
	return time.Now().Format(time.RFC3339)
}
