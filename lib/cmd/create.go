package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/samples"
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
	Use:     "create",
	Short:   "Create a HTML document from a text file",
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			data []byte
			err  error
			body = cmd.Flags().Lookup("body")
		)
		switch body.Changed {
		case true: // handle hidden --body flag
			data = []byte(body.Value.String())
		default:
			switch htmlArgs.Src {
			case "ascii":
				// internal example
				data, err = samples.Base64Decode(samples.LogoASCII)
				logs.ChkErr(logs.Err{Issue: "logoascii is invalid", Arg: htmlArgs.Src, Msg: err})
			case "":
				// no input (show help & exit)
				if cmd.Flags().NFlag() == 0 {
					fmt.Printf("%s\n\n", cmd.Short)
					err = cmd.Usage()
					logs.Check("create usage", err)
					os.Exit(0)
				}
				// show help and exit with no input otherwise a blank template will be shown
				err = cmd.Usage()
				logs.ReCheck(err)
				logs.FileMissingErr()
			default:
				data, err = filesystem.Read(htmlArgs.Src)
				logs.ChkErr(logs.Err{Issue: "file is invalid", Arg: htmlArgs.Src, Msg: err})
			}
		}
		htmlArgs.Cmd(data, "") // value should = arguments
		// check for a --server flag to serve the HTML
		if htmlArgs.HTTP {
			htmlArgs.Serve(&data)
		}
	},
}

type metaFlag struct {
	key   string   // configuration name
	strg  *string  // StringVarP(p) argument value
	boo   *bool    // BoolVarP(p) argument value
	i     *int     // IntVar(p) argument value
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
		1: {"create.server", nil, &htmlArgs.HTTP, nil, "server", "p", nil},
		2: {"create.server-port", nil, nil, &htmlArgs.Port, "port", "", nil},
		// main tag flags
		3: {"create.layout", &htmlArgs.Layout, nil, nil, "layout", "l", create.Options()},
		4: {"style.html", &htmlArgs.Syntax, nil, nil, "syntax-style", "c", nil},
		5: {"create.title", &htmlArgs.Title, nil, nil, "title", "t", nil},
		6: {"create.meta.description", &htmlArgs.Desc, nil, nil, "meta-description", "d", nil},
		7: {"create.meta.author", &htmlArgs.Author, nil, nil, "meta-author", "a", nil},
		// minor tag flags
		8:  {"create.meta.generator", nil, &htmlArgs.Generator, nil, "meta-generator", "g", nil},
		9:  {"create.meta.color-scheme", &htmlArgs.Author, nil, nil, "meta-color-scheme", "", nil},
		10: {"create.meta.keywords", &htmlArgs.Keys, nil, nil, "meta-keywords", "", nil},
		11: {"create.meta.referrer", &htmlArgs.Ref, nil, nil, "meta-referrer", "", nil},
		12: {"create.meta.theme-color", &htmlArgs.Scheme, nil, nil, "meta-theme-color", "", nil},
		// hidden flags
		0: {"create.body", &htmlArgs.Body, nil, nil, "body", "b", nil},
	}
	// create an ordered index for the flags
	var keys []int
	for k := range metaCfg {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// required flags
	createCmd.Flags().StringVarP(&htmlArgs.Src, "name", "n", "",
		str.Required("text file to parse")+"\n  run a built-in example "+str.Example("retrotxt create ascii")+"\n")
	createCmd.Flags().BoolVarP(&htmlArgs.SaveToFile, "save", "s", false, "save HTML to a file or ignore to print output")
	createCmd.Flags().BoolVarP(&htmlArgs.OW, "overwrite", "o", false, "overwrite any existing files when saving\n")
	// generate flags
	for i := range keys {
		c := metaCfg[i]
		var buf bytes.Buffer
		switch {
		case c.key == "create.body":
			fmt.Fprint(&buf, "override and inject a string into the HTML body element")
		case len(c.opts) == 0:
			fmt.Fprint(&buf, config.Hints[c.key])
		default:
			fmt.Fprint(&buf, str.Options(config.Hints[c.key], c.opts, true))
		}
		switch {
		case c.key == "create.server":
			createCmd.Flags().BoolVarP(c.boo, c.name, c.short, false, "serve HTML over an internal web server")
		case c.strg != nil:
			createCmd.Flags().StringVarP(c.strg, c.name, c.short, viper.GetString(c.key), buf.String())
		case c.boo != nil:
			createCmd.Flags().BoolVarP(c.boo, c.name, c.short, viper.GetBool(c.key), buf.String())
		case c.i != nil:
			createCmd.Flags().IntVar(c.i, c.name, viper.GetInt(c.key), buf.String())
		}
	}
	err = createCmd.Flags().MarkHidden("body")
	logs.ReCheck(err)
	createCmd.Flags().SortFlags = false
}
