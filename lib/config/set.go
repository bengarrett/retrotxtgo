package config

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/update"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

// List and print all the available configurations.
func List(w io.Writer) error {
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
	tw := tabwriter.NewWriter(w, minWidth, tabWidth, 0, ' ', 0)
	cmds := fmt.Sprintf(" %s config set ", meta.Bin)
	title := fmt.Sprintf("  Available %s configurations and settings", meta.Name)
	fmt.Fprintln(tw, "\n"+str.ColPri(title))
	fmt.Fprintln(tw, str.HR(len(title)))
	fmt.Fprintln(tw, tabs)
	fmt.Fprintf(tw, "Alias\t\tName\t\tHint\n")
	for i, key := range keys {
		tip := get.Tip()[key]
		fmt.Fprintln(tw, tabs)
		fmt.Fprintf(tw, " %d\t\t%s\t\t%s", i, key, suffix(capitalize(tip)))
		switch key {
		case get.LayoutTmpl:
			fmt.Fprintf(tw, "\n%schoices: %s (suggestion: %s)",
				tabs, str.ColPri(strings.Join(create.Layouts(), ", ")), str.Example("standard"))
		case get.Serve:
			fmt.Fprintf(tw, "\n%schoices: %s",
				tabs, input.PortInfo())
		}
		fmt.Fprint(tw, "\n")
	}
	fmt.Fprintln(tw, tabs)
	fmt.Fprintln(tw, str.HR(len(title)))
	fmt.Fprintln(tw, "\nEither the setting Name or the Alias can be used.")
	fmt.Fprintf(tw, "\n%s # To change the meta description setting\n",
		str.Example(cmds+get.Desc))
	fmt.Fprintf(tw, "%s # Will also change the meta description setting\n", str.Example(cmds+"6"))
	fmt.Fprintln(tw, "\nMultiple settings are supported.")
	fmt.Fprintf(tw, "\n%s\n", str.Example(cmds+"style.html style.info"))
	return tw.Flush()
}

// Set edits and saves a named setting within a configuration file.
// It also accepts numeric index values printed by List().
func Set(w io.Writer, name string) error {
	i, err := strconv.Atoi(name)
	namedSetting := err != nil
	switch {
	case namedSetting:
		return Update(w, name, false)
	case i >= 0 && i <= (len(get.Reset())-1):
		k := set.Keys()
		return Update(w, k[i], false)
	default:
		return Update(w, name, false)
	}
}

// Update edits and saves a named setting within a configuration file.
func Update(w io.Writer, name string, setup bool) error {
	if !set.Validate(name) {
		fmt.Fprintln(w, logs.Hint("config set --list", logs.ErrConfigName))
		return nil
	}
	if !setup {
		fmt.Fprint(w, Location())
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
		fmt.Fprint(w, update.Bool(b, name))
	}
	if s, ok := value.(string); ok {
		update.String(w, s, name, value.(string))
	}
	u := input.Update{Name: name, Setup: setup, Value: value}
	err := updatePrompt(w, u)
	switch {
	case errors.Is(err, prompt.ErrSkip):
		fmt.Fprintln(w, prompt.ErrSkip)
	case err != nil:
		return err
	}
	return nil
}

// updatePrompt prompts the user for input to a config file setting.
func updatePrompt(w io.Writer, u input.Update) error {
	switch u.Name {
	case "editor":
		return input.Editor(w, u)
	case get.SaveDir:
		return input.SaveDir(w, u)
	case get.Serve:
		return input.Serve(w, u)
	case get.Styleh:
		return input.StyleHTML(w, u)
	case get.Stylei:
		return input.StyleInfo(w, u)
	}
	return metaPrompts(w, u)
}

// metaPrompts prompts the user for a meta setting.
func metaPrompts(w io.Writer, u input.Update) error {
	switch u.Name {
	case get.FontEmbed:
		return set.FontEmbed(w, u.Value.(bool), u.Setup)
	case get.FontFamily:
		return set.Font(w, u.Value.(string), u.Setup)
	case get.LayoutTmpl:
		return input.Layout(w, u)
	case get.Author,
		get.Desc,
		get.Keywords:
		if err := input.PreviewMeta(w, u.Name, u.Value.(string)); err != nil {
			return err
		}
		return set.String(w, u.Name, u.Setup)
	case get.Theme:
		return theme(w, u)
	case get.Scheme:
		return input.ColorScheme(w, u)
	case get.Genr:
		return set.Generator(w, u.Value.(bool))
	case get.Notlate:
		return set.NoTranslate(w, u.Value.(bool), u.Setup)
	case get.Referr:
		return referr(w, u)
	case get.Bot:
		return bot(w, u)
	case get.Rtx:
		return set.RetroTxt(w, u.Value.(bool))
	case get.Title:
		return set.Title(w, u.Name, u.Value.(string), u.Setup)
	}
	return fmt.Errorf("%w: %s", prompt.ErrSkip, u.Name)
}

func theme(w io.Writer, u input.Update) error {
	if err := recommendMeta(w, u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	set.String(w, u.Name, u.Setup)
	return nil
}

func bot(w io.Writer, u input.Update) error {
	if err := recommendMeta(w, u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	cr := create.Robots()
	fmt.Fprintf(w, "%s\n  ", str.NumberizeKeys(cr[:]...))
	return set.Index(w, u.Name, u.Setup, cr[:]...)
}

func referr(w io.Writer, u input.Update) error {
	if err := recommendMeta(w, u.Name, u.Value.(string), ""); err != nil {
		return err
	}
	cr := create.Referrer()
	fmt.Fprintf(w, "%s\n  ", str.NumberizeKeys(cr[:]...))
	return set.Index(w, u.Name, u.Setup, cr[:]...)
}

func recommendMeta(w io.Writer, name, value, suggest string) error {
	if err := input.PrintMeta(w, name, value); err != nil {
		return fmt.Errorf("recommanded meta: %w", err)
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s\n  ", recommendPrompt(name, value, suggest))
	return nil
}

func recommendPrompt(name, value, suggest string) string {
	s := input.PreviewPromptS(name, value)
	return fmt.Sprintf("%s%s:", s, set.Recommend(suggest))
}
