package create

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
	"golang.org/x/text/encoding/charmap"
)

func ExampleHonk() {
	for i, c := range charmap.All {
		fmt.Println(i, c)
	}
	fmt.Printf("\n%v\n", charmap.All)
	fmt.Printf("\n%+v \n", charmap.CodePage437)
	fmt.Println(charmap.CodePage037.ID())
	// Output: ?
}

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
	l := Options()
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
	viper.SetDefault("create.title", "RetroTxt | example")

	args := Args{HTMLLayout: "standard"}
	w := "hello"
	d := []byte(w)
	got, _ := args.pagedata(d)
	if got.PreText != w {
		t.Errorf("pagedata().PreText = %v, want %v", got, w)
	}
	args.HTMLLayout = "mini"
	w = "RetroTxt | example"
	got, _ = args.pagedata(d)
	if got.PageTitle != w {
		t.Errorf("pagedata().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got, _ = args.pagedata(d)
	if got.MetaDesc != w {
		t.Errorf("pagedata().MetaDesc = %v, want %v", got, w)
	}
	args.HTMLLayout = "standard"
	w = ""
	got, _ = args.pagedata(d)
	if got.MetaAuthor != w {
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
	a := Args{HTMLLayout: "standard", Test: true}
	type args struct {
		data []byte
		name string
	}
	tmpFile := path.Join(os.TempDir(), "retrotxtgo_create_test.txt")
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
			if err := a.File(tt.args.data, tt.args.name); (err != nil) != tt.wantErr {
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
