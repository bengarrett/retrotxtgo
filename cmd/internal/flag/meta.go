package flag

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Meta struct {
	Key   string   // configuration name
	Strg  *string  // StringVarP(p) argument value
	Boo   *bool    // BoolVarP(p) argument value
	I     *uint    // UintVar(p) argument value
	Name  string   // flag long name
	Short string   // flag short name
	Opts  []string // flag choices for display in the usage string
}

// Init initializes the create command flags and their help.
func Init() map[int]Meta {
	const (
		serve = iota
		layout
		style
		title
		desc
		author
		retro
		gen
		cscheme
		kwords
		nolang
		refer
		bots
		themec
		fontf
		fonte
		body
		cache
	)
	return map[int]Meta{
		// main tag flags
		style:  {"style.html", &Build.Syntax, nil, nil, "syntax-style", "", nil},
		desc:   {"html.meta.description", &Build.Metadata.Description.Value, nil, nil, "meta-description", "d", nil},
		author: {"html.meta.author", &Build.Metadata.Author.Value, nil, nil, "meta-author", "a", nil},
		retro:  {"html.meta.retrotxt", nil, &Build.Metadata.RetroTxt, nil, "meta-retrotxt", "r", nil},
		// minor tag flags
		gen:     {"html.meta.generator", nil, &Build.Metadata.Generator, nil, "meta-generator", "g", nil},
		cscheme: {"html.meta.color_scheme", &Build.Metadata.ColorScheme.Value, nil, nil, "meta-color-scheme", "", nil},
		kwords:  {"html.meta.keywords", &Build.Metadata.Keywords.Value, nil, nil, "meta-keywords", "", nil},
		nolang:  {"html.meta.notranslate", nil, &Build.Metadata.NoTranslate, nil, "meta-notranslate", "", nil},
		refer:   {"html.meta.referrer", &Build.Metadata.Referrer.Value, nil, nil, "meta-referrer", "", nil},
		bots:    {"html.meta.robots", &Build.Metadata.Robots.Value, nil, nil, "meta-robots", "", nil},
		themec:  {"html.meta.theme_color", &Build.Metadata.ThemeColor.Value, nil, nil, "meta-theme-color", "", nil},
		fontf:   {"html.font.family", &Build.FontFamily.Value, nil, nil, "font-family", "f", nil},
		// hidden flags
		body: {"html.body", &Build.Source.HiddenBody, nil, nil, "body", "b", nil},
	}
}

// Sort creates an ordered index of the meta flags.
func Sort(flags map[int]Meta) []int {
	k := make([]int, len(flags))
	for i := range flags {
		k[i] = i
	}
	sort.Ints(k)
	return k
}

// Init initializes the public facing flags.
func (c *Meta) Init(cmd *cobra.Command, buf bytes.Buffer) bytes.Buffer {
	switch {
	case c.Key == "serve":
		fmt.Fprintf(&buf, "\ngive a 0 value, %s or %s, to use the default %d port",
			str.Example("-p0"), str.Example("--serve=0"), meta.WebPort)
		cmd.Flags().UintVarP(c.I, c.Name, c.Short, viper.GetUint(c.Key), buf.String())
	case c.Strg != nil:
		cmd.Flags().StringVarP(c.Strg, c.Name, c.Short, viper.GetString(c.Key), buf.String())
	case c.Boo != nil:
		cmd.Flags().BoolVarP(c.Boo, c.Name, c.Short, viper.GetBool(c.Key), buf.String())
	case c.I != nil:
		cmd.Flags().UintVarP(c.I, c.Name, c.Short, viper.GetUint(c.Key), buf.String())
	}
	return buf
}
