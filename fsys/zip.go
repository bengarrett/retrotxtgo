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
			return fmt.Errorf("zip walker failed with %q: %w", path, err)
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
		return err
	}
	return files.Zip(z.Writer, z.Name, z.Comment, z.Overwrite)
}

// Zip packages and compresses files to an archive using the provided name.
func (files *Files) Zip(w io.Writer, name, comment string, ow bool) error {
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
		n   string
		f   *os.File
	)
	switch ow {
	case true:
		n = name
		f, err = os.OpenFile(n, overwrite, readWriteAll)
		if err != nil {
			return fmt.Errorf("zip create %q: %w", n, err)
		}
		defer f.Close()
	default:
		n, err = UniqueName(name)
		if err != nil {
			return fmt.Errorf("zip name %q: %w", name, err)
		}
		w, err = os.OpenFile(n, mustNotExist, readWriteAll)
		if err != nil {
			return fmt.Errorf("zip create %q: %w", n, err)
		}
		defer f.Close()
	}
	zipper := zip.NewWriter(w)
	defer zipper.Close()
	if comment != "" {
		if err := zipper.SetComment(comment); err != nil {
			return fmt.Errorf("zip set comment %q: %w", comment, err)
		}
	}
	for _, fname := range *files {
		if err := InsertZip(zipper, fname); err != nil {
			return fmt.Errorf("add zip %q: %w", fname, err)
		}
	}
	if err := zipper.Close(); err != nil {
		return fmt.Errorf("zip close: %w", err)
	}
	s, err := os.Stat(n)
	if err != nil {
		return fmt.Errorf("zip could not stat %q: %w", n, err)
	}
	abs, err := filepath.Abs(s.Name())
	if err != nil {
		return fmt.Errorf("zip abs %q: %w", s.Name(), err)
	}
	fmt.Fprintln(w, "created zip file:", abs,
		humanize.Decimal(s.Size(), language.AmericanEnglish))
	return nil
}

// InsertZip adds the named file to a zip archive.
func InsertZip(z *zip.Writer, name string) error {
	if z == nil {
		return ErrWriter
	}
	s, err := os.Stat(name)
	if err != nil {
		return err
	}
	fh, err := zip.FileInfoHeader(s)
	if err != nil {
		return fmt.Errorf("file info header: %w", err)
	}
	f, err := z.CreateHeader(fh)
	if err != nil {
		return fmt.Errorf("create header: %w", err)
	}
	b, err := Read(name)
	if err != nil {
		return err
	}
	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("io writer: %w", err)
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
		return name, err
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
			return "", fmt.Errorf("unique name aborted after %d attempts: %w", maxAttempts, ErrMax)
		}
	}
}
