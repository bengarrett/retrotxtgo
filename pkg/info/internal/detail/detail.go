package detail

import (
	"archive/zip"
	"bytes"
	"text/tabwriter"

	//nolint:gosec
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io/fs"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/term"
	"github.com/bengarrett/sauce"
	"github.com/bengarrett/sauce/humanize"
	gookit "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var ErrFmt = errors.New("format is not known")

// Detail of a file.
//
//nolint:musttag
type Detail struct {
	XMLName    xml.Name     `json:"-"          xml:"file"`
	Name       string       `json:"filename"   xml:"name"`
	Unicode    string       `json:"unicode"    xml:"unicode,attr"`
	LineBreak  LineBreaks   `json:"lineBreak"  xml:"line_break"`
	Count      Stats        `json:"counts"     xml:"counts"`
	Size       Sizes        `json:"size"       xml:"size"`
	Lines      int          `json:"lines"      xml:"lines"`
	Width      int          `json:"width"      xml:"width"`
	Modified   ModDates     `json:"modified"   xml:"last_modified"`
	Sums       Checksums    `json:"checksums"  xml:"checksums"`
	Mime       Content      `json:"mime"       xml:"mime"`
	Slug       string       `json:"slug"       xml:"id,attr"`
	Sauce      sauce.Record `json:"sauce"      xml:"sauce"`
	ZipComment string       `json:"zipComment" xml:"zip_comment"`
	UTF8       bool
	sauceIndex int
}

// Checksums and hashes of the file.
type Checksums struct {
	CRC32  string `json:"crc32"  xml:"crc32"`
	CRC64  string `json:"crc64"  xml:"crc64"`
	MD5    string `json:"md5"    xml:"md5"`
	SHA256 string `json:"sha256" xml:"sha256"`
}

// Content metadata from either MIME content type and magic file data.
type Content struct {
	Type  string `json:"-"        xml:"-"`
	Media string `json:"media"    xml:"media"`
	Sub   string `json:"subMedia" xml:"sub_media"`
	Commt string `json:"comment"  xml:"comment"`
}

// LineBreaks for new line toggles.
type LineBreaks struct {
	Abbr     string  `json:"string"   xml:"string,attr"`
	Escape   string  `json:"escape"   xml:"-"`
	Decimals [2]rune `json:"decimals" xml:"decimal"`
}

// ModDates is the file last modified dates in multiple output formats.
type ModDates struct {
	Time  time.Time `json:"iso"   xml:"date"`
	Epoch int64     `json:"epoch" xml:"epoch,attr"`
}

// Sizes of the file in multiples.
type Sizes struct {
	Bytes   int64  `json:"bytes"   xml:"bytes"`
	Decimal string `json:"decimal" xml:"decimal,attr"`
	Binary  string `json:"binary"  xml:"binary,attr"`
}

// Stats are the text file content statistics and counts.
type Stats struct {
	Chars    int `json:"characters"   xml:"characters"`
	Controls int `json:"ansiControls" xml:"ansi_controls"`
	Words    int `json:"words"        xml:"words"`
}

// Format of the text to output.
type Format int

const (
	// ColorText is ANSI colored text.
	ColorText Format = iota
	// PlainText is standard text.
	PlainText
	// JSON data-interchange format.
	JSON
	// JSONMin is JSON data minified.
	JSONMin
	// XML markup data.
	XML
)

const (
	ans        = "ANSI controls"
	cmmt       = "comment"
	txt        = "text"
	zipComment = "zip comment"
	lf         = 10
	cr         = 13
	nl         = 21
	nel        = 133
	uc8        = "UTF-8"
)

const (
	octetStream = "application/octet-stream"
	zipType     = "application/zip"
)

// lang returns the English Language tag used for numeric syntax formatting.
func lang() language.Tag {
	return language.English
}

// Ctrls counts the number of ANSI escape controls in the named file.
func (d *Detail) Ctrls(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	var cnt int
	if cnt, err = fsys.Controls(f); err != nil {
		return err
	}
	d.Count.Controls = cnt
	return f.Close()
}

// LineTotals counts the totals lines in the named file.
func (d *Detail) LineTotals(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	var l int
	if l, err = fsys.Lines(f, d.LineBreak.Decimals); err != nil {
		return err
	}
	d.Lines = l
	return f.Close()
}

// Marshal the Detail data to a text format syntax.
func (d *Detail) Marshal(f Format) ([]byte, error) {
	var err error
	var b []byte
	switch f {
	case ColorText:
		return d.printMarshal(true)
	case PlainText:
		return d.printMarshal(false)
	case JSON:
		b, err = json.MarshalIndent(d, "", "    ")
		if err != nil {
			return nil, fmt.Errorf("detail json indent marshal: %w", err)
		}
	case JSONMin:
		b, err = json.Marshal(d)
		if err != nil {
			return nil, fmt.Errorf("detail json marshal: %w", err)
		}
	case XML:
		b, err = xml.MarshalIndent(d, "", "\t")
		if err != nil {
			return nil, fmt.Errorf("detail xml marshal: %w", err)
		}
	default:
		return nil, fmt.Errorf("detail marshal %q: %w", f, ErrFmt)
	}
	return b, nil
}

// MimeUnknown detects non-Standard legacy data.
func (d *Detail) MimeUnknown() {
	if d.Mime.Commt != "unknown" {
		return
	}
	if d.Count.Controls > 0 {
		d.Mime.Commt = "Text document with ANSI controls"
		return
	}
	switch d.LineBreak.Decimals {
	case [2]rune{21}, [2]rune{133}:
		d.Mime.Commt = "EBCDIC encoded text document"
		return
	}
	if d.Mime.Type == octetStream {
		if !d.UTF8 && d.Count.Words > 0 {
			d.Mime.Commt = "US-ASCII encoded text document"
			return
		}
	}
}

// Parse the file and the raw data content.
func (d *Detail) Parse(name string, stat os.FileInfo, data ...byte) error {
	const routines = 8
	var wg sync.WaitGroup
	wg.Add(routines)
	go func() {
		defer wg.Done()
		d.sauceIndex = sauce.Index(data)
		if d.sauceIndex > 0 {
			d.Sauce = sauce.Decode(data)
		}
	}()
	go func() {
		defer wg.Done()
		d.mime(name, data...)
	}()
	go func() {
		defer wg.Done()
		d.input(len(data), stat)
	}()
	go func() {
		defer wg.Done()
		crc32sum := crc32.ChecksumIEEE(data)
		d.Sums.CRC32 = fmt.Sprintf("%x", crc32sum)
	}()
	go func() {
		defer wg.Done()
		crc64sum := crc64.Checksum(data, crc64.MakeTable(crc64.ECMA))
		d.Sums.CRC64 = fmt.Sprintf("%x", crc64sum)
	}()
	go func() {
		defer wg.Done()
		md5sum := md5.Sum(data) //nolint:gosec
		d.Sums.MD5 = fmt.Sprintf("%x", md5sum)
	}()
	go func() {
		defer wg.Done()
		shasum := sha256.Sum256(data)
		d.Sums.SHA256 = fmt.Sprintf("%x", shasum)
	}()
	go func() {
		defer wg.Done()
		d.UTF8 = utf8.Valid(data)
		d.Unicode = unicode(&data, d.UTF8)
	}()
	wg.Wait()
	return nil
}

func (d *Detail) mime(name string, data ...byte) {
	mm := mimemagic.MatchMagic(data)
	d.Mime.Media = mm.Media
	d.Mime.Sub = mm.Subtype
	d.Mime.Type = fmt.Sprintf("%s/%s", mm.Media, mm.Subtype)
	d.Mime.Commt = mm.Comment
	if d.Mime.Commt == "plain text document" {
		reader := bytes.NewReader(data)
		if s := bbs.Find(reader).Name(); s != "" {
			d.Mime.Commt += fmt.Sprintf(" with %s BBS color codes", s)
		}
	}
	if d.ValidText() {
		var err error
		b := bytes.NewBuffer(data)
		if d.Count.Chars, err = fsys.Runes(b); err != nil {
			fmt.Fprintf(os.Stdout, "mine sniffer failure, %s\n", err)
		}
		return
	}
	if d.Mime.Type == zipType {
		r, e := zip.OpenReader(name)
		if e != nil {
			fmt.Fprintf(os.Stdout, "open zip file failure: %s\n", e)
		}
		defer r.Close()
		d.ZipComment = r.Comment
	}
}

func unicode(b *[]byte, uni8 bool) string {
	UTF8Bom := []byte{0xEF, 0xBB, 0xBF}
	// little endianness, x86, ARM
	UTF16LEBom := []byte{0xFF, 0xFE}
	UTF32LEBom := []byte{0xFF, 0xFE, 0x00, 0x00}
	// big endianness, legacy mainframes, RISC
	UTF16BEBom := []byte{0xFE, 0xFF}
	UTF32BEBom := []byte{0x00, 0x00, 0xFE, 0xFF}
	switch {
	case bytes.HasPrefix(*b, UTF8Bom):
		return uc8
	case bytes.HasPrefix(*b, UTF16LEBom):
		return "UTF-16 LE"
	case bytes.HasPrefix(*b, UTF16BEBom):
		return "UTF-16 BE"
	case bytes.HasPrefix(*b, UTF32LEBom):
		return "UTF-32 LE"
	case bytes.HasPrefix(*b, UTF32BEBom):
		return "UTF-32 BE"
	default:
		if uni8 {
			return "UTF-8 compatible"
		}
		return "no"
	}
}

// input parses simple statistical data on the file.
func (d *Detail) input(data int, stat fs.FileInfo) {
	if stat != nil {
		b := stat.Size()
		d.Size.Bytes = b
		d.Size.Binary = humanize.Binary(b, lang())
		d.Size.Decimal = humanize.Decimal(b, lang())
		d.Name = stat.Name()
		d.Modified.Time = stat.ModTime().UTC()
		d.Modified.Epoch = stat.ModTime().Unix()
		d.Slug = slugify.Slugify(stat.Name())
		return
	}
	b := int64(data)
	d.Size.Bytes = b
	d.Size.Binary = humanize.Binary(b, lang())
	d.Size.Decimal = humanize.Decimal(b, lang())
	d.Name = "n/a (stdin)"
	d.Slug = "n/a"
	d.Modified.Time = time.Now()
	d.Modified.Epoch = time.Now().Unix()
}

// printMarshal returns the marshaled detail data as plain or color text.
func (d *Detail) printMarshal(color bool) ([]byte, error) {
	const padding, width = 10, 80
	info := func(s string) string {
		return fmt.Sprintf("%s\t", s)
	}
	gookit.Enable = color
	data := d.printMarshalData()
	l := len(fmt.Sprintf(" filename%s%s", strings.Repeat(" ", padding), data[0].v))
	const tabWidth = 8
	b := &bytes.Buffer{}
	w := tabwriter.NewWriter(b, 0, tabWidth, 0, '\t', 0)
	if _, err := term.Head(w, width, "File information"); err != nil {
		return nil, err
	}
	for _, x := range data {
		if !d.marshalDataValid(x.k, x.v) {
			continue
		}
		if x.k == zipComment {
			if x.v != "" {
				fmt.Fprintln(w, term.HR(l))
				fmt.Fprintln(w, x.v)
				if d.sauceIndex <= 0 {
					break
				}
				// divider for sauce metadata
				fmt.Fprintln(w, term.HR(l))
				continue
			}
			if d.sauceIndex <= 0 {
				break
			}
			// divider for sauce metadata
			fmt.Fprintln(w, term.HR(l))
			continue
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
		if x.k == cmmt {
			if d.Sauce.Comnt.Count <= 0 {
				break
			}
		}
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// marshalDataValid returns true if the key and value data validates.
func (d *Detail) marshalDataValid(k, v string) bool {
	if !d.ValidText() {
		switch k {
		case uc8, "line break", "characters", ans, "words", "lines", "width":
			return false
		}
	} else if k == ans {
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

// LineBreaks determines the new lines characters found in the rune pair.
func (d *Detail) LineBreaks(r [2]rune) {
	a, e := "", ""
	switch r {
	case [2]rune{lf}:
		a = "lf"
		e = "\n"
	case [2]rune{cr}:
		a = "cr"
		e = "\r"
	case [2]rune{cr, lf}:
		a = "crlf"
		e = "\r\n"
	case [2]rune{lf, cr}:
		a = "lfcr"
		e = "\n\r"
	case [2]rune{nl}, [2]rune{nel}:
		a = "nl"
		e = "\025"
	}
	d.LineBreak.Decimals = r
	d.LineBreak.Abbr = strings.ToUpper(a)
	d.LineBreak.Escape = e
}

// printMarshalData returns the data structure used for print marshaling.
func (d *Detail) printMarshalData() []struct{ k, v string } {
	const (
		noBreakSpace     = "\u00A0"
		symbolForNewline = "\u2424"
	)
	p := message.NewPrinter(lang())
	data := []struct {
		k, v string
	}{
		{k: "slug", v: d.Slug},
		{k: "filename", v: d.Name},
		{k: "filetype", v: d.Mime.Commt},
		{k: "Unicode", v: d.Unicode},
		{k: "line break", v: fsys.LineBreak(d.LineBreak.Decimals, true)},
		{k: "characters", v: p.Sprint(d.Count.Chars)},
		{k: ans, v: p.Sprint(d.Count.Controls)},
		{k: "words", v: p.Sprint(d.Count.Words)},
		{k: "size", v: d.Size.Decimal},
		{k: "lines", v: p.Sprint(d.Lines)},
		{k: "width", v: p.Sprint(d.Width)},
		{k: "modified", v: humanize.DMY.Format(d.Modified.Time.UTC())},
		{k: "media mime type", v: d.Mime.Type},
		{k: "SHA256 checksum", v: d.Sums.SHA256},
		{k: "CRC64 ECMA", v: d.Sums.CRC64},
		{k: "CRC32", v: d.Sums.CRC32},
		{k: "MD5", v: d.Sums.MD5},
		{k: zipComment, v: d.ZipComment},
		// sauce data
		{k: "title", v: d.Sauce.Title},
		{k: "author", v: d.Sauce.Author},
		{k: "group", v: d.Sauce.Group},
		{k: "date", v: humanize.DMY.Format(d.Modified.Time.UTC())},
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
			comment.k = cmmt
		}
		data = append(data, comment)
	}
	return data
}

// Read and parse the named file and content.
func (d *Detail) Read(name string) error {
	// Get the file details
	stat, err := os.Stat(name)
	if err != nil {
		return err
	}
	// Read file content
	data, err := fsys.ReadAllBytes(name)
	if err != nil {
		return err
	}
	return d.Parse(name, stat, data...)
}

// ValidText returns true if the MIME content-type value is valid for text files.
func (d *Detail) ValidText() bool {
	s := strings.Split(d.Mime.Type, "/")
	const req = 2
	if len(s) != req {
		return false
	}
	if s[0] == txt {
		return true
	}
	if d.Mime.Type == octetStream {
		return true
	}
	return false
}

// Len counts the number of characters used per line in the named file.
func (d *Detail) Len(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	var w int
	if w, err = fsys.Columns(f, d.LineBreak.Decimals); err != nil {
		return err
	} else if w < 0 {
		w = d.Count.Chars
	}
	d.Width = w
	return f.Close()
}

// Words counts the number of words used in the named file.
func (d *Detail) Words(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	var w int
	switch d.LineBreak.Decimals {
	case [2]rune{nl}, [2]rune{nel}:
		if w, err = fsys.WordsEBCDIC(f); err != nil {
			return err
		}
	default:
		if w, err = fsys.Words(f); err != nil {
			return err
		}
	}
	d.Count.Words = w
	return f.Close()
}
