package cmd

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type rootFlags struct {
	config string
}

var rootFlag = rootFlags{}

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
	rootCmd.SilenceErrors = true // set to false to debug
	if err := rootCmd.Execute(); err != nil {
		if len(os.Args) < 2 {
			logs.Check("rootcmd.usage", rootCmd.Usage())
		}
		rootErr := logs.CmdErr{Args: os.Args[1:], Err: err}
		fmt.Println(rootErr.Error().String())
	}
}

func init() {
	// OnInitialize will not run if there is no provided command.
	cobra.OnInitialize(initConfig)
	// TODO: get viper to flag file autocomplete
	rootCmd.PersistentFlags().StringVar(&rootFlag.config, "config", "",
		"optional config file location")
}

func defaultConfig() string {
	named := viper.GetViper().ConfigFileUsed()
	if named == "" {
		named = config.Path()
	}
	return named
}

// initConfig reads in the config file and ENV variables if set.
// this does not run when rootCmd is in use.
func initConfig() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	if err := config.SetConfig(rootFlag.config); err != nil {
		logs.Check(fmt.Sprintf("config file %q", viper.ConfigFileUsed()), err)
		os.Exit(1)
	}
}

// checkUsage will print the help and exit when no arguments are supplied.
func checkUse(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cmd.Help()
		logs.Check("cmd.help", err)
		os.Exit(0)
	}
}

type internalPack struct {
	// choices: d convert.Dump, t convert.Text (default when blank)
	convert string
	// default character encoding for the packed data
	encoding string
	// package name used in internal/pack/blob.go
	name string
}

var internalPacks = map[string]internalPack{
	"437.cr":        {"d", "", "text/cp437-cr.txt"},
	"437.crlf":      {"d", "", "text/cp437-crlf.txt"},
	"437.lf":        {"d", "", "text/cp437-lf.txt"},
	"865":           {"", "ibm865", "text/cp865.txt"},
	"1252":          {"", "cp1252", "text/cp1252.txt"},
	"ascii":         {"", "cp437", "text/retrotxt.asc"},
	"ansi":          {"", "", "text/retrotxt.ans"},
	"ansi.aix":      {"", "", "text/ansi-aixterm.ans"},
	"ansi.blank":    {"", "", "text/ansi-blank"},
	"ansi.cp":       {"", "", "text/ansi-cp.ans"},
	"ansi.cpf":      {"", "", "text/ansi-cpf.ans"},
	"ansi.hvp":      {"", "", "text/ansi-hvp.ans"},
	"ansi.proof":    {"", "", "text/ansi-proof.ans"},
	"ansi.rgb":      {"", "cp437", "text/ansi-rgb.ans"},
	"ansi.setmodes": {"", "", "text/ansi-setmodes.ans"},
	"iso-1":         {"", "1", "text/iso-8859-1.txt"},
	"iso-15":        {"", "15", "text/iso-8859-15.txt"},
	"sauce":         {"", "", "text/sauce.txt"},
	"shiftjis":      {"", "shift-jis", "text/shiftjis.txt"},
	"us-ascii":      {"", "cp1252", "text/us-ascii.txt"},
	"utf8":          {"", "", "text/utf-8.txt"},
	"utf8.bom":      {"", "", "text/utf-8-bom.txt"},
	"utf16.be":      {"", "utf-16be", "text/utf-16-be.txt"},
	"utf16.le":      {"", "utf-16le", "text/utf-16-le.txt"},
}
