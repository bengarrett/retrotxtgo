package cmd

import (
	"github.com/bengarrett/retrotxtgo/pkg/str"
	"github.com/gookit/color"
)

// init is always called by the Cobra library to be used for global flags and commands.
//
//nolint:gochecknoinits
func init() {
	const highColor, basicColor = "COLORTERM", "TERM"
	if str.Term(str.GetEnv(highColor), str.GetEnv(basicColor)) == "none" {
		// disable all color output
		color.Enable = false
	}
}
