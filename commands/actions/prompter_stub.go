package actions

// PrompterStub is a stub for a Prompter, used for tests.
type PrompterStub struct {
	IncludeLogs             bool
	IncludeSystem           bool
	IncludeConfiguration    bool
	IncludeThreadDump       bool
	IncludeLogsErr          error
	IncludeSystemErr        error
	IncludeConfigurationErr error
	IncludeThreadDumpErr    error
}

// AskIncludeLogs tells if logs must be included.
func (s *PrompterStub) AskIncludeLogs() (bool, error) {
	return s.IncludeLogs, s.IncludeLogsErr
}

// AskIncludeSystem tells if system info must be included.
func (s *PrompterStub) AskIncludeSystem() (bool, error) {
	return s.IncludeSystem, s.IncludeSystemErr
}

// AskIncludeConfiguration tells if configuration must be included.
func (s *PrompterStub) AskIncludeConfiguration() (bool, error) {
	return s.IncludeConfiguration, s.IncludeConfigurationErr
}

// AskThreadDump tells if thread dumps must be included.
func (s *PrompterStub) AskThreadDump() (bool, error) {
	return s.IncludeThreadDump, s.IncludeThreadDumpErr
}
