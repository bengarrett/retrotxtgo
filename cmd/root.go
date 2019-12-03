/*
Copyright © 2019 Ben Garrett <code.by.ben@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

const (
	// Ver is the application version
	Ver string = "0.0.4"
	// Www is the application domain name
	Www string = "retrotxt.com"
	// FileDate is a non-standard date format for file modifications
	FileDate string = "2 Jan 15:04 2006"
)

// PageData holds template data used by the HTML layouts.
type PageData struct {
	BuildVersion    string
	BuildDate       time.Time
	CacheRefresh    string
	MetaAuthor      string
	MetaColorScheme string
	MetaDesc        string
	MetaGenerator   bool
	MetaKeywords    string
	MetaReferrer    string
	MetaThemeColor  string
	PageTitle       string
	PreText         string
}

// ErrorFmt is an interface for error messages
type ErrorFmt struct {
	Issue string
	Arg   string
	Msg   error
}

var (
	// Layout template data
	Layout      PageData
	cfgFile     string
	suppressCfg bool = false

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "retrotxtgo",
		Short: cp("RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML"),
		Long: color.Info.Sprint(`Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`),
	}
)

// color aliases
var (
	alert = func() string {
		return color.Error.Sprint("ERROR:")
	}
	cc = func(t string) string {
		return color.Comment.Sprint(t)
	}
	cf = func(t string) string {
		return color.OpFuzzy.Sprint(t)
	}
	ci = func(t string) string {
		return color.OpItalic.Sprint(t)
	}
	cp = func(t string) string {
		return color.Primary.Sprint(t)
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
	ia := color.OpItalic.Sprintf("%s %s", e.Issue, e.Arg)
	m := color.OpFuzzy.Sprintf(" %v", e.Msg)
	return color.Sprintf("%s %s%s", alert(), ia, m)
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

// errorPrint returns a coloured invalid flag message.
func (e *ErrorFmt) errorFlag() string {
	a := fmt.Sprintf("\"--%s %s\"", e.Issue, e.Arg)
	m := color.OpFuzzy.Sprintf(" valid %s values: %v", e.Issue, e.Msg)
	return color.Sprintf("\n%s %s %s%s\n", alert(), ci("invalid flag"), a, m)
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
	i := ci("missing the FILE argument")
	m := cf("you need to provide a path to a text file")
	color.Printf("\n%s %s %s\n", alert(), i, m)
	os.Exit(1)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	Check(ErrorFmt{"execute", "cobra", rootCmd.Execute()})
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file "+cf("(default is $HOME/.retrotxtgo.yaml)"))
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
			fmt.Printf("%s\n", err)
		}
	} else if suppressCfg == false {
		fmt.Println("Using config file:", cf(viper.ConfigFileUsed()))
	}
}
