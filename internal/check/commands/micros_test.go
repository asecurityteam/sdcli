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

func TestMicrosCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Output         []byte
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name:   "success",
			Output: []byte("6.1.1"),
		},
		{
			Name:           "old version",
			Output:         []byte("6.0.0"),
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:        "bad output",
			Output:      []byte("bad output"),
			ExpectError: true,
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
			mockRunner.EXPECT().Run("micros", "version").Return(tt.Output, tt.Error)
			var checker = NewMicrosChecker(mockRunner)
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

func TestMicrosName(t *testing.T) {
	var checker = NewMicrosChecker(nil)
	if checker.Name() != "micros" {
		t.Fatalf("Expected checker name to be micros, but was %s", checker.Name())
	}
}
