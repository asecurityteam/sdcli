package lint

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type testRunner struct {
	exitCode     int
	output       string
	unknownError error
}

func (r *testRunner) Run(_ string, args ...string) ([]byte, error) {
	if r.unknownError != nil {
		return nil, r.unknownError
	}
	cs := []string{"-test.run=TestExec", "--"}
	cs = append(cs, args...)
	cs = append(cs, fmt.Sprintf("%d", r.exitCode), r.output)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_HELPER_FUNC_EXEC=1"}
	out, err := cmd.CombinedOutput()
	return out, err
}

func (r *testRunner) RunEnv(_ []string, _ string, _ ...string) ([]byte, error) {
	return nil, nil
}

func TestLint(t *testing.T) {
	tc := []struct {
		Name           string
		ExpectedError  bool
		ExpectedOutput string
		ExitCode       int
		Error          error
	}{
		{
			Name: "linter-pass",
		},
		{
			Name:           "linter-fail",
			ExpectedOutput: "some feedback from running linter",
			ExitCode:       1,
		},
		{
			Name:           "other-exit-error",
			ExpectedOutput: "error running linter: exit status 2",
			ExitCode:       2,
			ExpectedError:  true,
		},
		{
			Name:           "unknown-error",
			ExpectedOutput: "error running linter: oops",
			Error:          errors.New("oops"),
			ExpectedError:  true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r := &testRunner{exitCode: tt.ExitCode, output: tt.ExpectedOutput, unknownError: tt.Error}
			buff := &bytes.Buffer{}
			cmd := NewCommand(r)
			cmd.SetOutput(buff)

			err := cmd.Execute()
			if tt.ExpectedError {
				assert.NotNil(t, err, "expected error to be non-nil")
				assert.Equal(t, "Error: "+tt.ExpectedOutput+"\n", buff.String())
			}
		})
	}
}

func TestExec(t *testing.T) {
	// This test function should only get run as a subprocess of go test, and exec'd by testRunner
	if os.Getenv("GO_HELPER_FUNC_EXEC") != "1" {
		return
	}

	nArgs := len(os.Args)
	output := os.Args[nArgs-1]
	exitCodeString := os.Args[nArgs-2]
	exitCode, _ := strconv.Atoi(exitCodeString)
	defer os.Exit(exitCode)

	fmt.Printf(output)
}
