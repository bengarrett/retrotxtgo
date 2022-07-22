package set

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/colorise"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

const RemoveChr = "-"

var (
	ErrBreak     = errors.New("break of loop")
	ErrDataEmpty = errors.New("data argument cannot be empty")
	ErrSaveType  = errors.New("save value type is unsupported")
	ErrSkip      = errors.New("skipped, no change")
	ErrUnused    = errors.New("is unused")
	ErrVal       = errors.New("value is unsupported")
	ErrWriter    = errors.New("writer argument cannot be nil")
)

func skip(w io.Writer, name string, setup bool, value interface{}) error {
	err := SkipWrite(name, value)
	switch {
	case errors.Is(err, ErrSkip):
	case errors.Is(err, ErrUnused):
		fmt.Fprint(w, skipSet(setup))
		return nil
	case err != nil:
		return err
	}
	return nil
}

// Write the value of the named setting to the configuration file.
func Write(w io.Writer, name string, setup bool, value interface{}) error {
	if w == nil {
		return ErrWriter
	}
	if name == "" {
		return fmt.Errorf("save: %w", logs.ErrNameNil)
	}
	if !Validate(name) {
		return fmt.Errorf("save %q: %w", name, logs.ErrConfigName)
	}
	switch v := value.(type) {
	case string:
		if v == RemoveChr {
			value = ""
			break
		}
		if err := skip(w, name, setup, value); err != nil {
			return err
		}
		if value == "" {
			return nil
		}
	case any:
		if err := skip(w, name, setup, value); err != nil {
			return err
		}
	}
	if value == nil {
		return ErrVal
	}
	viper.Set(name, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			fmt.Fprintf(w, "  %s %s\n", str.ColSuc(name), ErrUnused)
			return nil
		}
	default:
	}
	fmt.Fprintf(w, "  %s is set to \"%v\"\n", str.ColSuc(name), value)
	return nil
}

// SkipWrite returns an error if the named value doesn't need updating.
func SkipWrite(name string, value interface{}) error {
	if viper.Get(name) == nil {
		return fmt.Errorf("setting name: %s, type: %T, %w", name, nil, logs.ErrConfigName)
	}
	switch v := value.(type) {
	case bool:
		if viper.Get(name).(bool) == v {
			return ErrSkip
		}
		return nil
	case string:
		if viper.Get(name).(string) == v {
			return ErrSkip
		}
		if value.(string) == "" {
			return ErrUnused
		}
		return nil
	case uint:
		i := viper.Get(name).(int)
		if uint(i) == uint(v) {
			return ErrSkip
		}
		if name == get.Serve && v == 0 {
			return ErrUnused
		}
		return nil
	case int:
		if viper.Get(name) == v {
			return ErrSkip
		}
		if name == get.Serve && v == 0 {
			return ErrUnused
		}
		return nil
	}
	return fmt.Errorf("setting: %s, type: %T, %w", name, value, ErrSaveType)
}

func skipSet(setup bool) string {
	if !setup {
		return fmt.Sprintf("  %s\n", ErrSkip.Error())
	}
	return str.ColSuc("\r  "+ErrSkip.Error()) + "\n"
}

// Directory prompts and checks directory path for save.
func Directory(w io.Writer, name string, setup bool) error {
	if w == nil {
		return ErrWriter
	}
	if name == "" {
		return fmt.Errorf("set directory: %w", logs.ErrNameNil)
	}
	s := DirExpansion(prompt.String(w))
	if s == "" {
		fmt.Fprint(w, skipSet(setup))
		return ErrBreak
	}
	if s == RemoveChr {
		if err := Write(w, name, setup, RemoveChr); err != nil {
			return err
		}
		return ErrBreak
	}
	if _, err := os.Stat(s); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Fprintf(w, "%s The directory does not exist: %s\n", str.Alert(), s)
			return nil
		}
		fmt.Fprintf(w, "%s the directory is invalid: %s\n", str.Alert(), errors.Unwrap(err))
		return nil
	}
	if err := Write(w, name, setup, s); err != nil {
		return err
	}
	return ErrBreak
}

// DirExpansion traverses the named directory to apply shell-like expansions.
// It supports limited Bash tilde, shell dot and double dot syntax.
func DirExpansion(name string) string {
	const sep, homeDir, currentDir, parentDir = string(os.PathSeparator), "~", ".", ".."
	if name == "" || name == sep {
		return name
	}
	// Bash tilde expension http://www.gnu.org/software/bash/manual/html_node/Tilde-Expansion.html
	dir, r, paths := "", bool(name[0:1] == sep), strings.Split(name, sep)
	var err error
	for i, s := range paths {
		var p string
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

// Editor checks the existence of given text editor location
// and saves it as a configuration regardless of the result.
func Editor(w io.Writer, name string, setup bool) error {
	if w == nil {
		return ErrWriter
	}
	if name == "" {
		return fmt.Errorf("set editor: %w", logs.ErrNameNil)
	}
	s := prompt.String(w)
	switch s {
	case RemoveChr:
		if err := Write(w, name, setup, RemoveChr); err != nil {
			return err
		}
		return nil
	case "":
		fmt.Fprint(w, skipSet(setup))
		return nil
	}
	if _, err := exec.LookPath(s); err != nil {
		fmt.Fprintf(w, "%s%s\nThe %s editor is not usable by %s.\n",
			str.Alert(), errors.Unwrap(err), s, meta.Name)
	}
	return Write(w, name, setup, s)
}

// Font previews and saves a default font setting.
func Font(w io.Writer, value string, setup bool) error {
	f := create.Family(value)
	if f == create.Automatic {
		f = create.VGA
	}
	s := fmt.Sprintf("%s\n%s\n%s\n%s\n",
		"  @font-face {",
		fmt.Sprintf("    font-family: \"%s\";", f),
		fmt.Sprintf("    src: url(\"%s.woff2\") format(\"woff2\");", f),
		"  }")
	if err := colorise.CSS(w, s); err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n%s%s %s: ",
		str.ColFuz("  About font families: https://developer.mozilla.org/en-US/docs/Web/CSS/font-family"),
		"  Choose a font, ",
		str.UnderlineKeys(create.Fonts()...),
		fmt.Sprintf("(suggestion: %s)", str.Example("automatic")))
	return ShortStrings(w, get.FontFamily, setup, create.Fonts()...)
}

// Font previews and saves the embedded Base64 font setting.
func FontEmbed(w io.Writer, value, setup bool) error {
	const name = get.FontEmbed
	elm := fmt.Sprintf("  %s\n  %s\n  %s\n",
		"@font-face{",
		"  font-family: vga8;",
		"  src: url(data:font/woff2;base64,[a large font binary will be embedded here]...) format('woff2');",
	)
	if err := colorise.CSS(w, elm); err != nil {
		return err
	}
	q := fmt.Sprintf("%s\n%s\n%s",
		"  The use of this setting not recommended,",
		"  unless you always want large, self-contained HTML files for distribution.",
		"  Embed the font as Base64 text within the HTML")
	if value {
		q = "  Keep the embedded font option?"
	}
	q += Recommend("no")
	b := prompt.YesNo(w, q, viper.GetBool(name))
	return Write(w, name, setup, b)
}

// Generator prompts for and previews the custom program generator meta tag.
func Generator(w io.Writer, value bool) error {
	const name = get.Genr
	elm := fmt.Sprintf("  %s\n    %s\n  %s\n",
		"<head>",
		fmt.Sprintf("<meta name=\"generator\" content=\"%s, %s\">",
			meta.Name, meta.App.Date),
		"</head>")
	if !meta.IsGoBuild() {
		elm = fmt.Sprintf("  %s\n    %s\n  %s\n",
			"<head>",
			fmt.Sprintf("<meta name=\"generator\" content=\"%s %s, %s\">",
				meta.Name, meta.Print(), meta.App.Date),
			"</head>")
	}
	if err := colorise.HTML(w, elm); err != nil {
		return err
	}
	p := "Enable the generator element?"
	if value {
		p = "Keep the generator element?"
	}
	p = fmt.Sprintf("  %s%s", p, Recommend("yes"))
	b := prompt.YesNo(w, p, viper.GetBool(name))
	return Write(w, name, true, b)
}

// Index prompts for a value from a list of valid choices and saves the result.
func Index(w io.Writer, name string, setup bool, data ...string) error {
	if name == "" {
		return logs.ErrNameNil
	}
	if len(data) == 0 {
		return ErrDataEmpty
	}
	s := prompt.IndexStrings(w, &data, setup)
	data = append(data, RemoveChr)
	data = append(data, "")
	sort.Strings(data)
	// validate s against data
	i := sort.Search(len(data), func(i int) bool { return data[i] >= s })
	if i < len(data) && data[i] == s {
		return Write(w, name, setup, s)
	}
	return fmt.Errorf("config set %s %w: %q", name, ErrVal, s)
}

// Keys list all the available configuration setting names sorted alphabetically.
func Keys() []string {
	keys := make([]string, len(get.Reset()))
	i := 0
	for key := range get.Reset() {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// NoTranslate previews and prompts for the notranslate HTML attribute
// and Google meta elemenet.
func NoTranslate(w io.Writer, value, setup bool) error {
	const name = get.Notlate
	elm := fmt.Sprintf("  %s\n    %s\n  %s\n  %s\n",
		"<head>",
		"<meta name=\"google\" content=\"notranslate\">",
		"</head>",
		"<body class=\"notranslate\">")
	if err := colorise.HTML(w, elm); err != nil {
		return err
	}
	q := "Enable the no translate option?"
	if value {
		q = "Keep the translate option?"
	}
	q = fmt.Sprintf("  %s%s", q, Recommend("no"))
	b := prompt.YesNo(w, q, viper.GetBool(name))
	return Write(w, name, setup, b)
}

// Port prompts for and saves HTTP port.
func Port(w io.Writer, name string, setup bool) error {
	if name == "" {
		return logs.ErrNameNil
	}
	u := prompt.Port(w, true, setup)
	if u == prompt.PortReset {
		fmt.Fprint(w, skipSet(setup))
		return nil
	}
	return Write(w, name, setup, u)
}

// Recommend uses the s value as a user input suggestion.
func Recommend(s string) string {
	if s == "" {
		return fmt.Sprintf(" (suggestion: %s)", str.Example("do not use"))
	}
	return fmt.Sprintf(" (suggestion: %s)", str.Example(s))
}

// RetroTxt previews and prompts the custom retrotxt meta tag.
func RetroTxt(w io.Writer, value bool) error {
	const name = get.Rtx
	elm := fmt.Sprintf("%s\n%s\n%s\n",
		"  <head>",
		"    <meta name=\"retrotxt\" content=\"encoding: IBM437; linebreak: CRLF; length: 50; width: 80; name: file.txt\">",
		"  </head>")
	if err := colorise.HTML(w, elm); err != nil {
		return err
	}
	p := "Enable the retrotxt element?"
	if value {
		p = "Keep the retrotxt element?"
	}
	p = fmt.Sprintf("  %s%s", p, Recommend("yes"))
	b := prompt.YesNo(w, p, viper.GetBool(name))
	return Write(w, name, true, b)
}

// ShortStrings prompts and saves setting values that support 1 character aliases.
func ShortStrings(w io.Writer, name string, setup bool, data ...string) error {
	if name == "" {
		return fmt.Errorf("set short string: %w", logs.ErrNameNil)
	}
	s := prompt.ShortStrings(w, &data)
	return Write(w, name, setup, s)
}

// String prompts and saves a single word setting value.
func String(w io.Writer, name string, setup bool) error {
	if name == "" {
		return fmt.Errorf("set string: %w", logs.ErrNameNil)
	}
	s := prompt.String(w)
	return Write(w, name, setup, s)
}

// Strings prompts and saves a string of text setting value.
func Strings(w io.Writer, name string, setup bool, data ...string) error {
	if name == "" {
		return fmt.Errorf("set strings: %w", logs.ErrNameNil)
	}
	s := prompt.Strings(w, &data, setup)
	return Write(w, name, setup, s)
}

// Title prompts for and previews a HTML title element value.
func Title(w io.Writer, name, value string, setup bool) error {
	elm := fmt.Sprintf("%s\n%s\n%s\n",
		"  <head>",
		fmt.Sprintf("    <title>%s</title>", value),
		"  </head>")
	if err := colorise.HTML(w, elm); err != nil {
		return err
	}
	fmt.Fprintf(w, "%s\n%s\n  ",
		str.ColFuz("  About this value: https://developer.mozilla.org/en-US/docs/Web/HTML/Element/title"),
		fmt.Sprintf("  Choose a new %s:", get.Tip()[name]))
	return String(w, name, setup)
}

// Validate the existence of the key in a list of settings.
func Validate(key string) bool {
	keys := Keys()
	// var i must be sorted in ascending order.
	if i := sort.SearchStrings(keys, key); i == len(keys) || keys[i] != key {
		return false
	}
	return true
}
