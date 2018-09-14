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

func TestLinterCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name: "success",
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
			var checker = NewLinterChecker(mockRunner)

			mockRunner.EXPECT().Run("golangci-lint", "linters").Return(nil, tt.Error)

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

func TestLinterName(t *testing.T) {
	var checker = NewLinterChecker(nil)
	if checker.Name() != "golangci-lint" {
		t.Fatalf("Expected checker name to be golangci-lint, but was %s", checker.Name())
	}
}
