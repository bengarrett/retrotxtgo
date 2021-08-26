// nolint:goconst
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
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

type update struct {
	name  string
	setup bool
	value interface{}
}

// ColorCSS returns the element colored using CSS syntax highlights.
func ColorCSS(elm string) string {
	style := viper.GetString(styleh)
	return colorElm(elm, "css", style, true)
}

// ColorHTML returns the element colored using HTML syntax highlights.
func ColorHTML(elm string) string {
	style := viper.GetString(styleh)
	return colorElm(elm, "html", style, true)
}

// List and print all the available configurations.
func List() error {
	capitalize := func(s string) string {
		return strings.Title(s[:1]) + s[1:]
	}
	suffix := func(s string) string {
		if strings.HasSuffix(s, "?") {
			return s
		}
		return fmt.Sprintf("%s.", s)
	}
	keys := Keys()
	const minWidth, tabWidth, tabs = 2, 2, "\t\t\t\t"
	w := tabwriter.NewWriter(os.Stdout, minWidth, tabWidth, 0, ' ', 0)
	cmds := fmt.Sprintf(" %s config set ", meta.Bin)
	title := fmt.Sprintf("  Available %s configurations and settings", meta.Name)
	fmt.Fprintln(w, "\n"+str.ColPri(title))
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintln(w, tabs)
	fmt.Fprintf(w, "Alias\t\tName\t\tHint\n")
	for i, key := range keys {
		tip := Tip()[key]
		fmt.Fprintln(w, tabs)
		fmt.Fprintf(w, " %d\t\t%s\t\t%s", i, key, suffix(capitalize(tip)))
		switch key {
		case layoutTmpl:
			fmt.Fprintf(w, "\n%schoices: %s (suggestion: %s)",
				tabs, str.ColPri(strings.Join(create.Layouts(), ", ")), str.Example("standard"))
		case serve:
			fmt.Fprintf(w, "\n%schoices: %s",
				tabs, portInfo())
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprintln(w, tabs)
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintln(w, "\nEither the setting Name or the Alias can be used.")
	fmt.Fprintf(w, "\n%s # To change the meta description setting\n",
		str.Example(cmds+desc))
	fmt.Fprintf(w, "%s # Will also change the meta description setting\n", str.Example(cmds+"6"))
	fmt.Fprintln(w, "\nMultiple settings are supported.")
	fmt.Fprintf(w, "\n%s\n", str.Example(cmds+"style.html style.info"))
	return w.Flush()
}

// Names returns the chroma style names.
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
		fmt.Println(logs.Hint("config set --list", logs.ErrConfigName))
		return
	}
	if !setup {
		fmt.Print(Location())
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

// Validate the existence of the key in a list of settings.
func Validate(key string) (ok bool) {
	keys := Keys()
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return false
	}
	return true
}

// updateBool the boolean value of the named setting.
func updateBool(b bool, name string) {
	switch b {
	case true:
		fmt.Printf("\n  The %s is enabled.\n", str.ColFuz(name))
	default:
		fmt.Printf("\n  The %s is not in use.\n", str.ColFuz(name))
	}
}

// updateString the string value of the named setting.
func updateString(s, name, value string) {
	const sd = saveDir
	switch s {
	case "":
		fmt.Printf("\n  The empty %s setting is not in use.\n",
			str.ColFuz(name))
		if name == sd {
			fmt.Printf("  Files created by %s will always be saved to the active directory.\n",
				meta.Name)
		}
	default:
		fmt.Printf("\n  The %s is set to %q.", str.ColFuz(name), value)
		// print the operating system's ability to use the existing set values
		// does the 'editor' exist in the env path, does the save-directory exist?
		switch name {
		case "editor":
			_, err := exec.LookPath(value)
			fmt.Print(" ", str.Bool(err == nil))
		case sd:
			f := false
			if _, err := os.Stat(value); !os.IsNotExist(err) {
				f = true
			}
			fmt.Print(" ", str.Bool(f))
		default:
		}
		fmt.Println()
	}
}

// recommend uses the s value as a user input suggestion.
func recommend(s string) string {
	if s == "" {
		return fmt.Sprintf(" (suggestion: %s)", str.Example("do not use"))
	}
	return fmt.Sprintf(" (suggestion: %s)", str.Example(s))
}

// updatePrompt prompts the user for input to a config file setting.
func updatePrompt(u update) {
	switch u.name {
	case "editor":
		promptEditor(u)
	case saveDir:
		promptSaveDir(u)
	case serve:
		promptServe(u)
	case styleh:
		promptStyleHTML(u)
	case stylei:
		promptStyleInfo(u)
	default:
		metaPrompts(u)
	}
}

// metaPrompts prompts the user for a meta setting.
func metaPrompts(u update) {
	switch u.name {
	case fontEmbed:
		setFontEmbed(u.value.(bool), u.setup)
	case fontFamily:
		setFont(u.value.(string), u.setup)
	case layoutTmpl:
		promptLayout(u)
	case author,
		desc,
		keywords:
		previewMeta(u.name, u.value.(string))
		setString(u.name, u.setup)
	case theme:
		recommendMeta(u.name, u.value.(string), "")
		setString(u.name, u.setup)
	case scheme:
		promptColorScheme(u)
	case genr:
		setGenerator(u.value.(bool))
	case notlate:
		setNoTranslate(u.value.(bool), u.setup)
	case referr:
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Referrer()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		setIndex(u.name, u.setup, cr[:]...)
	case bot:
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Robots()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		setIndex(u.name, u.setup, cr[:]...)
	case rtx:
		setRetroTxt(u.value.(bool))
	case title:
		setTitle(u.name, u.value.(string), u.setup)
	default:
		log.Fatalln("config is not configured:", u.name)
	}
}

// promptColorScheme prompts the user for the color scheme setting.
func promptColorScheme(u update) {
	previewMeta(u.name, u.value.(string))
	c := create.ColorScheme()
	prints := make([]string, len(c[:]))
	copy(prints, c[:])
	fmt.Printf("%s%s: ",
		str.UnderlineKeys(prints...), recommend(""))
	setShortStrings(u.name, u.setup, c[:]...)
}

// promptEditor prompts the user for the editor setting.
func promptEditor(u update) {
	s := fmt.Sprint("  Set a " + Tip()[u.name])
	if u.value.(string) != "" {
		s = fmt.Sprint(s, " or use a dash [-] to remove")
	} else if ed := Editor(); ed != "" {
		s = fmt.Sprintf("  Instead %s found %s and will use this editor.\n\n%s",
			meta.Name, str.ColPri(ed), s)
	}
	fmt.Printf("%s:\n  ", s)
	setEditor(u.name, u.setup)
}

// promptLayout prompts the user for the layout setting.
func promptLayout(u update) {
	fmt.Printf("\n%s\n%s\n%s\n%s\n",
		"  Standard: Recommended, uses external CSS, JS and woff2 fonts and is the recommended layout for online hosting.",
		"  Inline:   Not recommended as it includes both the CSS and JS as inline elements that cannot be cached.",
		"  Compact:  The same as the standard layout but without any <meta> tags.",
		"  None:     No template is used and instead only the generated markup is returned.")
	fmt.Printf("\n%s%s%s ",
		"  Choose a ", str.Options(Tip()[u.name], true, false, create.Layouts()...),
		fmt.Sprintf(" (suggestion: %s):", str.Example("standard")))
	setShortStrings(u.name, u.setup, create.Layouts()...)
}

// promptSaveDir prompts the user for the a save destination directory setting.
func promptSaveDir(u update) {
	fmt.Printf("  Choose a new %s.\n\n  Directory aliases, use:", Tip()[u.name])
	if home, err := os.UserHomeDir(); err == nil {
		fmt.Printf("\n   %s (tilde) to save to your home directory: %s",
			str.Example("~"), str.Path(home))
	}
	if wd, err := os.Getwd(); err == nil {
		fmt.Printf("\n   %s (period or full stop) to always save to this directory: %s",
			str.Example("."), str.Path(wd))
	}
	fmt.Printf("\n   %s (hyphen-minus) to disable the setting and always use the active directory.\n  ",
		str.Example("-"))
	for {
		if setDirectory(u.name, u.setup) {
			break
		}
	}
}

// promptServe prompts the user for a HTTP server port setting.
func promptServe(u update) {
	reset := func() {
		var p uint
		if u, ok := Reset()[serve].(uint); ok {
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
	check := str.Bool(create.Port(p))
	fmt.Printf("\n  Internal HTTP server port number: %s%d %s\n",
		str.ColSec("http://localhost:"), p, check)
	fmt.Printf("\t\t\t\t    %s%d %s\n",
		str.ColSec("http://127.0.0.1:"), p, check)
	fmt.Printf("\n  Port %s is reserved, port numbers less than %s are not recommended.\n",
		str.Example("0"), str.Example("1024"))
	fmt.Printf("  Set a HTTP port number, choices %s: ", portInfo())
	setPort(u.name, u.setup)
}

// promptStyleHTML prompts the user for the a HTML and CSS style setting.
func promptStyleHTML(u update) {
	d := ""
	if s, ok := Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new HTML syntax style%s: ",
		str.Italic(Names("css")), recommend(d))
	setStrings(u.name, u.setup, styles.Names()...)
}

// promptStyleInfo prompts the user for the a JS style setting.
func promptStyleInfo(u update) {
	d := ""
	if s, ok := Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new %s syntax style%s: ",
		str.Italic(Names("json")), str.Example("config info"), recommend(d))
	setStrings(u.name, u.setup, styles.Names()...)
}

// colorElm applies color syntax to an element.
func colorElm(elm, lexer, style string, color bool) string {
	if elm == "" {
		return ""
	}
	var b bytes.Buffer
	_ = io.Writer(&b)
	if err := str.HighlightWriter(&b, elm, lexer, style, color); err != nil {
		logs.FatalMark(fmt.Sprint("html ", lexer), logs.ErrHighlight, err)
	}
	return fmt.Sprintf("\n%s\n", b.String())
}

type names []string

// string lists and applies the named themes for the HighlightWriter.
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
	s := make([]string, len(n))
	split := (len(n) / 2)
	const space = 2
	for i, name := range n {
		b, t := bytes.Buffer{}, ""
		pad := maxWidth - len(fmt.Sprintf("%s=%q", name, name))
		// prints a sequential list of styles
		if lexer == "json" {
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }%s", name, name, strings.Repeat(" ", pad+space))
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.FatalMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
			if split+i >= len(n) {
				break
			}
			b = bytes.Buffer{}
			t = fmt.Sprintf("{ %q:%q }\n", n[split+i], n[split+i])
			if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
				logs.FatalMark(name, logs.ErrHighlight, err)
			}
			s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
			continue
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>%s", name, name, strings.Repeat(" ", pad+space))
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", i, b.String()))
		if split+i >= len(n) {
			break
		}
		b = bytes.Buffer{}
		t = fmt.Sprintf("<%s=%q>\n", n[split+i], n[split+i])
		if err := str.HighlightWriter(&b, t, lexer, name, true); err != nil {
			logs.FatalMark(name, logs.ErrHighlight, err)
		}
		s = append(s, fmt.Sprintf("%2d %s", split+i, b.String()))
	}
	return strings.Join(s, "")
}

// dirExpansion traverses the named directory to apply shell-like expansions.
// It supports limited Bash tilde, shell dot and double dot syntax.
func dirExpansion(name string) (dir string) {
	const sep, homeDir, currentDir, parentDir = string(os.PathSeparator), "~", ".", ".."
	if name == "" || name == sep {
		return name
	}
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	r, paths := bool(name[0:1] == sep), strings.Split(name, sep)
	var err error
	for i, s := range paths {
		p := ""
		switch s {
		case homeDir:
			p, err = os.UserHomeDir()
			if err != nil {
				logs.FatalSave(err)
			}
		case currentDir:
			if i != 0 {
				continue
			}
			p, err = os.Getwd()
			if err != nil {
				logs.FatalSave(err)
			}
		case parentDir:
			if i != 0 {
				dir = filepath.Dir(dir)
				continue
			}
			wd, err := os.Getwd()
			if err != nil {
				logs.FatalSave(err)
			}
			p = filepath.Dir(wd)
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

// portInfo returns recommended and valid HTTP port values.
func portInfo() string {
	type ports struct {
		max uint
		min uint
		rec uint
	}
	port := ports{
		max: prompt.PortMax,
		min: prompt.PortMin,
		rec: meta.WebPort,
	}
	pm, px, pr :=
		strconv.Itoa(int(port.min)),
		strconv.Itoa(int(port.max)),
		strconv.Itoa(int(port.rec))
	return fmt.Sprintf("%s-%s (suggestion: %s)",
		str.Example(pm), str.Example(px), str.Example(pr))
}

// previewMeta previews and prompts for a meta element content value.
func previewMeta(name, value string) {
	previewMetaPrint(name, value)
	fmt.Printf("\n%s \n  ", previewPrompt(name, value))
}

func previewMetaPrint(name, value string) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("preview meta: %w", logs.ErrNameNil))
	}
	if !Validate(name) {
		logs.FatalSave(fmt.Errorf("preview meta %q: %w", name, logs.ErrConfigName))
	}
	s := strings.Split(name, ".")
	const splits = 3
	switch {
	case len(s) != splits, s[0] != "html", s[1] != "meta":
		return
	}
	element := func() string {
		v := value
		if v == "" {
			v = metaDefaults(name)
		}
		return fmt.Sprintf("%s\n%s\n%s\n",
			"  <head>",
			fmt.Sprintf("    <meta name=\"%s\" value=\"%s\">", s[2], v),
			"  </head>")
	}
	fmt.Print(ColorHTML(element()))
	h := strings.Split(Tip()[name], " ")
	fmt.Printf("%s\n  %s %s.",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"),
		strings.Title(h[0]), strings.Join(h[1:], " "))
}

func metaDefaults(name string) string {
	switch name {
	case author:
		return "Your name goes here"
	case scheme:
		return "normal"
	case desc:
		return "A brief description of the page could go here."
	case keywords:
		return "some, keywords, go here"
	case referr:
		return "same-origin"
	case bot:
		return "noindex"
	case theme:
		return "ghostwhite"
	case title:
		return "A page title foes here."
	}
	return ""
}

// setTitle prompts for and previews a HTML title element value.
func setTitle(name, value string, setup bool) {
	elm := fmt.Sprintf("%s\n%s\n%s\n",
		"  <head>",
		fmt.Sprintf("    <title>%s</title>", value),
		"  </head>")
	fmt.Print(ColorHTML(elm))
	fmt.Printf("%s\n%s\n  ",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"),
		fmt.Sprintf("  Choose a new %s:", Tip()[name]))
	setString(name, setup)
}

// previewPrompt returns the available options for the named setting.
func previewPrompt(name, value string) string {
	return fmt.Sprintf("%s:", previewPromptPrint(name, value))
}

func previewPromptPrint(name, value string) string {
	p := "Set a new value"
	if name == keywords {
		p = "Set some comma-separated keywords"
		if value != "" {
			p = "Replace the current keywords"
		}
	}
	if value != "" {
		return fmt.Sprintf("  %s, leave blank to keep as-is or use a dash [-] to remove", p)
	}
	return fmt.Sprintf("  %s or leave blank to keep it unused", p)
}

func recommendMeta(name, value, suggest string) {
	previewMetaPrint(name, value)
	fmt.Printf("\n%s\n  ", recommendPrompt(name, value, suggest))
}

func recommendPrompt(name, value, suggest string) string {
	p := previewPromptPrint(name, value)
	return fmt.Sprintf("%s%s:", p, recommend(suggest))
}

// save the value of the named setting to the configuration file.
func save(name string, setup bool, value interface{}) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("save: %w", logs.ErrNameNil))
	}
	if !Validate(name) {
		logs.FatalSave(fmt.Errorf("save %q: %w", name, logs.ErrConfigName))
	}
	if skipSave(name, value) {
		fmt.Print(skipSet(setup))
		return
	}
	switch v := value.(type) {
	case string:
		if v == "-" {
			value = ""
		}
	default:
	}
	viper.Set(name, value)
	if err := UpdateConfig("", false); err != nil {
		logs.FatalSave(err)
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			fmt.Printf("  %s is now unused\n",
				str.ColSuc(name))
			if !setup {
				os.Exit(0)
			}
			return
		}
	default:
	}
	fmt.Printf("  %s is set to \"%v\"\n",
		str.ColSuc(name), value)
	if !setup {
		os.Exit(0)
	}
}

// skipSave returns true if the named value doesn't need updating.
func skipSave(name string, value interface{}) bool {
	switch v := value.(type) {
	case bool:
		if viper.Get(name).(bool) == v {
			return true
		}
	case string:
		if viper.Get(name).(string) == v {
			return true
		}
		if value.(string) == "" {
			return true
		}
	case uint:
		if viper.Get(name).(int) == int(v) {
			return true
		}
		if name == serve && v == 0 {
			return true
		}
	default:
		logs.FatalSave(fmt.Errorf("name: %s, type: %T, %w", name, value, ErrSaveType))
	}
	return false
}

func skipSet(setup bool) string {
	if !setup {
		return ""
	}
	return str.ColSuc("\r  skipped setting")
}

// SetDirectory prompts for, checks and saves the directory path.
func setDirectory(name string, setup bool) (ok bool) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set directory: %w", logs.ErrNameNil))
	}
	s := dirExpansion(prompt.String())
	if s == "" {
		fmt.Print(skipSet(true))
		if setup {
			return true
		}
		os.Exit(0)
	}
	if s == "-" {
		save(name, setup, "-")
		return true
	}
	if _, err := os.Stat(s); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%s The directory does not exist: %s\n", str.Alert(), s)
			return false
		}
		fmt.Printf("%s the directory is invalid: %s\n", str.Alert(), errors.Unwrap(err))
		return false
	}
	save(name, setup, s)
	return true
}

// setEditor checks the existence of given text editor location
// and saves it as a configuration regardless of the result.
func setEditor(name string, setup bool) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set editor: %w", logs.ErrNameNil))
	}
	s := prompt.String()
	switch s {
	case "-":
		save(name, setup, "-")
		return
	case "":
		fmt.Print(skipSet(setup))
		return
	}
	if _, err := exec.LookPath(s); err != nil {
		fmt.Printf("%s%s\nThe %s editor is not usable by %s.\n",
			str.Alert(), errors.Unwrap(err), s, meta.Name)
	}
	save(name, setup, s)
}

// SetFont previews and saves a default font setting.
func setFont(value string, setup bool) {
	b, f := bytes.Buffer{}, create.Family(value)
	if f == create.Automatic {
		f = create.VGA
	}
	fmt.Fprintf(&b, "%s\n%s\n%s\n%s\n",
		"  @font-face {",
		fmt.Sprintf("    font-family: \"%s\";", f),
		fmt.Sprintf("    src: url(\"%s.woff2\") format(\"woff2\");", f),
		"  }")
	fmt.Print(ColorCSS(b.String()))
	fmt.Printf("%s\n%s%s %s: ",
		str.ColFuz("  About font families: https://developer.mozilla.org/en-US/docs/Web/CSS/font-family"),
		"  Choose a font, ",
		str.UnderlineKeys(create.Fonts()...),
		fmt.Sprintf("(suggestion: %s)", str.Example("automatic")))
	setShortStrings(fontFamily, setup, create.Fonts()...)
}

// setFont previews and saves the embedded Base64 font setting.
func setFontEmbed(value, setup bool) {
	const name = fontEmbed
	elm := fmt.Sprintf("  %s\n  %s\n  %s\n",
		"@font-face{",
		"  font-family: vga8;",
		"  src: url(data:font/woff2;base64,[a large font binary will be embedded here]...) format('woff2');",
	)
	fmt.Print(ColorCSS(elm))
	q := fmt.Sprintf("%s\n%s\n%s",
		"  The use of this setting not recommended,",
		"  unless you always want large, self-contained HTML files for distribution.",
		"  Embed the font as Base64 text within the HTML")
	if value {
		q = "  Keep the embedded font option?"
	}
	q += recommend("no")
	b := prompt.YesNo(q, viper.GetBool(name))
	save(name, setup, b)
}

// SetGenerator prompts for and previews the custom program generator meta tag.
func setGenerator(value bool) {
	const name = genr
	elm := fmt.Sprintf("  %s\n    %s\n  %s\n",
		"<head>",
		fmt.Sprintf("<meta name=\"generator\" content=\"%s %s, %s\">",
			meta.Name, meta.Print(), meta.App.Date),
		"</head>")
	fmt.Print(ColorHTML(elm))
	p := "Enable the generator element?"
	if value {
		p = "Keep the generator element?"
	}
	p = fmt.Sprintf("  %s%s", p, recommend("yes"))
	b := prompt.YesNo(p, viper.GetBool(name))
	save(name, true, b)
}

// setIndex prompts for a value from a list of valid choices and saves the result.
func setIndex(name string, setup bool, data ...string) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set index: %w", logs.ErrNameNil))
	}
	s := prompt.IndexStrings(&data, setup)
	save(name, setup, s)
}

// setNoTranslate previews and prompts for the notranslate HTML attribute
// and Google meta elemenet.
func setNoTranslate(value, setup bool) {
	const name = notlate
	elm := fmt.Sprintf("  %s\n    %s\n  %s\n  %s\n",
		"<head>",
		"<meta name=\"google\" content=\"notranslate\">",
		"</head>",
		"<body class=\"notranslate\">")
	fmt.Print(ColorHTML(elm))
	q := "Enable the no translate option?"
	if value {
		q = "Keep the translate option?"
	}
	q = fmt.Sprintf("  %s%s", q, recommend("no"))
	b := prompt.YesNo(q, viper.GetBool(name))
	save(name, setup, b)
}

// SetPort prompts for and saves HTTP port.
func setPort(name string, setup bool) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set port: %w", logs.ErrNameNil))
	}
	u := prompt.Port(true, setup)
	save(name, setup, u)
}

// setRetroTxt previews and prompts the custom retrotxt meta tag.
func setRetroTxt(value bool) {
	const name = rtx
	elm := fmt.Sprintf("%s\n%s\n%s\n",
		"  <head>",
		"    <meta name=\"retrotxt\" content=\"encoding: IBM437; linebreak: CRLF; length: 50; width: 80; name: file.txt\">",
		"  </head>")
	fmt.Print(ColorHTML(elm))
	p := "Enable the retrotxt element?"
	if value {
		p = "Keep the retrotxt element?"
	}
	p = fmt.Sprintf("  %s%s", p, recommend("yes"))
	b := prompt.YesNo(p, viper.GetBool(name))
	save(name, true, b)
}

// setShortStrings prompts and saves setting values that support 1 character aliases.
func setShortStrings(name string, setup bool, data ...string) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set short string: %w", logs.ErrNameNil))
	}
	s := prompt.ShortStrings(&data)
	save(name, setup, s)
}

// setString prompts and saves a single word setting value.
func setString(name string, setup bool) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set string: %w", logs.ErrNameNil))
	}
	s := prompt.String()
	save(name, setup, s)
}

// setStrings prompts and saves a string of text setting value.
func setStrings(name string, setup bool, data ...string) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("set strings: %w", logs.ErrNameNil))
	}
	s := prompt.Strings(&data, setup)
	save(name, setup, s)
}
