// Package flag provides the command flags handlers.
package flag

import (
	"errors"
	"fmt"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/spf13/cobra"
)

var ErrHide = errors.New("could not hide the flag")

// Command flags for the main cmd.
type Command struct {
	Tester bool // internal automated tester
}

// Cmd returns the flags for the main cmd.
var Cmd Command

// Info handles the info "format" flag.
var Info struct {
	Format string // output format
}

// Views handles the view command flags.
type Views struct {
	Input    string   // input character encoding used by the files
	Controls []string // control codes to implement
	Swap     []string // swap out these characters with Unicode control pictures
	Width    int      // maximum document character/column width
	Original bool     // output the sample's original character encoding to stdout
}

// View returns the Views struct with default values.
func View() Views {
	return Views{
		Input:    "CP437",
		Controls: []string{"eof", "tab"},
		Swap:     []string{"null", "bar"},
		Width:    0,
		Original: false,
	}
}

// Controls handles the "controls" flag.
func Controls(p *[]string, cc *cobra.Command) {
	//nolint:dupword
	cc.Flags().StringSliceVarP(p, "controls", "c", []string{},
		`implement these control codes (default "eof,tab")`+
			`
  separate multiple controls with commas
  eof    end of file mark
  tab    horizontal tab
  bell   bell or terminal alert
  cr     carriage return
  lf     line feed
  bs backspace, del delete character, esc escape character
  ff formfeed, vt vertical tab
`)
}

// Encode handles the "input" flag.
func Encode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "input", "i", "",
		fmt.Sprintf("character encoding used by the filename(s) (default \"CP437\")\n%s%s\n%s\n",
			"see the list of encode values ",
			term.Example(meta.Bin+" list codepages"),
			"this flag has no effect for the inbuilt samples"))
}

// SwapChars handles the "swap-chars" flag.
func SwapChars(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "swap-chars", "x", []string{},
		`swap out these characters with common alternatives (default "null,bar")
  separate multiple values with commas
  null	C null for a space
  bar	Unicode vertical bar | for the IBM broken pipe ¦
  house	IBM house ⌂ for the Greek capital delta Δ
  pipe	Box pipe │ for the Unicode integral extension ⎮
  root	Square root √ for the Unicode check mark ✓
  space	Space for the Unicode open box ␣
  `)
}

// FlagOG handles the hidden original flag.
func FlagOG(p *bool, cc *cobra.Command) error {
	const name = "original"
	cc.Flags().BoolVar(p, name, false,
		"use the original document character encoding to print to terminal"+
			"\nby default, "+meta.Bin+" view always prints texts as UTF-8"+
			"\nthis flag only works with the inbuilt samples\n")
	return nil
}

// Width handles the "width" flag.
func Width(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", View().Width,
		`maximum document character/column width
any horizontal tab characters are replaced with three spaces
any newline characters are replaced with a space
`)
}
