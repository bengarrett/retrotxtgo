package input

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/alecthomas/chroma/styles"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/colorise"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

var ErrMeta = errors.New("cannot use name as a meta element")

type Update struct {
	Name  string
	Setup bool
	Value interface{}
}

// ColorScheme prompts the user for the color scheme setting.
func ColorScheme(w io.Writer, u Update) error {
	if err := PreviewMeta(w, u.Name, u.Value.(string)); err != nil {
		return err
	}
	c := create.ColorScheme()
	prints := make([]string, len(c[:]))
	copy(prints, c[:])
	fmt.Fprintf(w, "%s%s: ", str.UnderlineKeys(prints...), set.Recommend(""))
	return set.ShortStrings(w, u.Name, u.Setup, c[:]...)
}

func Defaults(name string) string {
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

// Editor prompts the user for the editor setting.
func Editor(w io.Writer, u Update) error {
	s := fmt.Sprint("  Set a " + get.Tip()[u.Name])
	if u.Value.(string) != "" {
		s = fmt.Sprint(s, " or use a dash [-] to remove")
	} else if ed := get.TextEditor(w); ed != "" {
		s = fmt.Sprintf("  Instead %s found %s and will use this editor.\n\n%s",
			meta.Name, str.ColPri(ed), s)
	}
	fmt.Fprintf(w, "%s:\n  ", s)
	return set.Editor(w, u.Name, u.Setup)
}

// Layout prompts the user for the layout setting.
func Layout(w io.Writer, u Update) error {
	fmt.Fprintf(w, "\n%s\n%s\n%s\n%s\n",
		"  Standard: Recommended, uses external CSS, JS and woff2 fonts and is the recommended layout for online hosting.",
		"  Inline:   Not recommended as it includes both the CSS and JS as inline elements that cannot be cached.",
		"  Compact:  The same as the standard layout but without any <meta> tags.",
		"  None:     No template is used and instead only the generated markup is returned.")
	fmt.Fprintf(w, "\n%s%s%s ",
		"  Choose a ", str.Options(get.Tip()[u.Name], true, false, create.Layouts()...),
		fmt.Sprintf(" (suggestion: %s):", str.Example("standard")))
	return set.ShortStrings(w, u.Name, u.Setup, create.Layouts()...)
}

// PortInfo returns recommended and valid HTTP port values.
func PortInfo() string {
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

// PreviewMeta previews and prompts for a meta element content value.
func PreviewMeta(w io.Writer, name, value string) error {
	if err := PrintMeta(w, name, value); err != nil {
		return err
	}
	fmt.Fprintf(w, "\n%s \n  ", PreviewPrompt(name, value))
	return nil
}

// PreviewPrompt returns the available options for the named setting.
func PreviewPrompt(name, value string) string {
	return fmt.Sprintf("%s:", PreviewPromptS(name, value))
}

func PreviewPromptS(name, value string) string {
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

func PrintMeta(w io.Writer, name, value string) error {
	if name == "" {
		return fmt.Errorf("preview meta: %w", logs.ErrNameNil)
	}
	if !set.Validate(name) {
		return fmt.Errorf("preview meta %q: %w", name, logs.ErrConfigName)
	}
	s := strings.Split(name, ".")
	const splits = 3
	switch {
	case len(s) != splits, s[0] != "html", s[1] != "meta":
		return fmt.Errorf("preview meta %q: %w", name, ErrMeta)
	}
	element := func() string {
		v := value
		if v == "" {
			v = Defaults(name)
		}
		return fmt.Sprintf("%s\n%s\n%s\n",
			"  <head>",
			fmt.Sprintf("    <meta name=\"%s\" value=\"%s\">", s[2], v),
			"  </head>")
	}
	if err := colorise.HTML(w, element()); err != nil {
		return err
	}
	h := strings.Split(get.Tip()[name], " ")
	a := fmt.Sprintf("%s\n  %s %s.",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"),
		strings.Title(h[0]), strings.Join(h[1:], " "))
	fmt.Fprint(w, a)
	return nil
}

// Serve prompts the user for a HTTP server port setting.
func Serve(w io.Writer, u Update) error {
	reset := func() error {
		var p uint
		if u, ok := get.Reset()[get.Serve].(uint); ok {
			p = u
		}
		if err := set.Write(w, u.Name, false, p); err != nil {
			return err
		}
		return nil
	}
	var p uint
	switch v := u.Value.(type) {
	case uint:
		p = v
	case int:
		p = uint(v)
	default:
		if err := reset(); err != nil {
			return err
		}
	}
	if p > prompt.PortMax {
		if err := reset(); err != nil {
			return err
		}
	}
	check := str.Bool(create.Port(p))
	fmt.Fprintf(w, "\n  Internal HTTP server port number: %s%d %s\n",
		str.ColSec("http://localhost:"), p, check)
	fmt.Fprintf(w, "\t\t\t\t    %s%d %s\n",
		str.ColSec("http://127.0.0.1:"), p, check)
	fmt.Fprintf(w, "\n  Port %s is reserved, port numbers less than %s are not recommended.\n",
		str.Example("0"), str.Example("1024"))
	fmt.Fprintf(w, "  Set a HTTP port number, choices %s: ", PortInfo())
	return set.Port(w, u.Name, u.Setup)
}

// SaveDir prompts the user for the a save destination directory setting.
func SaveDir(w io.Writer, u Update) error {
	fmt.Fprintf(w, "  Choose a new %s.\n\n  Directory aliases, use:", get.Tip()[u.Name])
	if home, err := os.UserHomeDir(); err == nil {
		fmt.Fprintf(w, "\n   %s (tilde) to save to your home directory: %s",
			str.Example("~"), str.Path(home))
	}
	if wd, err := os.Getwd(); err == nil {
		fmt.Fprintf(w, "\n   %s (period or full stop) to always save to this directory: %s",
			str.Example("."), str.Path(wd))
	}
	fmt.Fprintf(w, "\n   %s (hyphen-minus) to disable the setting and always use the active directory.\n  ",
		str.Example("-"))
	// this will loop for all errors (dir does not exist etc.)
	// but will break when an empty string [Enter key press] is returned
	for {
		if err := set.Directory(w, u.Name, u.Setup); errors.Is(err, set.ErrBreak) {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

// StyleHTML prompts the user for the a HTML and CSS style setting.
func StyleHTML(w io.Writer, u Update) error {
	d := ""
	if s, ok := get.Reset()[u.Name].(string); ok {
		d = s
	}
	italic := new(bytes.Buffer)
	colorise.ChromaNames(italic, "css")
	fmt.Fprintf(w, "\n%s\n\n  Choose the number to set a new HTML syntax style%s: ",
		str.Italic(italic.String()), set.Recommend(d))
	return set.Strings(w, u.Name, u.Setup, styles.Names()...)
}

// StyleInfo prompts the user for the a JS style setting.
func StyleInfo(w io.Writer, u Update) error {
	d := ""
	if s, ok := get.Reset()[u.Name].(string); ok {
		d = s
	}
	italic := new(bytes.Buffer)
	colorise.ChromaNames(italic, "json")
	fmt.Fprintf(w, "\n%s\n\n  Choose the number to set a new %s syntax style%s: ",
		str.Italic(italic.String()), str.Example("config info"), set.Recommend(d))
	return set.Strings(w, u.Name, u.Setup, styles.Names()...)
}
