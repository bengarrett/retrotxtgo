package cmd

import (
	"github.com/bengarrett/retrotxtgo/cmd/hexa"
	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/spf13/cobra"
)

func Dec() *cobra.Command {
	s := "Convert decimal to hexadecimal numbers"
	l := `Rudimentary decimal to hexadecimal conversions.

Convert one or a series of decimal numbers to their hexadecimal, base 16 values.
For example the number "255" is converted to "FF".
The number "9" is returned as "9".

No prefixes or leading characters are added to the hexadecimal numbers.
Negative signs should not be used as they could be interpreted as
a command flag.
`
	return &cobra.Command{
		Use:     "dec",
		Aliases: []string{"d"},
		Short:   s,
		Long:    l,
		GroupID: IDcodepage,
		Example: `  retrotxt dec 0 255 106 161`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			const base = 10
			if flag.Hex.Raw {
				return hexa.Parser(cmd.OutOrStdout(), base, args...)
			}
			return hexa.Writer(cmd.ErrOrStderr(), base, args...)
		},
	}
}

func Hex() *cobra.Command {
	s := "Convert hexadecimal to decimal numbers"
	l := `Rudimentary hexadecimal to decimal conversions.

Convert one or a series of hexadecimal numbers to their decimal, base 10 values.
For example the hexadecimal number "FF" is converted to "255".
The hexadecimal number "55" is returned as "85".
The hexadecimal number "09" is converted to "9".

Common, case insensitive, hexadecimal prefixes are stripped.
  - x prefix (x00)
  - # hash prefix, found in CSS (#000)
  - $ dollar prefix, found in retro microcomputers ($00)
  - U+ unicode prefix, used in unicode (U+0000)
  - 0x zero prefix, found in linux and unix C syntax (0x00)
  - \x escape prefix, found in programming languages (\x00)

Numeric character reference (NCR) syntax is also supported.
  - &#000; decimal NCR syntax
  - &#x00; hexadecimal NCR syntax

Any signs are ignored. If a string is not a hexadecimal number then the 
value is printed as "invalid".
`
	return &cobra.Command{
		Use:     "hex",
		Aliases: []string{"h", "x"},
		Short:   s,
		Long:    l,
		GroupID: IDcodepage,
		Example: `  retrotxt hex 0x00 xff U+006A a1`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			const base = 16
			if flag.Hex.Raw {
				return hexa.Parser(cmd.OutOrStdout(), base, args...)
			}
			return hexa.Writer(cmd.ErrOrStderr(), base, args...)
		},
	}
}

func DecInit() *cobra.Command {
	const s = "raw output only returns the space separated results"
	d := Dec()
	d.Flags().BoolVarP(&flag.Hex.Raw, "raw", "r", false, s)
	return d
}

func HexInit() *cobra.Command {
	const s = "raw output only returns the space separated results"
	h := Hex()
	h.Flags().BoolVarP(&flag.Hex.Raw, "raw", "r", false, s)
	return h
}

func init() {
	Cmd.AddCommand(DecInit(), HexInit())
}
