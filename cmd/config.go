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
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
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
		cfg := viper.ConfigFileUsed()
		if cfg != "" {
			println("A config file already exists at ", cfg)
			os.Exit(1)
		}
		c := viper.AllSettings()
		bs, _ := yaml.Marshal(c)
		d, _ := os.UserHomeDir()
		err := ioutil.WriteFile(d+"/.retrotxtgo.yaml", bs, 0660)
		if err != nil {
			fmt.Printf("%s", err)
		}
	},
}
var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: cp("remove the default config file"),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := viper.ConfigFileUsed()
		if cfg == "" {
			configMissing(cmd.CommandPath(), "delete")
		}
		if _, err := os.Stat(cfg); os.IsNotExist(err) {
			configMissing(cmd.CommandPath(), "delete")
		}
		switch prompt("Confirm the file deletion", false) {
		case true:
			if err := os.Remove(cfg); err != nil {
				e := ErrorFmt{"Could not remove", cfg, err}
				e.GoErr()
			}
			fmt.Println("Deletion is done")
		}
	},
}

func prompt(query string, yesDefault bool) bool {
	var input string
	y := "Y"
	n := "n"
	if yesDefault == false {
		y = "y"
		n = "N"
	}
	fmt.Printf("%s? [%s/%s] ", query, color.Success.Sprint(y), color.Danger.Sprint(n))
	fmt.Scanln(&input)
	switch input {
	case "":
		if yesDefault == true {
			return true
		}
	case "yes", "y":
		return true
	}
	return false
}

func configMissing(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + "create"
	fmt.Printf("No config file is in use.\nTo create one run: %s\n", cp(cmd))
	os.Exit(1)
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: cp("edit the default config file"),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := viper.ConfigFileUsed()
		if cfg == "" {
			configMissing(cmd.CommandPath(), "edit")
		}
		var edit string
		if err := viper.BindEnv("editor", "EDITOR"); err != nil {
			editors := []string{"nano", "vim", "emacs"}
			if runtime.GOOS == "windows" {
				editors = append(editors, "notepad++.exe", "notepad.exe")
			}
			for _, editor := range editors {
				if _, err := exec.LookPath(editor); err == nil {
					edit = editor
					break
				}
			}
			if edit != "" {
				fmt.Printf("There is no %s environment variable set so using: %s\n", ci("EDITOR"), cp(edit))
			}
		} else {
			edit = viper.GetString("editor")
			if _, err := exec.LookPath(edit); err != nil {
				e := ErrorFmt{edit, "command not found", exec.ErrNotFound}
				e.GoErr()
			}
		}
		// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
		exe := exec.Command(edit, cfg)
		exe.Stdin = os.Stdin
		exe.Stdout = os.Stdout
		err := exe.Run()
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	},
}
var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: cp("view settings configured by the config"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cp("These are the default configurations used by the commands of RetroTxt when no flags are given.\n"))
		c := viper.AllSettings()
		bs, _ := yaml.Marshal(c)
		fmt.Printf("%s\n", bs)
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
