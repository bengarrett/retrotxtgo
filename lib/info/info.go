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

	c "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/humanize"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sauce"
	"retrotxt.com/retrotxt/lib/str"
)

// Detail of a file.
type Detail struct {
	XMLName    xml.Name     `json:"-" xml:"file"`
	Name       string       `json:"filename" xml:"name"`
	Utf8       bool         `json:"utf8" xml:"utf8,attr"`
	Newline    [2]rune      `json:"newline" xml:"newline"` // TODO: improve xml handling
	Count      Stats        `json:"counts" xml:"counts"`
	Size       Sizes        `json:"size" xml:"size"`
	Lines      int          `json:"lines" xml:"lines"`
	Width      int          `json:"width" xml:"width"`
	Modified   ModDates     `json:"modified" xml:"last_modified"`
	Sums       Checksums    `json:"checksums" xml:"checksums"`
	Mime       Content      `json:"mime" xml:"mime"`
	Slug       string       `json:"slug" xml:"id,attr"`
	Sauce      sauce.Record `json:"sauce" xml:"sauce"`
	index      int
	length     int
	sauceIndex int
}

// Stats are the text file content statistics and counts.
type Stats struct {
	CharCount int `json:"characters" xml:"characters"`
	CtrlCount int `json:"ansiControls" xml:"ansi_controls"`
	WordCount int `json:"words" xml:"words"`
}

// ModDates is the file last modified dates in multiple output formats.
type ModDates struct {
	Time  time.Time `json:"iso" xml:"date"`
	Epoch int64     `json:"epoch" xml:"epoch,attr"`
}

// Checksums and hashes of the file.
type Checksums struct {
	MD5    string `json:"MD5" xml:"md5"`
	SHA256 string `json:"SHA256" xml:"sha256"`
}

// Content metadata from either MIME content type and magic file data.
type Content struct {
	Type  string `json:"-" xml:"-"`
	Media string `json:"media" xml:"media"`
	Sub   string `json:"subMedia" xml:"sub_media"`
	Commt string `json:"comment" xml:"comment"`
}

// Sizes of the file in multiples.
type Sizes struct {
	Bytes   int64  `json:"bytes" xml:"bytes"`
	Decimal string `json:"decimal" xml:"decimal,attr"`
	Binary  string `json:"binary" xml:"binary,attr"`
}

// Names index and totals.
type Names struct {
	Index  int
	Length int
}

const (
	// DTFormat is the date-time format.
	DTFormat = "DMY24"
	// DFormat is the date format.
	DFormat = "DMY"
)

var (
	// ErrFmt format error.
	ErrFmt = errors.New("format is not known")
	// ErrNoName name cannot be empty.
	ErrNoName = errors.New("name cannot be empty")
	// ErrNoDir directories not usable with command.
	ErrNoDir = errors.New("directories are not usable with this command")
	// ErrNoFile file does not exist.
	ErrNoFile = errors.New("file does not exist")
)

// Info parses the named file and prints out its details in a specific syntax.
func (n Names) Info(name, format string) logs.Generic {
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
	} else if err := Print(name, format, n.Index, n.Length); err != nil {
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
func Print(filename, format string, i, length int) error {
	var d Detail
	if err := d.Read(filename); err != nil {
		return err
	}
	d.index, d.length = i, length
	if IsText(d.Mime.Type) {
		var err error
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
	var cnt int
	if cnt, err = filesystem.Controls(f); err != nil {
		return err
	}
	d.Count.CtrlCount = cnt
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
		w = d.Count.CharCount
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
	switch d.Newline {
	case [2]rune{21}, [2]rune{133}:
		if w, err = filesystem.WordsEBCDIC(f); err != nil {
			return err
		}
	default:
		if w, err = filesystem.Words(f); err != nil {
			return err
		}
	}
	d.Count.WordCount = w
	return f.Close()
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) error {
	var d Detail
	if err := d.parse(nil, b...); err != nil {
		return err
	}
	if IsText(d.Mime.Type) {
		var err error
		d.Newline = filesystem.Newlines(true, []rune(string(b))...)
		if d.Count.CtrlCount, err = filesystem.Controls(bytes.NewReader(b)); err != nil {
			return err
		}
		if d.Lines, err = filesystem.Lines(bytes.NewReader(b), d.Newline); err != nil {
			return err
		}
		if d.Width, err = filesystem.Columns(bytes.NewReader(b), d.Newline); err != nil {
			return err
		} else if d.Width < 0 {
			d.Width = d.Count.CharCount
		}
		if d.Count.WordCount, err = filesystem.Words(bytes.NewReader(b)); err != nil {
			return err
		}
	}
	return d.format(format)
}

func (d *Detail) format(format string) error {
	switch format {
	case "color", "c", "":
		fmt.Printf("%s", d.Text(true))
	case "json", "j":
		b, err := json.MarshalIndent(d, "", "    ")
		if err != nil {
			return fmt.Errorf("detail json indent format: %w", err)
		}
		fmt.Printf("%s\n", b)
	case "json.min", "jm":
		b, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("detail json format: %w", err)
		}
		fmt.Printf("%s\n", b)
	case "text", "t":
		fmt.Printf("%s", d.Text(false))
	case "xml", "x":
		b, err := xml.MarshalIndent(d, "", "\t")
		if err != nil {
			return fmt.Errorf("detail xml format: %w", err)
		}
		fmt.Printf("%s\n", b)
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
	const channels = 6
	ch := make(chan bool, channels)
	go func() {
		d.sauceIndex = sauce.Scan(data...)
		if d.sauceIndex > 0 {
			d.Sauce = sauce.Parse(data...)
		}
		ch <- true
	}()
	go func() {
		mm := mimemagic.MatchMagic(data)
		d.Mime.Media = mm.Media
		d.Mime.Sub = mm.Subtype
		d.Mime.Type = fmt.Sprintf("%s/%s", mm.Media, mm.Subtype)
		d.Mime.Commt = mm.Comment
		if IsText(d.Mime.Type) {
			if d.Count.CharCount, err = filesystem.Runes(bytes.NewBuffer(data)); err != nil {
				fmt.Printf("minesniffer errored, %s\n", err)
			}
		}
		ch <- true
	}()
	go func() {
		var standardInput os.FileInfo = nil
		if stat != standardInput {
			b := stat.Size()
			d.Size.Bytes = b
			d.Size.Binary = humanize.Binary(b, Language())
			d.Size.Decimal = humanize.Decimal(b, Language())
			d.Name = stat.Name()
			d.Modified.Time = stat.ModTime().UTC()
			d.Modified.Epoch = stat.ModTime().Unix()
			d.Slug = slugify.Slugify(stat.Name())
		} else {
			b := int64(len(data))
			d.Size.Bytes = b
			d.Size.Binary = humanize.Binary(b, Language())
			d.Size.Decimal = humanize.Decimal(b, Language())
			d.Name = "n/a (stdin)"
			d.Slug = "n/a"
			d.Modified.Time = time.Now()
			d.Modified.Epoch = time.Now().Unix()
		}
		ch <- true
	}()
	go func() {
		md5sum := md5.Sum(data)
		d.Sums.MD5 = fmt.Sprintf("%x", md5sum)
		ch <- true
	}()
	go func() {
		shasum := sha256.Sum256(data)
		d.Sums.SHA256 = fmt.Sprintf("%x", shasum)
		ch <- true
	}()
	go func() {
		d.Utf8 = utf8.Valid(data)
		ch <- true
	}()
	_, _, _, _, _, _ = <-ch, <-ch, <-ch, <-ch, <-ch, <-ch
	return err
}

// Text format and returns the details of a file.
func (d *Detail) Text(color bool) string {
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
		{k: "filetype", v: d.comment()},
		{k: "UTF-8", v: str.Bool(d.Utf8)},
		{k: "newline", v: filesystem.Newline(d.Newline, true)},
		{k: "characters", v: p.Sprint(d.Count.CharCount)},
		{k: "ANSI controls", v: p.Sprint(d.Count.CtrlCount)},
		{k: "words", v: p.Sprint(d.Count.WordCount)},
		{k: "size", v: d.Size.Decimal},
		{k: "lines", v: p.Sprint(d.Lines)},
		{k: "width", v: p.Sprint(d.Width)},
		{k: "modified", v: humanize.Datetime(DTFormat, d.Modified.Time.UTC())},
		{k: "MD5 checksum", v: d.Sums.MD5},
		{k: "SHA256 checksum", v: d.Sums.SHA256},
		{k: "media mime type", v: d.Mime.Type},
		{k: "slug", v: d.Slug},
		// sauce data
		{k: "title", v: d.Sauce.Title},
		{k: "author", v: d.Sauce.Author},
		{k: "group", v: d.Sauce.Group},
		{k: "date", v: humanize.Date(DFormat, d.Sauce.Date.Time.UTC())},
		{k: "original size", v: d.Sauce.FileSize.Decimal},
		{k: "file type", v: d.Sauce.File.Name},
		{k: "data type", v: d.Sauce.Data.Name},
		{k: "description", v: d.Sauce.Desc},
		{k: d.Sauce.Info.Info1.Info, v: fmt.Sprint(d.Sauce.Info.Info1.Value)},
		{k: d.Sauce.Info.Info2.Info, v: fmt.Sprint(d.Sauce.Info.Info2.Value)},
		{k: d.Sauce.Info.Info3.Info, v: fmt.Sprint(d.Sauce.Info.Info3.Value)},
		{k: "interpretation", v: d.Sauce.Info.Flags.String()},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	l := len(fmt.Sprintf(" filename%s%s", strings.Repeat(" ", 10), data[0].v))
	fmt.Fprint(w, hr(l))
	for _, x := range data {
		if !IsText(d.Mime.Type) {
			switch x.k {
			case "UTF-8", "newline", "characters", "ANSI controls", "words", "lines", "width":
				continue
			}
		} else if x.k == "ANSI controls" {
			if d.Count.CtrlCount == 0 {
				continue
			}
		}
		if x.k == "description" && x.v == "" {
			continue
		}
		if x.k == d.Sauce.Info.Info1.Info && d.Sauce.Info.Info1.Value == 0 {
			continue
		}
		if x.k == d.Sauce.Info.Info2.Info && d.Sauce.Info.Info2.Value == 0 {
			continue
		}
		if x.k == d.Sauce.Info.Info3.Info && d.Sauce.Info.Info3.Value == 0 {
			continue
		}
		if x.k == "interpretation" && x.v == "" {
			continue
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
		if x.k == "slug" {
			if d.sauceIndex <= 0 {
				break
			}
			fmt.Fprint(w, "\t \t   -───-\n")
		}
	}
	if d.index == d.length {
		fmt.Fprint(w, hr(l))
	}
	if err := w.Flush(); err != nil {
		logs.Fatal("flush of tab writer failed", "", err)
	}
	return buf.String()
}

// todo: update d.Mime.Cmmt
func (d *Detail) comment() string {
	if d.Mime.Commt != "unknown" {
		return d.Mime.Commt
	}
	if d.Count.CtrlCount > 0 {
		return "ANSI encoded text document"
	}
	switch d.Newline {
	case [2]rune{21}, [2]rune{133}:
		return "EBCDIC encoded text document"
	}
	if d.Mime.Type == "application/octet-stream" {
		if d.Count.WordCount > 0 {
			// todo !utf8 check
			return "US-ASCII encoded text document"
		}
	}
	return d.Mime.Commt
}
