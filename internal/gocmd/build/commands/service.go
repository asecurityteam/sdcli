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
	atlassianSecurityRegistry    = "docker.atl-paas.net/asecurityteam"
	atlassianSecuritySoxRegistry = "docker.atl-paas.net/asecurityteam/sox"
)

type globFunc func(string) ([]string, error)

type serviceGenerator struct {
	r        runner.Runner
	globFunc globFunc
	isDev    bool
}

// Service represents information required for deploying a service to micros
type Service struct {
	Image             string
	ImageName         string
	ImageTag          string
	ServiceDescriptor string
	ServiceName       string
}

func newServiceGenerator(r runner.Runner, isDev bool, g globFunc) *serviceGenerator {
	if g == nil {
		g = filepath.Glob
	}
	return &serviceGenerator{r: r, globFunc: g, isDev: isDev}
}

// NewService returns a new service, accepts a runner and globFunc for testing purposes, isDev
// dictates tagging procedures
func NewService(r runner.Runner, isDev bool, g globFunc) (*Service, error) {
	serviceGenerator := newServiceGenerator(r, isDev, g)
	return serviceGenerator.getService()
}

func (s *serviceGenerator) getService() (*Service, error) {
	serviceDescriptor, err := s.getServiceDescriptor()
	if err != nil {
		return nil, err
	}
	tag, err := s.getTag()
	if err != nil {
		return nil, err
	}
	serviceName := s.getServiceName(serviceDescriptor)
	imageName := s.getImageName(serviceName)
	image := fmt.Sprintf("%s:%s", imageName, tag)
	return &Service{
		Image:             image,
		ImageName:         imageName,
		ImageTag:          tag,
		ServiceDescriptor: serviceDescriptor,
		ServiceName:       serviceName,
	}, nil
}

func (s *serviceGenerator) getServiceDescriptor() (string, error) {
	serviceDescriptor, err := s.globFunc("*.sd.yml")
	if err != nil {
		return "", errors.Wrap(err, "error finding service descriptor")
	}
	if len(serviceDescriptor) != 1 {
		return "", errors.Errorf("found %d service descriptors", len(serviceDescriptor))
	}
	return serviceDescriptor[0], nil
}

func (s *serviceGenerator) getServiceName(serviceDescriptor string) string {
	return strings.TrimSuffix(filepath.Base(serviceDescriptor), ".sd.yml")
}

func (s *serviceGenerator) getImageName(serviceName string) string {
	registry := atlassianSecuritySoxRegistry
	if s.isDev {
		registry = atlassianSecurityRegistry
	}
	return fmt.Sprintf("%s/%s", registry, serviceName)
}

func (s *serviceGenerator) getTag() (string, error) {
	var hasUncommittedChanges bool
	_, err := s.r.Run("git", "diff", "--cached", "--quiet")
	if err != nil {
		hasUncommittedChanges = true
	}

	hash, err := s.r.Run("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return "", errors.Wrap(err, "error fetching hash")
	}

	tag := strings.TrimSpace(string(hash))

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
