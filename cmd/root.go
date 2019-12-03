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
	Ver string = "0.0.3"
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
	Layout  PageData
	cfgFile string

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
	println("InitDefaults")
	viper.SetDefault("version.format", "color")
	viper.SetDefault("info.format", "color")
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
}

// ErrorPrint returns a coloured error message.
func (e *ErrorFmt) ErrorPrint() string {
	ia := color.OpItalic.Sprintf("%s %s", e.Issue, e.Arg)
	m := color.OpFuzzy.Sprintf(" %v", e.Msg)
	return color.Sprintf("\n%s %s%s", alert(), ia, m)
}

// GoErr exits with a coloured error message.
func (e *ErrorFmt) GoErr() {
	println(e.ErrorPrint())
	os.Exit(1)
}

// FlagErr exits with an invalid command flag value.
func (e *ErrorFmt) FlagErr() {
	a := fmt.Sprintf("\"--%s %s\"", e.Issue, e.Arg)
	m := color.OpFuzzy.Sprintf(" valid %s values: %v", e.Issue, e.Msg)
	color.Printf("\n%s %s %s%s\n", alert(), ci("invalid flag"), a, m)
	os.Exit(1)
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Search config in home directory with name ".retrotxtgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".retrotxtgo")
	}

	// Read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	// todo: toggle a flag to hide this when CREATE XML/JSON/TEXT as it won't be able to be piped
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", cf(viper.ConfigFileUsed()))
	}
}
