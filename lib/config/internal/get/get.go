package get

import (
	"errors"

	"github.com/bengarrett/retrotxtgo/lib/logs"
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
	Scheme     = "html.meta.color-scheme"
	Desc       = "html.meta.description"
	Genr       = "html.meta.generator"
	Keywords   = "html.meta.keywords"
	Notlate    = "html.meta.notranslate"
	Referr     = "html.meta.referrer"
	Rtx        = "html.meta.retrotxt"
	Bot        = "html.meta.robots"
	Theme      = "html.meta.theme-color"
	Title      = "html.title"
	SaveDir    = "save-directory"
	Serve      = "serve"
	Stylei     = "style.info"
	Styleh     = "style.html"
)

// Defaults for setting keys.
type Defaults map[string]interface{}

// Reset configuration values.
// These are the default values whenever a setting is deleted,
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
func Bool(key string) bool {
	if v := viper.Get(key); v != nil {
		return v.(bool)
	}
	switch Reset()[key].(type) {
	case bool:
		return Reset()[key].(bool)
	default:
		logs.FatalMark(key, ErrBool, logs.ErrConfigName)
	}
	return false
}

// UInt returns the default integer value for the key.
func UInt(key string) uint {
	if v := viper.GetUint(key); v != 0 {
		return v
	}
	switch Reset()[key].(type) {
	case uint:
		return Reset()[key].(uint)
	default:
		logs.FatalMark(key, ErrUint, logs.ErrConfigName)
	}
	return 0
}

// String returns the default string value for the key.
func String(key string) string {
	if v := viper.GetString(key); v != "" {
		return v
	}
	switch Reset()[key].(type) {
	case string:
		return Reset()[key].(string)
	default:
		logs.FatalMark(key, ErrString, logs.ErrConfigName)
	}
	return ""
}
