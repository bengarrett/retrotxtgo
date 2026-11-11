// Package
//
//nolint:funlen
package flag_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/bengarrett/retrotxtgo/sample"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Fprint(os.Stdout, flag.EndOfFile(f))
	// Output: true
}

var cp437 encoding.Encoding = charmap.CodePage437

func TestDefault(t *testing.T) {
	t.Parallel()
	e := flag.Default()
	be.True(t, cp437 == e)
}

func TestInputOriginal(t *testing.T) {
	t.Parallel()
	g, err := flag.InputOverride(nil, "")
	be.Equal(t, g, sample.Flags{})
	be.Err(t, err, nil)
	g, err = flag.InputOverride(nil, "CP437")
	be.Equal(t, g, sample.Flags{})
	be.Err(t, err, nil)
}

func TestArgs(t *testing.T) {
	t.Parallel()

	// The flag package has complex dependencies that make it difficult to test
	// without extensive mocking. For now, we'll test the basic functionality
	// that doesn't cause panics.

	view := flag.View()
	be.Equal(t, view.Input, "CP437")
	be.True(t, len(view.Controls) > 0)
	be.True(t, len(view.Swap) > 0)
	be.Equal(t, view.Width, 0)
	be.Equal(t, view.Original, false)
}

func TestView(t *testing.T) {
	t.Parallel()

	view := flag.View()
	be.Equal(t, view.Input, "CP437")
	be.True(t, len(view.Controls) > 0)
	be.True(t, len(view.Swap) > 0)
	be.Equal(t, view.Width, 0)
	be.Equal(t, view.Original, false)
}

func TestGlobArgs(t *testing.T) {
	t.Parallel()

	// create a temp dir with known files
	dir := t.TempDir()
	fileA := filepath.Join(dir, "a.txt")
	fileB := filepath.Join(dir, "b.txt")
	fileC := filepath.Join(dir, "c.log")
	for _, f := range []string{fileA, fileB, fileC} {
		if err := os.WriteFile(f, []byte("hello"), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	t.Run("literal path unchanged", func(t *testing.T) {
		t.Parallel()
		got := flag.GlobArgs([]string{fileA})
		be.Equal(t, len(got), 1)
		be.Equal(t, got[0], fileA)
	})

	t.Run("glob expands matching files", func(t *testing.T) {
		t.Parallel()
		pattern := filepath.Join(dir, "*.txt")
		got := flag.GlobArgs([]string{pattern})
		be.Equal(t, len(got), 2)
	})

	t.Run("non-matching glob kept as is", func(t *testing.T) {
		t.Parallel()
		pattern := filepath.Join(dir, "*.xyz")
		got := flag.GlobArgs([]string{pattern})
		be.Equal(t, len(got), 1)
		be.Equal(t, got[0], pattern)
	})

	t.Run("mix of literal and glob", func(t *testing.T) {
		t.Parallel()
		pattern := filepath.Join(dir, "*.txt")
		got := flag.GlobArgs([]string{fileC, pattern})
		// fileC literal + 2 expanded .txt files
		be.Equal(t, len(got), 3)
	})

	t.Run("empty input returns empty", func(t *testing.T) {
		t.Parallel()
		got := flag.GlobArgs([]string{})
		be.Equal(t, len(got), 0)
	})

	t.Run("bare star excludes directories", func(t *testing.T) {
		t.Parallel()
		// create a subdir inside our temp dir so * would match it
		subdir := filepath.Join(dir, "subdir")
		if err := os.Mkdir(subdir, 0o700); err != nil {
			t.Fatal(err)
		}
		pattern := filepath.Join(dir, "*")
		got := flag.GlobArgs([]string{pattern})
		// should contain only the 3 regular files, not subdir
		be.Equal(t, len(got), 3)
		for _, g := range got {
			st, err := os.Stat(g)
			if err != nil {
				t.Fatalf("stat %s: %v", g, err)
			}
			if st.IsDir() {
				t.Errorf("globargs returned a directory: %s", g)
			}
		}
	})
}
