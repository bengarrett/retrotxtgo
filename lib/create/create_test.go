package create

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
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
	l := strings.Split(Layouts(), ",")
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
	args := Args{HTMLLayout: "standard"}
	w := filepath.Clean("../../static/html/standard.html")
	got, _ := args.filename(true)
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	w = filepath.Clean("static/html/standard.html")
	got, _ = args.filename(false)
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	args.HTMLLayout = "error"
	_, err := args.filename(false)
	if (err != nil) != true {
		t.Errorf("filename = %v, want %v", got, w)
	}
}

func Test_pagedata(t *testing.T) {
	args := Args{HTMLLayout: "standard"}
	w := "hello"
	d := []byte(w)
	got := args.pagedata(d).PreText
	if got != w {
		t.Errorf("pagedata().PreText = %v, want %v", got, w)
	}
	args.HTMLLayout = "mini"
	w = "RetroTxt | example"
	got = args.pagedata(d).PageTitle
	if got != w {
		t.Errorf("pagedata().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got = args.pagedata(d).MetaDesc
	if got != w {
		t.Errorf("pagedata().MetaDesc = %v, want %v", got, w)
	}
	args.HTMLLayout = "standard"
	w = ""
	got = args.pagedata(d).MetaAuthor
	if got != w {
		t.Errorf("pagedata().MetaAuthor = %v, want %v", got, w)
	}
}

func Test_serveFile(t *testing.T) {
	a := Args{HTMLLayout: "standard"}
	type args struct {
		data []byte
		port uint
		test bool
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
			if err := a.serveFile(tt.args.data, tt.args.port, tt.args.test); (err != nil) != tt.wantErr {
				t.Errorf("serveFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func Test_writeFile(t *testing.T) {
	a := Args{HTMLLayout: "standard"}
	type args struct {
		data []byte
		name string
		test bool
	}
	tmpFile := path.Join(os.TempDir(), "retrotxtgo_create_test.txt")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data", args{[]byte(""), "", true}, true},
		{"invalid", args{[]byte("abc"), "this-is-an-invalid-path", true}, true},
		{"tempDir", args{[]byte("abc"), tmpFile, true}, true},
		{"homeDir", args{[]byte("abc"), "~", true}, false},
		{"currentDir", args{[]byte("abc"), ".", true}, false},
		{"path as name", args{[]byte("abc"), os.TempDir(), true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := a.File(tt.args.data, tt.args.name, tt.args.test); (err != nil) != tt.wantErr {
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