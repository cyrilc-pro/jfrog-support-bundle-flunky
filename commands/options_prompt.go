package commands

import (
	"time"
)

type Prompter interface {
	AskIncludeLogs() (bool, error)
	AskIncludeSystem() (bool, error)
	AskIncludeConfiguration() (bool, error)
	AskThreadDump() (bool, error)
}

type PromptOptionsProvider struct {
	GetDate  func() time.Time
	Prompter Prompter
}

func newPromptOptionsProvider() optionsProvider {
	return &PromptOptionsProvider{
		GetDate:  time.Now,
		Prompter: &Terminal{},
	}
}

func (p *PromptOptionsProvider) GetOptions(caseNumber string) (SupportBundleCreationOptions, error) {
	options, err := (&DefaultOptionsProvider{GetDate: p.GetDate}).GetOptions(caseNumber)
	if err != nil {
		return options, err
	}
	options.Parameters = &SupportBundleParameters{Logs: &SupportBundleParametersLogs{}, ThreadDump: &SupportBundleParametersThreadDump{}}

	if options.Parameters.Logs.Include, err = p.Prompter.AskIncludeLogs(); err != nil {
		return options, err
	}
	if options.Parameters.Configuration, err = p.Prompter.AskIncludeConfiguration(); err != nil {
		return options, err
	}
	if options.Parameters.System, err = p.Prompter.AskIncludeSystem(); err != nil {
		return options, err
	}
	if askThreadDump, err := p.Prompter.AskThreadDump(); err != nil {
		return options, err
	} else if askThreadDump {
		options.Parameters.ThreadDump.Count = 1
		options.Parameters.ThreadDump.Interval = 0
	}

	now := p.GetDate()
	yesterday := now.Add(-24 * time.Hour)
	options.Parameters.Logs.StartDate = yesterday.Format("2006-01-02")
	options.Parameters.Logs.EndDate = now.Format("2006-01-02")

	return options, nil
}
