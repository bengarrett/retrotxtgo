package cmd

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/lib/codepage"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	humanize "github.com/labstack/gommon/bytes"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/gookit/color.v1"
)

const infoFormats = "color, json, json.min, text, xml"

// Detail of a file
type Detail struct {
	Bytes     int64
	CharCount int
	Name      string
	MD5       string
	Mime      string
	Modified  time.Time
	SHA256    string
	Slug      string
	Size      string
	Utf8      bool
}

var (
	infoFmt  string
	fileName string
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Information on a text file",
	Run: func(cmd *cobra.Command, args []string) {
		// only show Usage() with no errors if no flags .NFlags() are set
		if fileName == "" {
			if cmd.Flags().NFlag() == 0 {
				fmt.Printf("%s\n\n", cmd.Short)
				err := cmd.Usage()
				CheckErr(err)
				os.Exit(0)
			}
			err := cmd.Usage()
			CheckErr(err)
			FileMissingErr()
		}
		f, err := details(fileName)
		Check(ErrorFmt{"file is invalid", fileName, err})
		CheckFlag(f.infoSwitch(viper.GetString("info.format")))
	},
}

func init() {
	InitDefaults()
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().StringVarP(&fileName, "name", "n", "", cp("text file to analyse")+" (required)\n")
	infoCmd.Flags().StringVarP(&infoFmt, "format", "f", viper.GetString("info.format"), "output format \noptions: "+ci(infoFormats))
	err := viper.BindPFlag("info.format", infoCmd.Flags().Lookup("format"))
	CheckErr(err)
	infoCmd.Flags().SortFlags = false
}

func infoJSON(indent bool, f Detail) (js []byte) {
	var err error
	switch indent {
	case true:
		js, err = json.MarshalIndent(f, "", "    ")
	default:
		js, err = json.Marshal(f)
	}
	Check(ErrorFmt{"could not create", "json", err})
	return js
}

func (f Detail) infoSwitch(format string) (err ErrorFmt) {
	switch format {
	case "color", "c":
		fmt.Printf("%s", infoText(true, f))
	case "json", "j":
		fmt.Printf("%s\n", infoJSON(true, f))
	case "json.min", "jm":
		fmt.Printf("%s\n", infoJSON(false, f))
	case "text":
		fmt.Printf("%s", infoText(false, f))
	case "xml", "x":
		fmt.Printf("%s\n", infoXML(f))
	default:
		err = ErrorFmt{"format", infoFmt, fmt.Errorf(infoFormats)}
	}
	return err
}

func infoText(c bool, f Detail) string {
	color.Enable = c
	var info = func(t string) string {
		return cinf(fmt.Sprintf("%s\t", t))
	}
	var hr = func() string {
		return fmt.Sprintf("\t%s\n", cf(strings.Repeat("\u2015", 26)))
	}
	var data = []struct {
		d, v string
	}{
		{d: "filename", v: f.Name},
		{d: "UTF-8", v: fmt.Sprintf("%v", f.Utf8)},
		{d: "characters", v: fmt.Sprintf("%v", f.CharCount)},
		{d: "size", v: f.Size},
		{d: "modified", v: fmt.Sprintf("%v", f.Modified.Format(FileDate))},
		{d: "MD5 checksum", v: f.MD5},
		{d: "SHA256 checksum", v: f.SHA256},
		{d: "MIME type", v: f.Mime},
		{d: "slug", v: f.Slug},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, hr())
	for _, x := range data {
		fmt.Fprintf(w, "\t %s\t  %s\n", x.d, info(x.v))
	}
	fmt.Fprint(w, hr())
	w.Flush()
	return buf.String()
}

func infoXML(d Detail) []byte {
	type xmldetail struct {
		XMLName   xml.Name  `xml:"file"`
		ID        string    `xml:"id,attr"`
		Name      string    `xml:"name"`
		Mime      string    `xml:"content>mime"`
		Utf8      bool      `xml:"content>utf8"`
		Bytes     int64     `xml:"size>bytes"`
		Size      string    `xml:"size>value"`
		CharCount int       `xml:"size>character-count"`
		MD5       string    `xml:"checksum>md5"`
		SHA256    string    `xml:"checksum>sha256"`
		Modified  time.Time `xml:"modified"`
	}
	x := xmldetail{
		Bytes:     d.Bytes,
		CharCount: d.CharCount,
		ID:        d.Slug,
		MD5:       d.MD5,
		Mime:      d.Mime,
		Modified:  d.Modified,
		Name:      d.Name,
		SHA256:    d.SHA256,
		Size:      d.Size,
		Utf8:      d.Utf8,
	}
	xmlData, err := xml.MarshalIndent(x, "", "\t")
	Check(ErrorFmt{"could not create", "xml", err})
	return xmlData
}

func details(name string) (d Detail, err error) {
	// Get the file details
	stat, err := os.Stat(name)
	if err != nil {
		return d, err
	}
	// Read file content
	data, err := filesystem.ReadAllBytes(name)
	if err != nil {
		return d, err
	}
	return parse(data, stat)
}

func parse(data []byte, stat os.FileInfo) (d Detail, err error) {
	md5sum := md5.Sum(data)
	sha256 := sha256.Sum256(data)
	mime := mimesniffer.Sniff(data)
	// create a table of data
	d.Bytes = stat.Size()
	d.CharCount = runewidth.StringWidth(string(data))
	d.Name = stat.Name()
	d.MD5 = fmt.Sprintf("%x", md5sum)
	d.Modified = stat.ModTime()
	d.Slug = slugify.Slugify(stat.Name())
	d.SHA256 = fmt.Sprintf("%x", sha256)
	d.Utf8 = codepage.UTF8(data)
	if stat.Size() < 1000 {
		d.Size = fmt.Sprintf("%v bytes", stat.Size())
	} else {
		d.Size = fmt.Sprintf("%v (%v bytes)", humanize.Format(stat.Size()), stat.Size())
	}
	if strings.Contains(mime, ";") {
		d.Mime = strings.Split(mime, ";")[0]
	} else {
		d.Mime = mime
	}
	return d, err
}
