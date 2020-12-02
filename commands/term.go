package commands

import "github.com/AlecAivazis/survey/v2"

// Terminal is a Prompter that gets answers through questions to the user.
type Terminal struct {
}

// AskIncludeLogs tells if logs must be included.
func (t *Terminal) AskIncludeLogs() (bool, error) {
	return t.askBoolean("Include logs?")
}

// AskIncludeSystem tells if system info must be included.
func (t *Terminal) AskIncludeSystem() (bool, error) {
	return t.askBoolean("Include system info?")
}

// AskIncludeConfiguration tells if configuration must be included.
func (t *Terminal) AskIncludeConfiguration() (bool, error) {
	return t.askBoolean("Include configuration?")
}

// AskThreadDump tells if thread dumps must be included.
func (t *Terminal) AskThreadDump() (bool, error) {
	return t.askBoolean("Include thread dump?")
}

func (t *Terminal) askBoolean(question string) (bool, error) {
	answer := false
	confirm := &survey.Confirm{
		Message: question,
	}
	err := survey.AskOne(confirm, &answer)
	return answer, err
}
