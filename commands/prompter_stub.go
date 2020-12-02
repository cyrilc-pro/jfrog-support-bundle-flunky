package commands

type PrompterStub struct {
	IncludeLogs          bool
	IncludeSystem        bool
	IncludeConfiguration bool
	IncludeThreadDump    bool
	err                  error
}

func (s *PrompterStub) AskIncludeLogs() (bool, error) {
	return s.IncludeLogs, s.err
}
func (s *PrompterStub) AskIncludeSystem() (bool, error) {
	return s.IncludeSystem, s.err
}
func (s *PrompterStub) AskIncludeConfiguration() (bool, error) {
	return s.IncludeConfiguration, s.err
}
func (s *PrompterStub) AskThreadDump() (bool, error) {
	return s.IncludeThreadDump, s.err
}
