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

const depOutputTpl = `dep:
 version     : %s
 build date  : 2018-07-26
 git hash    : 224a564
 go version  : go1.10.3
 go compiler : gc
 platform    : darwin/amd64
 features    : ImportDuringSolve=false
`

func TestDepCheck(t *testing.T) {
	var tc = []struct {
		Name           string
		Output         []byte
		Error          error
		ExpectError    bool
		CheckerFailure bool
	}{
		{
			Name:   "success",
			Output: []byte(fmt.Sprintf(depOutputTpl, "0.5.0")),
		},
		{
			Name:           "old version",
			Output:         []byte(fmt.Sprintf(depOutputTpl, "0.4.0")),
			ExpectError:    true,
			CheckerFailure: true,
		},
		{
			Name:           "bad output",
			Output:         []byte("bad output"),
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
			mockRunner.EXPECT().Run("dep", "version").Return(tt.Output, tt.Error)
			var checker = NewDepChecker(mockRunner)
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

func TestDepName(t *testing.T) {
	var checker = NewDepChecker(nil)
	if checker.Name() != "dep" {
		t.Fatalf("Expected checker name to be dep, but was %s", checker.Name())
	}
}
