package filesystem

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
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
	// Quiet suppresses all non-error messages.
	Quiet bool
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
	return files.Zip(z.Name, z.Comment, z.Overwrite, z.Quiet)
}

// Zip packages and compresses files to an archive using the provided name.
func (files *Files) Zip(name, comment string, ow, quiet bool) error {
	const (
		overwrite    = os.O_RDWR | os.O_CREATE
		mustNotExist = os.O_RDWR | os.O_CREATE | os.O_EXCL
		readWriteAll = 0o666
	)
	var (
		err error
		n   string
		w   *os.File
	)
	switch ow {
	case true:
		n = name
		w, err = os.OpenFile(n, overwrite, readWriteAll)
		if err != nil {
			return fmt.Errorf("zip create %q: %w", n, err)
		}
		defer w.Close()
	default:
		n, err = UniqueName(name)
		if err != nil {
			return fmt.Errorf("zip name %q: %w", name, err)
		}
		w, err = os.OpenFile(n, mustNotExist, readWriteAll)
		if err != nil {
			return fmt.Errorf("zip create %q: %w", n, err)
		}
		defer w.Close()
	}
	z := zip.NewWriter(w)
	defer z.Close()
	if comment != "" {
		if err = z.SetComment(comment); err != nil {
			return fmt.Errorf("zip set comment %q: %w", comment, err)
		}
	}
	for _, f := range *files {
		err = AddZip(f, z)
		if err != nil {
			return fmt.Errorf("add zip %q: %w", f, err)
		}
	}
	err = z.Close()
	if err != nil {
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
	if !quiet {
		fmt.Println("created zip file:", abs, humanize.Decimal(s.Size(), language.AmericanEnglish))
	}
	return nil
}

func AddZip(name string, z *zip.Writer) error {
	s, err := os.Stat(name)
	if err != nil {
		fmt.Println("skipping file, could not stat", name)
		return nil //nolint:nilerr
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
		fmt.Println("skipping file, could not read", name)
		return nil //nolint:nilerr
	}
	_, err = f.Write(b)
	if err != nil {
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
		return "", fmt.Errorf("unique name is a directory %q: %w", name, logs.ErrFileSaveD)
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
			return "", fmt.Errorf("unique name aborted after %d attempts: %w", maxAttempts, logs.ErrMax)
		}
	}
}
