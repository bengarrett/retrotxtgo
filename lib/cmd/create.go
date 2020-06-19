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
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var htmlArgs create.Args

// createCmd makes create usage examples
var exampleCmd = func() string {
	var b bytes.Buffer
	s := string(os.PathSeparator)
	fmt.Fprint(&b, "  retrotxt create -n=ascii\n")
	fmt.Fprint(&b, `  retrotxt create -n=file.txt -t "Textfile" -d "Some text"`)
	fmt.Fprintf(&b, "\n  retrotxt create --name ~%sDownloads%stextfile.txt --save", s, s)
	return str.Cinf(b.String())
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create [filenames]",
	Aliases: []string{"c", "html"},
	Short:   "Create a HTML document from a text file",
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		// handle hidden --body flag that ignores all args
		if body := cmd.Flags().Lookup("body"); body.Changed {
			b := []byte(body.Value.String())
			createHTML(0, b)
			os.Exit(0)
		}
		checkUse(cmd, args)
		var b []byte
		for i, arg := range args {
			if b = createPackage(arg); b == nil {
				var err error
				b, err = filesystem.Read(arg)
				logs.ChkErr(logs.Err{Issue: "file is invalid", Arg: arg, Msg: err})
			}
			createHTML(i, b)
		}
	},
}

func createHTML(i int, b []byte) {
	htmlArgs.Cmd(b, "")
	// only serve the first file
	if i == 0 && htmlArgs.HTTP {
		htmlArgs.Serve(&b)
	}
}

func createPackage(name string) (b []byte) {
	var s = strings.ToLower(name)
	if _, err := os.Stat(s); !os.IsNotExist(err) {
		return nil
	}
	pkg, exist := internalPacks[s]
	println(fmt.Sprintf("%+v", pkg), exist)
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
		1: {"html.server", nil, &htmlArgs.HTTP, nil, "server", "p", nil},
		2: {"html.server-port", nil, nil, &htmlArgs.Port, "port", "", nil},
		// main tag flags
		3: {"html.layout", &htmlArgs.Layout, nil, nil, "layout", "l", create.Layouts()},
		4: {"style.html", &htmlArgs.Syntax, nil, nil, "syntax-style", "c", nil},
		5: {"html.title", &htmlArgs.Title, nil, nil, "title", "t", nil},
		6: {"html.meta.description", &htmlArgs.Desc, nil, nil, "meta-description", "d", nil},
		7: {"html.meta.author", &htmlArgs.Author, nil, nil, "meta-author", "a", nil},
		// minor tag flags
		8:  {"html.meta.generator", nil, &htmlArgs.Generator, nil, "meta-generator", "g", nil},
		9:  {"html.meta.color-scheme", &htmlArgs.Author, nil, nil, "meta-color-scheme", "", nil},
		10: {"html.meta.keywords", &htmlArgs.Keys, nil, nil, "meta-keywords", "", nil},
		11: {"html.meta.notranslate", nil, &htmlArgs.NoTranslate, nil, "meta-notranslate", "", nil},
		12: {"html.meta.referrer", &htmlArgs.Ref, nil, nil, "meta-referrer", "", nil},
		13: {"html.meta.robots", &htmlArgs.Robots, nil, nil, "meta-robots", "", nil},
		14: {"html.meta.theme-color", &htmlArgs.Scheme, nil, nil, "meta-theme-color", "", nil},
		// hidden flags
		0: {"html.body", &htmlArgs.Body, nil, nil, "body", "b", nil},
	}
	// create an ordered index for the flags
	var keys []int
	for k := range metaCfg {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// output flags
	createCmd.Flags().StringVarP(&htmlArgs.Enc, "encode", "e", "", "text encoding of the named text file\nwhen ignored, UTF8 encoding will be automatically detected\notherwise encode will assume default (default CP437)\nsee a list of encode values "+str.Example("retrotxt view codepages")+"\n")
	createCmd.Flags().BoolVarP(&htmlArgs.SaveToFile, "save", "s", false, "save HTML to a file or ignore to print output")
	createCmd.Flags().BoolVarP(&htmlArgs.OW, "overwrite", "o", false, "overwrite any existing files when saving\n")
	// html flags
	for i := range keys {
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
		case c.key == "create.server":
			createCmd.Flags().BoolVarP(c.boo, c.name, c.short, false, "serve HTML over an insecure local web server\n"+color.Warn.Sprint("please do not use on a production environment"))
		case c.strg != nil:
			createCmd.Flags().StringVarP(c.strg, c.name, c.short, viper.GetString(c.key), buf.String())
		case c.boo != nil:
			createCmd.Flags().BoolVarP(c.boo, c.name, c.short, viper.GetBool(c.key), buf.String())
		case c.i != nil:
			createCmd.Flags().UintVar(c.i, c.name, viper.GetUint(c.key), buf.String())
		}
	}
	err = createCmd.Flags().MarkHidden("body")
	logs.Check("create mark body hidden", err)
	createCmd.Flags().SortFlags = false
}
