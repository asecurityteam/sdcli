package internal

import (
	"bytes"
	"testing"
)

func TestDepCommand(t *testing.T) {
	var out bytes.Buffer
	depCommand := depCmd()

	depCommand.SetOutput(&out)
	_ = depCommand.Execute()

	result := out.String()
	if result != "dep called, no args" {
		t.Errorf("Expected %s, got: %s", "dep called", result)
	}
}

func TestDepCommandWithArg(t *testing.T) {
	var out bytes.Buffer
	depCommand := depCmd()

	depCommand.SetOutput(&out)
	depCommand.SetArgs([]string{"arg1"})
	_ = depCommand.Execute()

	result := out.String()
	expected := "dep called with arg1"
	if result != expected {
		t.Errorf("Expected %s, got: %s", expected, result)
	}
}
