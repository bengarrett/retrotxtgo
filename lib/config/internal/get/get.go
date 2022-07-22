package get

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
)

var (
	ErrBool   = errors.New("key is not a boolean (true/false) value")
	ErrString = errors.New("key is not a string (text) value")
	ErrUint   = errors.New("key is not a absolute number")
)

const (
	Editor     = "editor"
	FontEmbed  = "html.font.embed"
	FontFamily = "html.font.family"
	LayoutTmpl = "html.layout"
	Author     = "html.meta.author"
	Scheme     = "html.meta.color_scheme"
	Desc       = "html.meta.description"
	Genr       = "html.meta.generator"
	Keywords   = "html.meta.keywords"
	Notlate    = "html.meta.notranslate"
	Referr     = "html.meta.referrer"
	Rtx        = "html.meta.retrotxt"
	Bot        = "html.meta.robots"
	Theme      = "html.meta.theme_color"
	Title      = "html.title"
	SaveDir    = "save_directory"
	Serve      = "serve"
	Stylei     = "style.info"
	Styleh     = "style.html"
)

// Defaults for setting keys.
type Defaults map[string]interface{}

// Hints for configuration values.
type Hints map[string]string

// Tip provides brief help on the config file setting.
func Tip() Hints {
	return Hints{
		Editor:     "text editor to launch when using " + str.Example("config edit"),
		FontEmbed:  "encode and embed the font as Base64 binary-to-text within the CSS",
		FontFamily: "specifies the font to use with the HTML",
		LayoutTmpl: "HTML template for the layout of CSS, JS and fonts",
		Author:     "defines the name of the page authors",
		Scheme:     "specifies one or more color schemes with which the page is compatible",
		Desc:       "a short and accurate summary of the content of the page",
		Genr:       fmt.Sprintf("include the %s version and page generation date?", meta.Name),
		Keywords:   "words relevant to the page content",
		Notlate:    "used to declare that the page should not be translated by Google Translate",
		Referr:     "controls the Referer HTTP header attached to requests sent from the page",
		Rtx:        "include a custom tag containing the meta information of the source textfile",
		Bot:        "behavor that crawlers from Google, Bing and other engines should use with the page",
		Theme:      "indicates a suggested color that user agents should use to customize the display of the page",
		Title:      "page title that is shown in the browser tab",
		SaveDir:    fmt.Sprintf("directory to store HTML assets created by %s", meta.Name),
		Serve:      "serve files using an internal web server with this port",
		Stylei:     "syntax highlighter for the config info output",
		Styleh:     "syntax highlighter for html previews",
	}
}

// Reset returns the default configuration values.
// These values are used whenever a setting is deleted,
// or when a new configuration file is created.
func Reset() Defaults {
	return Defaults{
		Editor:     "",
		FontEmbed:  false,
		FontFamily: "vga",
		LayoutTmpl: "standard",
		Author:     "",
		Scheme:     "",
		Desc:       "",
		Genr:       true,
		Keywords:   "",
		Notlate:    false,
		Referr:     "",
		Rtx:        true,
		Bot:        "",
		Theme:      "",
		Title:      meta.Name,
		SaveDir:    "",
		Serve:      meta.WebPort,
		Stylei:     "dracula",
		Styleh:     "lovelace",
	}
}

// Bool returns the default boolean value for the key.
func Bool(key string) (bool, error) {
	if v := viper.Get(key); v != nil {
		return v.(bool), nil
	}
	switch Reset()[key].(type) {
	case bool:
		return Reset()[key].(bool), nil
	default:
		return false, fmt.Errorf("%w: %s", ErrBool, key)
	}
}

// String returns the default string value for the key.
func String(key string) (string, error) {
	if v := viper.GetString(key); v != "" {
		return v, nil
	}
	switch Reset()[key].(type) {
	case string:
		return Reset()[key].(string), nil
	default:
		return "", fmt.Errorf("%w: %s", ErrString, key)
	}
}

// UInt returns the default integer value for the key.
func UInt(key string) (uint, error) {
	if v := viper.GetUint(key); v != 0 {
		return v, nil
	}
	switch Reset()[key].(type) {
	case uint:
		return Reset()[key].(uint), nil
	default:
		return 0, fmt.Errorf("%w: %s", ErrUint, key)
	}
}

// Marshal default values for use in a YAML configuration file.
func Marshal() (interface{}, error) {
	sc := Settings{}
	for key := range Reset() {
		if err := sc.marshals(key); err != nil {
			return sc, err
		}
	}
	return sc, nil
}

// Settings types and names to be saved as YAML.
type Settings struct {
	Editor string
	HTML   struct {
		Font struct {
			Embed  bool   `yaml:"embed"`
			Family string `yaml:"family"`
			Size   string `yaml:"size"`
		}
		Layout string `yaml:"layout"`
		Meta   struct {
			Author      string `yaml:"author"`
			ColorScheme string `yaml:"colorScheme"`
			Description string `yaml:"description"`
			Keywords    string `yaml:"keywords"`
			Referrer    string `yaml:"referrer"`
			Robots      string `yaml:"robots"`
			ThemeColor  string `yaml:"themeColor"`
			Generator   bool   `yaml:"generator"`
			Notranslate bool   `yaml:"notranslate"`
			RetroTxt    bool   `yaml:"retrotxt"`
		}
		Title string `yaml:"title"`
	}
	SaveDirectory string `yaml:"saveDirectory"`
	ServerPort    uint   `yaml:"serve"`
	Style         struct {
		Info string `yaml:"info"`
		HTML string `yaml:"html"`
	}
}

// marshals sets the default value for the key.
func (sc *Settings) marshals(key string) error { // nolint:funlen
	switch key {
	case Editor:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.Editor = s
	case FontEmbed:
		b, err := Bool(key)
		if err != nil {
			return err
		}
		sc.HTML.Font.Embed = b
	case FontFamily:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Font.Family = s
	case LayoutTmpl:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Layout = s
	case Author:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Author = s
	case Scheme:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.ColorScheme = s
	case Desc:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Description = s
	case Genr:
		b, err := Bool(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Generator = b
	case Keywords:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Keywords = s
	case Notlate:
		b, err := Bool(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Notranslate = b
	case Referr:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Referrer = s
	case Rtx:
		b, err := Bool(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.RetroTxt = b
	case Bot:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.Robots = s
	case Theme:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Meta.ThemeColor = s
	case Title:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.HTML.Title = s
	case SaveDir:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.SaveDirectory = s
	case Serve:
		i, err := UInt(key)
		if err != nil {
			return err
		}
		sc.ServerPort = i
	case Stylei:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.Style.Info = s
	case Styleh:
		s, err := String(key)
		if err != nil {
			return err
		}
		sc.Style.HTML = s
	default:
		return fmt.Errorf("marshals %q: %w", key, logs.ErrConfigName)
	}
	return nil
}

// TextEditor returns the path of a configured or discovered text editor.
func TextEditor(w io.Writer) string {
	edit := viper.GetString("editor")
	_, err := exec.LookPath(edit)
	if err != nil {
		if edit != "" {
			fmt.Fprintf(w, "%s\nwill attempt to use the $EDITOR environment variable\n", err)
		}
		if err := viper.BindEnv("editor", "EDITOR"); err != nil {
			return DiscEditor()
		}
		edit = viper.GetString("editor")
		if _, err := exec.LookPath(edit); err != nil {
			return DiscEditor()
		}
	}
	return edit
}

// DiscEditor attempts to discover any known text editors on the host system.
func DiscEditor() string {
	editors := [5]string{"nano", "vim", "emacs"}
	if runtime.GOOS == "windows" {
		editors[3] = "notepad++.exe"
		editors[4] = "notepad.exe"
	}
	edit := ""
	for _, editor := range editors {
		if editor == "" {
			continue
		}
		if _, err := exec.LookPath(editor); err == nil {
			edit = editor
			break
		}
	}
	return edit
}
