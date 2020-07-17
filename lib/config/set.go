package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/chroma/styles"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/create"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/prompt"
	"retrotxt.com/retrotxt/lib/str"
	v "retrotxt.com/retrotxt/lib/version"
)

// ColorCSS prints colored CSS syntax highlighting.
func ColorCSS(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "css", style)
}

// ColorHTML prints colored syntax highlighting to HTML elements.
func ColorHTML(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "html", style)
}

// List and print all the available configurations.
func List() (err error) {
	keys := Keys()
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 0, ' ', 0)
	fmt.Fprintf(w, "\t\tname value\t\thint\n")
	for i, key := range keys {
		fmt.Fprintf(w, "%d\t\t%s\t\t%s", i, key, Hints[key])
		switch key {
		case "html.layout":
			fmt.Fprintf(w, ", choices: %s (recommend: %s)",
				str.Cp(strings.Join(create.Layouts(), ", ")), str.Cp("standard"))
		case "serve":
			fmt.Fprintf(w, ", choices: %s", portInfo())
		}
		fmt.Fprint(w, "\n")
	}
	return w.Flush()
}

// Names lists the names of chroma styles.
func Names() string {
	var s names = styles.Names()
	return s.string(true)
}

// Set edits and saves a named setting within a configuration file.
// It also accepts numeric index values printed by List().
func Set(name string) {
	if i, err := strconv.Atoi(name); err != nil {
		Update(name)
	} else if i >= 0 && i <= (len(Defaults)-1) {
		k := Keys()
		Update(k[i])
	} else {
		Update(name)
	}
}

// Update edits and saves a named setting within a configuration file.
func Update(name string) {
	if !Validate(name) {
		h := logs.Hint{
			Error: logs.Generic{
				Issue: "invalid setting name",
				Arg:   fmt.Sprintf("%q for config", name),
				Err:   errors.New(("to see a list of setting names")),
			},
			Hint: "config set --list",
		}
		fmt.Println(h.String())
		return
	}
	if !setupMode {
		PrintLocation()
	}
	// print the current status of the named setting
	value := viper.Get(name)
	switch value.(type) {
	case nil:
		// avoid potential panics from missing settings by implementing the default value
		viper.Set(name, Defaults[name])
		value = viper.Get(name)
	}
	if b, ok := value.(bool); ok {
		switch b {
		case true:
			fmt.Printf("\n%s is in use\n", str.Cf(name))
		default:
			fmt.Printf("\n%s is currently not in use\n", str.Cf(name))
		}
	}
	if s, ok := value.(string); ok {
		switch s {
		case "":
			fmt.Printf("\n%s is currently not in use\n", str.Cf(name))
		default:
			fmt.Printf("\n%s is set to %q", str.Cf(name), value)
			// print the operating system's ability to use the existing set values
			// does the 'editor' exist in the env path, does the save-directory exist?
			switch name {
			case "editor":
				_, err := exec.LookPath(value.(string))
				fmt.Print(" ", str.Bool(err == nil))
				fmt.Println()
			case "save-directory":
				f := false
				if _, err := os.Stat(value.(string)); !os.IsNotExist(err) {
					f = true
				}
				fmt.Print(" ", str.Bool(f))
				fmt.Println()
			default:
				fmt.Println()
			}
		}
	}
	// print the setting user input prompt
	switch name {
	case "editor":
		s := fmt.Sprint("Set a " + Hints[name])
		if value.(string) != "" {
			s = fmt.Sprint(s, " or use a dash [-] to remove")
		}
		fmt.Printf("%s:\n", s)
		setEditor(name)
	case "html.font.embed":
		setFontEmbed(value.(bool))
	case "html.font.family":
		setFont(value.(string))
	case "html.layout":
		fmt.Println("\nChoose a new " + str.Options(Hints[name], true, create.Layouts()...))
		fmt.Println("\n  standard: uses external CSS, JS and woff2 fonts and is the recommended layout for web servers")
		fmt.Println("  inline:   includes both the CSS and JS as inline elements but is not recommended")
		fmt.Println("  compact:  is the same as the standard layout but without any <meta> tags")
		fmt.Println("  none:     no template is used, instead only the generated markup is returned")
		setShortStrings(name, create.Layouts()...)
	case "html.meta.author",
		"html.meta.description",
		"html.meta.keywords",
		"html.meta.theme-color":
		previewMeta(name, value.(string))
		setString(name)
	case "html.meta.color-scheme":
		previewMeta(name, value.(string))
		var prints = make([]string, len(create.ColorScheme[:]))
		copy(prints, create.ColorScheme[:])
		fmt.Println(str.UnderlineKeys(prints...))
		setShortStrings(name, create.ColorScheme[:]...)
	case "html.meta.generator":
		setGenerator(value.(bool))
	case "html.meta.notranslate":
		setNoTranslate(value.(bool))
	case "html.meta.referrer":
		previewMeta(name, value.(string))
		fmt.Println(str.NumberizeKeys(create.Referrer[:]...))
		setIndex(name, create.Referrer[:]...)
	case "html.meta.robots":
		previewMeta(name, value.(string))
		fmt.Println(str.NumberizeKeys(create.Robots[:]...))
		setIndex(name, create.Robots[:]...)
	case "html.meta.retrotxt":
		setRetrotxt(value.(bool))
	case "html.title":
		previewTitle(value.(string))
		fmt.Println("Choose a new " + Hints[name] + ":")
		setString(name)
	case "save-directory":
		fmt.Println("Choose a new " + Hints[name] + ":")
		setDirectory(name)
	case "serve":
		p := Defaults["serve"].(uint)
		if u, ok := value.(uint); ok {
			p = u
		} else if i, ok := value.(int); ok {
			p = uint(i)
		}
		fmt.Printf("\n%slocalhost%s%d %s\n", str.Cb("http://"),
			str.Cb(":"), p, str.Bool(create.Port(p)))
		fmt.Printf("\nSet a new port value, to %s\nChoices %s:\n", Hints[name], portInfo())
		setPort(name)
	case "style.html":
		fmt.Printf("Set a new HTML syntax style:\n%s\n", str.Ci(Names()))
		setStrings(name, styles.Names()...)
	case "style.info":
		fmt.Printf("Set a new %s syntax style:\n%s\n", str.Example("config info"), str.Ci(Names()))
		setStrings(name, styles.Names()...)
	default:
		log.Fatalln("config is not configured:", name)
	}
}

// Validate the existence of a setting key name.
func Validate(key string) (ok bool) {
	ok = false
	keys := Keys()
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return ok
	}
	return true
}

func colorElm(elm, lexer, style string) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style); err != nil {
		logs.Fatal("logs", "colorhtml", err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
}

type names []string

func (n names) string(theme bool) string {
	maxWidth := 0
	for _, s := range n {
		if l := len(fmt.Sprintf("%s=%q", s, s)); l > maxWidth {
			maxWidth = l
		}
	}
	if !theme {
		return strings.Join(n, ", ")
	}
	var s = make([]string, len(n))
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
		if err := str.HighlightWriter(&b, t, "yaml", name); err != nil {
			logs.Println("yaml highlight writer failed", name, err)
		}
		s = append(s, b.String())
	}
	return strings.Join(s, "")
}

// dirExpansion traverses the named directory to apply shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func dirExpansion(name string) (dir string) {
	const sep = string(os.PathSeparator)
	if name == "" || name == sep {
		return name
	}
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	r, paths := bool(name[0:1] == sep), strings.Split(name, sep)
	for i, s := range paths {
		p := ""
		switch s {
		case "~":
			var err error
			p, err = os.UserHomeDir()
			if err != nil {
				logs.LogFatal(err)
			}
		case ".":
			if i != 0 {
				continue
			}
			var err error
			p, err = os.Getwd()
			if err != nil {
				logs.LogFatal(err)
			}
		case "..":
			if i == 0 {
				wd, err := os.Getwd()
				if err != nil {
					logs.LogFatal(err)
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
	if r {
		dir = filepath.Join(sep, dir)
	}
	return dir
}

func portInfo() string {
	type ports struct {
		max uint
		min uint
		rec uint
	}
	var port = ports{
		max: prompt.PortMax,
		min: prompt.PortMin,
		rec: prompt.PortRec,
	}
	pm, px, pr := strconv.Itoa(int(port.min)), strconv.Itoa(int(port.max)), strconv.Itoa(int(port.rec))
	return str.Cp(pm) + "-" + str.Cp(px) + fmt.Sprintf(" (recommend: %s)", str.Cp(pr))
}

func previewMeta(name, value string) {
	if name == "" {
		logs.LogFatal(errors.New("preview meta name cannot be empty"))
	}
	if !Validate(name) {
		logs.LogFatal(fmt.Errorf("preview meta setting name is not known: %q", name))
	}
	s := strings.Split(name, ".")
	switch {
	case len(s) != 3, s[0] != "html", s[1] != "meta":
		return
	}
	elm := fmt.Sprintf("<head>\n  <meta name=\"%s\" value=\"%s\">", s[2], value)
	fmt.Print(ColorHTML(elm))
	h := strings.Split(Hints[name], " ")
	fmt.Printf("\n%s %s.", strings.Title(h[0]), strings.Join(h[1:], " "))
	fmt.Println(str.Cf("\nAbout this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"))
	fmt.Printf("\n%s \n", previewPrompt(name, value))
}

func previewTitle(value string) {
	elm := fmt.Sprintf("<head>\n  <title>%s</title>", value)
	fmt.Print(ColorHTML(elm))
	fmt.Println(str.Cf("\nAbout this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"))
	fmt.Println()
}

func previewPrompt(name, value string) (p string) {
	p = "Set a new value"
	switch name {
	case "html.meta.keywords":
		p = "Set some comma-separated keywords"
		if value != "" {
			p = "Replace the current keywords"
		}
	}
	if value != "" {
		p += ", leave blank to keep as-is or use a dash [-] to remove:"
	} else {
		p += " or leave blank to keep it unused:"
	}
	return
}

func save(name string, value interface{}) {
	if name == "" {
		logs.LogFatal(errors.New("save name value cannot be empty"))
	}
	if !Validate(name) {
		logs.LogFatal(fmt.Errorf("save name is not known: %q", name))
	}
	// don't save unchanged input values
	if viper.Get(name) == fmt.Sprint(value) {
		if setupMode {
			return
		}
		os.Exit(0)
	}
	// save named value
	viper.Set(name, value)
	if err := UpdateConfig("", false); err != nil {
		logs.LogFatal(err)
	}
	fmt.Printf("%s %s is set to \"%v\"\n", str.Cs("✓"), str.Cs(name), value)
	if !setupMode {
		os.Exit(0)
	}
}

func setDirectory(name string) {
	if name == "" {
		logs.LogFatal(errors.New("set directory name cannot be empty"))
	}
	dir := dirExpansion(prompt.String())
	if setupMode && dir == "" {
		return
	}
	if dir == "-" {
		dir = ""
	}
	if _, err := os.Stat(dir); err != nil {
		es := fmt.Sprint(err)
		e := strings.Split(es, ":")
		if len(e) > 1 {
			es = strings.TrimSpace(strings.Join(e[1:], ""))
		}
		fmt.Printf("%s the directory is invalid: %s\n", str.Alert(), es)
	}
	save(name, dir)
}

func setEditor(name string) {
	if name == "" {
		logs.LogFatal(errors.New("set editor name cannot be empty"))
	}
	v := prompt.String()
	switch v {
	case "-":
		save(name, "")
		return
	case "":
		return
	}
	if _, err := exec.LookPath(v); err != nil {
		fmt.Printf("%s this editor choice is not accessible by RetroTxt\n%s\n",
			str.Info(), err.Error())
	}
	save(name, v)
}

func setFont(value string) {
	var (
		b bytes.Buffer
		f = create.Family(value)
	)
	if f == create.Automatic {
		f = create.VGA
	}
	fmt.Fprintln(&b, "@font-face {")
	fmt.Fprintf(&b, "  font-family: \"%s\";\n", f.String())
	fmt.Fprintf(&b, "  src: url(\"%s.woff2\") format(\"woff2\");\n", f.String())
	fmt.Fprintln(&b, "  font-display: swap;\n}")
	fmt.Print(ColorCSS(b.String()))
	fmt.Println(str.Cf("About font families: https://developer.mozilla.org/en-US/docs/Web/CSS/font-family"))
	fmt.Println("Choose a font (recommend: automatic):")
	fmt.Println(str.UnderlineKeys(create.Fonts()...))
	setShortStrings("html.font.family", create.Fonts()...)
}

func setFontEmbed(value bool) {
	var name = "html.font.embed"
	elm := "@font-face{\n  font-family: vga8;\n  src: url(data:font/woff2;base64,[a large font binary will be embedded here]...) format('woff2');\n}"
	fmt.Println(ColorCSS(elm))
	q := "This is not recommended unless you need to create self-contained HTML files for offline distribution.\nEmbed the font as base64 data in the HTML"
	if value {
		q = "Keep the embedded font option"
	}
	v := prompt.YesNo(q, viper.GetBool(name))
	save(name, v)
}

func setGenerator(value bool) {
	var name, ver = "html.meta.generator", v.Semantic(v.B.Version)
	elm := fmt.Sprintf("<head>\n  <meta name=\"generator\" content=\"RetroTxt v%s, %s\">",
		ver.String(), v.B.Date)
	fmt.Println(ColorHTML(elm))
	p := "Enable the generator element"
	if value {
		p = "Keep the generator element"
	}
	viper.Set(name, prompt.YesNo(p, viper.GetBool(name)))
	if err := UpdateConfig("", false); err != nil {
		logs.LogFatal(err)
	}
}

func setIndex(name string, data ...string) {
	if name == "" {
		logs.LogFatal(errors.New("setindex name string is empty"))
	}
	p := prompt.IndexStrings(&data)
	switch p {
	case "-":
		p = ""
	case "":
		return
	}
	save(name, p)
}

func setNoTranslate(value bool) {
	var name = "html.meta.notranslate"
	elm := "<html translate=\"no\">\n  <head>\n    <meta name=\"google\" content=\"notranslate\">"
	fmt.Println(ColorHTML(elm))
	q := "Enable the no translate option"
	if value {
		q = "Keep the translate option"
	}
	v := prompt.YesNo(q, viper.GetBool(name))
	save(name, v)
}

func setPort(name string) {
	if name == "" {
		logs.LogFatal(errors.New("set port name cannot be empty"))
	}
	v := prompt.Port(true)
	if setupMode && v == 0 {
		return
	}
	save(name, v)
}

func setRetrotxt(value bool) {
	var name = "html.meta.retrotxt"
	elm := "<head>\n  <meta name=\"retrotxt\" content=\"encoding: IBM437; newline: CRLF; length: 50; width: 80; name: file.txt\">"
	fmt.Println(ColorHTML(elm))
	p := "Enable the retrotxt element"
	if value {
		p = "Keep the retrotxt element"
	}
	viper.Set(name, prompt.YesNo(p, viper.GetBool(name)))
	if err := UpdateConfig("", false); err != nil {
		logs.LogFatal(err)
	}
}

func setShortStrings(name string, data ...string) {
	v := prompt.ShortStrings(&data)
	switch v {
	case "-":
		v = ""
	case "":
		return
	}
	save(name, v)
}

func setString(name string) {
	if name == "" {
		logs.LogFatal(errors.New("set string name cannot be empty"))
	}
	v := prompt.String()
	switch v {
	case "-":
		v = ""
	case "":
		return
	}
	save(name, v)
}

func setStrings(name string, data ...string) {
	if name == "" {
		logs.LogFatal(errors.New("set strings name cannot be empty"))
	}
	v := prompt.Strings(&data)
	switch v {
	case "-":
		v = ""
	case "":
		return
	}
	save(name, v)
}
