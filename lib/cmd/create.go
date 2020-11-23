package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"retrotxt.com/retrotxt/internal/pack"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/create"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

var html create.Args

// createCmd makes create usage examples.
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
		monitorFlags(cmd)
		// piped input from other programs
		if filesystem.IsPipe() {
			b, err := filesystem.ReadPipe()
			if err != nil {
				logs.Fatal("create", "read stdin", err)
			}
			if h := htmlServe(0, cmd, &b); !h {
				html.Create(&b)
			}
			os.Exit(0)
		}
		// hidden --body flag that ignores all args
		if body := cmd.Flags().Lookup("body"); body.Changed {
			b := []byte(body.Value.String())
			if h := htmlServe(0, cmd, &b); !h {
				html.Create(&b)
			}
			os.Exit(0)
		}
		// user arguments
		checkUse(cmd, args...)
		for i, arg := range args {
			html.SourceName = arg
			enc, font, b := createPackage(arg)
			if b == nil {
				var err error
				b, err = filesystem.Read(arg)
				if err != nil {
					logs.Fatal("file is invalid", arg, err)
				}
			}
			if cp := cmd.Flags().Lookup("encode"); !cp.Changed {
				html.Encoding = enc
			}
			if ff := cmd.Flags().Lookup("font-family"); !ff.Changed {
				if font == "" {
					font = "vga"
				}
				html.FontFamily.Value = font
			}
			if h := htmlServe(i, cmd, &b); !h {
				html.Create(&b)
			}
		}
	},
}

func htmlServe(i int, cmd *cobra.Command, b *[]byte) bool {
	if i != 0 {
		return false
		// only ever serve the first file given to the args.
		// in the future, when handling multiple files an index.html
		// could be generated with links to each of the htmls.
	}
	if serve := cmd.Flags().Lookup("serve"); serve.Changed {
		html.Serve(b)
		return true
	}
	return false
}

func createPackage(name string) (enc, font string, b []byte) {
	var s = strings.ToLower(name)
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return "", "", nil
	}
	pkg, exist := internalPacks[s]
	if !exist {
		return "", "", nil
	}
	b = pack.Get(pkg.name)
	if b == nil {
		return "", "", nil
	}
	fmt.Println("pack:", name, "enc", pkg.encoding, "font", pkg.font)
	return pkg.encoding, pkg.font, b
}

func monitorFlags(cmd *cobra.Command) {
	var changed = func(key string) bool {
		l := cmd.Flags().Lookup(key)
		if l == nil {
			return false
		}
		return l.Changed
	}
	// monitor string flag changes to allow three user states.
	// 1) flag not changed so use viper default.
	// 2) flag with new value to overwrite viper default.
	// 3) blank flag value to overwrite viper default with an empty/disable value.
	html.FontFamily.Flag = changed("font-family")
	html.Metadata.Author.Flag = changed("meta-author")
	html.Metadata.ColorScheme.Flag = changed("meta-color-scheme")
	html.Metadata.Description.Flag = changed("meta-description")
	html.Metadata.Keywords.Flag = changed("meta-keywords")
	html.Metadata.Referrer.Flag = changed("meta-referrer")
	html.Metadata.Robots.Flag = changed("meta-robots")
	html.Metadata.ThemeColor.Flag = changed("meta-theme-color")
	html.Title.Flag = changed("title")
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

func init() {
	var err error
	rootCmd.AddCommand(createCmd)
	// config must be initialized before getting saved default values
	initConfig()
	// init flags and their usage
	var metaCfg = metaConfig()
	// create an ordered index for the flags
	var keys = make([]int, len(metaCfg))
	for i := range metaCfg {
		keys[i] = i
	}
	sort.Ints(keys)
	// output flags
	flagEncode(&html.Encoding, createCmd)
	createCmd.Flags().BoolVarP(&html.SaveToFile, "save", "s", false,
		`save HTML and static files to a the save directory
or ignore to print (save directory: `+viper.GetString("save-directory")+")")
	createCmd.Flags().BoolVarP(&html.Compress, "compress", "c", false, "store and compress all files into an archive when saving")
	createCmd.Flags().BoolVarP(&html.OW, "overwrite", "o", false, "overwrite any existing files when saving")
	// html flags, the key int value must be used as the index
	// rather than the loop count, otherwise flags might be skipped
	for _, i := range keys {
		c := metaCfg[i]
		var buf bytes.Buffer
		switch {
		case c.key == "html.body":
			fmt.Fprint(&buf, "override and inject a string into the HTML body element")
		case len(c.opts) == 0:
			fmt.Fprint(&buf, config.Tip()[c.key])
		default:
			fmt.Fprint(&buf, str.Options(config.Tip()[c.key], true, c.opts...))
		}
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
	}
	if err = createCmd.Flags().MarkHidden("body"); err != nil {
		logs.Fatal("create mark", "body hidden", err)
	}
	if err = createCmd.Flags().MarkHidden("cache"); err != nil {
		logs.Fatal("create mark", "cache hidden", err)
	}
	createCmd.Flags().SortFlags = false
}

func metaConfig() map[int]metaFlag {
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
	// init flags and their usage
	return map[int]metaFlag{
		// output
		serve: {"serve", nil, nil, &html.Port, "serve", "p", nil},
		// main tag flags
		layout: {"html.layout", &html.Layout, nil, nil, "layout", "l", create.Layouts()},
		style:  {"style.html", &html.Syntax, nil, nil, "syntax-style", "", nil},
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
		body:  {"html.body", &html.Body, nil, nil, "body", "b", nil},
		cache: {"html.layout.cache", nil, &html.Cache, nil, "cache", "", nil},
	}
}
