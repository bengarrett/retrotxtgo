package create

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/viper"
)

func Test_Save(t *testing.T) {
	type args struct {
		data []byte
		name string
	}
	tmpFile := path.Join(os.TempDir(), "retrotxt_create_test.txt")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data", args{[]byte(""), ""}, true},
		{"invalid", args{[]byte("abc"), "this-is-an-invalid-path"}, true},
		{"tempDir", args{data: []byte("abc"), name: tmpFile}, true},
		{"homeDir", args{[]byte("abc"), "~"}, false},
		{"currentDir", args{[]byte("abc"), "."}, false},
		{"path as name", args{[]byte("abc"), os.TempDir()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := make(chan error)
			a := Args{Layout: "standard", Test: true}
			a.OW = true
			a.Dest = tt.args.name
			go a.savehtml(&tt.args.data, ch)
			err := <-ch
			if (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
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

func TestArgs_Stdout(t *testing.T) {
	var (
		a = Args{
			Layout: "standard",
		}
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

func Test_Layouts(t *testing.T) {
	l := Layouts()
	if got := len(l); got != 5 {
		t.Errorf("createTemplates() = %v, want %v", got, 5)
	}
	if got := createTemplates()["body"]; got != "body-content" {
		t.Errorf("createTemplates() = %v, want %v", got, "body-content")
	}
}

func Test_createTemplates(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"body", "body", "body-content"},
		{"standard", "standard", "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTemplates()[tt.key]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTemplates() = %v, want %v", got, tt.want)
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
		Layout: "standard",
		tmpl:   tmpFile.Name(),
	}
	if err = a.templateSave(); err != nil {
		t.Errorf("templateSave() created an error: %s", err)
	}
}

func Test_pagedata(t *testing.T) {
	viper.SetDefault("html.title", "RetroTxt | example")

	args := Args{Layout: "standard"}
	w := "hello"
	d := []byte(w)
	got, _ := args.pagedata(&d)
	if got.PreText != w {
		t.Errorf("pagedata().PreText = %v, want %v", got, w)
	}
	args.Layout = "mini"
	w = "RetroTxt | example"
	got, _ = args.pagedata(&d)
	if got.PageTitle != w {
		t.Errorf("pagedata().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got, _ = args.pagedata(&d)
	if got.MetaDesc != w {
		t.Errorf("pagedata().MetaDesc = %v, want %v", got, w)
	}
	args.Layout = "standard"
	w = ""
	got, _ = args.pagedata(&d)
	if got.MetaAuthor != w {
		t.Errorf("pagedata().MetaAuthor = %v, want %v", got, w)
	}
}

func Test_destination(t *testing.T) {
	saved := viper.GetString("save-directory")
	wd, _ := os.Getwd()
	home, _ := os.UserHomeDir()
	spaces := filepath.Join(home, "some directory", "some file.html")
	root, _ := filepath.Abs("/")
	sub := filepath.Clean(filepath.Join(home, "/html/example.htm"))
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
			gotPath, err := destination(tt.args)
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
