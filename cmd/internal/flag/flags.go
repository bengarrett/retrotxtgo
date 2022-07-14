package flag

import (
	"errors"
	"fmt"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var ErrHide = errors.New("could not hide the flag")

var Build create.Args

type Creates struct {
	Controls []string // character encoding used by the filename
	Encode   string   // use these control codes
	Swap     []string // swap out these characters with UTF8 alternatives
}

func Create() Creates {
	return Creates{
		Controls: []string{"eof", "tab"},
		Encode:   EncodingDefault,
		Swap:     []string{"null", "bar"},
	}
}

type Commands struct {
	Config string
}

var Command Commands

var Info struct {
	Format string
}

type Views struct {
	Controls []string
	Encode   string
	Swap     []string
	To       string
	Width    int
}

var View = views()

func views() Views {
	return Views{
		Controls: []string{"eof", "tab"},
		Encode:   "CP437",
		Swap:     []string{"null", "bar"},
		To:       "",
		Width:    0,
	}
}

// Controls handles the --controls flag.
func Controls(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "controls", "c", []string{},
		`implement these control codes (default "eof,tab")
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

// Encode handles the --encode flag.
func Encode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		fmt.Sprintf("character encoding used by the filename(s) (default \"CP437\")\n%s\n%s%s\n",
			color.Info.Sprint("this flag has no effect for Unicode and EBCDIC samples"),
			"see the list of encode values ",
			str.Example(meta.Bin+" list codepages")))
}

// SwapChars handles the --swap-chars flag.
func SwapChars(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "swap-chars", "x", []string{},
		`swap out these characters with UTF8 alternatives (default "null,bar")
separate multiple values with commas
  null	C null for a space
  bar	Unicode vertical bar | for the IBM broken pipe ¦
  house	IBM house ⌂ for the Greek capital delta Δ
  pipe	Box pipe │ for the Unicode integral extension ⎮
  root	Square root √ for the Unicode check mark ✓
  space	Space for the Unicode open box ␣
  `)
}

// HiddenTo handles the hidden --to flag.
func HiddenTo(p *string, cc *cobra.Command) error {
	const name = "to"
	cc.Flags().StringVar(p, name, "",
		"alternative character encoding to print to stdout\nthis flag is unreliable and not recommended")
	if err := cc.Flags().MarkHidden(name); err != nil {
		return fmt.Errorf("%w, %s: %s", ErrHide, name, err)
	}
	return nil
}

// Width handles the --width flag.
func Width(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", View.Width,
		"maximum document character/column width")
}
