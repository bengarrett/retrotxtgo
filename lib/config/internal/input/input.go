package input

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/color"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
)

type Update struct {
	Name  string
	Setup bool
	Value interface{}
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

// PreviewMeta previews and prompts for a meta element content value.
func PreviewMeta(name, value string) {
	PrintMeta(name, value)
	fmt.Printf("\n%s \n  ", PreviewPrompt(name, value))
}

func PrintMeta(name, value string) {
	if name == "" {
		logs.FatalSave(fmt.Errorf("preview meta: %w", logs.ErrNameNil))
	}
	if !set.Validate(name) {
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
			v = Defaults(name)
		}
		return fmt.Sprintf("%s\n%s\n%s\n",
			"  <head>",
			fmt.Sprintf("    <meta name=\"%s\" value=\"%s\">", s[2], v),
			"  </head>")
	}
	fmt.Print(color.ColorHTML(element()))
	h := strings.Split(get.Tip()[name], " ")
	fmt.Printf("%s\n  %s %s.",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name"),
		strings.Title(h[0]), strings.Join(h[1:], " "))
}

// ColorScheme prompts the user for the color scheme setting.
func ColorScheme(u Update) {
	PreviewMeta(u.Name, u.Value.(string))
	c := create.ColorScheme()
	prints := make([]string, len(c[:]))
	copy(prints, c[:])
	fmt.Printf("%s%s: ",
		str.UnderlineKeys(prints...), set.Recommend(""))
	set.ShortStrings(u.Name, u.Setup, c[:]...)
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

// Layout prompts the user for the layout setting.
func Layout(u Update) {
	fmt.Printf("\n%s\n%s\n%s\n%s\n",
		"  Standard: Recommended, uses external CSS, JS and woff2 fonts and is the recommended layout for online hosting.",
		"  Inline:   Not recommended as it includes both the CSS and JS as inline elements that cannot be cached.",
		"  Compact:  The same as the standard layout but without any <meta> tags.",
		"  None:     No template is used and instead only the generated markup is returned.")
	fmt.Printf("\n%s%s%s ",
		"  Choose a ", str.Options(get.Tip()[u.Name], true, false, create.Layouts()...),
		fmt.Sprintf(" (suggestion: %s):", str.Example("standard")))
	set.ShortStrings(u.Name, u.Setup, create.Layouts()...)
}

// Editor prompts the user for the editor setting.
func Editor(u Update) {
	s := fmt.Sprint("  Set a " + get.Tip()[u.Name])
	if u.Value.(string) != "" {
		s = fmt.Sprint(s, " or use a dash [-] to remove")
	} else if ed := get.TextEditor(); ed != "" {
		s = fmt.Sprintf("  Instead %s found %s and will use this editor.\n\n%s",
			meta.Name, str.ColPri(ed), s)
	}
	fmt.Printf("%s:\n  ", s)
	set.Editor(u.Name, u.Setup)
}

// Serve prompts the user for a HTTP server port setting.
func Serve(u Update) {
	reset := func() {
		var p uint
		if u, ok := get.Reset()[get.Serve].(uint); ok {
			p = u
		}
		set.Write(u.Name, false, p)
	}
	var p uint
	switch v := u.Value.(type) {
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
	fmt.Printf("  Set a HTTP port number, choices %s: ", PortInfo())
	set.Port(u.Name, u.Setup)
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

// SaveDir prompts the user for the a save destination directory setting.
func SaveDir(u Update) {
	fmt.Printf("  Choose a new %s.\n\n  Directory aliases, use:", get.Tip()[u.Name])
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
		if set.Directory(u.Name, u.Setup) {
			break
		}
	}
}
