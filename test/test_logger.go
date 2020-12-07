package test

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"testing"
)

type testLogger struct {
	t *testing.T
}

func (l *testLogger) GetLogLevel() log.LevelType {
	return log.DEBUG
}
func (l *testLogger) SetLogLevel(_ log.LevelType) {}
func (l *testLogger) SetOutputWriter(_ io.Writer) {}
func (l *testLogger) SetLogsWriter(_ io.Writer)   {}
func (l *testLogger) Debug(a ...interface{}) {
	l.print("DEBUG", a)
}
func (l *testLogger) Info(a ...interface{}) {
	l.print("INFO ", a)
}
func (l *testLogger) Warn(a ...interface{}) {
	l.print("WARN ", a)
}
func (l *testLogger) Error(a ...interface{}) {
	l.print("ERROR", a)
}
func (l *testLogger) Output(a ...interface{}) {
	l.print("OUT  ", a)
}

func (l *testLogger) print(level string, msgParts ...interface{}) {
	msg := level
	for i := range msgParts {
		msg += " " + fmt.Sprintf("%v", msgParts[i])
	}
	l.t.Log(msg)
}
