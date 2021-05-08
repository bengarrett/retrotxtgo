// nolint:goconst
package config

import (
	"bytes"
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
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	v "github.com/bengarrett/retrotxtgo/lib/version"
	"github.com/spf13/viper"
)

type update struct {
	name  string
	setup bool
	value interface{}
}

const skipped = ` ↵ skipped`

// ColorCSS prints colored CSS syntax highlighting.
func ColorCSS(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "css", style, true)
}

// ColorHTML prints colored syntax highlighting to HTML elements.
func ColorHTML(elm string) string {
	style := viper.GetString("style.html")
	return colorElm(elm, "html", style, true)
}

// List and print all the available configurations.
func List() (err error) {
	keys := Keys()
	w := tabwriter.NewWriter(os.Stdout, 2, 2, 0, ' ', 0)
	const title = "  Available RetroTxt configurations and settings"
	fmt.Fprintln(w, "\n"+str.Cp(title))
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintf(w, "Alias\t\tName value\t\tHint\n")
	for i, key := range keys {
		fmt.Fprintf(w, " %d\t\t%s\t\t%s", i, key, Tip()[key])
		switch key {
		case "html.layout":
			fmt.Fprintf(w, ", choices: %s (suggestion: %s)",
				str.Cp(strings.Join(create.Layouts(), ", ")), str.Cp("standard"))
		case "serve":
			fmt.Fprintf(w, ", choices: %s", portInfo())
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintln(w, "\nEither the Name value or the Alias can be used as the setting name")
	fmt.Fprintln(w, "\n"+str.Example(" retrotxt config set html.meta.description")+
		" to change the meta description setting")
	fmt.Fprintln(w, str.Example(" retrotxt config set 6")+
		" will also change the meta description setting")
	fmt.Fprintln(w, "\nMultiple settings are supported")
	fmt.Fprintln(w, "\n"+str.Example(" retrotxt config set style.html style.info"))
	fmt.Fprint(w, "\n")
	return w.Flush()
}

// Names lists the names of chroma styles.
func Names(lexer string) string {
	var s names = styles.Names()
	return s.string(true, lexer)
}

// Set edits and saves a named setting within a configuration file.
// It also accepts numeric index values printed by List().
func Set(name string) {
	i, err := strconv.Atoi(name)
	switch {
	case err != nil:
		Update(name, false)
	case i >= 0 && i <= (len(Reset())-1):
		k := Keys()
		Update(k[i], false)
	default:
		Update(name, false)
	}
}

// Update edits and saves a named setting within a configuration file.
func Update(name string, setup bool) {
	if !Validate(name) {
		fmt.Println(logs.Hint("config set --list", logs.ErrCfgName))
		return
	}
	if !setup {
		PrintLocation()
	}
	// print the current status of the named setting
	value := viper.Get(name)
	switch value.(type) {
	case nil:
		// avoid potential panics from missing settings by implementing the default value
		viper.Set(name, Reset()[name])
		value = viper.Get(name)
	default:
		// everything ok
	}
	if b, ok := value.(bool); ok {
		updateBool(b, name)
	}
	if s, ok := value.(string); ok {
		updateString(s, name, value.(string))
	}
	updatePrompt(update{name, setup, value})
}

func updateBool(b bool, name string) {
	switch b {
	case true:
		fmt.Printf("\n  %s is in use\n", str.Cf(name))
	default:
		fmt.Printf("\n  %s is currently not in use\n", str.Cf(name))
	}
}

func updateString(s, name, value string) {
	switch s {
	case "":
		fmt.Printf("\n  %s is currently not in use\n", str.Cf(name))
	default:
		fmt.Printf("\n  %s is set to %q", str.Cf(name), value)
		// print the operating system's ability to use the existing set values
		// does the 'editor' exist in the env path, does the save-directory exist?
		switch name {
		case "editor":
			_, err := exec.LookPath(value)
			fmt.Print(" ", str.Bool(err == nil))
			fmt.Println()
		case "save-directory":
			f := false
			if _, err := os.Stat(value); !os.IsNotExist(err) {
				f = true
			}
			fmt.Print(" ", str.Bool(f))
			fmt.Println()
		default:
			fmt.Println("skipped")
		}
	}
}

func recommend(s string) string {
	if s == "" {
		return " (suggestion: do not use)"
	}
	return fmt.Sprintf(" (suggestion: %s)", str.Cp(s))
}

// UpdatePrompt prompts the user for a config setting input.
func updatePrompt(u update) {
	switch u.name {
	case "editor":
		promptEditor(u)
	case "save-directory":
		promptSaveDir(u)
	case "serve":
		promptServe(u)
	case "style.html":
		promptStyleHTML(u)
	case "style.info":
		promptStyleInfo(u)
	default:
		metaPrompts(u)
	}
}

// MetaPrompts prompts the user for a meta setting.
func metaPrompts(u update) {
	switch u.name {
	case "html.font.embed":
		setFontEmbed(u.value.(bool), u.setup)
	case "html.font.family":
		setFont(u.value.(string), u.setup)
	case "html.layout":
		promptLayout(u)
	case "html.meta.author",
		"html.meta.description",
		"html.meta.keywords":
		previewMeta(u.name, u.value.(string))
		setString(u.name, u.setup)
	case "html.meta.theme-color":
		recommendMeta(u.name, u.value.(string), "")
		setString(u.name, u.setup)
	case "html.meta.color-scheme":
		promptColorScheme(u)
	case "html.meta.generator":
		setGenerator(u.value.(bool))
	case "html.meta.notranslate":
		setNoTranslate(u.value.(bool), u.setup)
	case "html.meta.referrer":
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Referrer()
		fmt.Println(str.NumberizeKeys(cr[:]...))
		setIndex(u.name, u.setup, cr[:]...)
	case "html.meta.robots":
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Robots()
		fmt.Println(str.NumberizeKeys(cr[:]...))
		setIndex(u.name, u.setup, cr[:]...)
	case "html.meta.retrotxt":
		setRetrotxt(u.value.(bool))
	case "html.title":
		previewTitle(u.value.(string))
		fmt.Println(" Choose a new " + Tip()[u.name] + ":")
		setString(u.name, u.setup)
	default:
		log.Fatalln("config is not configured:", u.name)
	}
}

// PromptColorScheme prompts the user for the color scheme setting.
func promptColorScheme(u update) {
	previewMeta(u.name, u.value.(string))
	ccc := create.ColorScheme()
	var prints = make([]string, len(ccc[:]))
	copy(prints, ccc[:])
	fmt.Println(" " + str.UnderlineKeys(prints...) + recommend(""))
	setShortStrings(u.name, u.setup, ccc[:]...)
}

// PromptEditor prompts the user for the editor setting.
func promptEditor(u update) {
	s := fmt.Sprint("Set a " + Tip()[u.name])
	if u.value.(string) != "" {
		s = fmt.Sprint(s, " or use a dash [-] to remove")
	} else if Editor() != "" {
		fmt.Printf("instead the %s editor will be run\n\n", str.Cp(Editor()))
	}
	fmt.Printf(" %s:\n", s)
	setEditor(u.name, u.setup)
}

// PromptLayout prompts the user for the layout setting.
func promptLayout(u update) {
	fmt.Println("\n  standard: uses external CSS, JS and woff2 fonts and is the recommended layout for web servers")
	fmt.Println("  inline:   includes both the CSS and JS as inline elements but is not recommended")
	fmt.Println("  compact:  is the same as the standard layout but without any <meta> tags")
	fmt.Println("  none:     no template is used, instead only the generated markup is returned")
	fmt.Println("\n Choose a " + str.Options(Tip()[u.name], true, create.Layouts()...) + " (suggestion: " + str.Cp("standard") + ")")
	setShortStrings(u.name, u.setup, create.Layouts()...)
}

// PromptSaveDir prompts the user for the a save destination directory setting.
func promptSaveDir(u update) {
	fmt.Println(" Choose a new " + Tip()[u.name] + ":")
	if home, err := os.UserHomeDir(); err == nil {
		fmt.Printf("\n  Use %s to save the home directory %s", str.Example("~"), str.Cb(home))
	}
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("\n      %s to save this current directory %s", str.Example("."), str.Cb(wd))
	}
	fmt.Printf("\n      %s to disable and always use the user's current directory\n\n", str.Example("-"))
	setDirectory(u.name, u.setup)
}

// PromptServe prompts the user for a HTTP server port setting.
func promptServe(u update) {
	var reset = func() {
		var p uint
		if u, ok := Reset()["serve"].(uint); ok {
			p = u
		}
		save(u.name, false, p)
	}
	var p uint
	switch v := u.value.(type) {
	case uint:
		p = v
	case int:
		p = uint(v)
	default:
		reset()
	}
	if p > prompt.PortMax {
		reset()
	}
	fmt.Printf("\n  HTTP server: %slocalhost%s%d %s\n", "http://",
		str.Cb(":"), p, str.Bool(create.Port(p)))
	fmt.Printf("\n Port %s is reserved, ", str.Example("0"))
	fmt.Printf("while the ports below %s are normally restricted\n by the operating system and are not recommended\n", str.Example("1024"))
	fmt.Printf("\n Set a HTTP port value, to %s\n Choices %s:\n", Tip()[u.name], portInfo())
	setPort(u.name, u.setup)
}

// PromptStyleHTML prompts the user for the a HTML and CSS style setting.
func promptStyleHTML(u update) {
	var d string
	if s, ok := Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n Set a new HTML syntax style%s:\n", str.Ci(Names("css")), recommend(d))
	setStrings(u.name, u.setup, styles.Names()...)
}

// PromptStyleInfo prompts the user for the a JS style setting.
func promptStyleInfo(u update) {
	var d string
	if s, ok := Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n Set a new %s syntax style%s:\n", str.Ci(Names("json")), str.Example("config info"), recommend(d))
	setStrings(u.name, u.setup, styles.Names()...)
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

// ColorElm applies color syntax to an element.
func colorElm(elm, lexer, style string, color bool) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style, color); err != nil {
		logs.ProblemMarkFatal(fmt.Sprint("html ", lexer), logs.ErrHighlight, err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
}

type names []string

// String lists and applies the named themes for the HighlightWriter.
func (n names) string(theme bool, lexer string) string {
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
	const space = 2
	for i, name := range n {
		var (
			b bytes.Buffer
			t string
		)
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		if lexer == "json" {
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }%s", name, name, strings.Repeat(" ", pad+space))
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.ProblemMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
			if split+i >= len(n) {
				break
			}
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }\n", n[split+i], n[split+i])
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.ProblemMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
			continue
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>%s", name, name, strings.Repeat(" ", pad+space))
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.ProblemMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>\n", n[split+i], n[split+i])
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.ProblemMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
	}
	return strings.Join(s, "")
}

// DirExpansion traverses the named directory to apply shell-like expansions.
// It currently supports limited Bash tilde, shell dot and double dot syntax.
func dirExpansion(name string) (dir string) {
	const sep, homeDir, currentDir, parentDir = string(os.PathSeparator), "~", ".", ".."
	if name == "" || name == sep {
		return name
	}
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	r, paths := bool(name[0:1] == sep), strings.Split(name, sep)
	for i, s := range paths {
		p := ""
		switch s {
		case homeDir:
			var err error
			p, err = os.UserHomeDir()
			if err != nil {
				logs.SaveFatal(err)
			}
		case currentDir:
			if i != 0 {
				continue
			}
			var err error
			p, err = os.Getwd()
			if err != nil {
				logs.SaveFatal(err)
			}
		case parentDir:
			if i == 0 {
				wd, err := os.Getwd()
				if err != nil {
					logs.SaveFatal(err)
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

// PortInfo returns valid and recommended HTTP port values.
func portInfo() string {
	type ports struct {
		max uint
		min uint
		rec uint
	}
	port := ports{
		max: prompt.PortMax,
		min: prompt.PortMin,
		rec: prompt.PortRec,
	}
	pm, px, pr := strconv.Itoa(int(port.min)), strconv.Itoa(int(port.max)), strconv.Itoa(int(port.rec))
	return str.Cp(pm) + "-" + str.Cp(px) + fmt.Sprintf(" (suggestion: %s)", str.Cp(pr))
}

// PreviewMeta previews and prompts for a meta element content value.
func previewMeta(name, value string) {
	previewMetaPrint(name, value)
	fmt.Printf("\n%s \n", previewPrompt(name, value))
}

func previewMetaPrint(name, value string) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("preview meta: %w", logs.ErrNameNil))
	}
	if !Validate(name) {
		logs.SaveFatal(fmt.Errorf("preview meta %q: %w", name, logs.ErrCfgName))
	}
	s := strings.Split(name, ".")
	const req = 3
	switch {
	case len(s) != req, s[0] != "html", s[1] != "meta":
		return
	}
	elm := fmt.Sprintf("<head>\n  <meta name=\"%s\" value=\"%s\">", s[2], value)
	fmt.Print(ColorHTML(elm))
	h := strings.Split(Tip()[name], " ")
	fmt.Println(str.Cf("\nAbout this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"))
	fmt.Printf("\n %s %s.", strings.Title(h[0]), strings.Join(h[1:], " "))
}

// PreviewTitle previews and prompts for the title element.
func previewTitle(value string) {
	elm := fmt.Sprintf("<head>\n  <title>%s</title>", value)
	fmt.Print(ColorHTML(elm))
	fmt.Println(str.Cf("\nAbout this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"))
	fmt.Println()
}

func previewPrompt(name, value string) string {
	return fmt.Sprintf("%s:", previewPromptPrint(name, value))
}

// PreviewPromptPrint returns the available input options.
func previewPromptPrint(name, value string) (p string) {
	p = "Set a new value"
	if name == "html.meta.keywords" {
		p = "Set some comma-separated keywords"
		if value != "" {
			p = "Replace the current keywords"
		}
	}
	if value != "" {
		return fmt.Sprintf(" %s, leave blank to keep as-is or use a dash [-] to remove", p)
	}
	return fmt.Sprintf(" %s or leave blank to keep it unused", p)
}

func recommendMeta(name, value, suggest string) {
	previewMetaPrint(name, value)
	fmt.Printf("\n%s\n", recommendPrompt(name, value, suggest))
}

func recommendPrompt(name, value, suggest string) string {
	p := previewPromptPrint(name, value)
	return fmt.Sprintf("%s%s:", p, recommend(suggest))
}

// Save value to the named configuration.
func save(name string, setup bool, value interface{}) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("save: %w", logs.ErrNameNil))
	}
	if !Validate(name) {
		logs.SaveFatal(fmt.Errorf("save %q: %w", name, logs.ErrCfgName))
	}
	// don't save unchanged input values
	if viper.Get(name) == fmt.Sprint(value) {
		if setup {
			return
		}
		os.Exit(0)
	}
	// save named value
	viper.Set(name, value)
	if err := UpdateConfig("", false); err != nil {
		logs.SaveFatal(err)
	}
	fmt.Printf(" %s %s is set to \"%v\"\n", str.Cs("✓"), str.Cs(name), value)
	if !setup {
		os.Exit(0)
	}
}

// SetDirectory checks the existence of a directory
// and saves the path as a configuration regardless of the result.
func setDirectory(name string, setup bool) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set directory: %w", logs.ErrNameNil))
	}
	dir := dirExpansion(prompt.String(setup))
	if setup && dir == "" {
		return
	}
	if dir == "-" {
		dir = ""
		save(name, setup, dir)
		return
	}
	if _, err := os.Stat(dir); err != nil {
		es := fmt.Sprint(err)
		e := strings.Split(es, ":")
		if len(e) > 1 {
			es = strings.TrimSpace(strings.Join(e[1:], ""))
		}
		fmt.Printf("%s the directory is invalid: %s\n", str.Alert(), es)
	}
	save(name, setup, dir)
}

// SetEditor checks the existence of given text editor location
// and saves it as a configuration regardless of the result.
func setEditor(name string, setup bool) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set editor: %w", logs.ErrNameNil))
	}
	val := prompt.String(setup)
	switch val {
	case "-":
		save(name, setup, "")
		return
	case "":
		fmt.Println(skipped)
		return
	}
	if _, err := exec.LookPath(val); err != nil {
		fmt.Printf("%s this editor choice is not accessible by RetroTxt\n%s\n",
			str.Info(), err.Error())
	}
	save(name, setup, val)
}

// SetFont previews and saves a default font setting.
func setFont(value string, setup bool) {
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
	fmt.Println(str.Cf("About font families: https://developer.mozilla.org/en-US/docs/Web/CSS/font-family") + "\n")
	fmt.Println(" Choose a font:")
	fmt.Println(" " + str.UnderlineKeys(create.Fonts()...) + " (suggestion: " + str.Cp("automatic") + ")")
	setShortStrings("html.font.family", setup, create.Fonts()...)
}

// SetFont previews and saves the embed Base64 font setting.
func setFontEmbed(value, setup bool) {
	name := "html.font.embed"
	elm := `@font-face{
  font-family: vga8;
  src: url(data:font/woff2;base64,[a large font binary will be embedded here]...) format('woff2');
}`
	fmt.Println(ColorCSS(elm))
	q := "This is not recommended, unless you need self-contained files for distribution.\n Embed the font as base64 data in the HTML"
	if value {
		q = "Keep the embedded font option"
	}
	q += recommend("no")
	val := prompt.YesNo(q, viper.GetBool(name))
	save(name, setup, val)
}

// SetGenerator previews and prompts the custom RetroTxt generator meta tag.
func setGenerator(value bool) {
	name, ver := "html.meta.generator", v.Semantic(v.Release.Version)
	elm := fmt.Sprintf("<head>\n  <meta name=\"generator\" content=\"RetroTxt v%s, %s\">",
		ver.String(), v.Release.Date)
	fmt.Println(ColorHTML(elm))
	p := "Enable the generator element"
	if value {
		p = "Keep the generator element"
	}
	p += recommend("yes")
	viper.Set(name, prompt.YesNo(p, viper.GetBool(name)))
	if err := UpdateConfig("", false); err != nil {
		logs.SaveFatal(err)
	}
}

// SetIndex prompts for a value from a list of valid choices and saves the result.
func setIndex(name string, setup bool, data ...string) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set index: %w", logs.ErrNameNil))
	}
	p := prompt.IndexStrings(&data, setup)
	save(name, setup, p)
}

// SetNoTranslate previews and prompts for the notranslate HTML attribute
// and Google meta elemenet.
func setNoTranslate(value, setup bool) {
	name := "html.meta.notranslate"
	elm := "<html translate=\"no\">\n  <head>\n    <meta name=\"google\" content=\"notranslate\">"
	fmt.Println(ColorHTML(elm))
	q := "Enable the no translate option"
	if value {
		q = "Keep the translate option"
	}
	q += recommend("no")
	val := prompt.YesNo(q, viper.GetBool(name))
	save(name, setup, val)
}

// SetPort prompts for and saves HTTP port.
func setPort(name string, setup bool) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set port: %w", logs.ErrNameNil))
	}
	val := prompt.Port(true, setup)
	save(name, setup, val)
}

// SetRetrotxt previews and prompts the custom retrotxt meta tag.
func setRetrotxt(value bool) {
	name := "html.meta.retrotxt"
	elm := "<head>\n  <meta name=\"retrotxt\" content=\"encoding: IBM437; linebreak: CRLF; length: 50; width: 80; name: file.txt\">"
	fmt.Println(ColorHTML(elm))
	p := "Enable the retrotxt element"
	if value {
		p = "Keep the retrotxt element"
	}
	p += recommend("yes")
	viper.Set(name, prompt.YesNo(p, viper.GetBool(name)))
	if err := UpdateConfig("", false); err != nil {
		logs.SaveFatal(err)
	}
}

// SetShortStrings prompts and saves setting values that support 1 character aliases.
func setShortStrings(name string, setup bool, data ...string) {
	val := prompt.ShortStrings(&data)
	switch val {
	case "-":
		val = ""
	case "":
		fmt.Println(skipped)
		return
	}
	save(name, setup, val)
}

// SetString prompts and saves a single word setting value.
func setString(name string, setup bool) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set string: %w", logs.ErrNameNil))
	}
	val := prompt.String(setup)
	switch val {
	case "-":
		val = ""
	case "":
		return
	}
	save(name, setup, val)
}

// SetStrings prompts and saves a string of text setting value.
func setStrings(name string, setup bool, data ...string) {
	if name == "" {
		logs.SaveFatal(fmt.Errorf("set strings: %w", logs.ErrNameNil))
	}
	val := prompt.Strings(&data, setup)
	switch val {
	case "-":
		val = ""
	case "":
		if !setup {
			fmt.Println(prompt.NoChange)
		}
		return
	}
	save(name, setup, val)
}
