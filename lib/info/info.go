// Package info that extracts file statistics and metadata.
package info

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/language"
	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/sauce"
)

const (
	// DTFormat is the date-time format.
	DTFormat = "DMY24"
	// DFormat is the date format.
	DFormat     = "DMY"
	octetStream = "application/octet-stream"
	zipType     = "application/zip"
)

// Detail of a file.
type Detail struct {
	XMLName    xml.Name     `json:"-" xml:"file"`
	Name       string       `json:"filename" xml:"name"`
	Utf8       bool         `json:"utf8" xml:"utf8,attr"`
	LineBreak  LineBreaks   `json:"lineBreak" xml:"line_break"`
	Count      Stats        `json:"counts" xml:"counts"`
	Size       Sizes        `json:"size" xml:"size"`
	Lines      int          `json:"lines" xml:"lines"`
	Width      int          `json:"width" xml:"width"`
	Modified   ModDates     `json:"modified" xml:"last_modified"`
	Sums       Checksums    `json:"checksums" xml:"checksums"`
	Mime       Content      `json:"mime" xml:"mime"`
	Slug       string       `json:"slug" xml:"id,attr"`
	Sauce      sauce.Record `json:"sauce" xml:"sauce"`
	ZipComment string       `json:"zipComment" xml:"zip_comment"`
	index      int
	length     int
	sauceIndex int
}

// LineBreaks for new line toggles.
type LineBreaks struct {
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
	CRC32  string `json:"CRC32" xml:"CRC32"`
	CRC64  string `json:"CRC64" xml:"CRC64"`
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

// Format of the text to output.
type Format uint

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
	f, err := output(format)
	if err != nil {
		gen.Err = ErrFmt
		return gen
	}
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		gen.Err = ErrNoFile
		return gen
	}
	if err != nil {
		gen.Err = err
		return gen
	}
	if s.IsDir() {
		const walkMode = -1
		// godirwalk.Walk is more performant than the standard library filepath.Walk
		err := godirwalk.Walk(name, &godirwalk.Options{
			Callback: func(osPathname string, de *godirwalk.Dirent) error {
				if skip, err := de.IsDirOrSymlinkToDir(); err != nil {
					return err
				} else if skip {
					return nil
				}
				return Marshal(osPathname, f, walkMode, walkMode)
			},
			ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
				return godirwalk.SkipNode
			},
			Unsorted: true, // set true for faster yet non-deterministic enumeration
		})
		if err != nil {
			gen.Issue = "info.print.directory"
			gen.Arg = format
			gen.Err = err
			return gen
		}
		return logs.Generic{}
	}
	if err := Marshal(name, f, n.Index, n.Length); err != nil {
		gen.Issue = "info.print"
		gen.Arg = format
		gen.Err = err
		return gen
	}
	return gen
}

// Language tag used for numeric syntax formatting.
func lang() language.Tag {
	return language.English
}

// output converts the --format argument value to a format type.
func output(argument string) (f Format, err error) {
	switch argument {
	case "color", "c", "":
		return ColorText, nil
	case "text", "t":
		return PlainText, nil
	case "json", "j":
		return JSON, nil
	case "json.min", "jm":
		return JSONMin, nil
	case "xml", "x":
		return XML, nil
	}
	return f, ErrFmt
}

// Marshal the meta and operating system details of a file.
func Marshal(filename string, f Format, i, length int) error {
	var d Detail
	if err := d.read(filename); err != nil {
		return err
	}
	d.index, d.length = i, length
	if d.validText() {
		var g errgroup.Group
		g.Go(func() error {
			var err error
			if d.LineBreak.Decimals, err = filesystem.ReadLineBreaks(filename); err != nil {
				return err
			}
			d.linebreaks(d.LineBreak.Decimals)
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
	if m, err = d.marshal(f); err != nil {
		return err
	}
	printf(f, m...)
	return nil
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) error {
	var d Detail
	f, e := output(format)
	if e != nil {
		return e
	}
	if err := d.parse("", nil, b...); err != nil {
		return err
	}
	if d.validText() {
		var g errgroup.Group
		g.Go(func() error {
			d.linebreaks(filesystem.LineBreaks(true, []rune(string(b))...))
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
			if d.Lines, err = filesystem.Lines(bytes.NewReader(b), d.LineBreak.Decimals); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Width, err = filesystem.Columns(bytes.NewReader(b), d.LineBreak.Decimals); err != nil {
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
	var m []byte
	if m, e = d.marshal(f); e != nil {
		return e
	}
	printf(f, m...)
	return nil
}

func printf(f Format, b ...byte) {
	switch f {
	case ColorText, PlainText:
		fmt.Printf("%s", b)
	case JSON, JSONMin, XML:
		fmt.Printf("%s\n", b)
	}
}
