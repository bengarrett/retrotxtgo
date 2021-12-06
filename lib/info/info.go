// Package info that extracts file statistics and metadata.
package info

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/info/internal/detail"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

type Detail detail.Detail

// Names index and totals.
type Names struct {
	Index  int
	Length int
}

// Info parses the named file and prints out its details in a specific syntax.
func (n Names) Info(name, format string) error {
	err1 := fmt.Sprintf("info on %s failed", name)
	if name == "" {
		return fmt.Errorf("%s: %w", err1, logs.ErrNameNil)
	}
	f, err := output(format)
	if err != nil {
		return fmt.Errorf("%s: %w", err1, logs.ErrFmt)
	}
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", err1, logs.ErrFileName)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", err1, err)
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
			return fmt.Errorf("info could not walk directory: %w", err)
		}
		return nil
	}
	if err := Marshal(name, f, n.Index, n.Length); err != nil {
		return fmt.Errorf("info on %s could not marshal: %w", name, err)
	}
	return nil
}

// output converts the --format argument value to a format type.
func output(argument string) (f detail.Format, err error) {
	switch argument {
	case "color", "c", "":
		return detail.ColorText, nil
	case "text", "t":
		return detail.PlainText, nil
	case "json", "j":
		return detail.JSON, nil
	case "json.min", "jm":
		return detail.JSONMin, nil
	case "xml", "x":
		return detail.XML, nil
	}
	return f, logs.ErrFmt
}

// Marshal the metadata and system details of a named file.
func Marshal(name string, f detail.Format, i, length int) error {
	var d detail.Detail
	if err := d.Read(name); err != nil {
		return err
	}
	if d.ValidText() {
		var err error
		// get the required linebreaks chars before running the multiple tasks
		if d.LineBreak.Decimals, err = filesystem.ReadLineBreaks(name); err != nil {
			return err
		}
		d.LineBreaks(d.LineBreak.Decimals)
		var g errgroup.Group
		g.Go(func() error {
			return d.Ctrls(name)
		})
		g.Go(func() error {
			return d.Len(name)
		})
		g.Go(func() error {
			return d.LineTotals(name)
		})
		g.Go(func() error {
			return d.Len(name)
		})
		g.Go(func() error {
			return d.Words(name)
		})
		if err := g.Wait(); err != nil {
			return err
		}
		d.MimeUnknown()
	}
	var (
		m   []byte
		err error
	)
	if m, err = d.Marshal(f); err != nil {
		return err
	}
	printf(f, m...)
	return nil
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) error {
	var d detail.Detail
	f, e := output(format)
	if e != nil {
		return e
	}
	if err := d.Parse("", nil, b...); err != nil {
		return err
	}
	if d.ValidText() {
		d.LineBreaks(filesystem.LineBreaks(true, []rune(string(b))...))
		var g errgroup.Group
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
		d.MimeUnknown()
	}
	var m []byte
	if m, e = d.Marshal(f); e != nil {
		return e
	}
	printf(f, m...)
	return nil
}

// printf prints the bytes as text and appends a newline to JSON and XML text.
func printf(f detail.Format, b ...byte) {
	switch f {
	case detail.ColorText, detail.PlainText:
		fmt.Printf("%s", b)
	case detail.JSON, detail.JSONMin, detail.XML:
		fmt.Printf("%s\n", b)
	}
}
