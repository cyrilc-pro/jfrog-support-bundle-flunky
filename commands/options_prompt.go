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

type promptOptionsProvider struct {
	getDate  func() time.Time
	prompter Prompter
}

func (p *promptOptionsProvider) GetOptions(caseNumber string) (SupportBundleCreationOptions, error) {
	options, err := (&defaultOptionsProvider{getDate: p.getDate}).GetOptions(caseNumber)
	if err != nil {
		return options, err
	}
	options.Parameters = &SupportBundleParameters{Logs: &SupportBundleParametersLogs{}, ThreadDump: &SupportBundleParametersThreadDump{}}

	if options.Parameters.Logs.Include, err = p.prompter.AskIncludeLogs(); err != nil {
		return options, err
	}
	if options.Parameters.Configuration, err = p.prompter.AskIncludeConfiguration(); err != nil {
		return options, err
	}
	if options.Parameters.System, err = p.prompter.AskIncludeSystem(); err != nil {
		return options, err
	}
	if askThreadDump, err := p.prompter.AskThreadDump(); err != nil {
		return options, err
	} else if askThreadDump {
		options.Parameters.ThreadDump.Count = 1
		options.Parameters.ThreadDump.Interval = 0
	}

	now := p.getDate()
	yesterday := now.Add(-24 * time.Hour)
	options.Parameters.Logs.StartDate = yesterday.Format("2006-01-02")
	options.Parameters.Logs.EndDate = now.Format("2006-01-02")

	return options, nil
}
