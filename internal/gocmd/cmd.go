// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gocmd

import (
	"bitbucket.org/asecurityteam/sdcli/internal/gocmd/build"
	"bitbucket.org/asecurityteam/sdcli/internal/gocmd/test"
	"github.com/spf13/cobra"
)

// goCmd represents the go command
func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "go",
		Short: "tools for building, testing, and deploying golang projects",
	}
	command.AddCommand(build.NewCommand())
	command.AddCommand(test.NewCommand())
	return command
}
