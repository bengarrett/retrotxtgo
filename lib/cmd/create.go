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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// create command flag
var (
	createFileName  string
	metaAuthor      string
	metaColorScheme string
	metaDesc        string
	metaGenerator   bool
	metaKeywords    string
	metaReferrer    string
	metaThemeColor  string
	pageTitle       string
	preText         string
	saveToFiles     string
)

var createArgs = create.Args{}

// createCmd makes create usage examples
var exampleCmd = func() string {
	var b bytes.Buffer
	s := string(os.PathSeparator)
	fmt.Fprint(&b, `  retrotxt create -n textfile.txt -t "Text file" -d "Some random text file"`)
	fmt.Fprintf(&b, "\n  retrotxt create --name ~%sDownloads%stextfile.txt --layout mini --save .%shtml", s, s, s)
	return str.Cinf(b.String())
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a HTML document from a text file",
	Example: exampleCmd(),
	Run: func(cmd *cobra.Command, args []string) {
		createArgs.HTMLLayout = viper.GetString("create.layout")
		var data []byte
		var err error
		// --body is a hidden flag to test this cmd without providing a file
		b := cmd.Flags().Lookup("body")
		switch b.Changed {
		case true:
			data = []byte(b.Value.String())
		default:
			if createFileName == "" {
				if cmd.Flags().NFlag() == 0 {
					fmt.Printf("%s\n\n", cmd.Short)
					err = cmd.Usage()
					logs.Check("create usage", err)
					os.Exit(0)
				}
				err = cmd.Usage()
				logs.ReCheck(err)
				logs.FileMissingErr()
			}
			data, err = filesystem.Read(createFileName)
			logs.ChkErr(logs.Err{Issue: "file is invalid", Arg: createFileName, Msg: err})
		}
		// check for a --save flag to save to files
		// otherwise output is sent to stdout
		s := cmd.Flags().Lookup("save")
		createArgs.Save(data, s.Value.String(), s.Changed)
		// check for a --server flag to serve the HTML
		if createArgs.ServerFiles {
			createArgs.Serve(data)
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
		// main flags
		0: {"create.layout", &createArgs.HTMLLayout, nil, nil, "layout", "l", create.Options()},
		1: {"style.html", &createArgs.Styles, nil, nil, "syntax-style", "c", nil},
		2: {"create.title", &pageTitle, nil, nil, "title", "t", nil},
		3: {"create.meta.description", &metaDesc, nil, nil, "meta-description", "d", nil},
		4: {"create.meta.author", &metaDesc, nil, nil, "meta-author", "a", nil},
		// minor flags
		5: {"create.meta.generator", nil, &metaGenerator, nil, "meta-generator", "g", nil},
		6: {"create.meta.color-scheme", &metaColorScheme, nil, nil, "meta-color-scheme", "", nil},
		7: {"create.meta.keywords", &metaKeywords, nil, nil, "meta-keywords", "", nil},
		8: {"create.meta.referrer", &metaReferrer, nil, nil, "meta-referrer", "", nil},
		9: {"create.meta.theme-color", &metaThemeColor, nil, nil, "meta-theme-color", "", nil},
		// output
		10: {"create.save-directory", &saveToFiles, nil, nil, "save", "s", nil},
		11: {"create.server", nil, &createArgs.ServerFiles, nil, "server", "p", nil},
		12: {"create.server-port", nil, nil, &createArgs.ServerPort, "port", "", nil},
		// hidden flags
		13: {"create.body", &preText, nil, nil, "body", "b", nil},
		// TODO: add sample flag to generate the RetroTxt ANSI/ascii logo as HTML?
	}
	// create an ordered index for the flags
	var keys []int
	for k := range metaCfg {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	// required flags
	createCmd.Flags().StringVarP(&createFileName, "name", "n", "",
		str.Required("text file to parse")+"\n")
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
