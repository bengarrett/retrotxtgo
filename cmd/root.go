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
	Ver string = "0.0.1"
	// Www is the website domain name
	Www string = "retrotxt.com"
)

//PageData contains template data used by standard.html
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

// Layout data
var Layout PageData

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retrotxtgo",
	Short: "RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
	Long: `Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

//ErrorFmt is
type ErrorFmt struct {
	Issue string
	Arg   string
	Msg   error
}

//ErrorPrint is
func (e *ErrorFmt) ErrorPrint() string {
	r := color.Error.Sprint("ERROR:")
	u := color.OpUnderscore.Sprintf("%s %s", e.Issue, e.Arg)
	f := color.OpFuzzy.Sprintf(" %v", e.Msg)
	return color.Sprintf("\n%s %s%s", r, u, f)
}

//UsageErr is
func (e *ErrorFmt) UsageErr(cmd *cobra.Command) {
	println("Usage:\n" + "  retrotxtgo " + cmd.Use)
	println("\nExamples:\n" + cmd.Example)
	println(e.ErrorPrint())
	os.Exit(1)
}

//GoErr is
func (e *ErrorFmt) GoErr() {
	println(e.ErrorPrint())
	os.Exit(1)
}

//LayoutDefault xxx
func LayoutDefault() PageData {
	l := Layout
	l.BuildVersion = Ver
	l.BuildDate = time.Now()
	l.CacheRefresh = fmt.Sprintf("?v=%v", Ver)
	l.MetaDesc = "A textfile example"
	l.MetaGenerator = true
	l.PageTitle = "RetroTxt | example"
	l.PreText = "Hello world."
	return l
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.retrotxtgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
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

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
