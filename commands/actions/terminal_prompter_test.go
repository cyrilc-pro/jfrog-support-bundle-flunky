package actions

import (
	"bytes"
	"errors"
	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

var cursorPositionANSISequence = []byte{0x1B, '[', '0', ';', '0', 'R'}

type fileReaderStub struct {
	r io.Reader
}

func (s *fileReaderStub) Read(p []byte) (n int, err error) {
	read, err := s.r.Read(p)
	if errors.Is(err, io.EOF) {
		// Once the answer has been provided, always provide the cursor position
		return bytes.NewReader(cursorPositionANSISequence).Read(p)
	}
	return read, err
}
func (s *fileReaderStub) Fd() uintptr {
	return os.Stdin.Fd()
}

type fileWriterStub struct{}

func (s *fileWriterStub) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}
func (s *fileWriterStub) Fd() uintptr {
	return os.Stdout.Fd()
}

func TestTerminalPrompter_AskIncludeLogs(t *testing.T) {
	test(t, func(prompter TerminalPrompter) (bool, error) {
		return prompter.AskIncludeLogs()
	})
}

func TestTerminalPrompter_AskIncludeSystem(t *testing.T) {
	test(t, func(prompter TerminalPrompter) (bool, error) {
		return prompter.AskIncludeSystem()
	})
}

func TestTerminalPrompter_AskIncludeConfiguration(t *testing.T) {
	test(t, func(prompter TerminalPrompter) (bool, error) {
		return prompter.AskIncludeConfiguration()
	})
}

func TestTerminalPrompter_AskThreadDump(t *testing.T) {
	test(t, func(prompter TerminalPrompter) (bool, error) {
		return prompter.AskThreadDump()
	})
}

func test(t *testing.T, ask func(TerminalPrompter) (bool, error)) {
	tests := map[string]bool{
		"y":   true,
		"yes": true,
		"n":   false,
		"no":  false,
		"":    true,
	}

	for input, expect := range tests {
		input += "\n"
		inputs := []string{
			input,
			"bad input\n" + input,
		}
		if input != "\n" {
			inputs = append(inputs, strings.ToUpper(input), strings.ToUpper("ZZZ\n"+input))
		}

		for _, in := range inputs {
			testIn := in
			testExpect := expect
			t.Run(testIn, func(t *testing.T) {
				readerStub := &fileReaderStub{strings.NewReader(testIn)}
				prompter := TerminalPrompter{opts: []survey.AskOpt{survey.WithStdio(readerStub, &fileWriterStub{}, os.Stderr)}}
				b, err := ask(prompter)
				require.NoError(t, err)
				assert.Equal(t, testExpect, b)
			})
		}
	}
}
