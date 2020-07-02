package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"retrotxt.com/retrotxt/lib/config"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

type rootFlags struct {
	config string
}

var rootFlag = rootFlags{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "retrotxt",
	Short: "RetroTxt is the tool that turns ANSI, ASCII, NFO text into browser ready HTML",
	Long: `Turn many pieces of ANSI text art and ASCII/NFO plain text into HTML5 text
using RetroTxt. The operating system agnostic tool that takes retro text
files and stylises them into a more pleasing, useful format to view and
copy in a web browser.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.SilenceErrors = true // set to false to debug
	if err := rootCmd.Execute(); err != nil {
		if len(os.Args) < 2 {
			logs.Check("rootcmd.usage", rootCmd.Usage())
		}
		rootErr := logs.CmdErr{Args: os.Args[1:], Err: err}
		fmt.Println(rootErr.Error().String())
	}
}

func init() {
	// OnInitialize will not run if there is no provided command.
	cobra.OnInitialize(initConfig)
	// TODO: get viper to flag file autocomplete
	rootCmd.PersistentFlags().StringVar(&rootFlag.config, "config", "",
		"optional config file location")
}

func defaultConfig() string {
	named := viper.GetViper().ConfigFileUsed()
	if named == "" {
		named = config.Path()
	}
	return named
}

// initConfig reads in the config file and ENV variables if set.
// this does not run when rootCmd is in use.
func initConfig() {
	// read in environment variables
	viper.SetEnvPrefix("env")
	viper.AutomaticEnv()
	// configuration file
	if err := config.SetConfig(rootFlag.config); err != nil {
		logs.Check(fmt.Sprintf("config file %q", viper.ConfigFileUsed()), err)
		os.Exit(1)
	}
}

// checkUsage will print the help and exit when no arguments are supplied.
func checkUse(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		err := cmd.Help()
		logs.Check("cmd.help", err)
		os.Exit(0)
	}
}

func flagEncode(p *string, cc *cobra.Command) {
	cc.Flags().StringVarP(p, "encode", "e", "",
		`characture encoding used by the filenames
when ignored, UTF8 encoding is detected
if that fails the default is used (default CP437)
see the list of encode values `+str.Example("retrotxt list codepages")+"\n")
}

func flagTab(p *bool, cc *cobra.Command) {
	cc.Flags().BoolVar(p, "tab", true, "parse horizontal tab control characters")
}

func flagWidth(p *int, cc *cobra.Command) {
	cc.Flags().IntVarP(p, "width", "w", viewFlag.width, "document column character width")
}

func flagControls(p *[]string, cc *cobra.Command) {
	cc.Flags().StringSliceVarP(p, "controls", "x", []string{}, "always use these control codes\n"+
		"  tab    horizontal tab\n"+
		"  bell   bell or terminal alert\n"+
		"  cr     carriage return\n"+
		"  lf     line feed\n"+
		"  bs backspace, del delete character, esc escape character, ff formfeed, vt vertical tab\n"+
		"separate multiple controls with commas\n"+
		str.Example("-x tab,bell")+" or "+str.Example("--controls=tab,bell")+"\n")
}

type internalPack struct {
	// choices: d convert.Dump, t convert.Text (default when blank)
	convert string
	// default character encoding for the packed data
	encoding string
	// package name used in internal/pack/blob.go
	name string
	// package description
	description string
}

var internalPacks = map[string]internalPack{
	"437.cr":        {"d", "cp437", "text/cp437-cr.txt", "CP-437 all characters test using CR (carriage return)"}, //
	"437.crlf":      {"d", "cp437", "text/cp437-crlf.txt", "CP-437 all characters test using Windows newline"},    //
	"437.lf":        {"d", "cp437", "text/cp437-lf.txt", "CP-437 all characters test using LF (line feed)"},       //
	"865":           {"", "ibm865", "text/cp865.txt", "CP-865 and CP-860 Nordic test"},                            //
	"1252":          {"", "cp1252", "text/cp1252.txt", "Windows-1252 English test"},                               //
	"ascii":         {"", "cp437", "text/retrotxt.asc", "RetroTxt ASCII logos"},                                   //
	"ansi":          {"", "cp437", "text/retrotxt.ans", "RetroTxt 256 color ANSI logo"},                           //
	"ansi.aix":      {"", "cp437", "text/ansi-aixterm.ans", "IBM AIX terminal colours"},                           //
	"ansi.blank":    {"", "cp437", "text/ansi-blank.ans", "Empty file test"},                                      //
	"ansi.cp":       {"", "cp437", "text/ansi-cp.ans", "ANSI cursor position tests"},                              //
	"ansi.cpf":      {"", "cp437", "text/ansi-cpf.ans", "ANSI cursor forward tests"},                              //
	"ansi.hvp":      {"", "cp437", "text/ansi-hvp.ans", "ANSI horizontal and vertical cursor positioning"},        //
	"ansi.proof":    {"", "cp437", "text/ansi-proof.ans", "ANSI formatting proof sheet"},                          //
	"ansi.rgb":      {"", "cp437", "text/ansi-rgb.ans", "ANSI RGB 24-bit color sheet"},                            //
	"ansi.setmodes": {"", "cp437", "text/ansi-setmodes.ans", "MS-DOS ANSI.SYS Set Mode examples"},                 //
	"iso-1":         {"", "1", "text/iso-8859-1.txt", "ISO 8859-1 select characters"},                             //
	"iso-15":        {"", "15", "text/iso-8859-15.txt", "ISO 8859-15 select characters"},                          //
	"sauce":         {"", "", "text/sauce.txt", "SAUCE metadata test"},                                            // TODO
	"shiftjis":      {"", "shift-jis", "text/shiftjis.txt", "Shift-JIS and Mona font test"},                       // TODO
	"us-ascii":      {"", "cp1252", "text/us-ascii.txt", "US-ASCII controls test"},                                //
	"utf8":          {"", "", "text/utf-8.txt", "UTF-8 test with no Byte Order Mark"},                             //
	"utf8.bom":      {"", "", "text/utf-8-bom.txt", "UTF-8 test with a Byte Order Mark"},                          //
	"utf16.be":      {"", "utf-16be", "text/utf-16-be.txt", "UTF-16 Big Endian test"},                             //
	"utf16.le":      {"", "utf-16le", "text/utf-16-le.txt", "UTF-16 Little Endian test"},                          //
}

func examples() *bytes.Buffer {
	keys := make([]string, 0, len(internalPacks))
	for k := range internalPacks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buf bytes.Buffer
	var flags uint = 0 //tabwriter.AlignRight | tabwriter.Debug
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', flags)
	fmt.Fprintln(w, str.Cp(" Packaged example text and ANSI files to test and play with RetroTxt"))
	fmt.Fprintln(w, strings.Repeat("-", 69))
	for _, k := range keys {
		fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t", k, internalPacks[k].description))
	}
	fmt.Fprintln(w, "\nAny of these packaged examples will work with both the", str.Example("create"), "and", str.Example("view"), "commands")
	fmt.Fprintln(w, "\n"+str.Example(" retrotxt view 1252"), "will print the Windows-1252 English test to the terminal")
	fmt.Fprintln(w, str.Example(" retrotxt view 1252 > file.txt"), "will convert and save the Windows-1252 English test to UTF-8 encoding")
	fmt.Fprintln(w, str.Example(" retrotxt save 1252 > file.txt"), "will save the Windows-1252 English test with its original encoding") // TODO
	fmt.Fprintln(w, str.Example(" retrotxt save 1252 | retrotxt info"), "displays statistics and information")                           // TODO
	fmt.Fprintln(w, str.Example(" retrotxt create 1252"), "creates a HTML document from the Windows-1252 English test")
	fmt.Fprintln(w, str.Example(" retrotxt create 1252 -p0"), "serves the Windows-1252 English test over a local web server")
	fmt.Fprintln(w, "\nMultiple examples are supported")
	fmt.Fprintln(w, str.Example(" retrotxt view ansi ascii ansi.rgb"))
	w.Flush()
	return &buf
}
