package config

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/update"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

var ErrNotCfged = errors.New("config is not configured")

// List and print all the available configurations.
func List() (*bytes.Buffer, error) {
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
	b := new(bytes.Buffer)
	w := tabwriter.NewWriter(b, minWidth, tabWidth, 0, ' ', 0)
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
	return b, nil
}

// Set edits and saves a named setting within a configuration file.
// It also accepts numeric index values printed by List().
func Set(name string) error {
	i, err := strconv.Atoi(name)
	namedSetting := err != nil
	switch {
	case namedSetting:
		return Update(name, false)
	case i >= 0 && i <= (len(get.Reset())-1):
		k := set.Keys()
		return Update(k[i], false)
	default:
		return Update(name, false)
	}
}

// Update edits and saves a named setting within a configuration file.
func Update(name string, setup bool) error {
	if !set.Validate(name) {
		fmt.Println(logs.Hint("config set --list", logs.ErrConfigName))
		return nil
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
		update.Bool(b, name)
	}
	if s, ok := value.(string); ok {
		update.String(s, name, value.(string))
	}
	if err := updatePrompt(input.Update{
		Name:  name,
		Setup: setup,
		Value: value,
	}); err != nil {
		return err
	}
	return nil
}

// updatePrompt prompts the user for input to a config file setting.
func updatePrompt(u input.Update) error {
	switch u.Name {
	case "editor":
		input.Editor(u)
	case get.SaveDir:
		input.SaveDir(u)
	case get.Serve:
		input.Serve(u)
	case get.Styleh:
		input.StyleHTML(u)
	case get.Stylei:
		input.StyleInfo(u)
	}
	return metaPrompts(u)
}

// metaPrompts prompts the user for a meta setting.
func metaPrompts(u input.Update) error {
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
		return theme(u)
	case get.Scheme:
		input.ColorScheme(u)
	case get.Genr:
		set.Generator(u.Value.(bool))
	case get.Notlate:
		set.NoTranslate(u.Value.(bool), u.Setup)
	case get.Referr:
		return referr(u)
	case get.Bot:
		return bot(u)
	case get.Rtx:
		set.RetroTxt(u.Value.(bool))
	case get.Title:
		set.Title(u.Name, u.Value.(string), u.Setup)
	}
	return fmt.Errorf("%w: %s", ErrNotCfged, u.Name)
}

func theme(u input.Update) error {
	if err := recommendMeta(u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	set.String(u.Name, u.Setup)
	return nil
}

func bot(u input.Update) error {
	if err := recommendMeta(u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	cr := create.Robots()
	fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
	set.Index(u.Name, u.Setup, cr[:]...)
	return nil
}

func referr(u input.Update) error {
	if err := recommendMeta(u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	cr := create.Referrer()
	fmt.Printf("%s\n  ", str.NumberizeKeys(cr[:]...))
	set.Index(u.Name, u.Setup, cr[:]...)
	return nil
}

func recommendMeta(name, value, suggest string) error {
	s, err := input.PrintMeta(name, value)
	if err != nil {
		return fmt.Errorf("recommanded meta: %w", err)
	}
	fmt.Printf("%s\n%s\n  ", s, recommendPrompt(name, value, suggest))
	return err
}

func recommendPrompt(name, value, suggest string) string {
	s := input.PreviewPromptS(name, value)
	return fmt.Sprintf("%s%s:", s, set.Recommend(suggest))
}
