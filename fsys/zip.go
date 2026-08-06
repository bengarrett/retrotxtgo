package fsys

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bengarrett/sauce/humanize"
	"golang.org/x/text/language"
)

// Files to zip.
type Files []string

// Zip archive details.
type Zip struct {
	// Zip path and filename.
	Name string
	// Root path of the directory to archive.
	Root string
	// Comment to embed.
	Comment string
	// Overwrite an existing named zip file if encountered.
	Overwrite bool
	// Writer for all the non-error messages, or use io.Discard to suppress.
	Writer io.Writer
}

// Create zip packages and compresses files contained the root directory into an archive using the provided name.
func (z *Zip) Create() error {
	const dotFile = "."
	files := Files{}
	walker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			const format = "zip walker failed with %q: %w"
			return fmt.Errorf(format, path, err)
		}
		if info.IsDir() && info.Name() != filepath.Base(path) {
			return filepath.SkipDir
		}
		// ignore directories because there is no recursive walking
		if info.IsDir() {
			return nil
		}
		// stop recursive walking
		if filepath.Base(filepath.Dir(path)) != filepath.Base(z.Root) {
			return nil
		}
		// ignore posix hidden files
		if info.Name()[:1] == dotFile {
			return nil
		}
		// ignore 0-byte files
		if info.Size() == 0 {
			return nil
		}
		files = append(files, path)
		return nil
	}
	if err := filepath.Walk(z.Root, walker); err != nil {
		const format = "zip create: %w"
		return fmt.Errorf(format, err)
	}
	return files.Zip(z.Writer, z.Name, z.Comment, z.Overwrite)
}

// Zip packages and compresses files to an archive using the provided name.
func (files *Files) Zip(w io.Writer, name, comment string, ow bool) error {
	const format = "zip %s %s: %w"
	if w == nil {
		w = io.Discard
	}
	const (
		overwrite    = os.O_RDWR | os.O_CREATE
		mustNotExist = os.O_RDWR | os.O_CREATE | os.O_EXCL
		readWriteAll = 0o666
	)
	var (
		err error
		unq string
		f   *os.File
	)
	switch ow {
	case true:
		f, err = os.OpenFile(name, overwrite, readWriteAll)
		if err != nil {
			return fmt.Errorf(format, "open", name, err)
		}
		defer f.Close()
	default:
		unq, err = UniqueName(name)
		if err != nil {
			return fmt.Errorf(format, "unique name", name, err)
		}
		w, err = os.OpenFile(unq, mustNotExist, readWriteAll)
		if err != nil {
			return fmt.Errorf(format, "create", name, err)
		}
		defer f.Close()
	}
	zipper := zip.NewWriter(w)
	defer zipper.Close()
	if comment != "" {
		if err := zipper.SetComment(comment); err != nil {
			return fmt.Errorf(format, "set comment", comment, err)
		}
	}
	for _, fname := range *files {
		if err := InsertZip(zipper, fname); err != nil {
			return fmt.Errorf(format, "zip", fname, err)
		}
	}
	if err := zipper.Close(); err != nil {
		return fmt.Errorf(format, "close", "zip writer", err)
	}
	s, err := os.Stat(unq)
	if err != nil {
		return fmt.Errorf(format, "stat", unq, err)
	}
	abs, err := filepath.Abs(s.Name())
	if err != nil {
		return fmt.Errorf(format, "filepath abs", s.Name(), err)
	}
	fmt.Fprintln(w, "created zip file:", abs,
		humanize.Decimal(s.Size(), language.AmericanEnglish))
	return nil
}

// InsertZip adds the named file to a zip archive.
func InsertZip(z *zip.Writer, name string) error {
	const format = "insert zip %s: %w"
	if z == nil {
		return ErrWriter
	}
	s, err := os.Stat(name)
	if err != nil {
		return fmt.Errorf(format, "stat", err)
	}
	fh, err := zip.FileInfoHeader(s)
	if err != nil {
		return fmt.Errorf(format, "file info header", err)
	}
	f, err := z.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf(format, "create header", err)
	}
	b, err := Read(name)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf(format, "io writer", err)
	}
	return nil
}

// UniqueName confirms the file name doesn't conflict with an existing file.
// If there is a conflict, a new incremental name will be returned.
func UniqueName(name string) (string, error) {
	const (
		maxAttempts = 9999
		macOS       = "darwin"
		windows     = "windows"
	)
	s, err := os.Stat(name)
	if os.IsNotExist(err) {
		return name, nil
	}
	if err != nil {
		const format = "zip unique name: %w"
		return name, fmt.Errorf(format, err)
	}
	if s.IsDir() {
		return "", fmt.Errorf("%q: %w", name, ErrName)
	}
	i := 1
	for {
		dir, file := path.Split(name)
		e := path.Ext(file)
		b := strings.TrimSuffix(file, e)
		var n string
		switch runtime.GOOS {
		case macOS:
			n = fmt.Sprintf("%s %d%s", b, i, e)
		case windows:
			n = fmt.Sprintf("%s (%d)%s", b, i, e)
		default:
			n = fmt.Sprintf("%s_%d%s", b, i, e)
		}
		p := filepath.Join(dir, n)
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			return p, nil
		}
		i++
		if i > maxAttempts {
			const format = "unique name aborted after %d attempts: %w"
			return "", fmt.Errorf(format, maxAttempts, ErrMax)
		}
	}
}
