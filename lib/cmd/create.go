// nolint: gochecknoglobals,gochecknoinits
package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/convert"
	"retrotxt.com/retrotxt/lib/create"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sample"
	"retrotxt.com/retrotxt/lib/sauce"
	"retrotxt.com/retrotxt/lib/str"
)

type createFlags struct {
	controls []string // character encoding used by the filename
	encode   string   // use these control codes
	swap     []int    // swap out these characters with UTF8 alternatives
}

type metaFlag struct {
	key   string   // configuration name
	strg  *string  // StringVarP(p) argument value
	boo   *bool    // BoolVarP(p) argument value
	i     *uint    // UintVar(p) argument value
	name  string   // flag long name
	short string   // flag short name
	opts  []string // flag choices for display in the usage string
}

// createFlag contain default values.
var createFlag = createFlags{
	controls: []string{eof, tab},
	encode:   "CP437",
	swap:     []int{null, verticalBar},
}

// flags container.
var html create.Args

// exampleCmd returns help usage examples.
var exampleCmd = func() string {
	var b bytes.Buffer
	tmpl := `  retrotxt create file.txt -t "A text file" -d "Some text goes here"
  retrotxt create file1.txt file2.asc --save
  retrotxt create ~{{.}}Downloads{{.}}file.txt --archive
  retrotxt create file.txt --serve=8080
  cat file.txt | retrotxt create`
	t := template.Must(template.New("example").Parse(tmpl))
	err := t.Execute(&b, string(os.PathSeparator))
	if err != nil {
		log.Fatal(err)
	}
	return str.Cinf(b.String())
}

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:     "create [filenames]",
	Aliases: []string{"c", "html"},
	Short:   "Create a HTML document from a text file",
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		f := convert.Flags{
			Controls:  createFlag.controls,
			SwapChars: createFlag.swap,
		}
		// handle defaults, use these control codes
		if c := cmd.Flags().Lookup("controls"); !c.Changed {
			f.Controls = []string{eof, tab}
		}
		// handle defaults, swap out these characters with UTF8 alternatives
		if s := cmd.Flags().Lookup("swap-chars"); !s.Changed {
			f.SwapChars = []int{null, verticalBar}
		}
		// handle the defaults for most other flags
		stringFlags(cmd)
		// handle standard input (stdio)
		if filesystem.IsPipe() {
			parsePipe(cmd, f)
		}
		// handle the hidden --body flag value,
		// used for debugging, it ignores most other flags and
		// overrides the <pre></pre> content before exiting
		if body := cmd.Flags().Lookup("body"); body.Changed {
			parseBody(cmd)
		}
		// print help if no flags are supplied
		checkUse(cmd, args...)
		// parse the flags to create the HTML
		parseFiles(cmd, f, args...)
	},
}

// stringFlags handles the defaults for flags that accept strings.
// These flags are parse to three different states.
// 1) the flag is unchanged, so use the configured viper default.
// 2) the flag has a new value to overwrite viper default.
// 3) a blank flag value is given to overwrite viper default with an empty/disable value.
func stringFlags(cmd *cobra.Command) {
	var changed = func(key string) bool {
		l := cmd.Flags().Lookup(key)
		if l == nil {
			return false
		}
		return l.Changed
	}
	html.FontFamily.Flag = changed("font-family")
	html.Metadata.Author.Flag = changed("meta-author")
	html.Metadata.ColorScheme.Flag = changed("meta-color-scheme")
	html.Metadata.Description.Flag = changed("meta-description")
	html.Metadata.Keywords.Flag = changed("meta-keywords")
	html.Metadata.Referrer.Flag = changed("meta-referrer")
	html.Metadata.Robots.Flag = changed("meta-robots")
	html.Metadata.ThemeColor.Flag = changed("meta-theme-color")
	html.Title.Flag = changed("title")
	ff := cmd.Flags().Lookup("font-family")
	if !ff.Changed {
		html.FontFamily.Value = "vga"
	}
	if html.FontFamily.Value == "" {
		html.FontFamily.Value = ff.Value.String()
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	// root config must be initialized before getting saved default values
	initConfig()
	// output flags
	flagEncode(&createFlag.encode, createCmd)
	flagControls(&createFlag.controls, createCmd)
	flagRunes(&createFlag.swap, createCmd)
	dir := saveDir()
	createCmd.Flags().BoolVarP(&html.Save.AsFiles, "save", "s", false,
		"save HTML and static files to a the save directory\nor ignore to print (save directory: "+dir+")")
	createCmd.Flags().BoolVarP(&html.Save.Compress, "compress", "z", false,
		"store and compress all files into an archive when saving")
	createCmd.Flags().BoolVarP(&html.Save.OW, "overwrite", "o", false,
		"overwrite any existing files when saving")
	// meta and html related flags.
	flags := initFlags()
	keys := index(flags)
	for _, i := range keys {
		c := flags[i]
		var buf bytes.Buffer
		buf = c.initBodyFlag(buf)
		buf = c.initFlags(buf)
	}
	createCmd.Flags().BoolVarP(&html.SauceData.Use, "sauce", "", true, "use any found SAUCE metadata as HTML meta tags")
	if err := createCmd.Flags().MarkHidden("body"); err != nil {
		logs.Fatal("create mark", "body hidden", err)
	}
	if err := createCmd.Flags().MarkHidden("cache"); err != nil {
		logs.Fatal("create mark", "cache hidden", err)
	}
	createCmd.Flags().SortFlags = false
}

// initBodyFlag initializes the hidden body flag.
func (c *metaFlag) initBodyFlag(buf bytes.Buffer) bytes.Buffer {
	switch {
	case c.key == "html.body":
		fmt.Fprint(&buf, "override and inject a string into the HTML body element")
	case len(c.opts) == 0:
		fmt.Fprint(&buf, config.Tip()[c.key])
	default:
		fmt.Fprint(&buf, str.Options(config.Tip()[c.key], true, c.opts...))
	}
	return buf
}

// initFlags initializes the public facing flags.
func (c *metaFlag) initFlags(buf bytes.Buffer) bytes.Buffer {
	switch {
	case c.key == "serve":
		fmt.Fprint(&buf, "\nsupply a 0 value to use the default port, "+str.Example("-p0")+" or "+str.Example("--serve=0"))
		createCmd.Flags().UintVarP(c.i, c.name, c.short, viper.GetUint(c.key), buf.String())
	case c.strg != nil:
		createCmd.Flags().StringVarP(c.strg, c.name, c.short, viper.GetString(c.key), buf.String())
	case c.boo != nil:
		createCmd.Flags().BoolVarP(c.boo, c.name, c.short, viper.GetBool(c.key), buf.String())
	case c.i != nil:
		createCmd.Flags().UintVarP(c.i, c.name, c.short, viper.GetUint(c.key), buf.String())
	}
	return buf
}

// saveDir returns the directory the created HTML and other files will be saved to.
func saveDir() string {
	var err error
	s := viper.GetString("save-directory")
	if s == "" {
		s, err = os.Getwd()
		if err != nil {
			fmt.Printf("current working directory error: %v\n", err)
		}
	}
	return s
}

// index creates an ordered index of the meta flags.
func index(cfg map[int]metaFlag) []int {
	k := make([]int, len(cfg))
	for i := range cfg {
		k[i] = i
	}
	sort.Ints(k)
	return k
}

// initFlags initializes the create command flags and their help.
func initFlags() map[int]metaFlag {
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
	return map[int]metaFlag{
		// output
		serve: {"serve", nil, nil, &html.Port, "serve", "p", nil},
		// main tag flags
		style:  {"style.html", &html.Syntax, nil, nil, "syntax-style", "", nil},
		layout: {"html.layout", &html.Layout, nil, nil, "layout", "l", create.Layouts()},
		title:  {"html.title", &html.Title.Value, nil, nil, "title", "t", nil},
		desc:   {"html.meta.description", &html.Metadata.Description.Value, nil, nil, "meta-description", "d", nil},
		author: {"html.meta.author", &html.Metadata.Author.Value, nil, nil, "meta-author", "a", nil},
		retro:  {"html.meta.retrotxt", nil, &html.Metadata.RetroTxt, nil, "meta-retrotxt", "r", nil},
		// minor tag flags
		gen:     {"html.meta.generator", nil, &html.Metadata.Generator, nil, "meta-generator", "g", nil},
		cscheme: {"html.meta.color-scheme", &html.Metadata.ColorScheme.Value, nil, nil, "meta-color-scheme", "", nil},
		kwords:  {"html.meta.keywords", &html.Metadata.Keywords.Value, nil, nil, "meta-keywords", "", nil},
		nolang:  {"html.meta.notranslate", nil, &html.Metadata.NoTranslate, nil, "meta-notranslate", "", nil},
		refer:   {"html.meta.referrer", &html.Metadata.Referrer.Value, nil, nil, "meta-referrer", "", nil},
		bots:    {"html.meta.robots", &html.Metadata.Robots.Value, nil, nil, "meta-robots", "", nil},
		themec:  {"html.meta.theme-color", &html.Metadata.ThemeColor.Value, nil, nil, "meta-theme-color", "", nil},
		fontf:   {"html.font.family", &html.FontFamily.Value, nil, nil, "font-family", "f", nil},
		fonte:   {"html.font.embed", nil, &html.FontEmbed, nil, "font-embed", "", nil},
		// hidden flags
		body:  {"html.body", &html.Source.HiddenBody, nil, nil, "body", "b", nil},
		cache: {"html.layout.cache", nil, &html.Save.Cache, nil, "cache", "", nil},
	}
}

// parseBody is a hidden function used for debugging.
// It takes the supplied text and uses it as the content of the generated HTML <pre></pre> elements.
func parseBody(cmd *cobra.Command) {
	// hidden --body flag ignores most other args
	if body := cmd.Flags().Lookup("body"); body.Changed {
		b := []byte(body.Value.String())
		serve := cmd.Flags().Lookup("serve").Changed
		if h := serveBytes(0, serve, &b); !h {
			i, a, err := html.Create(&b)
			if err != nil {
				logs.Fatal(i, a, err)
			}
		}
		os.Exit(0)
	}
}

// parseFiles parses the flags to create the HTML document or website.
// The generated HTML and associated files will either be served, saved or printed.
func parseFiles(cmd *cobra.Command, flags convert.Flags, args ...string) {
	conv := convert.Convert{
		Flags: flags,
	}
	f, ff := sample.Flags{}, cmd.Flags().Lookup("font-family")
	serve := cmd.Flags().Lookup("serve").Changed
	for i, arg := range args {
		src, cont := staticTextfile(f, &conv, arg, ff.Changed)
		if cont {
			continue
		}
		b := createHTML(cmd, flags, &src)
		if b == nil {
			continue
		}
		// serve the HTML over HTTP?
		if h := serveBytes(i, serve, &b); !h {
			i, a, err := html.Create(&b)
			if err != nil {
				logs.Fatal(i, a, err)
			}
		}
	}
}

// parsePipe creates HTML content using the standard input (stdio) of the operating system.
func parsePipe(cmd *cobra.Command, flags convert.Flags) {
	src, err := filesystem.ReadPipe()
	if err != nil {
		logs.Fatal("create", "read stdin", err)
	}
	b := createHTML(cmd, flags, &src)
	serve := cmd.Flags().Lookup("serve").Changed
	if h := serveBytes(0, serve, &b); !h {
		i, a, err := html.Create(&b)
		if err != nil {
			logs.Fatal(i, a, err)
		}
	}
	os.Exit(0)
}

// createHTML applies a HTML template to src text.
func createHTML(cmd *cobra.Command, flags convert.Flags, src *[]byte) []byte {
	var err error
	conv := convert.Convert{
		Flags: flags,
	}
	f := sample.Flags{}
	conv.Output = convert.Output{}
	// encode and convert the source text
	if cp := cmd.Flags().Lookup("encode"); cp.Changed {
		if f.From, err = convert.Encoding(cp.Value.String()); err != nil {
			logs.Fatal("encoding not known or supported", "createHTML", err)
		}
		conv.Source.E = f.From
	}
	// obtain any appended SAUCE metadata
	appendSAUCE(src)
	// convert the source text into web friendly UTF8
	var r []rune
	if endOfFile(conv.Flags) {
		r, err = conv.Text(src)
	} else {
		r, err = conv.Dump(src)
	}
	if err != nil {
		logs.Println("convert text", "createHTML", err)
		return nil
	}
	return []byte(string(r))
}

// appendSAUCE parses any embedded SAUCE metadata.
func appendSAUCE(src *[]byte) {
	if html.SauceData.Use {
		if index := sauce.Scan(*src...); index > 0 {
			s := sauce.Parse(*src...)
			html.SauceData.Title = s.Title
			html.SauceData.Author = s.Author
			html.SauceData.Group = s.Group
			html.SauceData.Description = s.Desc
			html.SauceData.Width = uint(s.Info.Info1.Value)
			html.SauceData.Lines = uint(s.Info.Info2.Value)
		}
	}
}

// staticTextfile fetches a static text file from `/static/text`
// and uses it as the input source text.
func staticTextfile(f sample.Flags, conv *convert.Convert, arg string, changed bool) (src []byte, cont bool) {
	var err error
	if ok := sample.Valid(arg); ok {
		var p sample.File
		p, err = f.Open(arg, conv)
		if err != nil {
			logs.Println("sample", arg, err)
			return nil, true
		}
		src = create.Normalize(p.Encoding, p.Runes...)
		if changed {
			// only apply the sample font when the --font-family flag is unused
			html.FontFamily.Value = p.Font.String()
		}
	}
	// read file
	if src == nil {
		if src, err = filesystem.Read(arg); err != nil {
			logs.Fatal("file is invalid", arg, err)
		}
	}
	return src, false
}

// serveBytes hosts the HTML using an internal HTTP server.
func serveBytes(i int, changed bool, b *[]byte) bool {
	if i != 0 {
		return false
		// only ever serve the first file given to the args.
		// in the future, when handling multiple files a dynamic
		// index.html could be generated with links to each of the htmls.
	}
	if changed {
		i, a, err := html.Serve(b)
		if err != nil {
			logs.Fatal(i, a, err)
		}
		return true
	}
	return false
}
