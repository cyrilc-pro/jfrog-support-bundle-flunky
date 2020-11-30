package test

import (
	"fmt"
	"github.com/jfrog/jfrog-client-go/utils/log"
	"io"
	"testing"
)

type testLog struct {
	t *testing.T
}

func (l *testLog) GetLogLevel() log.LevelType {
	return log.DEBUG
}
func (l *testLog) SetLogLevel(_ log.LevelType) {}
func (l *testLog) SetOutputWriter(_ io.Writer) {}
func (l *testLog) SetLogsWriter(_ io.Writer)   {}
func (l *testLog) Debug(a ...interface{}) {
	l.print("DEBUG", a)
}
func (l *testLog) Info(a ...interface{}) {
	l.print("INFO ", a)
}
func (l *testLog) Warn(a ...interface{}) {
	l.print("WARN ", a)
}
func (l *testLog) Error(a ...interface{}) {
	l.print("ERROR", a)
}
func (l *testLog) Output(a ...interface{}) {
	l.print("OUT  ", a)
}

func (l *testLog) print(level string, msgParts ...interface{}) {
	msg := level
	for i := range msgParts {
		msg += " " + fmt.Sprintf("%v", msgParts[i])
	}
	l.t.Log(msg)
}
