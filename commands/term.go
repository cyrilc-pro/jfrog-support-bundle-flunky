package commands

import "github.com/AlecAivazis/survey/v2"

type Terminal struct {
}

func (t *Terminal) AskIncludeLogs() (bool, error) {
	return t.askBoolean("Include logs?")
}
func (t *Terminal) AskIncludeSystem() (bool, error) {
	return t.askBoolean("Include system info?")
}
func (t *Terminal) AskIncludeConfiguration() (bool, error) {
	return t.askBoolean("Include configuration?")
}
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
