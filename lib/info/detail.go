package info

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
	"unicode/utf8"

	gookit "github.com/gookit/color"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/text/message"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/humanize"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sauce"
	"retrotxt.com/retrotxt/lib/str"
)

const zipComment = "zip comment"

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
	if l, err = filesystem.Lines(f, d.LineBreak.Decimals); err != nil {
		return err
	}
	d.Lines = l
	return f.Close()
}

func (d *Detail) marshal(f Format) (b []byte, err error) {
	switch f {
	case ColorText:
		return d.printMarshal(true), nil
	case PlainText:
		return d.printMarshal(false), nil
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

func (d *Detail) mimeUnknown() {
	if d.Mime.Commt == "unknown" {
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
			if !d.Utf8 && d.Count.Words > 0 {
				d.Mime.Commt = "US-ASCII encoded text document"
				return
			}
		}
	}
}

func (d *Detail) parse(name string, stat os.FileInfo, data ...byte) (err error) {
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
				fmt.Printf("mine sniffer failure, %s\n", err)
			}
			return
		}
		if d.Mime.Type == zipType {
			r, e := zip.OpenReader(name)
			if e != nil {
				fmt.Printf("open zip file failure: %s\n", e)
			}
			defer r.Close()
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
	gookit.Enable = color
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
		if x.k == zipComment {
			if x.v != "" {
				x.v = fmt.Sprintf("\n%s", x.v)
			}
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
		if x.k == zipComment || x.k == "slug" {
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
	if d.length > -1 && d.index == d.length {
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
		case "UTF-8", "line break", "characters", "ANSI controls", "words", "lines", "width":
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
	if k == zipComment && v == "" {
		return false
	}
	return true
}

func (d *Detail) linebreaks(r [2]rune) {
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
	d.LineBreak.Decimals = r
	d.LineBreak.Abbr = strings.ToUpper(a)
	d.LineBreak.Escape = e
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
		{k: "line break", v: filesystem.LineBreak(d.LineBreak.Decimals, true)},
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
		{k: zipComment, v: d.ZipComment},
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
	return d.parse(name, stat, data...)
}

// validText checks the MIME content-type value for valid text files.
func (d *Detail) validText() bool {
	s := strings.Split(d.Mime.Type, "/")
	const req = 2
	if len(s) != req {
		return false
	}
	if s[0] == "text" {
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
	if w, err = filesystem.Columns(f, d.LineBreak.Decimals); err != nil {
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
	switch d.LineBreak.Decimals {
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
