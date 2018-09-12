package common

import (
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"

	"github.com/fatih/color"
)

var whitespaceOnlyRx = regexp.MustCompile(`(?m)^[ \t]+$`)
var leadingWhitespaceRx = regexp.MustCompile(`(?m)^[ \t]*(?:[^ \t\n])`)

// Check outputs a message with a preceding green checkmark
func Check(out io.Writer, msg string, params ...interface{}) {
	fmt.Fprintf(out, color.GreenString("✓")+" "+msg+"\n", params...)
}

// Fail outputs a message with a preceding red X
func Fail(out io.Writer, msg string, params ...interface{}) {
	fmt.Fprintf(out, color.RedString("✗")+" "+msg+"\n", params...)
}

// Dedent removes the common indentation from a block of text by first scanning
// for leading whitespace and then stripping that amount of whitespace from
// each subsequent line
func Dedent(s string) string {
	// Stores the current margin
	var margin string

	// First, remove whitespace from all whitespace only lines
	text := whitespaceOnlyRx.ReplaceAllString(s, "")

	// Next, find the indents
	indents := leadingWhitespaceRx.FindAllString(text, -1)
	for _, indent := range indents {
		// Each indent will have the first non-space char at the end, so let's
		// slice that out
		indent = indent[:len(indent)-1]

		// If we don't have a margin yet, set it to this indent
		if margin == "" {
			margin = indent
		} else if strings.HasPrefix(indent, margin) {
			// This line is more deeply indented than the previous winner, so
			// leave the previous winner alone
			continue
		} else if strings.HasPrefix(margin, indent) {
			// Current line is consistent with and no deeper than the previous
			// winner, so it's the new winner
			margin = indent
		} else {
			// Otherwise, let's find the largest common whitespace between this
			// line and the previous winner
			margin = largestCommonWhitespace(margin, indent)
		}
	}

	if len(margin) > 0 {
		// If we have a margin then let's pull it out of each line
		rx := regexp.MustCompile(fmt.Sprintf(`(?m)^%s`, margin))
		text = rx.ReplaceAllString(text, "")
	}
	return text
}

func largestCommonWhitespace(margin string, indent string) string {
	var sz = int(math.Max(float64(len(margin)), float64(len(indent))))
	for i := 0; i < sz; i++ {
		for _, x := range margin {
			for _, y := range indent {
				if x != y {
					return margin[:i]
				}
			}
		}
	}
	return margin[:len(indent)]
}
