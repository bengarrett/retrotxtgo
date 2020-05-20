package config

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	v "github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/spf13/viper"
)

// Set edits and saves a setting within a configuration file.
func Set(name string) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, name); i == len(keys) || keys[i] != name {
		err := fmt.Errorf("to see a list of usable settings, run: retrotxt config info")
		logs.ChkErr(logs.Err{Issue: "invalid flag", Arg: fmt.Sprintf("--name %s", name), Msg: err})
	}
	value := viper.GetString(name)
	switch value {
	case "":
		fmt.Printf("\n%s is currently disabled\n", logs.Cp(name))
	default:
		fmt.Printf("\n%s is currently set to %q\n", logs.Cp(name), value)
	}
	switch name {
	case "create.layout":
		fmt.Printf("Set a new value, choice: %s\n",
			logs.Ci(createTemplates().String()))
		setStrings(name, createTemplates().Strings())
	case "info.format":
		fmt.Printf("Set a new value, choice: %s\n",
			logs.Ci(Format.String("info")))
		setStrings(name, Format.Info)
	case "version.format":
		fmt.Printf("Set a new value, choice: %s\n",
			logs.Ci(Format.String("version")))
		setStrings(name, Format.Version)
	case "create.server-port":
		fmt.Printf("Set a new HTTP port, choices: %d-%d (recommended: %d)\n",
			port.min, port.max, port.rec)
		setPort(name)
	case "create.meta.generator":
		setGenerator()
	default:
		q := "Set a new value or leave blank to keep it disabled:"
		setMeta(name, value)
		if value != "" {
			q = "Set a new value, leave blank to keep as-is or use a dash [-] to disable:"
		}
		fmt.Printf("\n%s \n", q)
		setString(value)
	}
}

// createTemplates creates a map of the template filenames used in conjunction with the layout flag.
func createTemplates() files {
	f := make(files)
	f["body"] = "body-content"
	f["full"] = "standard"
	f["mini"] = "standard"
	f["pre"] = "pre-content"
	f["standard"] = "standard"
	return f
}

// String method returns the files keys as a comma separated list.
func (f files) String() string {
	s := []string{}
	for key := range createTemplates() {
		s = append(s, key)
	}
	sort.Strings(s)
	return strings.Join(s, ", ")
}

// Strings method returns the files keys as a sorted slice.
func (f files) Strings() []string {
	s := []string{}
	for key := range createTemplates() {
		s = append(s, key)
	}
	sort.Strings(s)
	return s
}

func setGenerator() {
	var name = "create.meta.generator"
	// v{{.BuildVersion}}; {{.BuildDate}}
	elm := fmt.Sprintf("<head>\n  <meta name=\"generator\" content=\"RetroTxt v%s, %s\">",
		v.B.Version, v.B.Date)
	logs.ColorHTML(elm)
	viper.Set(name, logs.PromptYN("Enable this element", viper.GetBool(name)))
}

func setMeta(name, value string) {
	s := strings.Split(name, ".")
	switch {
	case len(s) != 3, s[0] != "create", s[1] != "meta":
		return
	}
	elm := fmt.Sprintf("<head>\n  <meta name=\"%s\" value=\"%s\">", s[2], value)
	logs.ColorHTML(elm)
	fmt.Println(logs.Cf("About this element: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta#attr-name"))
}

func setPort(name string) {
	save(name, logs.PromptPort(true))
}

func setString(name string) {
	save(name, logs.PromptString())
}

func setStrings(name string, data []string) {
	save(name, logs.PromptStrings(&data))
}

func save(name string, value interface{}) {
	viper.Set(name, value)
	fmt.Printf("%s %s is now set to \"%v\"\n", logs.Cs("✓"), logs.Cp(name), value)
	UpdateConfig(name, false)
	os.Exit(0)
}
