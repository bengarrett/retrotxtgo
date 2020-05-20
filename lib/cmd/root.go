package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// ErrorFmt is an interface for error messages
type ErrorFmt struct {
	Issue string
	Arg   string
	Msg   error
}

var configName string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retrotxt",
	Short: "RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
	Long: `Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`,
}

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

// CheckErr prints the error and exits.
func CheckErr(err error) {
	if err != nil {
		fmt.Printf("error: %s", err)
		os.Exit(1)
	}
}

// FileMissingErr exits with a missing FILE error.
func FileMissingErr() {
	i := logs.Ci("missing the --name flag")
	m := logs.Cf("you need to provide a path to a text file")
	fmt.Printf("\n%s %s %s\n", logs.Alert(), i, m)
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceErrors = true
	if err := rootCmd.Execute(); err != nil {
		rootCmd.Usage()
		rootErr := logs.CmdErr{Args: os.Args[1:], Err: err}
		fmt.Println(rootErr.Error().String())
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configName, "config", "", "config file "+
		logs.Cf(fmt.Sprintf("(default is %s)", defaultConfig())))
}

func defaultConfig() string {
	named := viper.GetViper().ConfigFileUsed()
	if named == "" {
		named = config.Filepath()
	}
	return named
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("yaml")
	if configName != "" {
		viper.SetConfigFile(configName)
	} else {
		viper.SetConfigFile(config.Filepath())
	}
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			config.Create(false)
		} else {
			logs.ChkErr(logs.Err{Issue: "config file failed", Arg: viper.ConfigFileUsed(), Msg: err})
		}
	}
}
