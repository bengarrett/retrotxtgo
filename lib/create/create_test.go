package create

import (
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
		data    []byte
		value   string
		changed bool
	}
	tests := []struct {
		name string
		args args
	}{
		// {"empty", args{nil, "", false}},
		// {"empty", args{nil, "", true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Save(tt.args.data, tt.args.value, tt.args.changed)
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

func Test_filename(t *testing.T) {
	args := Args{Layout: "standard"}
	w := filepath.Clean("../../static/html/standard.html")
	got, _ := args.filename()
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	w = filepath.Clean("static/html/standard.html")
	got, _ = args.filename()
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	args.Layout = "error"
	_, err := args.filename()
	if (err != nil) != true {
		t.Errorf("filename = %v, want %v", got, w)
	}
}

func Test_pagedata(t *testing.T) {
	viper.SetDefault("create.title", "RetroTxt | example")

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

func Test_serveFile(t *testing.T) {
	a := Args{Layout: "standard"}
	type args struct {
		data *[]byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// not tested as it requires operating system permissions
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := a.serveFile(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("serveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Test_writeFile(t *testing.T) {
	a := Args{Layout: "standard", Test: true}
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
		{"tempDir", args{[]byte("abc"), tmpFile}, true},
		{"homeDir", args{[]byte("abc"), "~"}, false},
		{"currentDir", args{[]byte("abc"), "."}, false},
		{"path as name", args{[]byte("abc"), os.TempDir()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := a.Save(&tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// clean-up
	if wd, err := os.Getwd(); err == nil {
		p := filepath.Join(wd, "index.html")
		if err := os.Remove(p); err != nil {
			t.Error(err)
		}
	}
}

func Test_writeStdout(t *testing.T) {
	type args struct {
		data []byte
		test bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data", args{[]byte(""), true}, false},
		{"some data", args{[]byte("hello world"), true}, false},
		{"nil data", args{nil, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if err := Stdout(tt.args.data, tt.args.test); (err != nil) != tt.wantErr {
			// 	t.Errorf("writeStdout() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}

func TestDest(t *testing.T) {
	saved := viper.GetString("create.save-directory")
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
			gotPath, err := Dest(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Dest() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
}
