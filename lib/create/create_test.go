package create_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/assets"
	"github.com/bengarrett/retrotxtgo/lib/create/internal/layout"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleColorScheme() {
	fmt.Print(create.ColorScheme()[0])
	// Output: normal
}
func ExampleReferrer() {
	fmt.Print(create.Referrer()[1])
	// Output: origin
}

func ExampleRobots() {
	fmt.Print(create.Robots()[2])
	// Output: follow
}

func TestSaveAssets(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := filepath.Join(os.TempDir(), "retrotxt_example_save_assets")
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			t.Errorf("saveAssets make temp dir: %q: %w", tmpDir, err)
		}
		defer os.RemoveAll(tmpDir)
		// Initialize
		a := create.Args{
			Test:    true,
			Layouts: layout.Compact,
		}
		a.Save.Destination = tmpDir
		// Save files
		b := []byte("hello")
		if err := a.SaveAssets(&b); err != nil {
			t.Errorf("saveAssets: %w", err)
		}
		// Count the saved files in the temporary directory
		files, err := ioutil.ReadDir(tmpDir)
		if err != nil {
			t.Errorf("saveAssets read dir: %q: %w", tmpDir, err)
		}
		const zero = 0

		if got := len(files); got == zero {
			t.Errorf("SaveAssets() file count = %v", got)
		}
	})
}

func TestZipAssets(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := filepath.Join(os.TempDir(), "retrotxt_example_save_assets")
		os.RemoveAll(tmpDir)
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			t.Errorf("saveAssets make temp dir: %w", err)
		}
		defer os.RemoveAll(tmpDir)
		// Initialize
		a := create.Args{
			Layouts: layout.Standard,
			Test:    true,
		}
		a.Save.Destination = tmpDir
		// Create a zip file
		name := filepath.Join(os.TempDir(), layout.ZipName)
		b := []byte("hello")
		a.ZipAssets(os.TempDir(), &b)
		defer os.Remove(name)
		// Print the filename of the new zip file
		file, err := os.Stat(name)
		if err != nil {
			t.Errorf("stat file: %w", err)
		}
		const want = "retrotxt.zip"
		if got := file.Name(); got != want {
			t.Errorf("ZipAssets() filename = %v, want %v", got, want)
		}
	})
}

func TestSave(t *testing.T) {
	type args struct {
		data []byte
		name string
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Save() user home error: %w", err)
	}
	tmpDir := os.TempDir()
	tmpFile := path.Join(tmpDir, "retrotxt_create_test.txt")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data", args{[]byte(""), ""}, true},
		{"invalid", args{[]byte("abc"), "this-is-an-invalid-path"}, true},
		{"tempDir", args{data: []byte("abc"), name: tmpFile}, true},
		{"homeDir", args{[]byte("abc"), homeDir}, false},
		{"currentDir", args{[]byte("abc"), "."}, false},
		{"path as name", args{[]byte("abc"), tmpDir}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan error)
			a := create.Args{Layouts: layout.Standard, Test: true}
			a.Save.OW = true
			a.Save.Destination = tt.args.name
			go a.SaveHTML(&tt.args.data, ch)
			err := <-ch
			if (err != nil) != tt.wantErr {
				fmt.Println("TestSave dir:", tmpDir)
				t.Errorf("Save(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
	// clean-up
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, create.HtmlFn.Write())
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}
}

func TestArgsStdout(t *testing.T) {
	var (
		a  = create.Args{Layouts: layout.Standard}
		b  = []byte("")
		hi = []byte("hello world")
	)
	tests := []struct {
		name    string
		args    create.Args
		b       *[]byte
		wantErr bool
	}{
		{"no data", a, &b, false},
		{"hi", a, &hi, false},
		{"nil", a, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.args.Stdout(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("Args.Stdout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	l := create.Layouts()
	if got := len(l); got != 4 {
		t.Errorf("Templates().Keys() = %v, want %v", got, 4)
	}
	if got := create.Layouts()[3]; got != "none" {
		t.Errorf("Templates().Keys() = %v, want %v", got, "none")
	}
}

func TestTemplateSave(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmplsave")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())
	a := create.Args{
		Layouts: layout.Standard,
		Tmpl:    tmpFile.Name(),
	}
	if err = a.TemplateSave(); err != nil {
		t.Errorf("TemplateSave() created an error: %w", err)
	}
}

func TestArgs_Marshal(t *testing.T) {
	ex := fmt.Sprintf("%s | example", meta.Name)
	viper.SetDefault("html.title", ex)
	args := create.Args{Layouts: layout.Standard}
	w := "hello"
	d := []byte(w)
	got, _ := args.Marshal(&d)
	if got.PreText != w {
		t.Errorf("Marshal().PreText = %v, want %v", got, w)
	}
	args.Layouts = layout.Compact
	w = ex
	got, _ = args.Marshal(&d)
	if got.PageTitle != w {
		t.Errorf("Marshal().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got, _ = args.Marshal(&d)
	if got.MetaDesc != w {
		t.Errorf("Marshal().MetaDesc = %v, want %v", got, w)
	}
	args.Layouts = layout.Standard
	w = ""
	got, _ = args.Marshal(&d)
	if got.MetaAuthor != w {
		t.Errorf("Marshal().MetaAuthor = %v, want %v", got, w)
	}
	args.Layouts = layout.Inline
	w = ""
	got, _ = args.Marshal(&d)
	if got.MetaAuthor != w {
		t.Errorf("Marshal().MetaAuthor = %v, want %v", got, w)
	}
}

func TestDestination(t *testing.T) {
	saved := viper.GetString("save-directory")
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	spaces := filepath.Join(home, "some directory", "some file.html")
	root, _ := filepath.Abs("/")
	sub := filepath.Clean(filepath.Join(home, "html", "example.htm"))
	winI, winO := "/", "/"
	if runtime.GOOS == "windows" {
		winI = "c:\\"
		winO = "\\"
	}
	tests := []struct {
		name     string
		args     []string
		wantPath string
		wantErr  bool
	}{
		{"empty", []string{}, saved, false},
		{"cwd", []string{"."}, wd, false},
		{"home", []string{"~"}, home, false},
		{"root", []string{"/"}, root, false},
		{"file", []string{"./example.html"}, "example.html", false},
		{"subdir", []string{"~/html/example.htm"}, sub, false},
		{"spaces", []string{"~/some directory/some file.html"}, spaces, false},
		{"windows", []string{winI}, winO, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := assets.Destination(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Destination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Destination() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}

func TestNormalize(t *testing.T) {
	hi := []rune("hello world")
	cp437 := charmap.CodePage437
	empty := []byte("")
	type args struct {
		e encoding.Encoding
		r []rune
	}
	tests := []struct {
		name  string
		args  args
		wantB []byte
	}{
		{"empty", args{nil, nil}, empty},
		{"no enc", args{nil, hi}, []byte(string(hi))},
		{"enc", args{cp437, hi}, []byte(string(hi))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotB := create.Normalize(tt.args.e, tt.args.r...); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("Normalize() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
