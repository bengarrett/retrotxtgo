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
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const shells string = "bash, powershell, zsh"

var (
	configShell   string
	configSetName string
	fileOverwrite bool
	infoStyles    string
	shellPreview  bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure RetroTxt defaults",
	// no Run: means the help will be display
}

var configCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new config file",
	Run: func(cmd *cobra.Command, args []string) {
		suppressCfg = true // todo: not working
		if cfg := viper.ConfigFileUsed(); cfg != "" && fileOverwrite != true {
			configExists(cmd.CommandPath(), "create")
		}
		writeConfig(false)
	},
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove the default config file",
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
	Short: "Edit the default config file",
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
	Short: "View settings configured by the config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cp("These are the default configurations used by the commands of RetroTxt when no flags are given.\n"))
		sets, err := yaml.Marshal(viper.AllSettings())
		Check(ErrorFmt{"read configuration", "yaml", err})
		quick.Highlight(os.Stdout, string(sets), "yaml", "terminal256", infoStyles)
		println()
	},
}

var configShellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Include retrotxt autocompletion to the terminal shell",
	Run: func(cmd *cobra.Command, args []string) {
		var buf bytes.Buffer
		var lexer string
		switch configShell {
		case "bash":
			lexer = "bash"
			cmd.GenBashCompletion(&buf)
		case "powershell", "posh":
			lexer = "powershell"
			cmd.GenPowerShellCompletion(&buf)
		case "zsh":
			lexer = "bash"
			cmd.GenZshCompletion(&buf)
		default:
		}
		quick.Highlight(os.Stdout, buf.String(), lexer, "terminal256", "monokai")
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Change a configuration",
	//todo add long with information on how to view settings
	Example: `--name create.meta.description # to change the meta description setting
--name version.format          # to set the version command output format`,
	Run: func(cmd *cobra.Command, args []string) {
		var name = configSetName
		keys := viper.AllKeys()
		sort.Strings(keys)
		// sort.SearchStrings() - The slice must be sorted in ascending order.
		var i = sort.SearchStrings(keys, name)
		if i == len(keys) || keys[i] != name {
			err := fmt.Errorf("retrotxt config info")
			CheckFlag(ErrorFmt{"name", name, err})
		}
		s := viper.GetString(name)
		switch s {
		case "":
			fmt.Printf("\n%s is currently disabled\n", cp(name))
		default:
			fmt.Printf("\n%s is currently set to %q\n", cp(name), s)
		}
		switch {
		case name == "create.layout":
			fmt.Printf("Set a new value, choices: %s\n", ci(createLayouts()))
			promptSlice(createLayouts())
		case name == "info.format":
			fmt.Printf("Set a new value, choice: %s\n", ci(infoFormats))
			promptSlice(infoFormats)
		case name == "version.format":
			fmt.Printf("Set a new value, choices: %s\n", ci(versionFormats))
			promptSlice(versionFormats)
		case name == "create.server-port":
			fmt.Printf("Set a new HTTP port, choices: %v-%v (recommended: 8080)\n", portMin, portMax)
			promptPort()
		case name == "create.meta.generator":
			fmt.Printf("<meta name=\"generator\" content=\"RetroTxt v%s\">\nEnable this element? [y/n]\n", Ver)
			promptBool()
		case s == "":
			promptMeta(s)
			fmt.Printf("\nSet a new value or leave blank to keep it disabled: \n")
			promptString(s)
		default:
			promptMeta(s)
			fmt.Printf("\nSet a new value, leave blank to keep as-is or use a dash [-] to disable: \n")
			promptString(s)
		}
	},
}

func promptMeta(val string) {
	s := strings.Split(configSetName, ".")
	switch {
	case len(s) != 3, s[0] != "create", s[1] != "meta":
		return
	}
	elm := fmt.Sprintf("<meta name=\"%s\" value=\"%s\">", s[2], val)
	var buf bytes.Buffer
	err := quick.Highlight(&buf, elm, "html", "terminal256", "lovelace")
	if err != nil {
		fmt.Printf("\n%s\n", elm)
	} else {
		fmt.Printf("\n%v\n", buf.String())
	}
	fmt.Println(cf("About this element: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta#attr-name"))
}

func promptString(keep string) {
	// allow multiple word user input
	scanner := bufio.NewScanner(os.Stdin)
	var save string
	for scanner.Scan() {
		txt := scanner.Text()
		switch txt {
		case "":
			os.Exit(0)
		case "-":
			save = ""
		default:
			save = txt
		}
		viper.Set(configSetName, save)
		fmt.Printf("%s %s is now set to \"%v\"\n", cs("✓"), cp(configSetName), save)
		writeConfig(true)
		os.Exit(0)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}

}

func promptBool() {
	var input string
	cnt := 0
	for {
		input = ""
		cnt++
		fmt.Scanln(&input)
		if input == "" {
			promptCheck(cnt)
			continue
		}
		switch input {
		case "n", "no", "f", "false":
			viper.Set(configSetName, false)
			fmt.Printf("%s %s is now disabled\n", cs("✓"), cp(configSetName))
		case "y", "yes", "t", "true":
			viper.Set(configSetName, true)
			fmt.Printf("%s %s is now enabled\n", cs("✓"), cp(configSetName))
		default:
			fmt.Printf("%s %v\n", ce("✗"), input)
			promptCheck(cnt)
			continue
		}
		writeConfig(true)
		os.Exit(0)
	}
}

func promptSlice(s string) {
	keys := strings.Split(s, ", ")
	sort.Strings(keys)
	var input string
	cnt := 0
	for {
		input = ""
		cnt++
		fmt.Scanln(&input)
		if input == "" {
			promptCheck(cnt)
			continue
		}
		var i = sort.SearchStrings(keys, input)
		if i >= len(keys) || keys[i] != input {
			fmt.Printf("%s %v\n", ce("✗"), input)
			promptCheck(cnt)
			continue
		}
		viper.Set(configSetName, input)
		fmt.Printf("%s %s is now set to \"%v\"\n", cs("✓"), cp(configSetName), input)
		writeConfig(true)
		os.Exit(0)
	}
}

// ValidPort checks that p can be used as a network port value
func ValidPort(p int) bool {
	if p < portMin || p > portMax {
		return false
	}
	return true
}

func promptCheck(cnt int) {
	switch {
	case cnt == 2:
		fmt.Println("Ctrl+C to keep the existing port")
	case cnt >= 4:
		os.Exit(1)
	}
}

func promptPort() {
	var input string
	cnt := 0
	for {
		input = ""
		cnt++
		fmt.Scanln(&input)
		if input == "" {
			promptCheck(cnt)
			continue
		}
		i, err := strconv.ParseInt(input, 10, 0)
		if err != nil && input != "" {
			fmt.Printf("%s %v\n", ce("✗"), input)
			promptCheck(cnt)
			continue
		}
		// check that the input a valid port
		if v := ValidPort(int(i)); v == false {
			fmt.Printf("%s %v, is out of range\n", ce("✗"), input)
			promptCheck(cnt)
			continue
		}
		viper.Set(configSetName, i)
		fmt.Printf("%s %s is now set to \"%v\"\n", cs("✓"), cp(configSetName), i)
		writeConfig(true)
		os.Exit(0)
	}
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configCreateCmd)
	configCreateCmd.Flags().BoolVarP(&fileOverwrite, "overwrite", "y", false, "overwrite any existing config file")
	configCmd.AddCommand(configDeleteCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configInfoCmd)
	configInfoCmd.Flags().StringVarP(&infoStyles, "syntax-style", "c", "monokai", "config syntax highligher, use "+ci("none")+" to disable")
	configCmd.AddCommand(configShellCmd)
	configShellCmd.Flags().StringVarP(&configShell, "interpreter", "i", "", "user shell to receive retrotxtgo auto-completions\nchoices: "+ci(shells))
	//configShellCmd.Flags().BoolVar(&shellPreview, "preview", false, "prints the shell completion instead of applying it")
	configShellCmd.MarkFlagRequired("interpreter")
	configShellCmd.SilenceErrors = true
	configCmd.AddCommand(configSetCmd)
	configSetCmd.Flags().StringVarP(&configSetName, "name", "n", "", `the configuration path to edit in dot syntax (see examples)
to see a list of names run: retrotxt config info`)
	configSetCmd.MarkFlagRequired("name")
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

func writeConfig(update bool) {
	bs, err := yaml.Marshal(viper.AllSettings())
	Check(ErrorFmt{"could not create", "settings", err})
	d, err := os.UserHomeDir()
	Check(ErrorFmt{"could not use", "user home directory", err})
	err = ioutil.WriteFile(d+"/.retrotxtgo.yaml", bs, 0660)
	Check(ErrorFmt{"could not write", "settings", err})
	s := "Created a new"
	if update == true {
		s = "Updated the"
	}
	fmt.Println(s+" config file at:", cf(d+"/.retrotxtgo.yaml"))
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
