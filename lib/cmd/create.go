package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/bengarrett/retrotxtgo/internal/pack"
	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var html create.Args

// createCmd makes create usage examples
var exampleCmd = func() string {
	var b bytes.Buffer
	s := string(os.PathSeparator)
	fmt.Fprint(&b, `  retrotxt create file.txt -t "A text file" -d "Some text goes here"`)
	fmt.Fprint(&b, "\n  retrotxt create file1.txt file2.asc --save")
	fmt.Fprintf(&b, "\n  retrotxt create ~%sDownloads%sfile.txt --archive", s, s)
	fmt.Fprint(&b, "\n  retrotxt create file.txt --serve=8080")
	fmt.Fprint(&b, "\n  cat file.txt | retrotxt create")
	return str.Cinf(b.String())
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create [filenames]",
	Aliases: []string{"c", "html"},
	Short:   "Create a HTML document from a text file",
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		// piped input from other programs
		if filesystem.IsPipe() {
			b, err := filesystem.ReadPipe()
			logs.Check("piped.create", err)
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
		checkUse(cmd, args)
		var b []byte
		for i, arg := range args {
			if b = createPackage(arg); b == nil {
				var err error
				b, err = filesystem.Read(arg)
				logs.ChkErr(logs.Err{Issue: "file is invalid", Arg: arg, Msg: err})
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

func createPackage(name string) (b []byte) {
	var s = strings.ToLower(name)
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return nil
	}
	pkg, exist := internalPacks[s]
	if !exist {
		return nil
	}
	b = pack.Get(pkg.name)
	if b == nil {
		return nil
	}
	return b
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
	var metaCfg = map[int]metaFlag{
		// output
		1: {"serve", nil, nil, &html.Port, "serve", "p", nil},
		// main tag flags
		3: {"html.layout", &html.Layout, nil, nil, "layout", "l", create.Layouts()},
		4: {"style.html", &html.Syntax, nil, nil, "syntax-style", "c", nil},
		5: {"html.title", &html.Title, nil, nil, "title", "t", nil},
		6: {"html.meta.description", &html.Desc, nil, nil, "meta-description", "d", nil},
		7: {"html.meta.author", &html.Author, nil, nil, "meta-author", "a", nil},
		// minor tag flags
		8:  {"html.meta.generator", nil, &html.Generator, nil, "meta-generator", "g", nil},
		9:  {"html.meta.color-scheme", &html.Author, nil, nil, "meta-color-scheme", "", nil},
		10: {"html.meta.keywords", &html.Keys, nil, nil, "meta-keywords", "", nil},
		11: {"html.meta.notranslate", nil, &html.NoTranslate, nil, "meta-notranslate", "", nil},
		12: {"html.meta.referrer", &html.Ref, nil, nil, "meta-referrer", "", nil},
		13: {"html.meta.robots", &html.Robots, nil, nil, "meta-robots", "", nil},
		14: {"html.meta.theme-color", &html.Scheme, nil, nil, "meta-theme-color", "", nil},
		15: {"html.font.family", &html.FontFamily, nil, nil, "font.family", "", nil},
		17: {"html.font.embed", nil, &html.FontEmbed, nil, "font.embed", "", nil},
		// hidden flags
		0: {"html.body", &html.Body, nil, nil, "body", "b", nil},
	}
	// create an ordered index for the flags
	var keys []int
	for k := range metaCfg {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// output flags
	createCmd.Flags().StringVarP(&html.Enc, "encode", "e", "", "text encoding of the named text file\nwhen ignored, UTF8 encoding will be automatically detected\notherwise encode will assume default (default CP437)\nsee a list of encode values "+str.Example("retrotxt view codepages")+"\n")
	createCmd.Flags().BoolVarP(&html.SaveToFile, "save", "s", false, "save HTML to a file or ignore to print output")
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
			fmt.Fprint(&buf, config.Hints[c.key])
		default:
			fmt.Fprint(&buf, str.Options(config.Hints[c.key], c.opts, true))
		}
		switch {
		case c.key == "serve":
			fmt.Fprint(&buf, "\nsupply a 0 value to use the default, "+str.Example("-p0")+" or "+str.Example("--serve=0"))
			createCmd.Flags().UintVarP(c.i, c.name, c.short, viper.GetUint(c.key), buf.String())
		case c.strg != nil:
			createCmd.Flags().StringVarP(c.strg, c.name, c.short, viper.GetString(c.key), buf.String())
		case c.boo != nil:
			createCmd.Flags().BoolVarP(c.boo, c.name, c.short, viper.GetBool(c.key), buf.String())
		case c.i != nil:
			createCmd.Flags().UintVarP(c.i, c.name, c.short, viper.GetUint(c.key), buf.String())
		}
	}
	err = createCmd.Flags().MarkHidden("body")
	logs.Check("create mark body hidden", err)
	createCmd.Flags().SortFlags = false
}
