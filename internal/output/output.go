package output

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

// Check outputs a message with a preceding green checkmark
func Check(out io.Writer, msg string, params ...interface{}) {
	fmt.Fprintf(out, color.GreenString("✓")+" "+msg+"\n", params...)
}

// Fail outputs a message with a preceding red X
func Fail(out io.Writer, msg string, params ...interface{}) {
	fmt.Fprintf(out, color.RedString("✗")+" "+msg+"\n", params...)
}
