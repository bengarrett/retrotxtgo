package flag

import "github.com/spf13/cobra"

// Runes handles the --swap-chars flag.
func Runes(p *[]string, cc *cobra.Command) {
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
