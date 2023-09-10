// Package info that extracts file statistics and metadata.
package info

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/info/internal/detail"
	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

var (
	ErrName = errors.New("name value cannot be empty")
	ErrFmt  = errors.New("format is not known")
)

type Detail detail.Detail

// Names index and totals.
type Names struct {
	Index  int
	Length int
}

// Info parses the named file and prints out its details in a specific syntax.
func (n Names) Info(name, format string) (string, error) {
	failure := fmt.Sprintf("info on %s failed", name)
	if name == "" {
		return "", ErrName
	}
	f, err := output(format)
	if err != nil {
		return "", err
	}
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s: %w", failure, err)
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", failure, err)
	}
	if !s.IsDir() {
		res, err := Marshal(name, f)
		if err != nil {
			return "", fmt.Errorf("%s: %w", failure, err)
		}
		return res, nil
	}
	// godirwalk.Walk is more performant than the standard library filepath.Walk
	err = godirwalk.Walk(name, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if skip, err := de.IsDirOrSymlinkToDir(); err != nil {
				return err
			} else if skip {
				return nil
			}
			_, err := Marshal(osPathname, f)
			return err
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
		Unsorted: true, // set true for faster yet non-deterministic enumeration
	})
	if err != nil {
		return "", fmt.Errorf("info could not walk directory: %w", err)
	}
	return "", nil
}

// output converts the --format argument value to a format type.
func output(argument string) (detail.Format, error) {
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
	return -1, fmt.Errorf("%w: %s", ErrFmt, argument)
}

// Marshal the metadata and system details of a named file.
func Marshal(name string, f detail.Format) (string, error) {
	var d detail.Detail
	if err := d.Read(name); err != nil {
		return "", err
	}
	if d.ValidText() {
		var err error
		// get the required linebreaks chars before running the multiple tasks
		if d.LineBreak.Decimals, err = fsys.ReadLineBreaks(name); err != nil {
			return "", err
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
			return "", err
		}
		d.MimeUnknown()
	}
	var (
		m   []byte
		err error
	)
	if m, err = d.Marshal(f); err != nil {
		return "", err
	}
	return printf(f, m...), nil
}

// Stdin parses piped data and prints out the details in a specific syntax.
func Stdin(format string, b ...byte) (string, error) {
	var d detail.Detail
	f, e := output(format)
	if e != nil {
		return "", e
	}
	if err := d.Parse("", nil, b...); err != nil {
		return "", err
	}
	if d.ValidText() { //nolint:nestif
		d.LineBreaks(fsys.LineBreaks(true, []rune(string(b))...))
		var g errgroup.Group
		g.Go(func() error {
			var err error
			if d.Count.Controls, err = fsys.Controls(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Count.Controls, err = fsys.Controls(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Lines, err = fsys.Lines(bytes.NewReader(b), d.LineBreak.Decimals); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Width, err = fsys.Columns(bytes.NewReader(b), d.LineBreak.Decimals); err != nil {
				return err
			} else if d.Width < 0 {
				d.Width = d.Count.Chars
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Count.Words, err = fsys.Words(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		if err := g.Wait(); err != nil {
			return "", err
		}
		d.MimeUnknown()
	}
	var m []byte
	if m, e = d.Marshal(f); e != nil {
		return "", e
	}
	return printf(f, m...), nil
}

// printf prints the bytes as text and appends a newline to JSON and XML text.
func printf(f detail.Format, b ...byte) string {
	switch f {
	case detail.ColorText, detail.PlainText:
		return string(b)
	case detail.JSON, detail.JSONMin, detail.XML:
		return string(b) + "\n"
	}
	return ""
}
