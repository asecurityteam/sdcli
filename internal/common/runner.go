package common

import "os/exec"

// Runner is a function which runs the command with the supplied parameters.
// The output from the command is is returned as a []byte, or an error if the command
// failed or was not found.
type Runner interface {
	Run(string, ...string) ([]byte, error)
}

// ExecRunner executes cmd with args using exec.Command.
type ExecRunner struct{}

// Run runs the ExecRunner
func (r ExecRunner) Run(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}
