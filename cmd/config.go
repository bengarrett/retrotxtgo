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

	"github.com/spf13/cobra"
)

var (
	fileOverwrite bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   cp("Configure RetroTxt defaults"),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("config called")
	},
}
var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: cp("create a new config file"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create config file")
	},
}
var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: cp("remove the default config file"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("delete config file")
	},
}
var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: cp("edit the default config file"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("edit config file")
	},
}
var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: cp("view settings configured by the config"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("into config file")
	},
}

// todo detect Windows to
// + replace terminal shell with command prompt
// + automatically use powershell
var configShellCmd = &cobra.Command{
	Use:   "shell",
	Short: cp("include retrotxt autocompletion to the terminal shell"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shell config file")
	},
}
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: cp("change a configuration"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set config file")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCreateCmd)
	configCreateCmd.Flags().BoolVarP(&fileOverwrite, "overwrite", "y", false, "overwrite any existing config file")
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInfoCmd)
	configCmd.AddCommand(configShellCmd)
	configCmd.AddCommand(configSetCmd)

	// todo
	// include text on how to use
	// rootCmd.GenBashCompletionFile("hi.sh")
	// rootCmd.GenBashCompletionFile("hi.sh")
	// rootCmd.GenPowershellCompletionFile("hi.sh")
	// rootCmd.GenZshCompletionFile("hi.sh")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

/*
config create
--overwrite=true/false
config edit
> look for EDITOR env var
config show
> output using viper.Sub("create.meta")
*/
/*

possible options:

font choice (family)
font size
font format
> base64
> woff2

codepage?

input (override for internal use)
> ascii
> ansi
> etc

output
> stout
> unique filename
> same filename `index.html`

css
> embed
> embed+minified
> link
> link+assets
> none

template (or css?)
> basic (raw)
> standard

title
quiet

*/
