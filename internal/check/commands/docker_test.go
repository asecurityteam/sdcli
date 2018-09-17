package commands

import (
	"errors"
	"fmt"
	"os/exec"
	"testing"

	"bitbucket.org/asecurityteam/sdcli/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const dockerVersionOutputTpl = `{"Client":{"Platform":{"Name":""},"Version":"%s","ApiVersion":"1.35","DefaultAPIVersion":"1.35","GitCommit":"c97c6d6","GoVersion":"go1.9.2","Os":"darwin","Arch":"amd64","BuildTime":"Wed Dec 27 20:03:51 2017"},"Server":{"Platform":{"Name":""},"Components":[{"Name":"Engine","Version":"17.12.0-ce","Details":{"ApiVersion":"1.35","Arch":"amd64","BuildTime":"Wed Dec 27 20:12:29 2017","Experimental":"true","GitCommit":"c97c6d6","GoVersion":"go1.9.2","KernelVersion":"4.9.60-linuxkit-aufs","MinAPIVersion":"1.12","Os":"linux"}}],"Version":"%s","ApiVersion":"1.35","MinAPIVersion":"1.12","GitCommit":"c97c6d6","GoVersion":"go1.9.2","Os":"linux","Arch":"amd64","KernelVersion":"4.9.60-linuxkit-aufs","Experimental":true,"BuildTime":"2017-12-27T20:12:29.000000000+00:00"}}`

func TestDockerCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Output         []byte
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name:   "success",
			Output: []byte(fmt.Sprintf(dockerVersionOutputTpl, "17.12.0-ce", "17.12.0-ce")),
		},
		{
			Name:           "old version",
			Output:         []byte(fmt.Sprintf(dockerVersionOutputTpl, "17.10.0-ce", "17.10.0-ce")),
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:        "bad output",
			Output:      []byte("bad output"),
			ExpectError: true,
		},
		{
			Name:           "daemon error",
			Error:          &exec.ExitError{},
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name: "not installed",
			Error: &exec.Error{
				Err: exec.ErrNotFound,
			},
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:        "unkown error",
			Error:       errors.New("oops"),
			ExpectError: true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			var ctrl = gomock.NewController(t)
			defer ctrl.Finish()

			var mockRunner = mocks.NewMockRunner(ctrl)
			mockRunner.EXPECT().Run("docker", "version", "-f", "{{json .}}").Return(tt.Output, tt.Error)
			var checker = NewDockerChecker(mockRunner)
			var e = checker.Check()
			if !tt.ExpectError {
				require.Nil(t, e)
			}
			if tt.ExpectError {
				require.NotNil(t, e)
			}
			if tt.CheckerFailure {
				var _, ok = e.(*CheckerFailure)
				require.True(t, ok, fmt.Sprintf("Expected error to be CheckerFailure but was %T", e))
			}
		})
	}
}

func TestDockerName(t *testing.T) {
	var checker = NewDockerChecker(nil)
	if checker.Name() != "docker" {
		t.Fatalf("Expected checker name to be docker, but was %s", checker.Name())
	}
}
