package commands

import (
	"errors"
	"fmt"
	"os/user"
	"testing"

	"bitbucket.org/asecurityteam/sdcli/internal/mocks"
	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetServiceAndImageName(t *testing.T) {
	tc := []struct {
		Name                string
		IsDev               bool
		ServiceDescriptor   string
		ExpectedImageName   string
		ExpectedServiceName string
		ExpectedErr         bool
	}{
		{
			Name:                "production success",
			ServiceDescriptor:   "my-service.sd.yml",
			ExpectedImageName:   "docker.atl-paas.net/asecurityteam/sox/my-service",
			ExpectedServiceName: "my-service",
		},
		{
			Name:                "dev success",
			IsDev:               true,
			ServiceDescriptor:   "my-service.sd.yml",
			ExpectedImageName:   "docker.atl-paas.net/asecurityteam/my-service",
			ExpectedServiceName: "my-service",
		},
	}

	for _, test := range tc {
		t.Run(test.Name, func(tt *testing.T) {
			runner := runner.ExecRunner{}
			serviceGenerator := newServiceGenerator(runner, test.IsDev, nil)

			serviceName := serviceGenerator.GetServiceName(test.ServiceDescriptor)
			require.Equal(tt, test.ExpectedServiceName, serviceName)

			imageName := serviceGenerator.GetImageName(serviceName)
			require.Equal(tt, test.ExpectedImageName, imageName)
		})
	}
}

func TestGetServiceDescriptor(t *testing.T) {
	tc := []struct {
		Name                      string
		GlobFuncResults           []string
		GlobFuncErr               error
		ExpectedServiceDescriptor string
		ExpectedErr               bool
	}{
		{
			Name:                      "production success",
			GlobFuncResults:           []string{"my-service.sd.yml"},
			ExpectedServiceDescriptor: "my-service.sd.yml",
		},
		{
			Name:            "file not found",
			GlobFuncResults: nil,
			GlobFuncErr:     errors.New("file not found"),
			ExpectedErr:     true,
		},
		{
			Name:            "too many results",
			GlobFuncResults: []string{"sd-1.sd.yml", "sd-2.sd.yml"},
			ExpectedErr:     true,
		},
	}

	for _, test := range tc {
		t.Run(test.Name, func(tt *testing.T) {
			globFunc := func(s string) ([]string, error) {
				return test.GlobFuncResults, test.GlobFuncErr
			}

			runner := runner.ExecRunner{}
			serviceGenerator := newServiceGenerator(runner, false, globFunc)
			actualServiceDescriptor, err := serviceGenerator.GetServiceDescriptor()

			require.Equal(tt, test.ExpectedServiceDescriptor, actualServiceDescriptor)
			if test.ExpectedErr {
				require.NotNil(tt, err)
			}
		})
	}
}

func TestGetTag(t *testing.T) {
	tc := []struct {
		Name                  string
		IsDev                 bool
		UncommittedChangesErr error
		Hash                  string
		HashErr               error
		ExpectedTag           string
		ExpectedErr           bool
	}{
		{
			Name:        "production build without uncommitted changes",
			Hash:        "a1b2c3d",
			ExpectedTag: "a1b2c3d",
		},
		{
			Name:                  "production build with uncommitted changes",
			UncommittedChangesErr: errors.New("changes"),
			Hash:                  "a1b2c3d",
			ExpectedTag:           "",
			ExpectedErr:           true,
		},
		{
			Name:        "dev build without uncommitted changes",
			IsDev:       true,
			Hash:        "a1b2c3d",
			ExpectedTag: "a1b2c3d",
			ExpectedErr: false,
		},
	}

	for _, test := range tc {
		t.Run(test.Name, func(tt *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRunner := mocks.NewMockRunner(ctrl)
			serviceGenerator := newServiceGenerator(mockRunner, test.IsDev, nil)

			mockRunner.EXPECT().Run("git", "diff", "--cached", "--quiet").Return(nil, test.UncommittedChangesErr)
			mockRunner.EXPECT().Run("git", "rev-parse", "--short", "HEAD").Return([]byte(test.Hash), test.HashErr)

			actualTag, err := serviceGenerator.GetTag()

			require.Equal(tt, test.ExpectedTag, actualTag)
			if test.ExpectedErr {
				require.NotNil(t, err)
			}
		})
	}
}

func TestGetTagForDevBuildWithUncommitedChanges(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunner := mocks.NewMockRunner(ctrl)
	serviceGenerator := newServiceGenerator(mockRunner, true, nil)

	currentUser, err := user.Current()
	if err != nil {
		t.Fatal(err.Error())
	}
	hash := "a1b2c3d"
	expectedTag := fmt.Sprintf("%s-%s", hash, currentUser.Username)

	mockRunner.EXPECT().Run("git", "diff", "--cached", "--quiet").Return(nil, errors.New("changes"))
	mockRunner.EXPECT().Run("git", "rev-parse", "--short", "HEAD").Return([]byte(hash), nil)

	actualTag, err := serviceGenerator.GetTag()

	require.Equal(t, expectedTag, actualTag)
	require.Nil(t, err)
}
