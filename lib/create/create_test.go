package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleColorScheme() {
	fmt.Print(ColorScheme()[0])
	// Output: normal
}
func ExampleReferrer() {
	fmt.Print(Referrer()[1])
	// Output: origin
}

func ExampleRobots() {
	fmt.Print(Robots()[2])
	// Output: follow
}

func Test_saveAssets(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := filepath.Join(os.TempDir(), "retrotxt_example_save_assets")
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			t.Errorf("saveAssets make temp dir: %q: %w", tmpDir, err)
		}
		defer os.RemoveAll(tmpDir)
		// Initialize
		a := Args{}
		a.Save.Destination = tmpDir
		a.Test = true
		// Save files
		b := []byte("hello")
		if err := a.saveAssets(&b); err != nil {
			t.Errorf("saveAssets: %w", err)
		}
		// Count the saved files in the temporary directory
		files, err := ioutil.ReadDir(tmpDir)
		if err != nil {
			t.Errorf("saveAssets read dir: %q: %w", tmpDir, err)
		}
		const zero = 0

		if got := len(files); got == zero {
			t.Errorf("saveAssets() file count = %v", got)
		}
	})
}

func Test_zipAssets(t *testing.T) {
	t.Run("comment", func(t *testing.T) {
		// Create a temporary directory
		tmpDir := filepath.Join(os.TempDir(), "retrotxt_example_save_assets")
		os.RemoveAll(tmpDir)
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			t.Errorf("saveAssets make temp dir: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		// Initialize
		a := Args{}
		a.layout = Standard
		a.Save.Destination = tmpDir
		a.Test = true

		// Create a zip file
		name := filepath.Join(os.TempDir(), zipName)
		b := []byte("hello")
		a.zipAssets(os.TempDir(), &b)
		defer os.Remove(name)

		// Print the filename of the new zip file
		file, err := os.Stat(name)
		if err != nil {
			t.Errorf("stat file: %w", err)
		}

		const want = "retrotxt.zip"
		if got := file.Name(); got != want {
			t.Errorf("zipAssets() filename = %v, want %v", got, want)
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
			a := Args{layout: Standard, Test: true}
			a.Save.OW = true
			a.Save.Destination = tt.args.name
			go a.saveHTML(&tt.args.data, ch)
			err := <-ch
			if (err != nil) != tt.wantErr {
				fmt.Println("dir:", tmpDir)
				t.Errorf("Save(%s) error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
	// clean-up
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, "index.html")
		if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
			t.Error(err)
		}
	}
}

func TestArgsStdout(t *testing.T) {
	var (
		a  = Args{layout: Standard}
		b  = []byte("")
		hi = []byte("hello world")
	)
	tests := []struct {
		name    string
		args    Args
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
	l := Layouts()
	if got := len(l); got != 4 {
		t.Errorf("Templates().Keys() = %v, want %v", got, 4)
	}
	if got := Layouts()[3]; got != "none" {
		t.Errorf("Templates().Keys() = %v, want %v", got, "none")
	}
}

func TestTemplates(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", "unknown"},
		{"none", "none", "none"},
		{"standard", "standard", "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, _ := layout(tt.key)
			if got := l.Pack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("layout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_templateSave(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmplsave")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tmpFile.Name())
	a := Args{
		layout: Standard,
		tmpl:   tmpFile.Name(),
	}
	if err = a.templateSave(); err != nil {
		t.Errorf("templateSave() created an error: %w", err)
	}
}

func Test_marshal(t *testing.T) {
	ex := fmt.Sprintf("%s | example", meta.Name)
	viper.SetDefault("html.title", ex)
	args := Args{layout: Standard}
	w := "hello"
	d := []byte(w)
	got, _ := args.marshal(&d)
	if got.PreText != w {
		t.Errorf("marshal().PreText = %v, want %v", got, w)
	}
	args.layout = Compact
	w = ex
	got, _ = args.marshal(&d)
	if got.PageTitle != w {
		t.Errorf("marshal().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got, _ = args.marshal(&d)
	if got.MetaDesc != w {
		t.Errorf("marshal().MetaDesc = %v, want %v", got, w)
	}
	args.layout = Standard
	w = ""
	got, _ = args.marshal(&d)
	if got.MetaAuthor != w {
		t.Errorf("marshal().MetaAuthor = %v, want %v", got, w)
	}
	args.layout = Inline
	w = ""
	got, _ = args.marshal(&d)
	if got.MetaAuthor != w {
		t.Errorf("marshal().MetaAuthor = %v, want %v", got, w)
	}
}

func Test_destination(t *testing.T) {
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
			gotPath, err := destination(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("destination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("destination() = %v, want %v", gotPath, tt.wantPath)
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
			if gotB := Normalize(tt.args.e, tt.args.r...); !reflect.DeepEqual(gotB, tt.wantB) {
				t.Errorf("Normalize() = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}
