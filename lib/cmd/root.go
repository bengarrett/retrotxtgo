// Package cmd handles the terminal interface, user flags and arguments.
// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	Use:   meta.Bin,
	Short: fmt.Sprintf("%s is the tool that turns ANSI, ASCII, NFO text into browser ready HTML", meta.Name),
	Long: fmt.Sprintf(`Turn many pieces of ANSI art, ASCII and NFO texts into HTML5 using %s.
It is the platform agnostic tool that takes nostalgic text files and stylises
them into a more modern, useful format to view or copy in a web browser.`, meta.Name),
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing other than print the help.
		// This func must remain otherwise root command flags are ignored by Cobra.
		printUsage(cmd)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SilenceErrors = silence
	rootCmd.Version = meta.Print()
	rootCmd.SetVersionTemplate(version())
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

func version() string {
	const tabWidth, copyright, years = 8, "\u00A9", "2020-21"
	exe, err := self()
	if err != nil {
		exe = err.Error()
	}
	newVer, v := chkRelease()
	appDate := ""
	if meta.App.Date != meta.Placeholder {
		appDate = fmt.Sprintf(" (%s)", meta.App.Date)
	}
	var b bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&b, 0, tabWidth, 0, '\t', 0)
	fmt.Fprintf(w, "%s %s\n", meta.Name, meta.Print())
	fmt.Fprintf(w, "%s %s Ben Garrett\n", copyright, years)
	fmt.Fprintln(w, color.Primary.Sprint(meta.URL))
	fmt.Fprintf(w, "\n%s\t%s %s%s\n", color.Secondary.Sprint("build:"), runtime.Compiler, meta.App.BuiltBy, appDate)
	fmt.Fprintf(w, "%s\t%s/%s\n", color.Secondary.Sprint("platform:"), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("go:"), strings.Replace(runtime.Version(), "go", "v", 1))
	fmt.Fprintf(w, "%s\t%s\n", color.Secondary.Sprint("path:"), exe)
	if newVer {
		fmt.Fprintf(w, "\n%s\n", newRelease(meta.App.Version, v))
	}
	w.Flush()
	return b.String()
}

// Self returns the path to this dupers executable file.
func self() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("self error: %w", err)
	}
	return exe, nil
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&rootFlag.config, "config", "",
		"optional config file location")
	rootCmd.LocalNonPersistentFlags().BoolP("version", "v", false, "")
}

// initConfig reads in the config file and ENV variables if set.
// This init might be run twice due to the Cobra initializer registers.
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
	scanner, s := bufio.NewScanner(&b), ""
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

// flagEncode handles the encode flag.
func flagEncode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		fmt.Sprintf(`character encoding used by the filename(s)
this flag is silently ignored if Unicode text is detected
otherwise the default is used (default "CP437")
see the list of encode values %s%s`,
			str.Example(meta.Bin+" list codepages"), "\n"))
}

// flagRunes handles the swap-chars flag.
func flagRunes(p *[]int, cc *cobra.Command) {
	cc.Flags().IntSliceVarP(p, "swap-chars", "x", []int{},
		`swap out these characters with UTF8 alternatives (default "0,124")
separate multiple values with commas
  0    C null for a space
  124  Unicode vertical bar | for the IBM broken pipe ¦
  127  IBM house ⌂ for the Greek capital delta Δ
  178  Box pipe │ for the Unicode integral extension ⎮
  251  Square root √ for the Unicode check mark ✓
`)
}

// flagTo handles the to flag.
func flagTo(p *string, cc *cobra.Command) {
	cc.Flags().StringVar(p, "to", "",
		fmt.Sprintf(`alternative character encoding to print to stdout
modern terminals and %s use UTF8 encoding
this flag is unreliable and not recommended
see the list of usable values %s%s`,
			meta.Name, str.Example(meta.Bin+" list codepages"), "\n"))
}

// flagWidth handles the width flag.
func flagWidth(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", viewFlag.width,
		"maximum document character/column width")
}
