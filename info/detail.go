package info

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io"
	"io/fs"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/nl"
	"github.com/bengarrett/retrotxtgo/term"
	"github.com/bengarrett/sauce"
	"github.com/bengarrett/sauce/humanize"
	gookit "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Detail is the exported file details.
type Detail struct {
	XMLName    xml.Name     `json:"-"          xml:"file"`
	Name       string       `json:"filename"   xml:"name"`          // Name is the file name.
	Unicode    string       `json:"unicode"    xml:"unicode,attr"`  // Unicode is the file encoding if in Unicode.
	LineBreak  nl.LineBreak `json:"lineBreak"  xml:"line_break"`    // LineBreak is the line break used in the file.
	Count      Stats        `json:"counts"     xml:"counts"`        // Count is the file content statistics.
	Size       Sizes        `json:"size"       xml:"size"`          // Size is the file size in multiples.
	Lines      int          `json:"lines"      xml:"lines"`         // Lines is the number of lines in the file.
	Width      int          `json:"width"      xml:"width"`         // Width is the number of characters per line in the file, this may be inaccurate.
	Modified   ModDates     `json:"modified"   xml:"last_modified"` // Modified is the last modified date of the file.
	Sums       Checksums    `json:"checksums"  xml:"checksums"`     // Sums are the checksums of the file.
	Mime       Content      `json:"mime"       xml:"mime"`          // Mime is the file content metadata.
	Slug       string       `json:"slug"       xml:"id,attr"`       // Slug is the file name slugified.
	Sauce      sauce.Record `json:"sauce"      xml:"sauce"`         // Sauce is the SAUCE metadata.
	ZipComment string       `json:"zipComment" xml:"zip_comment"`   // ZipComment is the zip file comment.
	UTF8       bool         `json:"-"          xml:"-"`             // UTF8 is true if the file is UTF-8 encoded.
	sauceIndex int          // sauceIndex is the index of the SAUCE record in the file.
}

// Checksums act as a fingerprint of the file for uniqueness and data corruption checks.
type Checksums struct {
	CRC32  string `json:"crc32"  xml:"crc32"`  // CRC32 is a cyclic redundancy check of the file.
	CRC64  string `json:"crc64"  xml:"crc64"`  // CRC64 is a cyclic redundancy check of the file.
	MD5    string `json:"md5"    xml:"md5"`    // MD5 is a weak cryptographic hash function.
	SHA256 string `json:"sha256" xml:"sha256"` // SHA256 is a strong cryptographic hash function.
}

// Content metadata from either MIME content type and magic file data.
type Content struct {
	Type  string `json:"-"        xml:"-"`
	Media string `json:"media"    xml:"media"`     // Media is the MIME media type.
	Sub   string `json:"subMedia" xml:"sub_media"` // Sub is the MIME sub type.
	Commt string `json:"comment"  xml:"comment"`   // Commt is the MIME comment.
}

// ModDates is the file last modified dates in multiple output formats.
type ModDates struct {
	Time  time.Time `json:"iso"   xml:"date"`       // Time is the last modified date of the file.
	Epoch int64     `json:"epoch" xml:"epoch,attr"` // Epoch is the last modified date of the file in seconds since the Unix epoch.
}

// Sizes of the file in multiples.
type Sizes struct {
	Bytes   int64  `json:"bytes"   xml:"bytes"`        // Bytes is the size of the file in bytes.
	Decimal string `json:"decimal" xml:"decimal,attr"` // Decimal is the size of the file with decimal units.
	Binary  string `json:"binary"  xml:"binary,attr"`  // Binary is the size of the file with binary units.
}

// Stats are the text file content statistics and counts.
type Stats struct {
	Chars    int `json:"characters"   xml:"characters"`    // Chars is the number of characters in the file.
	Controls int `json:"ansiControls" xml:"ansi_controls"` // Controls is the number of ANSI escape controls in the file.
	Words    int `json:"words"        xml:"words"`         // Words is the number of words in the file, this may be inaccurate.
}

// Format of the output text.
type Format int

const (
	ColorText Format = iota // ColorText is ANSI colored text.
	PlainText               // PlainText is standard text.
	JSON                    // JSON data-interchange format.
	JSONMin                 // JSONMin is JSON data minified.
	XML                     // XML markup data.
)

const (
	uc8         = "UTF-8"
	ans         = "ANSI controls"
	cmmt        = "comment"
	txt         = "text"
	zipComment  = "zip comment"
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
	cnt, err := fsys.Controls(f)
	if err != nil {
		return err
	}
	d.Count.Controls = cnt
	return nil
}

// Marshal writes the Detail data in a given format syntax.
func (d *Detail) Marshal(w io.Writer, f Format) error {
	if w == nil {
		w = io.Discard
	}
	const jsTab = "    "
	const xmlTab = "\t"
	switch f {
	case ColorText:
		return d.marshal(w, true)
	case PlainText:
		return d.marshal(w, false)
	case JSON:
		b, err := json.MarshalIndent(d, "", jsTab)
		if err != nil {
			return fmt.Errorf("detail json indent marshal: %w", err)
		}
		_, err = w.Write(b)
		return err
	case JSONMin:
		b, err := json.Marshal(d)
		if err != nil {
			return fmt.Errorf("detail json marshal: %w", err)
		}
		_, err = w.Write(b)
		return err
	case XML:
		b, err := xml.MarshalIndent(d, "", xmlTab)
		if err != nil {
			return fmt.Errorf("detail xml marshal: %w", err)
		}
		_, err = w.Write(b)
		return err
	default:
		return fmt.Errorf("detail marshal %q: %w", f, ErrFmt)
	}
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
	switch d.LineBreak.Decimal {
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
func (d *Detail) Parse(name string, data ...byte) error {
	const routines = 8
	wg := sync.WaitGroup{}
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
		stat, _ := os.Stat(name)
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
		d.UTF8 = utf8.Valid(data)
		d.Unicode = unicode(d.UTF8, data...)
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
	if ValidText(d.Mime.Type) {
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

func unicode(uni bool, b ...byte) string {
	UTF8Bom := []byte{0xEF, 0xBB, 0xBF}
	// little endianness, x86, ARM
	UTF16LEBom := []byte{0xFF, 0xFE}
	UTF32LEBom := []byte{0xFF, 0xFE, 0x00, 0x00}
	// big endianness, legacy mainframes, RISC
	UTF16BEBom := []byte{0xFE, 0xFF}
	UTF32BEBom := []byte{0x00, 0x00, 0xFE, 0xFF}
	switch {
	case bytes.HasPrefix(b, UTF8Bom):
		return uc8
	case bytes.HasPrefix(b, UTF16LEBom):
		return "UTF-16 LE"
	case bytes.HasPrefix(b, UTF16BEBom):
		return "UTF-16 BE"
	case bytes.HasPrefix(b, UTF32LEBom):
		return "UTF-32 LE"
	case bytes.HasPrefix(b, UTF32BEBom):
		return "UTF-32 BE"
	default:
		if uni {
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

// marshal returns the marshaled detail data as plain or color text.
func (d *Detail) marshal(w io.Writer, color bool) error {
	if w == nil {
		w = io.Discard
	}
	const padding, width = 10, 80
	info := func(s string) string {
		return fmt.Sprintf("%s\t", s)
	}
	gookit.Enable = color
	term.Head(w, width, "File information")
	fmt.Fprintln(w)
	data := d.marshalled()
	l := len(fmt.Sprintf(" filename%s%s", strings.Repeat(" ", padding), data[0].v))
	const tabWidth = 8
	tw := tabwriter.NewWriter(w, 0, tabWidth, 0, '\t', 0)
	for _, x := range data {
		if !d.validate(x) {
			continue
		}
		if x.k == zipComment {
			if x.v != "" {
				term.HR(tw, l)
				fmt.Fprintln(tw, x.v)
				if d.sauceIndex <= 0 {
					break
				}
				// divider for sauce metadata
				term.HR(tw, l)
				continue
			}
			if d.sauceIndex <= 0 {
				break
			}
			// divider for sauce metadata
			term.HR(tw, l)
			continue
		}
		fmt.Fprintf(tw, "\t %s\t  %s\n", x.k, info(x.v))
		if x.k == cmmt {
			if d.Sauce.Comnt.Count <= 0 {
				break
			}
		}
	}
	return tw.Flush()
}

// validate reports whether the key and value data validate.
func (d Detail) validate(x struct{ k, v string }) bool {
	if !ValidText(d.Mime.Type) {
		switch x.k {
		case uc8, "line break", "characters", ans, "words", "lines", "width":
			return false
		}
	} else if x.k == ans {
		if d.Count.Controls == 0 {
			return false
		}
	}
	if x.k == "description" && x.v == "" {
		return false
	}
	if x.k == d.Sauce.Info.Info1.Info && d.Sauce.Info.Info1.Value == 0 {
		return false
	}
	if x.k == d.Sauce.Info.Info2.Info && d.Sauce.Info.Info2.Value == 0 {
		return false
	}
	if x.k == d.Sauce.Info.Info3.Info && d.Sauce.Info.Info3.Value == 0 {
		return false
	}
	if x.k == "interpretation" && x.v == "" {
		return false
	}
	return true
}

// marshalled returns the data structure used for print marshaling.
func (d Detail) marshalled() []struct{ k, v string } {
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
		{k: "line break", v: fsys.LineBreak(d.LineBreak.Decimal, true)},
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
	// Read file content
	p, err := fsys.ReadAllBytes(name)
	if err != nil {
		return err
	}
	return d.Parse(name, p...)
}

// ValidText reports whether the MIME content-type value is valid for text files.
func ValidText(mime string) bool {
	s := strings.Split(mime, "/")
	const req = 2
	if len(s) != req {
		return false
	}
	if s[0] == txt {
		return true
	}
	if mime == octetStream {
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
	w, err := fsys.Columns(f, d.LineBreak.Decimal)
	if err != nil {
		return err
	}
	if w < 0 {
		w = d.Count.Chars
	}
	d.Width = w
	return nil
}

// Words counts the number of words used in the named file.
func (d *Detail) Words(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	switch d.LineBreak.Decimal {
	case [2]rune{nl.NL}, [2]rune{nl.NEL}:
		if d.Count.Words, err = fsys.WordsEBCDIC(f); err != nil {
			return err
		}
	default:
		if d.Count.Words, err = fsys.Words(f); err != nil {
			return err
		}
	}
	return nil
}
