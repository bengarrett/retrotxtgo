// nolint:goconst
package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/chroma/styles"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/upd"
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
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "css", style, true)
}

// ColorHTML returns the element colored using HTML syntax highlights.
func ColorHTML(elm string) string {
	style := viper.GetString(get.Styleh)
	return ColorElm(elm, "html", style, true)
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
	keys := set.Keys()
	const minWidth, tabWidth, tabs = 2, 2, "\t\t\t\t"
	w := tabwriter.NewWriter(os.Stdout, minWidth, tabWidth, 0, ' ', 0)
	cmds := fmt.Sprintf(" %s config set ", meta.Bin)
	title := fmt.Sprintf("  Available %s configurations and settings", meta.Name)
	fmt.Fprintln(w, "\n"+str.ColPri(title))
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintln(w, tabs)
	fmt.Fprintf(w, "Alias\t\tName\t\tHint\n")
	for i, key := range keys {
		tip := get.Tip()[key]
		fmt.Fprintln(w, tabs)
		fmt.Fprintf(w, " %d\t\t%s\t\t%s", i, key, suffix(capitalize(tip)))
		switch key {
		case get.LayoutTmpl:
			fmt.Fprintf(w, "\n%schoices: %s (suggestion: %s)",
				tabs, str.ColPri(strings.Join(create.Layouts(), ", ")), str.Example("standard"))
		case get.Serve:
			fmt.Fprintf(w, "\n%schoices: %s",
				tabs, portInfo())
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprintln(w, tabs)
	fmt.Fprintln(w, str.HR(len(title)))
	fmt.Fprintln(w, "\nEither the setting Name or the Alias can be used.")
	fmt.Fprintf(w, "\n%s # To change the meta description setting\n",
		str.Example(cmds+get.Desc))
	fmt.Fprintf(w, "%s # Will also change the meta description setting\n", str.Example(cmds+"6"))
	fmt.Fprintln(w, "\nMultiple settings are supported.")
	fmt.Fprintf(w, "\n%s\n", str.Example(cmds+"style.html style.info"))
	return w.Flush()
}

// ChromaNames returns the chroma style names.
func ChromaNames(lexer string) string {
	var s Names = styles.Names()
	return s.String(true, lexer)
}

// Set edits and saves a named setting within a configuration file.
// It also accepts numeric index values printed by List().
func Set(name string) {
	i, err := strconv.Atoi(name)
	switch {
	case err != nil:
		Update(name, false)
	case i >= 0 && i <= (len(get.Reset())-1):
		k := set.Keys()
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
		viper.Set(name, get.Reset()[name])
		value = viper.Get(name)
	default:
		// everything ok
	}
	if b, ok := value.(bool); ok {
		upd.Bool(b, name)
	}
	if s, ok := value.(string); ok {
		upd.String(s, name, value.(string))
	}
	updatePrompt(update{name, setup, value})
}

// Validate the existence of the key in a list of settings.
func Validate(key string) (ok bool) {
	keys := set.Keys()
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return false
	}
	return true
}

// Recommend uses the s value as a user input suggestion.
func Recommend(s string) string {
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
	case get.SaveDir:
		promptSaveDir(u)
	case get.Serve:
		promptServe(u)
	case get.Styleh:
		promptStyleHTML(u)
	case get.Stylei:
		promptStyleInfo(u)
	default:
		metaPrompts(u)
	}
}

// metaPrompts prompts the user for a meta setting.
func metaPrompts(u update) {
	switch u.name {
	case get.FontEmbed:
		set.FontEmbed(u.value.(bool), u.setup)
	case get.FontFamily:
		set.Font(u.value.(string), u.setup)
	case get.LayoutTmpl:
		promptLayout(u)
	case get.Author,
		get.Desc,
		get.Keywords:
		previewMeta(u.name, u.value.(string))
		set.String(u.name, u.setup)
	case get.Theme:
		recommendMeta(u.name, u.value.(string), "")
		set.String(u.name, u.setup)
	case get.Scheme:
		promptColorScheme(u)
	case get.Genr:
		set.Generator(u.value.(bool))
	case get.Notlate:
		set.NoTranslate(u.value.(bool), u.setup)
	case get.Referr:
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Referrer()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		set.Index(u.name, u.setup, cr[:]...)
	case get.Bot:
		recommendMeta(u.name, u.value.(string), "")
		cr := create.Robots()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		set.Index(u.name, u.setup, cr[:]...)
	case get.Rtx:
		set.RetroTxt(u.value.(bool))
	case get.Title:
		set.Title(u.name, u.value.(string), u.setup)
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
		str.UnderlineKeys(prints...), Recommend(""))
	set.ShortStrings(u.name, u.setup, c[:]...)
}

// promptEditor prompts the user for the editor setting.
func promptEditor(u update) {
	s := fmt.Sprint("  Set a " + get.Tip()[u.name])
	if u.value.(string) != "" {
		s = fmt.Sprint(s, " or use a dash [-] to remove")
	} else if ed := Editor(); ed != "" {
		s = fmt.Sprintf("  Instead %s found %s and will use this editor.\n\n%s",
			meta.Name, str.ColPri(ed), s)
	}
	fmt.Printf("%s:\n  ", s)
	set.Editor(u.name, u.setup)
}

// promptLayout prompts the user for the layout setting.
func promptLayout(u update) {
	fmt.Printf("\n%s\n%s\n%s\n%s\n",
		"  Standard: Recommended, uses external CSS, JS and woff2 fonts and is the recommended layout for online hosting.",
		"  Inline:   Not recommended as it includes both the CSS and JS as inline elements that cannot be cached.",
		"  Compact:  The same as the standard layout but without any <meta> tags.",
		"  None:     No template is used and instead only the generated markup is returned.")
	fmt.Printf("\n%s%s%s ",
		"  Choose a ", str.Options(get.Tip()[u.name], true, false, create.Layouts()...),
		fmt.Sprintf(" (suggestion: %s):", str.Example("standard")))
	set.ShortStrings(u.name, u.setup, create.Layouts()...)
}

// promptSaveDir prompts the user for the a save destination directory setting.
func promptSaveDir(u update) {
	fmt.Printf("  Choose a new %s.\n\n  Directory aliases, use:", get.Tip()[u.name])
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
		if set.Directory(u.name, u.setup) {
			break
		}
	}
}

// promptServe prompts the user for a HTTP server port setting.
func promptServe(u update) {
	reset := func() {
		var p uint
		if u, ok := get.Reset()[get.Serve].(uint); ok {
			p = u
		}
		set.Write(u.name, false, p)
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
	set.Port(u.name, u.setup)
}

// promptStyleHTML prompts the user for the a HTML and CSS style setting.
func promptStyleHTML(u update) {
	d := ""
	if s, ok := get.Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new HTML syntax style%s: ",
		str.Italic(ChromaNames("css")), Recommend(d))
	set.Strings(u.name, u.setup, styles.Names()...)
}

// promptStyleInfo prompts the user for the a JS style setting.
func promptStyleInfo(u update) {
	d := ""
	if s, ok := get.Reset()[u.name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new %s syntax style%s: ",
		str.Italic(ChromaNames("json")), str.Example("config info"), Recommend(d))
	set.Strings(u.name, u.setup, styles.Names()...)
}

// ColorElm applies color syntax to an element.
func ColorElm(elm, lexer, style string, color bool) string {
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

type Names []string

// String lists and applies the named themes for the HighlightWriter.
func (n Names) String(theme bool, lexer string) string {
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
	fmt.Printf("\n%s \n  ", PreviewPrompt(name, value))
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
	h := strings.Split(get.Tip()[name], " ")
	fmt.Printf("%s\n  %s %s.",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"),
		strings.Title(h[0]), strings.Join(h[1:], " "))
}

func metaDefaults(name string) string {
	switch name {
	case get.Author:
		return "Your name goes here"
	case get.Scheme:
		return "normal"
	case get.Desc:
		return "A brief description of the page could go here."
	case get.Keywords:
		return "some, keywords, go here"
	case get.Referr:
		return "same-origin"
	case get.Bot:
		return "noindex"
	case get.Theme:
		return "ghostwhite"
	case get.Title:
		return "A page title foes here."
	}
	return ""
}

// PreviewPrompt returns the available options for the named setting.
func PreviewPrompt(name, value string) string {
	return fmt.Sprintf("%s:", previewPromptPrint(name, value))
}

func previewPromptPrint(name, value string) string {
	p := "Set a new value"
	if name == get.Keywords {
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
	return fmt.Sprintf("%s%s:", p, Recommend(suggest))
}
