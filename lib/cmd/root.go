package cmd

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFlag string

// rootCmd represents the base command when called without any subcommands
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
	rootCmd.SilenceErrors = false // set to false to debug
	if err := rootCmd.Execute(); err != nil {
		//rootCmd.Usage()
		rootErr := logs.CmdErr{Args: os.Args[1:], Err: err}
		fmt.Println(rootErr.Error().String())
	}
}

func init() {
	// OnInitialize will not run if no command is provided.
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&configFlag,
		"config", "", "config file "+
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
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	config.SetConfig(configFlag)
}
