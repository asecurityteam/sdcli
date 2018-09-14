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

const dockerRegistryOutput = `{"auths":{"docker.atl-paas.net":{},"https://docker.atl-paas.net":{}},"HttpHeaders":{"User-Agent":"Docker-Client/17.12.0-ce (darwin)"},"credsStore":"%s"}`

func TestRegistryCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Output         []byte
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name:   "success",
			Output: []byte(fmt.Sprintf(dockerRegistryOutput, "foobar")),
		},
		{
			Name:        "bad output",
			Output:      []byte("bad output"),
			ExpectError: true,
		},
		{
			Name:           "not authenticated",
			Error:          &exec.ExitError{},
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:           "no auths",
			Output:         []byte("{}"),
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:           "no creds",
			Output:         []byte(fmt.Sprintf(dockerRegistryOutput, "")),
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
			mockRunner.EXPECT().Run("cat", gomock.Any()).Return(tt.Output, tt.Error)
			var checker = NewRegistryChecker(mockRunner)
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

func TestRegistryName(t *testing.T) {
	var checker = NewRegistryChecker(nil)
	if checker.Name() != "docker-registry" {
		t.Fatalf("Expected checker name to be docker-registry, but was %s", checker.Name())
	}
}
