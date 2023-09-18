package cmd

import (
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/gookit/color"
)

// init is always called by the Cobra library to be used for global flags and commands.
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if term.Term(term.GetEnv(highColor), term.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
}
