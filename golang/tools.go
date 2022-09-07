//go:build tools

package golang

import (
	_ "github.com/AlekSi/gocov-xml"
	_ "github.com/axw/gocov"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/wadey/gocovmerge"
)
