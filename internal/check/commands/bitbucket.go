package commands

import (
	"os/exec"

	"bitbucket.org/asecurityteam/sdcli/internal/runner"
	"github.com/pkg/errors"
)

const bitbucketFailure = `
SecDev uses Bitbucket Cloud nearly exclusively.

In order to start contributing, you'll need to register on
bitbucket.org, get access to the asecurityteam organization,
and add your SSH key.

  - Visit https://bitbucket.org/account/signin/ and use the
    "Log in with Google" option.  Use your @atlassian.com
    email address.

  - Follow the steps to get your SSH key added, creating one
    if you need to:
    https://bitbucket.org/account/user/{your_bb_username}/ssh-keys/

  - Finally, contact your manager or teammates to get added to
    asecurityteam.`

// BitbucketChecker verifies the git cli is installed, and that SSH access is configured
type BitbucketChecker struct {
	repo string
	name string
	r    runner.Runner
}

// NewBitbucketChecker creates and returns a new BitbucketChecker
func NewBitbucketChecker(r runner.Runner) *BitbucketChecker {
	return &BitbucketChecker{
		repo: "git@bitbucket.org:asecurityteam/sdcli.git",
		name: "bitbucket",
		r:    r,
	}
}

// Name returns the checker name
func (c *BitbucketChecker) Name() string {
	return c.name
}

// Check runs the checker
// It verifies the system has SSH access to Bitbucket cloud, and that
// the user has access to the asecurityteam organization
func (c *BitbucketChecker) Check() error {
	var err error
	if _, err = c.r.Run("git", "ls-remote", "--heads", "--quiet", c.repo, "master"); err == nil {
		return nil
	}
	if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
		return &CheckerFailure{
			Message: "It does not appear you have git installed. You should install git.\n",
		}
	}
	if _, ok := err.(*exec.ExitError); ok {
		// If it was an ExitError then it didn't return 0, which indicates that the connection failed.
		return &CheckerFailure{
			Message: bitbucketFailure,
		}
	}
	return errors.Wrap(err, "error communicating with bitbucket")
}
