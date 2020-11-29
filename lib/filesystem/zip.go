package filesystem

import (
	"archive/zip"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/language"
	"retrotxt.com/retrotxt/lib/humanize"
)

// Files to compress and archive.
type Files []string

// Zip packages and compresses files contained in root to an archive using the provided name.
func Zip(name, root, comment string) error {

	var files Files

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
		if filepath.Base(filepath.Dir(path)) != filepath.Base(root) {
			return nil
		}
		// ignore posix hidden files
		if info.Name()[:1] == "." {
			return nil
		}

		files = append(files, path)
		return nil
	}

	if err := filepath.Walk(root, walker); err != nil {
		return err
	}

	return files.Zip(name, comment)
}

// Zip packages and compresses files to an archive using the provided name.
func (files *Files) Zip(name, comment string) error {

	new, err := UniqueName(name)
	if err != nil {
		return fmt.Errorf("zip name %q: %w", name, err)
	}

	w, err := os.Create(new)
	if err != nil {
		return fmt.Errorf("zip create %q: %w", new, err)
	}
	defer w.Close()

	z := zip.NewWriter(w)
	defer z.Close()

	if comment != "" {
		if err = z.SetComment(comment); err != nil {
			return fmt.Errorf("zip set comment %q: %w", comment, err)
		}
	}

	for _, f := range *files {
		err = zipper(f, z)
		if err != nil {
			return fmt.Errorf("zipper %q: %w", f, err)
		}
	}

	err = z.Close()
	if err != nil {
		return fmt.Errorf("zip close: %w", err)
	}

	s, err := os.Stat(new)
	if err != nil {
		return fmt.Errorf("zip could not stat %q: %w", new, err)
	}
	abs, err := filepath.Abs(s.Name())
	if err != nil {
		return fmt.Errorf("zip abs %q: %w", s.Name(), err)
	}

	fmt.Println("created zip file:", abs, humanize.Decimal(s.Size(), language.AmericanEnglish))
	return nil
}

func zipper(name string, z *zip.Writer) error {

	s, err := os.Stat(name)
	if err != nil {
		fmt.Println("skipping file, could not stat", name)
		return nil
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
		return nil
	}

	_, err = f.Write(b)
	if err != nil {
		return fmt.Errorf("io writer: %w", err)
	}
	return nil
}

// UniqueName returns...
func UniqueName(name string) (string, error) {

	const (
		maxAttempts = 100
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
		return "", fmt.Errorf("increment name %q: is a directory", name)
	}

	i := 1
	for {

		dir, file := path.Split(name)
		e := path.Ext(file)
		b := strings.TrimSuffix(file, e)
		new := ""

		switch runtime.GOOS {
		case macOS:
			new = fmt.Sprintf("%s %d%s", b, i, e)
		case windows:
			new = fmt.Sprintf("%s (%d)%s", b, i, e)
		default:
			new = fmt.Sprintf("%s_%d%s", b, i, e)
		}

		path := filepath.Join(dir, new)
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return path, nil
		}

		i++
		if i > maxAttempts {
			return "", fmt.Errorf("increment name aborted after %d attempts", maxAttempts)
		}
	}
}
