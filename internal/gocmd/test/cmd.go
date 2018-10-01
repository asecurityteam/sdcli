package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	integrationFlag            = "integration"
	unitCoverageProfile        = ".coverage/unit.cover.out"
	integrationCoverageProfile = ".coverage/integration.cover.out"
	combinedCoverageProfile    = ".coverage/combined.cover.out"
	coverageDir                = ".coverage"
	allTestPattern             = "./..."
	integrationTestPattern     = "./tests/"
)

var baseTestArguments = [4]string{"test", "-race", "-v", "-cover"}

// NewCommand returns a new test command
func NewCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "test",
		Short: "run unit/integration tests and generate coverage reports",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := CreateCoverageDir(); err != nil {
				return err
			}

			integration, err := cmd.Flags().GetBool(integrationFlag)
			if err != nil {
				return errors.Wrap(err, "error getting integration flag")
			}

			var cmdOutput []byte
			if integration && HasIntegrationTests() {
				cmdOutput, err = RunTests(integrationCoverageProfile, integrationTestPattern)
			}
			cmdOutput, err = RunTests(unitCoverageProfile, allTestPattern)
			if err != nil {
				return err
			}
			cmd.Printf("%s\n", cmdOutput)

			return nil
		},
	}

	command.Flags().BoolP(integrationFlag, "i", false, "Run integration tests")
	command.AddCommand(CoverageCommand())

	return command
}

func CoverageCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "coverage",
		Aliases: []string{"cov"},
		Short:   "produce test coverage for unit and integration tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := CreateCoverageDir(); err != nil {
				return err
			}
			coverageFiles, err := filepath.Glob(".coverage/*.cover.out")
			if err != nil {
				return errors.Wrap(err, "error globbing coverage directory")
			}
			if len(coverageFiles) == 0 {
				return errors.New("no coverage files found")
			}
			gocovMergeOutput, err := exec.Command("gocovmerge", coverageFiles...).Output()
			if err != nil {
				return errors.Wrap(err, "error merging coverage")
			}
			mergedCoverage, err := os.Create(combinedCoverageProfile)
			if err != nil {
				return errors.Wrap(err, "error creating combined coverage file")
			}
			defer mergedCoverage.Close()
			if _, err := mergedCoverage.Write(gocovMergeOutput); err != nil {
				return err
			}
			report, _ := exec.Command("go", "tool", "cover", "-func", combinedCoverageProfile).CombinedOutput()
			cmd.Printf("%s\n", report)
			return nil
		},
	}
}

func CreateCoverageDir() error {
	if err := os.Mkdir(coverageDir, os.ModeDir|os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func HasIntegrationTests() bool {
	if _, err := os.Stat(integrationTestPattern); os.IsNotExist(err) {
		return false
	}
	return true
}

func RunTests(coverageProfile, testDir string) ([]byte, error) {
	testArgs := append(baseTestArguments[:], []string{"-coverprofile", coverageProfile, testDir}...)
	testOutput, err := exec.Command("go", testArgs...).CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "error running tests")
	}

	gocovConvert, err := exec.Command("gocov", "convert", coverageProfile).Output()
	if err != nil {
		return nil, errors.Wrap(err, "gocov: error converting coverage")
	}

	xmlFile := strings.Replace(coverageProfile, ".cover.out", ".xml", 1)
	gocovXML := exec.Command("gocov-xml")
	gocovXML.Stdin = bytes.NewBuffer(gocovConvert)
	xmlCoverage, err := gocovXML.Output()
	if err != nil {
		return nil, errors.Wrap(err, "gocov-xml: error converting coverage to xml")
	}
	xmlCoverageProfile, err := os.Create(xmlFile)
	if err != nil {
		return nil, errors.Wrap(err, "could not create xml coverage profile")
	}
	defer xmlCoverageProfile.Close()
	if _, err := xmlCoverageProfile.Write(xmlCoverage); err != nil {
		return nil, err
	}

	return testOutput, nil
}
