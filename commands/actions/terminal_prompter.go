package actions

import "github.com/AlecAivazis/survey/v2"

// TerminalPrompter is a Prompter that gets answers through questions to the user.
type TerminalPrompter struct {
}

// AskIncludeLogs tells if logs must be included.
func (t *TerminalPrompter) AskIncludeLogs() (bool, error) {
	return t.askBoolean("Include logs?")
}

// AskIncludeSystem tells if system info must be included.
func (t *TerminalPrompter) AskIncludeSystem() (bool, error) {
	return t.askBoolean("Include system info?")
}

// AskIncludeConfiguration tells if configuration must be included.
func (t *TerminalPrompter) AskIncludeConfiguration() (bool, error) {
	return t.askBoolean("Include configuration?")
}

// AskThreadDump tells if thread dumps must be included.
func (t *TerminalPrompter) AskThreadDump() (bool, error) {
	return t.askBoolean("Include thread dump?")
}

func (t *TerminalPrompter) askBoolean(question string) (bool, error) {
	answer := false
	confirm := &survey.Confirm{
		Message: question,
	}
	err := survey.AskOne(confirm, &answer)
	return answer, err
}
