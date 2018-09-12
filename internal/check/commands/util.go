package commands

import (
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

func checkMinimumVersion(version, constraint string) (bool, error) {
	var err error
	var c *semver.Constraints
	var v *semver.Version
	c, err = semver.NewConstraint(">= " + constraint)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("invalid constraint %s", constraint))
	}
	v, err = semver.NewVersion(version)
	if err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("invalid version %s", version))
	}

	var ok, _ = c.Validate(v)
	return ok, nil
}
