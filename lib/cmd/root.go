// Package cmd handles the terminal interface, user flags and arguments.
// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

type rootFlags struct {
	config string
}

const (
	eof         = "eof"
	tab         = "tab"
	null        = 0
	verticalBar = 124

	// silence can be set to false to debug cmd/flag feedback from Viper.
	silence = false
)

var rootFlag = rootFlags{}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "retrotxt",
	Short: "RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
	Long: `Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing other than print the help.
		// This func must remain otherwise root command flags are ignored by Cobra.
		printUsage(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceErrors = silence
	if err := rootCmd.Execute(); err != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if err1 := rootCmd.Usage(); err1 != nil {
				logs.ProblemMarkFatal("rootCmd", ErrUsage, err1)
			}
		}
		logs.Execute(err, os.Args[1:]...)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&rootFlag.config, "config", "",
		"optional config file location")
}

// initConfig reads in the config file and ENV variables if set.
// This init can be run twice due to the Cobra initializer registers.
func initConfig() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	if err := config.SetConfig(rootFlag.config); err != nil {
		logs.ProblemMarkFatal(viper.ConfigFileUsed(), logs.ErrCfgFile, err)
	}
}

// printUsage will print the help and exit when no arguments are supplied.
func printUsage(cmd *cobra.Command, args ...string) bool {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			logs.ProblemMarkFatal("checkUse", ErrHelp, err)
		}
		os.Exit(0)
	}
	return false
}

// endOfFile determines if an EOF marker should be obeyed.
func endOfFile(flags convert.Flags) bool {
	for _, c := range flags.Controls {
		if c == eof {
			return true
		}
	}
	return false
}

// exampleCmd returns help usage examples.
func exampleCmd(tmpl string) string {
	if tmpl == "" {
		return ""
	}
	var b bytes.Buffer
	// change example operating system path separator
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(&b, string(os.PathSeparator))
	if err != nil {
		log.Fatal(err)
	}
	// color the example text except text following
	// the last hash #, which is treated as a comment
	const cmmt, sentence = "#", 2
	var s string
	scanner := bufio.NewScanner(&b)
	for scanner.Scan() {
		ss := strings.Split(scanner.Text(), cmmt)
		l := len(ss)
		if l < sentence {
			s += str.Cinf(scanner.Text()) + "\n  "
			continue
		}
		// do not the last hash as a comment
		ex := strings.Join(ss[:l-1], cmmt)
		s += str.Cinf(ex)
		s += fmt.Sprintf("%s%s\n  ", cmmt, ss[l-1])
	}
	return s
}

// flagControls handles the controls flag.
func flagControls(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "controls", "c", []string{},
		`use these control codes
  eof    end of file mark
  tab    horizontal tab
  bell   bell or terminal alert
  cr     carriage return
  lf     line feed
  bs backspace, del delete character, esc escape character
  ff formfeed, vt vertical tab
(default eof,tab)
separate multiple controls with commas
`+str.Example("--controls=eof,tab,bell")+"\n")
}

// flagEncode handles the encode flag.
func flagEncode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		`character encoding used by the filename(s)
this flag is silently ignored if Unicode text is detected
otherwise the default is used (default CP437)
see the list of encode values `+str.Example("retrotxt list codepages")+"\n")
}

// flagRunes handles the swap-chars flag.
func flagRunes(p *[]int, cc *cobra.Command) {
	cc.Flags().IntSliceVarP(p, "swap-chars", "x", []int{},
		`swap out these characters with UTF8 alternatives
  0    C null for a space
  124  Unicode vertical bar | for the IBM broken pipe ¦
  127  IBM house ⌂ for the Greek capital delta Δ
  178  Box pipe │ for the Unicode integral extension ⎮
  251  Square root √ for the Unicode check mark ✓
(default 0,124)
separate multiple values with commas
`+str.Example("--swap-chars=0,124,127")+"\n")
}

// flagTo handles the to flag.
func flagTo(p *string, cc *cobra.Command) {
	cc.Flags().StringVar(p, "to", "", `alternative character encoding to print to stdout
modern terminals and RetroTxt use UTF8 encoding
this flag is unreliable and not recommended
see the list of usable values `+str.Example("retrotxt list codepages")+"\n")
}

// flagWidth handles the width flag.
func flagWidth(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", viewFlag.width, "maximum document character/column width")
}
