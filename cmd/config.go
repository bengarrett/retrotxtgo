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

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	fileOverwrite bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: cp("Configure RetroTxt defaults"),
	// no Run: means the help will be display
}

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: cp("create a new config file"),
	Run: func(cmd *cobra.Command, args []string) {
		suppressCfg = true // todo: not working
		if cfg := viper.ConfigFileUsed(); cfg != "" {
			configExists(cmd.CommandPath(), "create")
		}
		bs, err := yaml.Marshal(viper.AllSettings())
		Check(ErrorFmt{"could not create", "settings", err})
		d, err := os.UserHomeDir()
		Check(ErrorFmt{"could not use", "user home directory", err})
		err = ioutil.WriteFile(d+"/.retrotxtgo.yaml", bs, 0660)
		Check(ErrorFmt{"could not write", "settings", err})
		fmt.Println("Created a new config file at:", cf(d+"/.retrotxtgo.yaml"))
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
			Check(ErrorFmt{"Could not remove", cfg, os.Remove(cfg)})
			fmt.Println("The file is deleted")
		}
	},
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
				Check(ErrorFmt{edit, "command not found", exec.ErrNotFound})
			}
		}
		// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
		exe := exec.Command(edit, cfg)
		exe.Stdin = os.Stdin
		exe.Stdout = os.Stdout
		if err := exe.Run(); err != nil {
			fmt.Printf("%s\n", err)
		}
	},
}

var configInfoCmd = &cobra.Command{
	Use:   "info",
	Short: cp("view settings configured by the config"),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cp("These are the default configurations used by the commands of RetroTxt when no flags are given.\n"))
		sets, err := yaml.Marshal(viper.AllSettings())
		Check(ErrorFmt{"read configuration", "yaml", err})
		quick.Highlight(os.Stdout, string(sets), "yaml", "terminal256", "monokai")
		println()
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

func prompt(query string, yesDefault bool) bool {
	var input string
	y := "Y"
	n := "n"
	if yesDefault == false {
		y = "y"
		n = "N"
	}
	fmt.Printf("%s? [%s/%s] ", query, y, n)
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

func configExists(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix)
	fmt.Printf("A config file already is in use at: %s\n", cf(viper.ConfigFileUsed()))
	fmt.Printf("To edit it: %s\n", cp(cmd+"edit"))
	fmt.Printf("To delete:  %s\n", cp(cmd+"delete"))
	os.Exit(1)
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
