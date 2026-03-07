package zap

type SugaredLogger struct{}

func S() *SugaredLogger {
	return &SugaredLogger{}
}

func (s *SugaredLogger) Info(msg string, args ...any)  {}
func (s *SugaredLogger) Error(msg string, args ...any) {}
func (s *SugaredLogger) Warn(msg string, args ...any)  {}
func (s *SugaredLogger) Debug(msg string, args ...any) {}
func (s *SugaredLogger) Panic(msg string, args ...any) {}
func (s *SugaredLogger) Fatal(msg string, args ...any) {}
