package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

// ErrorFmt is an interface for error messages
type ErrorFmt struct {
	Issue string
	Arg   string
	Msg   error
}

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "retrotxtgo",
		Short: "RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
		Long: `Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`,
	}
)

// color aliases
var (
	alert = func() string {
		return color.Error.Sprint("problem:")
	}
	cc = func(t string) string {
		return color.Comment.Sprint(t)
	}
	ce = func(t string) string {
		return color.Warn.Sprint(t)
	}
	cf = func(t string) string {
		return color.OpFuzzy.Sprint(t)
	}
	ci = func(t string) string {
		return color.OpItalic.Sprint(t)
	}
	cinf = func(t string) string {
		return color.Info.Sprint(t)
	}
	cp = func(t string) string {
		return color.Primary.Sprint(t)
	}
	cs = func(t string) string {
		return color.Success.Sprint(t)
	}
)

// InitDefaults initialises flag and configuration defaults.
func InitDefaults() {
	viper.SetDefault("create.layout", "standard")
	viper.SetDefault("create.title", "RetroTxt | example")
	viper.SetDefault("create.meta.author", "")
	viper.SetDefault("create.meta.color-scheme", "")
	viper.SetDefault("create.meta.description", "An example")
	viper.SetDefault("create.meta.generator", true)
	viper.SetDefault("create.meta.keywords", "")
	viper.SetDefault("create.meta.referrer", "")
	viper.SetDefault("create.meta.theme-color", "")
	viper.SetDefault("create.save-directory", "")
	viper.SetDefault("create.server-port", 8080)
	viper.SetDefault("info.format", "color")
	viper.SetDefault("version.format", "color")
}

// errorPrint returns a coloured error message.
func (e *ErrorFmt) errorPrint() string {
	ia := ci(fmt.Sprintf("%s%s", e.Issue, e.Arg))
	m := cf(fmt.Sprintf(" %v", e.Msg))
	return fmt.Sprintf("%s %s%s", alert(), ia, m)
}

// Check parses the ErrorFmt and will exit with a message if an error is found.
// The ErrorFmt interface comprises of three fields.
// Issue is a summary of the problem.
// Arg is the argument, flag or item that triggered the error.
// Msg is the actual error generated.
func Check(e ErrorFmt) {
	if e.Msg != nil {
		println(e.errorPrint())
		os.Exit(1)
	}
}

// CheckCodePage ...
func CheckCodePage(e ErrorFmt) {
	if e.Msg != nil {
		e.Issue = "unsupported "
		println(e.errorPrint())
		println("         to see a list of supported code pages and aliases run: " + color.Primary.Sprint("retrotxtgo view codepages"))
		os.Exit(1)
	}
}

// CheckErr prints the error and exits.
func CheckErr(err error) {
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

// errorPrint returns a coloured invalid flag message.
func (e *ErrorFmt) errorFlag() string {
	a := fmt.Sprintf("\"--%s %s\"", e.Issue, e.Arg)
	m := cf(fmt.Sprintf(" valid %s values: %v", e.Issue, e.Msg))
	return fmt.Sprintf("\n%s %s %s %s\n", alert(), ci("invalid flag"), a, m)
}

// CheckFlag exits with an invalid command flag value.
func CheckFlag(e ErrorFmt) {
	if e.Msg != nil {
		println(e.errorFlag())
		os.Exit(1)
	}
}

// FileMissingErr exits with a missing FILE error.
func FileMissingErr() {
	i := ci("missing the --name flag")
	m := cf("you need to provide a path to a text file")
	fmt.Printf("\n%s %s %s\n", alert(), i, m)
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err == nil {
		return
	}
	msg := fmt.Sprintf("%s", err)
	m := strings.Split(msg, " ")
	switch {
	case len(msg) > 22 && msg[:22] == "unknown shorthand flag",
		len(msg) > 12 && msg[:12] == "unknown flag":
		Check(ErrorFmt{"invalid flag", m[len(m)-1], fmt.Errorf("is not a flag in use for this command")})
	case len(msg) > 22 && msg[:22] == "flag needs an argument":
		Check(ErrorFmt{"invalid flag", m[len(m)-1], fmt.Errorf("cannot be empty and requires a value")})
	case len(msg) > 17 && msg[:16] == "required flag(s)":
		Check(ErrorFmt{"a required flag missing", m[2], err})
	case len(msg) > 16 && msg[:15] == "unknown command":
		Check(ErrorFmt{"invalid command", m[2], err})
	case msg == "subcommand is required":
		return // ignored errors
	default:
		Check(ErrorFmt{"command", "execute", err})
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file "+cf(fmt.Sprint("(default is $HOME/.retrotxtgo.yaml)")))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		Check(ErrorFmt{"directory", "user home", err})
		// Search config in home directory with name ".retrotxtgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".retrotxtgo")
	}

	// Read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	// todo: toggle a flag to hide this when CREATE XML/JSON/TEXT as it won't be able to be piped
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			Check(ErrorFmt{"config file", viper.ConfigFileUsed(), err})
		}
	}
	//fmt.Println("Using config file:", cf(viper.ConfigFileUsed()))
}
