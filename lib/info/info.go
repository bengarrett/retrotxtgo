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
	"sync"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	c "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/sync/errgroup"
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
	Newline    Newlines     `json:"newline" xml:"newline"`
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

// Newlines or line endings.
type Newlines struct {
	Abbr     string  `json:"string" xml:"string,attr"`
	Escape   string  `json:"escape" xml:"-"`
	Decimals [2]rune `json:"decimals" xml:"decimal"`
}

// Stats are the text file content statistics and counts.
type Stats struct {
	Chars    int `json:"characters" xml:"characters"`
	Controls int `json:"ansiControls" xml:"ansi_controls"`
	Words    int `json:"words" xml:"words"`
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
	DFormat     = "DMY"
	octetStream = "application/octet-stream"
	text        = "text"
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
		// todo: directory walk
		gen.Issue = "info"
		gen.Err = ErrNoDir
	} else if err := Marshal(name, format, n.Index, n.Length); err != nil {
		gen.Issue = "info.print"
		gen.Arg = format
		gen.Err = err
	}
	return gen
}

// Language tag used for numeric syntax formatting.
func lang() language.Tag {
	return language.English
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
	d.Count.Controls = cnt
	return f.Close()
}

func (d *Detail) lines(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var l int
	if l, err = filesystem.Lines(f, d.Newline.Decimals); err != nil {
		return err
	}
	d.Lines = l
	return f.Close()
}

func (d *Detail) marshal(format string) (b []byte, err error) {
	switch format {
	case "color", "c", "":
		return d.printMarshal(true), nil
	case text, "t":
		return d.printMarshal(false), nil
	case "json", "j":
		b, err = json.MarshalIndent(d, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("detail json indent marshal: %w", err)
		}
	case "json.min", "jm":
		b, err = json.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("detail json marshal: %w", err)
		}
	case "xml", "x":
		b, err = xml.MarshalIndent(d, "", "\t")
		if err != nil {
			return nil, fmt.Errorf("detail xml marshal: %w", err)
		}
	default:
		return nil, fmt.Errorf("detail marshal %q: %w", format, ErrFmt)
	}
	return b, nil
}

func (d *Detail) mimeUnknown() {
	if d.Mime.Commt == "unknown" {
		if d.Count.Controls > 0 {
			d.Mime.Commt = "Text document with ANSI controls"
			return
		}
		switch d.Newline.Decimals {
		case [2]rune{21}, [2]rune{133}:
			d.Mime.Commt = "EBCDIC encoded text document"
			return
		}
		if d.Mime.Type == octetStream {
			if !d.Utf8 && d.Count.Words > 0 {
				d.Mime.Commt = "US-ASCII encoded text document"
				return
			}
		}
	}
}

func (d *Detail) parse(stat os.FileInfo, data ...byte) (err error) {
	const routines = 6
	var wg sync.WaitGroup
	wg.Add(routines)
	go func() {
		defer wg.Done()
		d.sauceIndex = sauce.Scan(data...)
		if d.sauceIndex > 0 {
			d.Sauce = sauce.Parse(data...)
		}
	}()
	go func() {
		defer wg.Done()
		mm := mimemagic.MatchMagic(data)
		d.Mime.Media = mm.Media
		d.Mime.Sub = mm.Subtype
		d.Mime.Type = fmt.Sprintf("%s/%s", mm.Media, mm.Subtype)
		d.Mime.Commt = mm.Comment
		if d.validText() {
			if d.Count.Chars, err = filesystem.Runes(bytes.NewBuffer(data)); err != nil {
				fmt.Printf("minesniffer errored, %s\n", err)
			}
		}
	}()
	go func() {
		defer wg.Done()
		var standardInput os.FileInfo = nil
		if stat != standardInput {
			b := stat.Size()
			d.Size.Bytes = b
			d.Size.Binary = humanize.Binary(b, lang())
			d.Size.Decimal = humanize.Decimal(b, lang())
			d.Name = stat.Name()
			d.Modified.Time = stat.ModTime().UTC()
			d.Modified.Epoch = stat.ModTime().Unix()
			d.Slug = slugify.Slugify(stat.Name())
		} else {
			b := int64(len(data))
			d.Size.Bytes = b
			d.Size.Binary = humanize.Binary(b, lang())
			d.Size.Decimal = humanize.Decimal(b, lang())
			d.Name = "n/a (stdin)"
			d.Slug = "n/a"
			d.Modified.Time = time.Now()
			d.Modified.Epoch = time.Now().Unix()
		}
	}()
	go func() {
		defer wg.Done()
		md5sum := md5.Sum(data)
		d.Sums.MD5 = fmt.Sprintf("%x", md5sum)
	}()
	go func() {
		defer wg.Done()
		shasum := sha256.Sum256(data)
		d.Sums.SHA256 = fmt.Sprintf("%x", shasum)
	}()
	go func() {
		defer wg.Done()
		d.Utf8 = utf8.Valid(data)
	}()
	wg.Wait()
	return err
}

func (d *Detail) printMarshal(color bool) []byte {
	c.Enable = color
	var (
		buf  bytes.Buffer
		info = func(t string) string {
			return str.Cinf(fmt.Sprintf("%s\t", t))
		}
		hr = func(l int) string {
			return fmt.Sprintf("\t%s\n", str.Cb(strings.Repeat("\u2500", l)))
		}
		data = d.printMarshalData()
		w    = new(tabwriter.Writer)
		l    = len(fmt.Sprintf(" filename%s%s", strings.Repeat(" ", 10), data[0].v))
	)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	fmt.Fprint(w, hr(l))
	for _, x := range data {
		if !d.marshalDataValid(x.k, x.v) {
			continue
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
		if x.k == "slug" {
			if d.sauceIndex <= 0 {
				break
			}
			fmt.Fprint(w, "\t \t   -───-\n")
		}
		if x.k == "comment" {
			if d.Sauce.Comnt.Count <= 0 {
				break
			}
		}
	}
	if d.index == d.length {
		fmt.Fprint(w, hr(l))
	}
	if err := w.Flush(); err != nil {
		logs.Fatal("flush of tab writer failed", "", err)
	}
	return buf.Bytes()
}

func (d *Detail) marshalDataValid(k, v string) bool {
	if !d.validText() {
		switch k {
		case "UTF-8", "newline", "characters", "ANSI controls", "words", "lines", "width":
			return false
		}
	} else if k == "ANSI controls" {
		if d.Count.Controls == 0 {
			return false
		}
	}
	if k == "description" && v == "" {
		return false
	}
	if k == d.Sauce.Info.Info1.Info && d.Sauce.Info.Info1.Value == 0 {
		return false
	}
	if k == d.Sauce.Info.Info2.Info && d.Sauce.Info.Info2.Value == 0 {
		return false
	}
	if k == d.Sauce.Info.Info3.Info && d.Sauce.Info.Info3.Value == 0 {
		return false
	}
	if k == "interpretation" && v == "" {
		return false
	}
	return true
}

func (d *Detail) newlines(r [2]rune) {
	a, e := "", ""
	switch r {
	case [2]rune{10}:
		a = "lf"
		e = "\n"
	case [2]rune{13}:
		a = "cr"
		e = "\r"
	case [2]rune{13, 10}:
		a = "crlf"
		e = "\r\n"
	case [2]rune{10, 13}:
		a = "lfcr"
		e = "\n\r"
	case [2]rune{21}, [2]rune{133}:
		a = "nl"
		e = "\025"
	}
	d.Newline.Decimals = r
	d.Newline.Abbr = strings.ToUpper(a)
	d.Newline.Escape = e
}

func (d *Detail) printMarshalData() (data []struct{ k, v string }) {
	const noBreakSpace, symbolForNewline = "\u00a0", "\u2424"
	p := message.NewPrinter(lang())
	data = []struct {
		k, v string
	}{
		{k: "filename", v: d.Name},
		{k: "filetype", v: d.Mime.Commt},
		{k: "UTF-8", v: str.Bool(d.Utf8)},
		{k: "newline", v: filesystem.Newline(d.Newline.Decimals, true)},
		{k: "characters", v: p.Sprint(d.Count.Chars)},
		{k: "ANSI controls", v: p.Sprint(d.Count.Controls)},
		{k: "words", v: p.Sprint(d.Count.Words)},
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
	// sauce comment
	for i, line := range d.Sauce.Comnt.Comment {
		comment := struct{ k, v string }{
			k: noBreakSpace, v: line,
		}
		if i == 0 {
			comment.k = "comment"
		}
		data = append(data, comment)
	}
	return data
}

func (d *Detail) read(name string) (err error) {
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

// validText checks the MIME content-type value for valid text files.
func (d *Detail) validText() bool {
	s := strings.Split(d.Mime.Type, "/")
	const req = 2
	if len(s) != req {
		return false
	}
	if s[0] == text {
		return true
	}
	if d.Mime.Type == octetStream {
		return true
	}
	return false
}

func (d *Detail) width(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	var w int
	if w, err = filesystem.Columns(f, d.Newline.Decimals); err != nil {
		return err
	} else if w < 0 {
		w = d.Count.Chars
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
	switch d.Newline.Decimals {
	case [2]rune{21}, [2]rune{133}:
		if w, err = filesystem.WordsEBCDIC(f); err != nil {
			return err
		}
	default:
		if w, err = filesystem.Words(f); err != nil {
			return err
		}
	}
	d.Count.Words = w
	return f.Close()
}

// Marshal the meta and operating system details of a file.
func Marshal(filename, format string, i, length int) error {
	var d Detail
	if err := d.read(filename); err != nil {
		return err
	}
	d.index, d.length = i, length
	if d.validText() {
		var g errgroup.Group
		g.Go(func() error {
			var err error
			if d.Newline.Decimals, err = filesystem.ReadNewlines(filename); err != nil {
				return err
			}
			d.newlines(d.Newline.Decimals)
			return nil
		})
		g.Go(func() error {
			if err := d.ctrls(filename); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			if err := d.width(filename); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			if err := d.lines(filename); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			if err := d.width(filename); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			if err := d.words(filename); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return err
		}
		d.mimeUnknown()
	}
	var (
		m   []byte
		err error
	)
	if m, err = d.marshal(format); err != nil {
		return err
	}
	print(format, m...)
	return nil
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) error {
	var d Detail
	if err := d.parse(nil, b...); err != nil {
		return err
	}
	if d.validText() {
		var g errgroup.Group
		g.Go(func() error {
			d.newlines(filesystem.Newlines(true, []rune(string(b))...))
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Count.Controls, err = filesystem.Controls(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Count.Controls, err = filesystem.Controls(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Lines, err = filesystem.Lines(bytes.NewReader(b), d.Newline.Decimals); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Width, err = filesystem.Columns(bytes.NewReader(b), d.Newline.Decimals); err != nil {
				return err
			} else if d.Width < 0 {
				d.Width = d.Count.Chars
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Count.Words, err = filesystem.Words(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return err
		}
		d.mimeUnknown()
	}
	var (
		m   []byte
		err error
	)
	if m, err = d.marshal(format); err != nil {
		return err
	}
	print(format, m...)
	return nil
}

func print(format string, b ...byte) {
	switch format {
	case "color", "c", "", text, "t":
		fmt.Printf("%s", b)
	default:
		fmt.Printf("%s\n", b)
	}
}
