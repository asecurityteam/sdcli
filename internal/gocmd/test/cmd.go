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
			if err := createCoverageDir(); err != nil {
				return err
			}

			integration, err := cmd.Flags().GetBool(integrationFlag)
			if err != nil {
				return errors.Wrap(err, "error getting integration flag")
			}

			var cmdOutput []byte
			if integration && hasIntegrationTests() {
				cmdOutput, err = runTests(integrationCoverageProfile, []string{integrationTestPattern})
			} else {
				allPackages, err := exec.Command("go", "list", allTestPattern).Output()
				if err != nil {
					return errors.Wrap(err, "error listing packages")
				}
				filterPackages := exec.Command("grep", "-v", "-e", "/inttest$")
				filterPackages.Stdin = bytes.NewBuffer(allPackages)

				filterPackagesOutput, err := filterPackages.Output()
				if err != nil {
					return errors.Wrap(err, "error excluding packages")
				}
				unitTestPackages := strings.Split(strings.TrimSpace(string(filterPackagesOutput)), "\n")
				cmdOutput, err = runTests(unitCoverageProfile, unitTestPackages)
			}
			if err != nil {
				return err
			}
			cmd.Printf("%s\n", cmdOutput)

			return nil
		},
	}

	command.Flags().BoolP(integrationFlag, "i", false, "Run integration tests")
	command.AddCommand(coverageCommand())

	return command
}

func coverageCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "coverage",
		Aliases: []string{"cov"},
		Short:   "produce test coverage for unit and integration tests",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := createCoverageDir(); err != nil {
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

			if err = createXMLCoverage(combinedCoverageProfile); err != nil {
				return err
			}
			report, _ := exec.Command("go", "tool", "cover", "-func", combinedCoverageProfile).CombinedOutput()
			cmd.Printf("%s\n", report)
			return nil
		},
	}
}

func createCoverageDir() error {
	if err := os.Mkdir(coverageDir, os.ModeDir|os.ModePerm); err != nil {
		if os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func hasIntegrationTests() bool {
	if _, err := os.Stat(integrationTestPattern); os.IsNotExist(err) {
		return false
	}
	return true
}

func createXMLCoverage(coverageProfile string) error {
	gocovConvert, err := exec.Command("gocov", "convert", coverageProfile).Output()
	if err != nil {
		return errors.Wrap(err, "gocov: error converting coverage")
	}

	xmlFile := strings.Replace(coverageProfile, ".cover.out", ".xml", 1)
	gocovXML := exec.Command("gocov-xml")
	gocovXML.Stdin = bytes.NewBuffer(gocovConvert)
	xmlCoverage, err := gocovXML.Output()
	if err != nil {
		return errors.Wrap(err, "gocov-xml: error converting coverage to xml")
	}
	xmlCoverageProfile, err := os.Create(xmlFile)
	if err != nil {
		return errors.Wrap(err, "could not create xml coverage profile")
	}
	defer xmlCoverageProfile.Close()
	if _, err := xmlCoverageProfile.Write(xmlCoverage); err != nil {
		return err
	}
	return nil
}

func runTests(coverageProfile string, testDirs []string) ([]byte, error) {
	testArgs := append(baseTestArguments[:], []string{"-coverprofile", coverageProfile}...)
	testArgs = append(testArgs, testDirs...)
	testOutput, err := exec.Command("go", testArgs...).CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "error running tests")
	}

	if err = createXMLCoverage(coverageProfile); err != nil {
		return nil, err
	}

	return testOutput, nil
}
