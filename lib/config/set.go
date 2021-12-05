// nolint:goconst
package config

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/alecthomas/chroma/styles"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/upd"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

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
				tabs, input.PortInfo())
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
	if !set.Validate(name) {
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
	updatePrompt(input.Update{name, setup, value})
}

// updatePrompt prompts the user for input to a config file setting.
func updatePrompt(u input.Update) {
	switch u.Name {
	case "editor":
		input.Editor(u)
	case get.SaveDir:
		input.SaveDir(u)
	case get.Serve:
		input.Serve(u)
	case get.Styleh:
		promptStyleHTML(u)
	case get.Stylei:
		promptStyleInfo(u)
	default:
		metaPrompts(u)
	}
}

// metaPrompts prompts the user for a meta setting.
func metaPrompts(u input.Update) {
	switch u.Name {
	case get.FontEmbed:
		set.FontEmbed(u.Value.(bool), u.Setup)
	case get.FontFamily:
		set.Font(u.Value.(string), u.Setup)
	case get.LayoutTmpl:
		input.Layout(u)
	case get.Author,
		get.Desc,
		get.Keywords:
		input.PreviewMeta(u.Name, u.Value.(string))
		set.String(u.Name, u.Setup)
	case get.Theme:
		recommendMeta(u.Name, u.Value.(string), "")
		set.String(u.Name, u.Setup)
	case get.Scheme:
		input.ColorScheme(u)
	case get.Genr:
		set.Generator(u.Value.(bool))
	case get.Notlate:
		set.NoTranslate(u.Value.(bool), u.Setup)
	case get.Referr:
		recommendMeta(u.Name, u.Value.(string), "")
		cr := create.Referrer()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		set.Index(u.Name, u.Setup, cr[:]...)
	case get.Bot:
		recommendMeta(u.Name, u.Value.(string), "")
		cr := create.Robots()
		fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
		set.Index(u.Name, u.Setup, cr[:]...)
	case get.Rtx:
		set.RetroTxt(u.Value.(bool))
	case get.Title:
		set.Title(u.Name, u.Value.(string), u.Setup)
	default:
		log.Fatalln("config is not configured:", u.Name)
	}
}

// promptStyleHTML prompts the user for the a HTML and CSS style setting.
func promptStyleHTML(u input.Update) {
	d := ""
	if s, ok := get.Reset()[u.Name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new HTML syntax style%s: ",
		str.Italic(ChromaNames("css")), set.Recommend(d))
	set.Strings(u.Name, u.Setup, styles.Names()...)
}

// promptStyleInfo prompts the user for the a JS style setting.
func promptStyleInfo(u input.Update) {
	d := ""
	if s, ok := get.Reset()[u.Name].(string); ok {
		d = s
	}
	fmt.Printf("\n%s\n\n  Choose the number to set a new %s syntax style%s: ",
		str.Italic(ChromaNames("json")), str.Example("config info"), set.Recommend(d))
	set.Strings(u.Name, u.Setup, styles.Names()...)
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

func recommendMeta(name, value, suggest string) {
	input.PrintMeta(name, value)
	fmt.Printf("\n%s\n  ", recommendPrompt(name, value, suggest))
}

func recommendPrompt(name, value, suggest string) string {
	p := input.PreviewPromptS(name, value)
	return fmt.Sprintf("%s%s:", p, set.Recommend(suggest))
}
