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
	"unicode/utf8"

	"github.com/aofei/mimesniffer"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/humanize"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	c "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Detail of a file
type Detail struct {
	Bytes     int64
	CharCount int
	CtrlCount int
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
	CtrlCount int       `xml:"size>ansi-control-count"`
	WordCount int       `xml:"size>word-count"`
	MD5       string    `xml:"checksum>md5"`
	SHA256    string    `xml:"checksum>sha256"`
	Modified  time.Time `xml:"modified"`
}

// Language tag used for numeric syntax formatting
var Language = language.English

// DTFormat is the datetime format
// DMY12, YMD12, MDY12, DMY24, YMD24, MDY24
var DTFormat = "DMY24"

// Info parses the named file and prints out its details in a specific syntax.
func Info(name, format string) (err logs.Err) {
	if name == "" {
		return logs.Err{Issue: "info", Arg: "name",
			Msg: errors.New("value cannot be empty")}
	}
	if s, err := os.Stat(name); os.IsNotExist(err) {
		return logs.Err{Issue: "info", Arg: name,
			Msg: errors.New("file does not exist")}
	} else if err != nil {
		return logs.Err{Issue: "info", Arg: name, Msg: err}
	} else if s.IsDir() {
		return logs.Err{Issue: "info", Arg: name,
			Msg: errors.New("directories are not usable with this command")}
	} else if e := Print(name, format); e != nil {
		return logs.Err{Issue: "info.print", Arg: format, Msg: e}
	}
	return err
}

// Print the meta and operating system details of a file.
func Print(filename, format string) (err error) {
	var d Detail
	if err := d.Read(filename); err != nil {
		return err
	}
	if IsText(d.Mime) {
		// TODO: display this as an info field
		newline, err := filesystem.ReadNewlines(filename)
		if err != nil {
			return err
		}
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer file.Close()
		if d.CtrlCount, err = filesystem.Controls(file); err != nil {
			return err
		}
		if d.Lines, err = filesystem.Lines(file); err != nil {
			return err
		}
		if d.Width, err = filesystem.Columns(file, newline); err != nil {
			return err
		} else if d.Width < 0 {
			d.Width = d.CharCount
		}
		if d.WordCount, err = filesystem.Words(file); err != nil {
			return err
		}
	}
	return d.format(format)
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(b []byte, format string) (err error) {
	var d Detail
	if err = d.parse(b, nil); err != nil {
		return err
	}
	if IsText(d.Mime) {
		// TODO: display this as an info field
		newline := filesystem.Newlines([]rune(string(b)))
		if d.CtrlCount, err = filesystem.Controls(bytes.NewReader(b)); err != nil {
			return err
		}
		if d.Lines, err = filesystem.Lines(bytes.NewReader(b)); err != nil {
			return err
		}
		if d.Width, err = filesystem.Columns(bytes.NewReader(b), newline); err != nil {
			return err
		} else if d.Width < 0 {
			d.Width = d.CharCount
		}
		if d.WordCount, err = filesystem.Words(bytes.NewReader(b)); err != nil {
			return err
		}
	}
	return d.format(format)
}

func (d Detail) format(format string) error {
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
		return errors.New("format invalid: " + format)
	}
	return nil
}

// IsText checks the MIME content-type value for valid text files.
func IsText(contentType string) bool {
	s := strings.Split(contentType, "/")
	if len(s) != 2 {
		return false
	}
	if s[0] == "text" {
		return true
	}
	if contentType == "application/octet-stream" {
		return true
	}
	return false
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
	return d.parse(data, stat)
}

// parse fileinfo and file content.
func (d *Detail) parse(data []byte, stat os.FileInfo) (err error) {
	p := message.NewPrinter(Language)
	ms := mimesniffer.Sniff(data)
	if strings.Contains(ms, ";") {
		d.Mime = strings.Split(ms, ";")[0]
	} else {
		d.Mime = ms
	}
	if IsText(d.Mime) {
		if d.CharCount, err = filesystem.Runes(bytes.NewBuffer(data)); err != nil {
			return err
		}
	}
	// create a table of data
	if stat != nil {
		d.Bytes = stat.Size()
		d.Name = stat.Name()
		d.Modified = stat.ModTime().UTC()
		d.Slug = slugify.Slugify(stat.Name())
		if stat.Size() < 1000 {
			d.Size = p.Sprintf("%v bytes", p.Sprint(stat.Size()))
		} else {
			d.Size = p.Sprintf("%v (%v bytes)", humanize.Bytes(stat.Size(), Language), p.Sprint(stat.Size()))
		}
	} else {
		d.Bytes = int64(len(data))
		d.Name = "n/a (stdin)"
		d.Slug = "n/a"
		d.Modified = time.Now()
		l := d.Bytes
		if l < 1000 {
			d.Size = p.Sprintf("%v bytes", p.Sprint(l))
		} else {
			d.Size = p.Sprintf("%v (%v bytes)", humanize.Bytes(l, Language), p.Sprint(l))
		}
	}

	md5sum := md5.Sum(data)
	d.MD5 = fmt.Sprintf("%x", md5sum)

	sha256 := sha256.Sum256(data)
	d.SHA256 = fmt.Sprintf("%x", sha256)

	d.Utf8 = utf8.Valid(data)

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
	p := message.NewPrinter(Language)
	c.Enable = color
	var info = func(t string) string {
		return str.Cinf(fmt.Sprintf("%s\t", t))
	}
	var hr = func(l int) string {
		return fmt.Sprintf("\t%s\n", str.Cb(strings.Repeat("\u2500", l)))
	}
	var data = []struct {
		k, v string
	}{
		{k: "filename", v: d.Name},
		{k: "UTF-8", v: str.Bool(d.Utf8)},
		{k: "characters", v: p.Sprint(d.CharCount)},
		{k: "ANSI controls", v: p.Sprint(d.CtrlCount)},
		{k: "words", v: p.Sprint(d.WordCount)},
		{k: "size", v: d.Size},
		{k: "lines", v: p.Sprint(d.Lines)},
		{k: "width", v: p.Sprint(d.Width)},
		{k: "modified", v: humanize.Datetime(DTFormat, d.Modified.UTC())},
		{k: "MD5 checksum", v: d.MD5},
		{k: "SHA256 checksum", v: d.SHA256},
		{k: "MIME type", v: d.Mime},
		{k: "slug", v: d.Slug},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	l := len(fmt.Sprintf(" filename%s%s", strings.Repeat(" ", 10), data[0].v))
	fmt.Fprint(w, hr(l))
	for _, x := range data {
		if !IsText(d.Mime) {
			switch x.k {
			case "UTF-8", "characters", "words", "lines", "width":
				continue
			}
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
	}
	fmt.Fprint(w, hr(l))
	w.Flush()
	return buf.String()
}

// XML formats and returns the details of a file.
func (d Detail) XML() ([]byte, error) {
	v := File{
		Bytes:     d.Bytes,
		CharCount: d.CharCount,
		CtrlCount: d.CtrlCount,
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
