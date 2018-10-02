package commands

import (
	"testing"

	"bitbucket.org/asecurityteam/sdcli/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestDeployCommand(t *testing.T) {
	service := &Service{
		Image:             "my-image-name:12345",
		ImageName:         "my-image-name",
		ImageTag:          "12345",
		ServiceName:       "my-service",
		ServiceDescriptor: "my-service.sd.yml",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRunner := mocks.NewMockRunner(ctrl)
	docker := NewDocker(mockRunner)
	deployCmd := NewDeployCommand(mockRunner, docker)

	mockRunner.EXPECT().Run("docker", "build", "-t", "my-image-name:12345", ".").Return(nil, nil)
	mockRunner.EXPECT().Run("docker", "push", "my-image-name:12345").Return(nil, nil)
	mockRunner.EXPECT().RunEnv(
		[]string{"DOCKER_IMAGE=my-image-name", "DOCKER_TAG=12345"},
		"micros",
		"service:deploy",
		"my-service",
		"-f",
		"my-service.sd.yml",
	).Return(nil, nil)

	e := deployCmd.deploy(service)
	require.Nil(t, e)
}
