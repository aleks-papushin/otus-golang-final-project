package collector

import "os/exec"

type CommandExecutor interface {
	Execute(name string, arg ...string) ([]byte, error)
}

type RealCommandExecutor struct{}

func (r *RealCommandExecutor) Execute(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).Output()
}

type MockCommandExecutor struct {
	Output []byte
	Err    error
}

func (m *MockCommandExecutor) Execute(name string, arg ...string) ([]byte, error) {
	return m.Output, m.Err
}
