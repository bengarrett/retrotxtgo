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
	c "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/humanize"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

// Detail of a file.
type Detail struct {
	Bytes     int64
	CharCount int
	CtrlCount int
	Lines     int
	Name      string
	MD5       string
	Mime      string
	Modified  time.Time
	Newline   [2]rune
	Newlines  string
	SHA256    string
	Slug      string
	Size      string
	Utf8      bool
	Width     int
	WordCount int
}

// File data for XML encoding.
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

// DTFormat is the datetime format.
// DMY12, YMD12, MDY12, DMY24, YMD24, MDY24.
const DTFormat = "DMY24"

var (
	ErrFmt    = errors.New("format is not known")
	ErrNoName = errors.New("name cannot be empty")
	ErrNoDir  = errors.New("directories are not usable with this command")
	ErrNoFile = errors.New("file does not exist")
)

// Info parses the named file and prints out its details in a specific syntax.
func Info(name, format string) logs.Generic {
	gen := logs.Generic{Issue: "info", Arg: name}
	if name == "" {
		gen.Issue = "name"
		gen.Err = ErrNoName
		return gen
	}
	if s, err := os.Stat(name); os.IsNotExist(err) {
		gen.Err = ErrNoFile
	} else if err != nil {
		gen.Err = err
	} else if s.IsDir() {
		gen.Issue = "info"
		gen.Err = ErrNoDir
	} else if err := Print(name, format); err != nil {
		gen.Issue = "info.print"
		gen.Arg = format
		gen.Err = err
	}
	return gen
}

// Language tag used for numeric syntax formatting.
func Language() language.Tag {
	return language.English
}

// Print the meta and operating system details of a file.
func Print(filename, format string) (err error) {
	var d Detail
	if err := d.Read(filename); err != nil {
		return err
	}
	if IsText(d.Mime) {
		d.Newline, err = filesystem.ReadNewlines(filename)
		if err != nil {
			return err
		}
		if err := d.ctrls(filename); err != nil {
			return err
		}
		if err := d.lines(filename); err != nil {
			return err
		}
		if err := d.width(filename); err != nil {
			return err
		}
		if err := d.words(filename); err != nil {
			return err
		}
	}
	return d.format(format)
}

func (d *Detail) ctrls(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var c int
	if c, err = filesystem.Controls(f); err != nil {
		return err
	}
	d.CtrlCount = c
	return f.Close()
}

func (d *Detail) lines(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var l int
	if l, err = filesystem.Lines(f, d.Newline); err != nil {
		return err
	}
	d.Lines = l
	return f.Close()
}

func (d *Detail) width(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var w int
	if w, err = filesystem.Columns(f, d.Newline); err != nil {
		return err
	} else if w < 0 {
		w = d.CharCount
	}
	d.Width = w
	return f.Close()
}

func (d *Detail) words(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var w int
	if w, err = filesystem.Words(f); err != nil {
		return err
	}
	d.WordCount = w
	return f.Close()
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) (err error) {
	var d Detail
	if err = d.parse(nil, b...); err != nil {
		return err
	}
	if IsText(d.Mime) {
		d.Newline = filesystem.Newlines(true, []rune(string(b))...)
		if d.CtrlCount, err = filesystem.Controls(bytes.NewReader(b)); err != nil {
			return err
		}
		if d.Lines, err = filesystem.Lines(bytes.NewReader(b), d.Newline); err != nil {
			return err
		}
		if d.Width, err = filesystem.Columns(bytes.NewReader(b), d.Newline); err != nil {
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
		data, err := d.XML()
		if err != nil {
			return fmt.Errorf("detail xml format: %w", err)
		}
		fmt.Printf("%s\n", data)
	default:
		return fmt.Errorf("detail format %q: %w", format, ErrFmt)
	}
	return nil
}

// IsText checks the MIME content-type value for valid text files.
func IsText(contentType string) bool {
	s := strings.Split(contentType, "/")
	const req = 2
	if len(s) != req {
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
	return d.parse(stat, data...)
}

// parse fileinfo and file content.
func (d *Detail) parse(stat os.FileInfo, data ...byte) (err error) {
	p := message.NewPrinter(Language())
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
	const kB = 1000
	// create a table of data
	if stat != nil {
		d.Bytes = stat.Size()
		d.Name = stat.Name()
		d.Modified = stat.ModTime().UTC()
		d.Slug = slugify.Slugify(stat.Name())
		if stat.Size() < kB {
			d.Size = p.Sprintf("%v bytes", p.Sprint(stat.Size()))
		} else {
			d.Size = p.Sprintf("%v (%v bytes)", humanize.Bytes(stat.Size(), Language()), p.Sprint(stat.Size()))
		}
	} else {
		d.Bytes = int64(len(data))
		d.Name = "n/a (stdin)"
		d.Slug = "n/a"
		d.Modified = time.Now()
		l := d.Bytes
		if l < kB {
			d.Size = p.Sprintf("%v bytes", p.Sprint(l))
		} else {
			d.Size = p.Sprintf("%v (%v bytes)", humanize.Bytes(l, Language()), p.Sprint(l))
		}
	}
	const channels = 3
	ch := make(chan bool, channels)
	go func() {
		md5sum := md5.Sum(data)
		d.MD5 = fmt.Sprintf("%x", md5sum)
		ch <- true
	}()
	go func() {
		sha256 := sha256.Sum256(data)
		d.SHA256 = fmt.Sprintf("%x", sha256)
		ch <- true
	}()
	go func() {
		d.Utf8 = utf8.Valid(data)
		ch <- true
	}()
	_, _, _ = <-ch, <-ch, <-ch
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
	if err != nil {
		logs.Fatal("info could not marshal", "json", err)
	}
	return js
}

// Text format and returns the details of a file.
func (d Detail) Text(color bool) string {
	p := message.NewPrinter(Language())
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
		{k: "newline", v: filesystem.Newline(d.Newline, true)},
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
	if err := w.Flush(); err != nil {
		logs.Fatal("flush of tab writer failed", "", err)
	}
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
