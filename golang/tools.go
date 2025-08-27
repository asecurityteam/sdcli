//go:build tools

package golang

import (
	_ "github.com/AlekSi/gocov-xml"
	_ "github.com/axw/gocov/gocov"                             //nb gocov/gocov for the tool binaries
// see Dockerfile, we should not install it like that	_ "github.com/golangci/golangci-lint/v2/cmd/golangci-lint" //nb postfix of cmd/golangci-lint for the tool binaries
	_ "github.com/wadey/gocovmerge"
)
