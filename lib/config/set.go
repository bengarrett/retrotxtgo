package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"text/template"

	"github.com/alecthomas/chroma/styles"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	v "github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/spf13/viper"
)

type files map[string]string

type names []string

func (n names) String(theme bool) string {
	maxWidth := 0
	for _, s := range n {
		if l := len(fmt.Sprintf("%s=%q", s, s)); l > maxWidth {
			maxWidth = l
		}
	}
	if !theme {
		return strings.Join(n, ", ")
	}
	var s []string
	split := (len(n) / 2)
	for i, name := range n {
		var t string
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		t = fmt.Sprintf(" %2d) %s=%q%s", i, name, name, strings.Repeat(" ", pad+2))
		if split+i < len(n) {
			t += fmt.Sprintf("%2d) %s=%q\n", split+i, n[split+i], n[split+i])
		} else {
			break
		}
		var b bytes.Buffer
		str.HighlightWriter(&b, t, "yaml", name)
		s = append(s, b.String())
	}
	return strings.Join(s, "")
}

// Hints provide brief help on the config file configurations.
var Hints = map[string]string{
	"create.layout":            "HTML output layout",
	"create.meta.author":       "defines the name of the page authors",
	"create.meta.color-scheme": "specifies one or more color schemes with which the page is compatible",
	"create.meta.description":  "a short and accurate summary of the content of the page",
	"create.meta.generator":    "include the RetroTxt version and page generation date?",
	"create.meta.keywords":     "words relevant to the page content",
	"create.meta.notranslate":  "used to declare that the page should not be translated by Google Translate",
	"create.meta.referrer":     "controls the Referer HTTP header attached to requests sent from the page",
	"create.meta.robots":       "behaviour that crawlers from Google, Bing and other engines should use with the page",
	"create.meta.theme-color":  "indicates a suggested color that user agents should use to customize the display of the page",
	"create.save-directory":    "directory to store RetroTxt created HTML files",
	"create.server":            "serve HTML over an internal web server",
	"create.server-port":       "port which the internal web server will use",
	"create.title":             "page title that is shown in a browser title bar or tab",
	"editor":                   "text editor to launch when using " + str.Example("config edit"),
	"style.html":               "syntax highlighter for html previews",
	"style.yaml":               "syntax highlighter for info and version commands",
}

var setupMode = false

type ports struct {
	max uint
	min uint
	rec uint
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
func (f files) String() (s string) {
	k := []string{}
	for key := range createTemplates() {
		k = append(k, key)
	}
	sort.Strings(k)
	// apply an ANSI underline to the first letter of each key
	t, err := template.New("underline").Parse("{{define \"TEXT\"}}\033[0m\033[4m{{.}}\033[0m{{end}}")
	if err != nil {
		logs.LogCont(err)
		return strings.Join(k, ", ")
	}
	for i, key := range k {
		if len(k) > 1 {
			var b bytes.Buffer
			err := t.ExecuteTemplate(&b, "TEXT", string(key[0]))
			logs.LogCont(err)
			k[i] = fmt.Sprintf("%s%s", b.String(), key[1:])
		}
	}
	return strings.Join(k, ", ")
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

func portInfo() string {
	var port = ports{
		max: prompt.PortMax,
		min: prompt.PortMin,
		rec: prompt.PortRec,
	}
	pm, px, pr := strconv.Itoa(int(port.min)), strconv.Itoa(int(port.max)), strconv.Itoa(int(port.rec))
	return str.Cp(pm) + "-" + str.Cp(px) + fmt.Sprintf(" (recommend: %s)", str.Cp(pr))
}

// List all the available configurations that can be passed to the --name flag.
func List() (err error) {
	keys := viper.AllKeys()
	sort.Strings(keys)
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 0, ' ', 0)
	fmt.Fprintf(w, "\t\tname value\t\thint\n")
	for i, key := range keys {
		fmt.Fprintf(w, "%d\t\t%s\t\t%s", i, key, Hints[key])
		switch key {
		case "create.layout":
			fmt.Fprintf(w, ", choices: %s (recommend: %s)", str.Cp(createTemplates().String()), str.Cp("standard"))
		case "create.server-port":
			fmt.Fprintf(w, ", choices: %s", portInfo())
		}
		fmt.Fprint(w, "\n")
	}
	return w.Flush()
}

// Names lists the names of chroma styles.
func Names() string {
	var s names = styles.Names()
	return s.String(true)
}

// Set edits and saves a setting within a configuration file.
func Set(name string) {
	if !Validate(name) {
		h := logs.Hint{
			Issue: "invalid value",
			Arg:   fmt.Sprintf("%q for --name", name),
			Msg:   errors.New(("to see a list of usable settings")),
			Hint:  "config info -c",
		}
		fmt.Println(h.String())
		return
	}
	if !setupMode {
		PrintLocation()
	}
	value := viper.GetString(name)
	switch value {
	case "":
		fmt.Printf("\n%s is currently disabled\n", str.Cf(name))
	default:
		fmt.Printf("\n%s is currently set to %q\n", str.Cf(name), value)
	}
	switch name {
	case "create.layout":
		fmt.Println("Choose a new " + str.Options(Hints[name], create.Layouts(), true))
		setShortStrings(name, createTemplates().Strings())
	case "create.meta.generator":
		setGenerator()
	case "create.meta.referrer":
		setMeta(name, value)
		fmt.Println(strings.Join(create.Referrer, ", "))
		setShortStrings(name, create.Referrer)
	case "create.meta.robots":
		setMeta(name, value)
		fmt.Println(strings.Join(create.Robots, ", "))
		setString(name)
	case "create.meta.color-scheme":
		setMeta(name, value)
		fmt.Println(strings.Join(create.ColorScheme, ", "))
		setShortStrings(name, create.ColorScheme)
	case "create.save-directory":
		fmt.Println("Choose a new " + Hints[name] + ":")
		setDirectory(name)
	case "create.server-port":
		fmt.Println("Set a new HTTP port to " + Hints[name])
		setPort(name)
	case "create.title":
		fmt.Println("Choose a new " + Hints[name] + ":")
		setString(name)
	case "editor":
		fmt.Println("Set a " + Hints[name] + ":")
		setEditor(name)
	case "style.html":
		fmt.Printf("Set a new HTML syntax style:\n%s\n", str.Ci(Names()))
		setStrings(name, styles.Names())
	case "style.yaml":
		fmt.Printf("Set a new YAML syntax style:\n%s\n", str.Ci(Names()))
		setStrings(name, styles.Names())
	default:
		setMeta(name, value)
		setString(name)
	}
}

// Setup walks through all the settings within a configuration file.
func Setup() {
	setupMode = true
	prompt.SetupMode = true
	keys := viper.AllKeys()
	sort.Strings(keys)
	PrintLocation()
	w := 42
	for i, key := range keys {
		if i == 0 {
			fmt.Printf("\n%s\n\n", str.Cb(enterKey()))
		}
		h := fmt.Sprintf("  %d/%d. RetroTxt Setup - %s",
			i+1, len(keys), key)
		if i == 0 {
			fmt.Println(hr(w))
		}
		fmt.Println(h)
		Set(key)
		fmt.Println(hr(w))
	}
}

func enterKey() string {
	if runtime.GOOS == "darwin" {
		return "Press ↩ return to skip the question or ⌃ control-c to quit"
	}
	return "Press ⏎ enter to skip the question or Ctrl-c to quit"
}

func hr(w int) string {
	return str.Cb(strings.Repeat("-", w))
}

// Validate the existence of a setting key name.
func Validate(key string) (ok bool) {
	ok = false
	keys := viper.AllKeys()
	sort.Strings(keys)
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return ok
	}
	return true
}

// dirExpansion traverses the named directory to apply shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func dirExpansion(name string) (dir string) {
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	var err error
	paths := strings.Split(name, string(os.PathSeparator))
	for i, s := range paths {
		p := ""
		switch s {
		case "~":
			p, err = os.UserHomeDir()
			if err != nil {
				logs.Log(err)
			}
		case ".":
			if i != 0 {
				continue
			}
			p, err = os.Getwd()
			if err != nil {
				logs.Log(err)
			}
		case "..":
			if i == 0 {
				wd, err := os.Getwd()
				if err != nil {
					logs.Log(err)
				}
				p = filepath.Dir(wd)
			} else {
				dir = filepath.Dir(dir)
				continue
			}
		default:
			p = s
		}
		dir = filepath.Join(dir, p)
	}
	return dir
}

func save(name string, value interface{}) {
	if name == "" {
		logs.Log(errors.New("save name string is empty"))
	}
	if !Validate(name) {
		logs.Log(errors.New("save name is an unknown setting: " + name))
	}
	if setupMode && fmt.Sprint(value) == "" {
		return
	}
	// don't save unchanged input values
	if viper.GetString(name) == fmt.Sprint(value) {
		if setupMode {
			return
		}
		os.Exit(0)
	}
	// save named value
	viper.Set(name, value)
	if err := UpdateConfig("", false); err != nil {
		logs.Log(err)
	}
	fmt.Printf("%s %s is set to \"%v\"\n", str.Cs("✓"), str.Cs(name), value)
	if !setupMode {
		os.Exit(0)
	}
}

func setDirectory(name string) {
	if name == "" {
		logs.Log(errors.New("setdirectory name string is empty"))
	}
	dir := dirExpansion(prompt.String())
	if setupMode && dir == "" {
		return
	}
	if _, err := os.Stat(dir); err != nil {
		es := fmt.Sprint(err)
		e := strings.Split(es, ":")
		if len(e) > 1 {
			es = fmt.Sprintf("%s", strings.TrimSpace(strings.Join(e[1:], "")))
		}
		fmt.Printf("%s this directory returned the following error: %s\n", str.Info(), es)
	}
	save(name, dir)
}

func setEditor(name string) {
	if name == "" {
		logs.Log(errors.New("setstring name string is empty"))
	}
	editor := prompt.String()
	if setupMode && editor == "" {
		return
	}
	if _, err := exec.LookPath(editor); err != nil {
		fmt.Printf("%s this editor choice is not accessible: %s\n", str.Info(), err)
	}
	save(name, editor)
}

func setGenerator() {
	var name = "create.meta.generator"
	elm := fmt.Sprintf("<head>\n  <meta name=\"generator\" content=\"RetroTxt v%s, %s\">",
		v.B.Version, v.B.Date)
	fmt.Println(logs.ColorHTML(elm))
	prmt := prompt.YesNo("Enable this element", viper.GetBool(name))
	viper.Set(name, prmt)
	if err := UpdateConfig("", false); err != nil {
		logs.Log(err)
	}
}

func setMeta(name, value string) {
	if name == "" {
		logs.Log(errors.New("setmeta name string is empty"))
	}
	if !Validate(name) {
		logs.Log(errors.New("setmeta name is an unknown setting: " + name))
	}
	s := strings.Split(name, ".")
	switch {
	case len(s) != 3, s[0] != "create", s[1] != "meta":
		return
	}
	elm := fmt.Sprintf("<head>\n  <meta name=\"%s\" value=\"%s\">", s[2], value)
	fmt.Print(logs.ColorHTML(elm))
	h := strings.Split(Hints[name], " ")
	fmt.Printf("\n%s %s.", strings.Title(h[0]), strings.Join(h[1:], " "))
	fmt.Println(str.Cf("\nAbout this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"))
	q := "Set a new value or leave blank to keep it disabled:"
	if value != "" {
		q = "Set a new value, leave blank to keep as-is or use a dash [-] to disable:"
	}
	fmt.Printf("\n%s \n", q)
}

func setPort(name string) {
	if name == "" {
		logs.Log(errors.New("setport name string is empty"))
	}
	p := prompt.Port(true)
	if setupMode && p == 0 {
		return
	}
	save(name, p)
}

func setShortStrings(name string, data []string) {
	if name == "" {
		logs.Log(errors.New("setstrings name string is empty"))
	}
	save(name, prompt.ShortStrings(&data))
}

func setString(name string) {
	if name == "" {
		logs.Log(errors.New("setstring name string is empty"))
	}
	save(name, prompt.String())
}

func setStrings(name string, data []string) {
	if name == "" {
		logs.Log(errors.New("setstrings name string is empty"))
	}
	save(name, prompt.Strings(&data))
}
