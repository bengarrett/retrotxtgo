// Package cmd is the terminal, command line interface.
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/pack"
	"retrotxt.com/retrotxt/lib/str"
)

type rootFlags struct {
	config string
}

var (
	// ErrIntpr no interpreter.
	ErrIntpr = errors.New("the interpreter is not supported")
	// ErrPackGet invalid pack name.
	ErrPackGet = errors.New("pack.get name is invalid")
	// ErrPackValue unknown pack value.
	ErrPackValue = errors.New("unknown package convert value")
	// ErrTempClose close temp file.
	ErrTempClose = errors.New("could not close temporary file")
	// ErrTempOpen open temp file.
	ErrTempOpen = errors.New("could not create temporary file")
	// ErrTempWrite write temp file.
	ErrTempWrite = errors.New("could not write to temporary file")
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
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceErrors = true // set to false to debug
	if err := rootCmd.Execute(); err != nil {
		const minArgs = 2
		if len(os.Args) < minArgs {
			if err1 := rootCmd.Usage(); err1 != nil {
				logs.Fatal("rootcmd", "usage", err1)
			}
		}
		logs.Execute(err, os.Args[1:]...)
	}
}

func init() {
	// OnInitialize will not run if there is no provided command.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&rootFlag.config, "config", "",
		"optional config file location")
}

// initConfig reads in the config file and ENV variables if set.
// this does not run when rootCmd is in use.
func initConfig() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	if err := config.SetConfig(rootFlag.config); err != nil {
		logs.Fatal("config file", viper.ConfigFileUsed(), err)
	}
}

// checkUsage will print the help and exit when no arguments are supplied.
func checkUse(cmd *cobra.Command, args ...string) {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			logs.Fatal("root", "cmd.help", err)
		}
		os.Exit(0)
	}
}

func flagControls(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "controls", "c", []string{},
		`use these control codes
  tab    horizontal tab
  bell   bell or terminal alert
  cr     carriage return
  lf     line feed
  bs backspace, del delete character, esc escape character
  ff formfeed, vt vertical tab
(default tab)
separate multiple controls with commas
`+str.Example("--controls=tab,bell")+"\n")
}

func flagEncode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		`character encoding used by the filenames
when ignored, UTF8 encoding is detected
if that fails the default is used (default CP437)
see the list of encode values `+str.Example("retrotxt list codepages")+"\n")
}

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

func flagTo(p *string, cc *cobra.Command) {
	cc.Flags().StringVar(p, "to", "", `alternative character encoding to print to stdout
modern terminals and RetroTxt use UTF8 encoding
this alternative option is unreliable and not recommended
see the list of usable values `+str.Example("retrotxt list codepages")+"\n")
}

func flagWidth(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", viewFlag.width, "maximum document character/column width")
}

func examples() *bytes.Buffer {
	m := pack.Map()
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	var flags uint = 0 // tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', flags)
	const title = " Packaged example text and ANSI files to test and play with RetroTxt "
	fmt.Fprintln(w, str.Cp(title))
	fmt.Fprintln(w, strings.Repeat("-", len(title)))
	for _, k := range keys {
		fmt.Fprintf(w, "%s\t%s\t\n", k, m[k].Description)
	}
	fmt.Fprintln(w, "\nAny of these packaged examples will work with both the",
		str.Example("create")+",", str.Example("info"), "and", str.Example("view"), "commands")
	fmt.Fprintln(w, "\n"+str.Example(" retrotxt view 1252"),
		"will print the Windows-1252 English test to the terminal")
	fmt.Fprintln(w, str.Example(" retrotxt view 1252 > file.txt"),
		"will convert and save the Windows-1252 English test to UTF-8 encoding")
	fmt.Fprintln(w, str.Example(" retrotxt view --to=cp1252 1252 > file.txt"),
		"will save the Windows-1252 English test with its original encoding")
	fmt.Fprintln(w, str.Example(" retrotxt view --to=cp1252 1252 | retrotxt info"), "displays statistics and information from a piped source")
	fmt.Fprintln(w, str.Example(" retrotxt info 1252"), "displays statistics and information from the Windows-1252 English test")
	fmt.Fprintln(w, str.Example(" retrotxt info sauce"), "displays statistics, information and SAUCE metadata from the SAUCE test")
	fmt.Fprintln(w, str.Example(" retrotxt create 1252"), "creates a HTML document from the Windows-1252 English test")
	fmt.Fprintln(w, str.Example(" retrotxt create 1252 -p0"), "serves the Windows-1252 English test over a local web server")
	fmt.Fprintln(w, "\nMultiple examples are supported")
	fmt.Fprintln(w, str.Example(" retrotxt view ansi ascii ansi.rgb"))
	if err := w.Flush(); err != nil {
		logs.Fatal("flush of tab writer failed", "", err)
	}
	return &buf
}
