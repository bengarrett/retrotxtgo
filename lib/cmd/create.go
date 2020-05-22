package cmd

import (
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
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
	s := string(os.PathSeparator)
	e := `  retrotxt create -n textfile.txt -t "Text file" -d "Some random text file"` +
		fmt.Sprintf("\n  retrotxt create --name ~%sDownloads%stextfile.txt --layout mini --save .%shtml", s, s, s)
	return logs.Cinf(e)
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

func init() {
	//config.InitDefaults()
	rootCmd.AddCommand(createCmd)
	// required flags
	createCmd.Flags().StringVarP(&createFileName, "name", "n", "",
		logs.Cp("text file to parse")+" (required)\n")
	// main flags
	createCmd.Flags().StringVarP(&createArgs.HTMLLayout, "layout", "l", def("create.layout"),
		"output HTML layout\noptions: "+logs.Ci(create.Layouts()))
	err := viper.BindPFlag("create.layout", createCmd.Flags().Lookup(("layout")))
	logs.ReCheck(err)
	createCmd.Flags().StringVarP(&createArgs.Styles, "syntax-style", "c", "lovelace",
		"HTML syntax highligher, use "+logs.Ci("none")+" to disable")
	createCmd.Flags().StringVarP(&pageTitle, "title", "t", def("create.title"),
		"defines the page title that is shown in a browser title bar or tab")
	err = viper.BindPFlag("create.title", createCmd.Flags().Lookup("title"))
	logs.ReCheck(err)
	createCmd.Flags().StringVarP(&metaDesc, "meta-description", "d", def("create.meta.description"),
		"a short and accurate summary of the content of the page")
	err = viper.BindPFlag("create.meta.description", createCmd.Flags().Lookup("meta-description"))
	logs.ReCheck(err)
	createCmd.Flags().StringVarP(&metaAuthor, "meta-author", "a", def("create.meta.author"),
		"defines the name of the page authors")
	err = viper.BindPFlag("create.meta.author", createCmd.Flags().Lookup("meta-author"))
	logs.ReCheck(err)
	// minor flags
	createCmd.Flags().BoolVarP(&metaGenerator, "meta-generator", "g", viper.GetBool("create.meta.generator"),
		"include the RetroTxt version and page generation date")
	err = viper.BindPFlag("create.meta.generator", createCmd.Flags().Lookup("meta-generator"))
	logs.ReCheck(err)
	createCmd.Flags().StringVar(&metaColorScheme, "meta-color-scheme", def("create.meta.color-scheme"),
		"specifies one or more color schemes with which the page is compatible")
	err = viper.BindPFlag("create.meta.color-scheme", createCmd.Flags().Lookup("meta-color-scheme"))
	logs.ReCheck(err)
	createCmd.Flags().StringVar(&metaKeywords, "meta-keywords", def("create.meta.keywords"),
		"words relevant to the page content")
	err = viper.BindPFlag("create.meta.keywords", createCmd.Flags().Lookup("meta-keywords"))
	logs.ReCheck(err)
	createCmd.Flags().StringVar(&metaReferrer, "meta-referrer", def("create.meta.referrer"),
		"controls the Referer HTTP header attached to requests sent from the page")
	err = viper.BindPFlag("create.meta.referrer", createCmd.Flags().Lookup("meta-referrer"))
	logs.ReCheck(err)
	createCmd.Flags().StringVar(&metaThemeColor, "meta-theme-color", def("create.meta.theme-color"),
		"indicates a suggested color that user agents should use to customize the display of the page")
	err = viper.BindPFlag("create.meta.theme-color", createCmd.Flags().Lookup("meta-theme-color"))
	logs.ReCheck(err)
	// output flags
	// todo: when using save-directory config setting, there is no way to stdout using flags
	// instead add an output flag with print, file|save
	createCmd.Flags().StringVarP(&saveToFiles, "save", "s", def("create.save-directory"),
		"save HTML as files to store this directory"+homeDir()+workingDir())
	err = viper.BindPFlag("create.save-directory", createCmd.Flags().Lookup("save"))
	logs.ReCheck(err)
	createCmd.Flags().BoolVarP(&createArgs.ServerFiles, "server", "p", false,
		"serve HTML over an internal web server")
	createCmd.Flags().IntVar(&createArgs.ServerPort, "port", viper.GetInt("create.server-port"),
		"port which the internet web server will listen")
	err = viper.BindPFlag("create.server-port", createCmd.Flags().Lookup("port"))
	logs.ReCheck(err)
	// hidden flags
	createCmd.Flags().StringVarP(&preText, "body", "b", "",
		"override and inject string content into the body element")
	// flag options
	err = createCmd.Flags().MarkHidden("body")
	logs.ReCheck(err)
	createCmd.Flags().SortFlags = false
}

func def(key string) string {
	return viper.GetString(key)
}

func homeDir() string {
	s := "\n" + logs.Ci("--save ~") + " saves to the home or user directory"
	d, err := os.UserHomeDir()
	if err != nil {
		return s
	}
	return s + " at " + logs.Cf(d)
}

func workingDir() string {
	s := "\n" + logs.Ci("--save .") + " saves to the current working directory"
	d, err := os.Getwd()
	if err != nil {
		return s
	}
	return s + " at " + logs.Cf(d)
}
