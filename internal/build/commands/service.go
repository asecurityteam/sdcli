package commands

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const (
	AtlassianSecurityRegistry    = "docker.atl-paas.net/asecurityteam"
	AtlassianSecuritySoxRegistry = "docker.atl-paas.net/asecurityteam/sox"
)

type GlobFunc func(string) ([]string, error)

type ServiceGenerator struct {
	r        runner.Runner
	globFunc GlobFunc
	isDev    bool
}

type Service struct {
	Image             string
	ImageName         string
	ImageTag          string
	ServiceDescriptor string
	ServiceName       string
}

func newServiceGenerator(r runner.Runner, isDev bool, g GlobFunc) *ServiceGenerator {
	if g == nil {
		g = filepath.Glob
	}
	return &ServiceGenerator{r: r, globFunc: g, isDev: isDev}
}

func NewService(r runner.Runner, isDev bool, g GlobFunc) (*Service, error) {
	serviceGenerator := newServiceGenerator(r, isDev, g)
	return serviceGenerator.GetService()
}

func (s *ServiceGenerator) GetService() (*Service, error) {
	serviceDescriptor, err := s.GetServiceDescriptor()
	if err != nil {
		return nil, err
	}
	tag, err := s.GetTag()
	if err != nil {
		return nil, err
	}
	serviceName := s.GetServiceName(serviceDescriptor)
	imageName := s.GetImageName(serviceName)
	image := fmt.Sprintf("%s:%s", imageName, tag)
	return &Service{
		Image:             image,
		ImageName:         imageName,
		ImageTag:          tag,
		ServiceDescriptor: serviceDescriptor,
		ServiceName:       serviceName,
	}, nil
}

func (s *ServiceGenerator) GetServiceDescriptor() (string, error) {
	serviceDescriptor, err := s.globFunc("*.sd.yml")
	if err != nil {
		return "", errors.Wrap(err, "error finding service descriptor")
	}
	if len(serviceDescriptor) != 1 {
		return "", errors.Errorf("found %d service descriptors", len(serviceDescriptor))
	}
	return serviceDescriptor[0], nil
}

func (s *ServiceGenerator) GetServiceName(serviceDescriptor string) string {
	return strings.TrimSuffix(filepath.Base(serviceDescriptor), ".sd.yml")
}

func (s *ServiceGenerator) GetImageName(serviceName string) string {
	registry := AtlassianSecuritySoxRegistry
	if s.isDev {
		registry = AtlassianSecurityRegistry
	}
	return fmt.Sprintf("%s/%s", registry, serviceName)
}

func (s *ServiceGenerator) GetTag() (string, error) {
	var hasUncommittedChanges bool
	var tag string
	_, err := s.r.Run("git", "diff", "--cached", "--quiet")
	if err != nil {
		hasUncommittedChanges = true
	}
	hash, err := s.r.Run("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", errors.Wrap(err, "error fetching hash")
	}
	tag = strings.TrimSpace(string(hash))

	if s.isDev {
		if hasUncommittedChanges {
			currentUser, err := user.Current()
			if err != nil {
				return "", errors.Wrap(err, "error getting current user")
			}
			tag = fmt.Sprintf("%s-%s", tag, currentUser.Username)
		}
	} else {
		if hasUncommittedChanges {
			return "", errors.New("uncommited changes")
		}
	}

	return tag, nil
}
