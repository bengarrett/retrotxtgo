package config

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	v "github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/spf13/viper"
)

type hints map[string]string

type files map[string]string

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

// TODO: improve descriptions and use them for create --flag hints
func list() hints {
	pm, px, pr := strconv.Itoa(int(port.min)), strconv.Itoa(int(port.max)), strconv.Itoa(int(port.rec))
	ports := logs.Cp(pm) + "-" + logs.Cp(px) + fmt.Sprintf(" (recommend: %s)", logs.Cp(pr))
	return hints{
		"create.layout": "HTML output layout, choices: " +
			logs.Cp(createTemplates().String()),
		"create.meta.author":       "defines the name of the page authors",
		"create.meta.color-scheme": "specifies one or more color schemes with which the page is compatible",
		"create.meta.description":  "a short and accurate summary of the content of the page",
		"create.meta.generator":    "include the RetroTxt version and page generation date?",
		"create.meta.keywords":     "words relevant to the page content",
		"create.meta.referrer":     "controls the Referer HTTP header attached to requests sent from the page",
		"create.meta.theme-color":  "indicates a suggested color that user agents should use to customize the display of the page",
		"create.save-directory":    "directory that all HTML files get saved to",
		"create.server-port":       "serve HTML over an internal web server, choices: " + ports,
		"create.title":             "page title that is shown in a browser title bar or tab",
		"editor":                   "set an text editor to launch when using " + logs.Example("config edit"),
		"style.html":               "syntax highlighter for html previews",
		"style.yaml":               "syntax highlighter for info and version commands",
	}
}

// List all the available configurations that can be passed to the --name flag.
func List() (err error) {
	hints := list()
	keys := viper.AllKeys()
	sort.Strings(keys)
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 0, ' ', 0)
	fmt.Fprintf(w, "\t\tname value\t\thint\n")
	for i, key := range keys {
		fmt.Fprintf(w, "%d\t\t%s\t\t%s\n", i, key, hints[key])
	}
	return w.Flush()
}

// Set edits and saves a setting within a configuration file.
func Set(name string) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, name); i == len(keys) || keys[i] != name {
		h := logs.Hint{
			Issue: "invalid value",
			Arg:   fmt.Sprintf("%q for --name", name),
			Msg:   errors.New(("to see a list of usable settings")),
			Hint:  "config info",
		}
		fmt.Println(h.String())
		return
	}
	//
	PrintLocation()
	value := viper.GetString(name)
	switch value {
	case "":
		fmt.Printf("\n%s is currently disabled\n", logs.Cp(name))
	default:
		fmt.Printf("\n%s is currently set to %q\n", logs.Cp(name), value)
	}
	hints := list()
	switch name {
	case "create.layout":
		fmt.Println("Choose a new " + hints[name])
		setStrings(name, createTemplates().Strings())
	case "create.meta.generator":
		setGenerator()
	case "create.save-directory":
		fmt.Println("Choose a new " + hints[name])
		setString(value) // TODO: setDirectory? check exist
	case "create.server-port":
		fmt.Println("Set a new HTTP port to " + hints[name])
		setPort(name)
	case "create.title":
		fmt.Println("Choose a new value " + hints[name])
		setString(value)
	case "style.html":
		fmt.Printf("Choose a new value, choice: %s\n",
			logs.Ci(Format.String("info")))
		// TODO sample HTML
		setStrings(name, Format.Info)
	case "style.yaml":
		fmt.Printf("Set a new value, choice: %s\n",
			logs.Ci(Format.String("version")))
		// logs.Ci(Format.String("info")))
		setStrings(name, Format.Version)
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
