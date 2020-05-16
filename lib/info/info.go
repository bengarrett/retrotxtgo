package info

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
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/gookit/color"
	"github.com/mattn/go-runewidth"
	"github.com/mozillazg/go-slugify"

	humanize "github.com/labstack/gommon/bytes"
)

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

// XMLData ...
type XMLData struct {
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

// FileDate is a non-standard date format for file modifications
const FileDate string = "2 Jan 15:04 2006"

// File ...
func File(name string) (d Detail, err error) {
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

// JSON ...
func (d Detail) JSON(indent bool) (js []byte) {
	var err error
	switch indent {
	case true:
		js, err = json.MarshalIndent(d, "", "    ")
	default:
		js, err = json.Marshal(d)
	}
	logs.Check(logs.Err{"could not create", "json", err})
	return js
}

// Text ...
func (d Detail) Text(c bool) string {
	color.Enable = c
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
		{k: "UTF-8", v: fmt.Sprintf("%v", d.Utf8)},
		{k: "characters", v: fmt.Sprintf("%v", d.CharCount)},
		{k: "size", v: d.Size},
		{k: "modified", v: fmt.Sprintf("%v", d.Modified.Format(FileDate))},
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

// XML ...
func (d Detail) XML() ([]byte, error) {
	v := XMLData{
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
	return xml.MarshalIndent(v, "", "\t")
}
