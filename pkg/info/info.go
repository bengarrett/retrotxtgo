// Package info that extracts and returns file statistics and metadata.
package info

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/nl"
	"github.com/karrick/godirwalk"
	"golang.org/x/sync/errgroup"
)

var (
	ErrFmt  = errors.New("format is not known")
	ErrName = errors.New("name value cannot be empty")
)

// Info parses the named file and writes the details in a formal syntax.
func Info(w io.Writer, name, format string) error {
	if w == nil {
		w = io.Discard
	}
	failure := fmt.Sprintf("info on %s failed", name)
	if name == "" {
		return ErrName
	}
	f, err := output(format)
	if err != nil {
		return err
	}
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: %w", failure, err)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", failure, err)
	}
	if !s.IsDir() {
		if err := Marshal(w, name, f); err != nil {
			return fmt.Errorf("%s: %w", failure, err)
		}
		return nil
	}
	// godirwalk.Walk is more performant than the standard library filepath.Walk
	err = godirwalk.Walk(name, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			if skip, err := de.IsDirOrSymlinkToDir(); err != nil {
				return err
			} else if skip {
				return nil
			}
			return Marshal(w, osPathname, f)
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

// output converts the argument string to a format type.
func output(arg string) (Format, error) {
	switch arg {
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
	return -1, fmt.Errorf("%w: %s", ErrFmt, arg)
}

// Marshal and write the metadata and system details of a named file.
func Marshal(w io.Writer, name string, f Format) error {
	if w == nil {
		w = io.Discard
	}
	var d Detail
	if err := d.Read(name); err != nil {
		return err
	}
	if ValidText(d.Mime.Type) {
		var err error
		// get the required linebreaks chars before running the multiple tasks
		if d.LineBreak.Decimal, err = fsys.ReadLineBreaks(name); err != nil {
			return err
		}
		d.LineBreak.Find(d.LineBreak.Decimal)
		g := errgroup.Group{}
		g.Go(func() error {
			return d.Ctrls(name)
		})
		g.Go(func() error {
			return d.Len(name)
		})
		g.Go(func() error {
			i, err := d.LineBreak.Total(name)
			if err != nil {
				return err
			}
			d.Lines = i
			return nil
		})
		g.Go(func() error {
			return d.Words(name)
		})
		if err := g.Wait(); err != nil {
			return err
		}
		d.MimeUnknown()
	}
	if err := d.Marshal(w, f); err != nil {
		return err
	}
	printnl(w, f)
	return nil
}

// Stream parses piped data and writes out the details in a specific syntax.
func Stream(w io.Writer, format string, b ...byte) error {
	if w == nil {
		w = io.Discard
	}
	var d Detail
	f, e := output(format)
	if e != nil {
		return e
	}
	if err := d.Parse("", b...); err != nil {
		return err
	}
	if ValidText(d.Mime.Type) { //nolint:nestif
		d.LineBreak.Find(fsys.LineBreaks(true, []rune(string(b))...))
		g := errgroup.Group{}
		g.Go(func() error {
			var err error
			if d.Count.Controls, err = fsys.Controls(bytes.NewReader(b)); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Lines, err = nl.Lines(bytes.NewReader(b), d.LineBreak.Decimal); err != nil {
				return err
			}
			return nil
		})
		g.Go(func() error {
			var err error
			if d.Width, err = fsys.Columns(bytes.NewReader(b), d.LineBreak.Decimal); err != nil {
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
			return err
		}
		d.MimeUnknown()
	}
	if err := d.Marshal(w, f); err != nil {
		return err
	}
	printnl(w, f)
	return nil
}

// printnl appends a newline to JSON and XML text.
func printnl(w io.Writer, f Format) {
	if w == nil {
		w = io.Discard
	}
	switch f {
	case ColorText, PlainText:
		return
	case JSON, JSONMin, XML:
		fmt.Fprintln(w)
		return
	}
}
