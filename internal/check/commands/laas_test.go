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

func TestLaasCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Output         []byte
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name:   "success",
			Output: []byte("cached credentials for user"),
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
			Name:           "unauthenticated",
			Output:         []byte("Not currently logged in"),
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
			var checker = NewLaasChecker(mockRunner)

			mockRunner.EXPECT().Run("laas", "whoami").Return(tt.Output, tt.Error)

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

func TestLaasName(t *testing.T) {
	var checker = NewLaasChecker(nil)
	if checker.Name() != "laas" {
		t.Fatalf("Expected checker name to be laas, but was %s", checker.Name())
	}
}
