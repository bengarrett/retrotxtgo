package info

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/lib/codepage"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	c "github.com/gookit/color"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"

	humanize "github.com/labstack/gommon/bytes"
)

// Detail of a file
type Detail struct {
	Bytes     int64
	CharCount int
	Lines     int
	Name      string
	MD5       string
	Mime      string
	Modified  time.Time
	SHA256    string
	Slug      string
	Size      string
	Utf8      bool
	Width     int
	WordCount int
}

// File data for XML encoding
type File struct {
	XMLName   xml.Name  `xml:"file"`
	ID        string    `xml:"id,attr"`
	Name      string    `xml:"name"`
	Mime      string    `xml:"content>mime"`
	Utf8      bool      `xml:"content>utf8"`
	Bytes     int64     `xml:"size>bytes"`
	Size      string    `xml:"size>value"`
	Lines     int       `xml:"size>lines"`
	Width     int       `xml:"size>width"`
	CharCount int       `xml:"size>character-count"`
	WordCount int       `xml:"size>word-count"`
	MD5       string    `xml:"checksum>md5"`
	SHA256    string    `xml:"checksum>sha256"`
	Modified  time.Time `xml:"modified"`
}

// FileDate is a non-standard date format for file modifications
const FileDate string = "2 Jan 15:04 2006"

// Info parses the named file and prints out its details in a specific syntax.
func Info(name, format string) (err logs.Err) {
	if name == "" {
		return logs.Err{Issue: "info", Arg: "name", Msg: errors.New("value cannot be empty")}
	}
	if s, err := os.Stat(name); os.IsNotExist(err) {
		return logs.Err{Issue: "info --name", Arg: name, Msg: errors.New("file does not exist")}
	} else if err != nil {
		return logs.Err{Issue: "info --name", Arg: name, Msg: err}
	} else if s.IsDir() {
		return logs.Err{Issue: "info --name", Arg: name, Msg: errors.New("directories are not usable with this command")}
	} else if e := Print(name, format); e != nil {
		return logs.Err{Issue: "info print", Arg: format, Msg: e}
	}
	return err
}

// Print the meta and operating system details of a file.
func Print(filename, format string) (err error) {
	var d Detail
	if err := d.Read(filename); err != nil {
		return err
	}
	if d.Mime == "text/plain" {
		if d.Lines, err = filesystem.Lines(filename); err != nil {
			return err
		}
		if d.Width, err = filesystem.Columns(filename); err != nil {
			return err
		}
		if d.WordCount, err = filesystem.Words(filename); err != nil {
			return err
		}
	}
	switch format {
	case "color", "c", "":
		fmt.Printf("%s", d.Text(true))
	case "json", "j":
		fmt.Printf("%s\n", d.JSON(true))
	case "json.min", "jm":
		fmt.Printf("%s\n", d.JSON(false))
	case "text", "t":
		fmt.Printf("%s", d.Text(false))
	case "xml", "x":
		data, _ := d.XML()
		fmt.Printf("%s\n", data)
	default:
		return errors.New("format:invalid")
	}
	return err
}

// Read returns the operating system and meta detail of a named file.
func (d *Detail) Read(name string) (err error) {
	// Get the file details
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}
	// Read file content
	data, err := filesystem.ReadAllBytes(name)
	if err != nil {
		return err
	}
	return d.parse(data, stat, name)
}

// parse fileinfo and file content.
func (d *Detail) parse(data []byte, stat os.FileInfo, name string) (err error) {
	md5sum := md5.Sum(data)
	sha256 := sha256.Sum256(data)
	mime := mimesniffer.Sniff(data)
	if strings.Contains(mime, ";") {
		d.Mime = strings.Split(mime, ";")[0]
	} else {
		d.Mime = mime
	}
	if d.Mime == "text/plain" {
		d.CharCount = runewidth.StringWidth(string(data))
	}
	// create a table of data
	d.Bytes = stat.Size()
	d.Name = stat.Name()
	d.MD5 = fmt.Sprintf("%x", md5sum)
	d.Modified = stat.ModTime().UTC()
	d.Slug = slugify.Slugify(stat.Name())
	d.SHA256 = fmt.Sprintf("%x", sha256)
	d.Utf8 = codepage.UTF8(data)
	if stat.Size() < 1000 {
		d.Size = fmt.Sprintf("%v bytes", stat.Size())
	} else {
		d.Size = fmt.Sprintf("%v (%v bytes)", humanize.Format(stat.Size()), stat.Size())
	}

	return err
}

// JSON format and returns the details of a file.
func (d Detail) JSON(indent bool) (js []byte) {
	var err error
	switch indent {
	case true:
		js, err = json.MarshalIndent(d, "", "    ")
	default:
		js, err = json.Marshal(d)
	}
	logs.ChkErr(logs.Err{Issue: "could not create", Arg: "json", Msg: err})
	return js
}

// Text format and returns the details of a file.
func (d Detail) Text(color bool) string {
	c.Enable = color
	var info = func(t string) string {
		return logs.Cinf(fmt.Sprintf("%s\t", t))
	}
	var hr = func() string {
		return fmt.Sprintf("\t%s\n", logs.Cf(strings.Repeat("\u2015", 26)))
	}
	var data = []struct {
		k, v string
	}{
		{k: "filename", v: d.Name},
		{k: "UTF-8", v: logs.Bool(d.Utf8)},
		{k: "characters", v: fmt.Sprint(d.CharCount)},
		{k: "words", v: fmt.Sprint(d.WordCount)},
		{k: "size", v: d.Size},
		{k: "lines", v: fmt.Sprint(d.Lines)},
		{k: "width", v: fmt.Sprint(d.Width)},
		{k: "modified", v: fmt.Sprintf("%v", d.Modified.UTC().Format(FileDate))},
		{k: "MD5 checksum", v: d.MD5},
		{k: "SHA256 checksum", v: d.SHA256},
		{k: "MIME type", v: d.Mime},
		{k: "slug", v: d.Slug},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, hr())
	for _, x := range data {
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
	}
	fmt.Fprint(w, hr())
	w.Flush()
	return buf.String()
}

// XML formats and returns the details of a file.
func (d Detail) XML() ([]byte, error) {
	v := File{
		Bytes:     d.Bytes,
		CharCount: d.CharCount,
		ID:        d.Slug,
		Lines:     d.Lines,
		MD5:       d.MD5,
		Mime:      d.Mime,
		Modified:  d.Modified,
		Name:      d.Name,
		SHA256:    d.SHA256,
		Size:      d.Size,
		Utf8:      d.Utf8,
		Width:     d.Width,
		WordCount: d.WordCount,
	}
	return xml.MarshalIndent(v, "", "\t")
}
