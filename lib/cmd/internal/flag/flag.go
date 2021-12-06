package flag

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/sauce"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
)

var ErrFilenames = errors.New("ignoring [filenames]")

type Configs struct {
	Configs bool
	Ow      bool
	Styles  bool
	Style   string
}

type RootFlags struct {
	Config string
}

type ViewFlags struct {
	Controls []string
	Encode   string
	Swap     []string
	To       string
	Width    int
}

var RootFlag RootFlags

var ViewFlag = ViewFlags{
	Controls: []string{"eof", "tab"},
	Encode:   "CP437",
	Swap:     []string{"null", "bar"},
	To:       "",
	Width:    0,
}

var Config Configs

var InfoFlag struct {
	Format string
}

// flags container.
var HTML create.Args

func ConfigInfo() (exit bool) {
	if Config.Configs {
		if err := config.List(); err != nil {
			logs.FatalFlag("config info", "list", err)
		}
	}
	if Config.Styles {
		fmt.Print(str.JSONStyles(fmt.Sprintf("%s info --style", meta.Bin)))
		return true
	}
	style := viper.GetString("style.info")
	if Config.Style != "" {
		style = Config.Style
	}
	if style == "" {
		style = "dracula"
	}
	if err := config.Info(style); err != nil {
		logs.Fatal(err)
	}
	return false
}

type Create struct {
	Controls []string // character encoding used by the filename
	Encode   string   // use these control codes
	Swap     []string // swap out these characters with UTF8 alternatives
}

const EncodingDefault = "CP437"

// createFlag
// CreateDefaults contain default values.
var CreateDefaults = Create{
	Controls: []string{"eof", "tab"},
	Encode:   EncodingDefault,
	Swap:     []string{"null", "bar"},
}

type Meta struct {
	Key   string   // configuration name
	Strg  *string  // StringVarP(p) argument value
	Boo   *bool    // BoolVarP(p) argument value
	I     *uint    // UintVar(p) argument value
	Name  string   // flag long name
	Short string   // flag short name
	Opts  []string // flag choices for display in the usage string
}

// initBodyFlag initializes the hidden body flag.
func (c *Meta) Body(buf bytes.Buffer) bytes.Buffer {
	switch {
	case c.Key == "html.body":
		fmt.Fprint(&buf, "override and inject a string into the HTML body element")
	case len(c.Opts) == 0:
		fmt.Fprint(&buf, config.Tip()[c.Key])
	default:
		fmt.Fprint(&buf, str.Options(config.Tip()[c.Key], true, true, c.Opts...))
	}
	return buf
}

// initFlags initializes the create command flags and their help.
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
		// output
		serve: {"serve", nil, nil, &HTML.Port, "serve", "p", nil},
		// main tag flags
		style:  {"style.html", &HTML.Syntax, nil, nil, "syntax-style", "", nil},
		layout: {"html.layout", &HTML.Layout, nil, nil, "layout", "l", create.Layouts()},
		title:  {"html.title", &HTML.Title.Value, nil, nil, "title", "t", nil},
		desc:   {"html.meta.description", &HTML.Metadata.Description.Value, nil, nil, "meta-description", "d", nil},
		author: {"html.meta.author", &HTML.Metadata.Author.Value, nil, nil, "meta-author", "a", nil},
		retro:  {"html.meta.retrotxt", nil, &HTML.Metadata.RetroTxt, nil, "meta-retrotxt", "r", nil},
		// minor tag flags
		gen:     {"html.meta.generator", nil, &HTML.Metadata.Generator, nil, "meta-generator", "g", nil},
		cscheme: {"html.meta.color-scheme", &HTML.Metadata.ColorScheme.Value, nil, nil, "meta-color-scheme", "", nil},
		kwords:  {"html.meta.keywords", &HTML.Metadata.Keywords.Value, nil, nil, "meta-keywords", "", nil},
		nolang:  {"html.meta.notranslate", nil, &HTML.Metadata.NoTranslate, nil, "meta-notranslate", "", nil},
		refer:   {"html.meta.referrer", &HTML.Metadata.Referrer.Value, nil, nil, "meta-referrer", "", nil},
		bots:    {"html.meta.robots", &HTML.Metadata.Robots.Value, nil, nil, "meta-robots", "", nil},
		themec:  {"html.meta.theme-color", &HTML.Metadata.ThemeColor.Value, nil, nil, "meta-theme-color", "", nil},
		fontf:   {"html.font.family", &HTML.FontFamily.Value, nil, nil, "font-family", "f", nil},
		fonte:   {"html.font.embed", nil, &HTML.FontEmbed, nil, "font-embed", "", nil},
		// hidden flags
		body:  {"html.body", &HTML.Source.HiddenBody, nil, nil, "body", "b", nil},
		cache: {"html.layout.cache", nil, &HTML.Save.Cache, nil, "cache", "", nil},
	}
}

// Sort creates an ordered index of the meta flags.
func Sort(cfg map[int]Meta) []int {
	k := make([]int, len(cfg))
	for i := range cfg {
		k[i] = i
	}
	sort.Ints(k)
	return k
}

// initFlags initializes the public facing flags.
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

// SAUCE parses any embedded SAUCE metadata.
func SAUCE(src *[]byte) {
	if HTML.SauceData.Use {
		sr := sauce.Decode(*src)
		if !sr.Valid() {
			return
		}
		HTML.SauceData.Title = sr.Title
		HTML.SauceData.Author = sr.Author
		HTML.SauceData.Group = sr.Group
		HTML.SauceData.Description = sr.Desc
		HTML.SauceData.Width = uint(sr.Info.Info1.Value)
		HTML.SauceData.Lines = uint(sr.Info.Info2.Value)
	}
}

// endOfFile determines if the end-of-file control flag was requested.
func EndOfFile(flags convert.Flag) bool {
	for _, c := range flags.Controls {
		if c == "eof" {
			return true
		}
	}
	return false
}

// initArgs initializes the command arguments and flags.
func InitArgs(cmd *cobra.Command, args ...string) ([]string, *convert.Convert, sample.Flags, error) {
	conv := convert.Convert{}
	conv.Flags = convert.Flag{
		Controls:  ViewFlag.Controls,
		SwapChars: ViewFlag.Swap,
		MaxWidth:  ViewFlag.Width,
	}
	l := len(args)

	if c := cmd.Flags().Lookup("controls"); !c.Changed {
		conv.Flags.Controls = []string{"eof", "tab"}
	}
	if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
		conv.Flags.SwapChars = []string{"null", "bar"}
	}
	if filesystem.IsPipe() {
		var err error
		if l > 0 {
			err = fmt.Errorf("%v;%w for piped text", err, ErrFilenames)
			args = []string{""}
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, logs.Sprint(err))
		}
	} else if err := PrintUsage(cmd, args...); err != nil {
		logs.Fatal(err)
	}
	if l == 0 {
		args = []string{""}
	}
	samp, err := InitEncodings(cmd, "")
	if err != nil {
		return nil, nil, samp, err
	}
	if conv.Input.Encoding == nil {
		conv.Input.Encoding = DfaultInput()
	}
	return args, &conv, samp, nil
}

// dfaultInput returns the default encoding when the --encoding flag is unused.
func DfaultInput() encoding.Encoding {
	if filesystem.IsPipe() {
		return unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	}
	return charmap.CodePage437
}

// readArg returns the content of argument supplied filepath, embed sample file or piped data.
func ReadArg(arg string, cmd *cobra.Command, c *convert.Convert, f sample.Flags) ([]byte, error) {
	var (
		b   []byte
		err error
	)
	// if no argument, then assume the source is piped via stdin
	if arg == "" {
		b, err = filesystem.ReadPipe()
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	// attempt to see if arg is a embed sample file request
	if b, err = openSample(arg, cmd, c, f); err != nil {
		return nil, err
	} else if b != nil {
		return b, nil
	}
	// the arg should be a filepath
	b, err = openFile(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// openFile returns the content of the named file given via an argument.
func openFile(arg string) ([]byte, error) {
	b, err := filesystem.Read(arg)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// openSample returns the content of the named embed sample file given via an argument.
func openSample(arg string, cmd *cobra.Command, c *convert.Convert, f sample.Flags) ([]byte, error) {
	if ok := sample.Valid(arg); !ok {
		return nil, nil
	}
	p, err := f.Open(arg, c)
	if err != nil {
		return nil, err
	}
	// handle flags
	if ff := cmd.Flags().Lookup("font-family"); ff != nil && !ff.Changed {
		// only apply the sample font when the --font-family flag is unused
		// html is a global flag, create.Args
		HTML.FontFamily.Value = p.Font.String()
	}
	return []byte(string(p.Runes)), nil
}

// printUsage will print the help and exit when no arguments are supplied.
func PrintUsage(cmd *cobra.Command, args ...string) error {
	if len(args) == 0 {
		if err := cmd.Help(); err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}

// initEncodings applies the public --encode and the hidden --to encoding values to embed sample data.
func InitEncodings(cmd *cobra.Command, dfault string) (sample.Flags, error) {
	parse := func(name string) (encoding.Encoding, error) {
		cp := cmd.Flags().Lookup(name)
		lookup := dfault
		if cp != nil && cp.Changed {
			lookup = cp.Value.String()
		} else if dfault == "" {
			return nil, nil
		}
		if lookup == "" {
			return nil, nil
		}
		return convert.Encoder(lookup)
	}
	var (
		frm encoding.Encoding
		to  encoding.Encoding
	)
	if cmd == nil {
		return sample.Flags{}, nil
	}
	// handle encode flag or apply the default
	frm, err := parse("encode")
	if err != nil {
		return sample.Flags{}, err
	}
	// handle the hidden reencode (--to) flag
	to, err = parse("to")
	if err != nil {
		return sample.Flags{}, err
	}
	return sample.Flags{From: frm, To: to}, err
}
