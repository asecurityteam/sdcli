package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoverPkg(t *testing.T) {
	pkg, err := getCoverPkg("bitbucket.org/asecurityteam/sdcli/...")
	assert.Nil(t, err)
	assert.NotContains(t, pkg, "bitbucket.org/asecurityteam/sdcli,")
}
