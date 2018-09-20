package check

import (
	"bytes"
	"errors"
	"regexp"
	"testing"

	"bitbucket.org/asecurityteam/sdcli/internal/check/commands"
	"bitbucket.org/asecurityteam/sdcli/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
)

func TestSuccessAll(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var mockCheckerA = mocks.NewMockchecker(ctrl)
	var mockCheckerB = mocks.NewMockchecker(ctrl)
	mockCheckerA.EXPECT().Name().Return("a")
	mockCheckerB.EXPECT().Name().Return("b")
	mockCheckerA.EXPECT().Check().Return(nil)
	mockCheckerB.EXPECT().Check().Return(nil)

	var checks = []checker{
		mockCheckerA,
		mockCheckerB,
	}

	var cmd = &cobra.Command{}

	var buff = &bytes.Buffer{}
	cmd.SetOutput(buff)
	runChecks(checks)(cmd, nil)

	var matcher = regexp.MustCompile("✓")
	var matches = matcher.FindAllString(buff.String(), -1)
	require.Equal(t, 2, len(matches))
}

func TestSuccessSubcommand(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var a = "a"

	var mockCheckerA = mocks.NewMockchecker(ctrl)
	mockCheckerA.EXPECT().Name().Return(a)
	mockCheckerA.EXPECT().Check().Return(nil)

	var cmd = &cobra.Command{}
	var checks = []checker{
		mockCheckerA,
	}

	var buff = &bytes.Buffer{}
	cmd.SetOutput(buff)
	runChecks(checks)(cmd, []string{a})

	var matcher = regexp.MustCompile("✓")
	var matches = matcher.FindAllString(buff.String(), -1)
	require.Equal(t, 1, len(matches))
}

func TestFailure(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var a = "a"

	var mockCheckerA = mocks.NewMockchecker(ctrl)
	mockCheckerA.EXPECT().Name().Return(a)
	mockCheckerA.EXPECT().Check().Return(&commands.CheckerFailure{})

	var cmd = &cobra.Command{}
	var checks = []checker{
		mockCheckerA,
	}

	var buff = &bytes.Buffer{}
	cmd.SetOutput(buff)
	runChecks(checks)(cmd, []string{a})

	var matcher = regexp.MustCompile("✗ failure")
	var matches = matcher.FindAllString(buff.String(), -1)
	require.Equal(t, 1, len(matches))
}

func TestError(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var a = "a"

	var mockCheckerA = mocks.NewMockchecker(ctrl)
	mockCheckerA.EXPECT().Name().Return(a)
	mockCheckerA.EXPECT().Check().Return(errors.New(""))

	var cmd = &cobra.Command{}
	var checks = []checker{
		mockCheckerA,
	}

	var buff = &bytes.Buffer{}
	cmd.SetOutput(buff)
	runChecks(checks)(cmd, []string{a})

	var matcher = regexp.MustCompile("✗ error")
	var matches = matcher.FindAllString(buff.String(), -1)
	require.Equal(t, 1, len(matches))
}
