package info

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
)

func ExampleDetail_XML() {
	tmp := sampleFile()
	millennia(tmp)
	var file Detail
	file.Read(tmp)
	data, _ := file.XML()
	filesystem.Clean(tmp)
	s := strings.ReplaceAll(string(data), "\t", "")
	ln := strings.Split(s, "\n")
	fmt.Println(ln[0])
	// Output: <file id="info-test-txt">
}

var sampleFile = func() string {
	path, err := filesystem.SaveTemp("info_test.txt", []byte(filesystem.Tabs))
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func TestIsText(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        bool
	}{
		{"empty", "", false},
		{"image", "image/jpeg", false},
		{"stream", "application/octet-stream", true},
		{"text", "text/plain", true},
		{"js", "text/javascript", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsText(tt.contentType); got != tt.want {
				t.Errorf("IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Print(t *testing.T) {
	tmp := sampleFile()
	type args struct {
		filename string
		format   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", ""}, true},
		{"file not exist", args{"notexistingfile", "json"}, true},
		{"invalid fmt", args{tmp, "yaml"}, true},
		{"color", args{tmp, ""}, false},
		{"json", args{tmp, "json"}, false},
		{"json.min", args{tmp, "jm"}, false},
		{"text", args{tmp, "t"}, false},
		{"xml", args{tmp, "xml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Print(tt.args.filename, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	filesystem.Clean(tmp)
}

func Test_File(t *testing.T) {
	tmp := sampleFile()
	var got Detail
	err := got.Read(tmp)
	if err != nil {
		t.Errorf("Read() = %v, want %v", err, nil)
	}
	if got.Bytes != 57 {
		t.Errorf("Read() = %v, want %v", got.Bytes, 57)
	}
	if got.Name != "info_test.txt" {
		t.Errorf("Read() = %v, want %v", got.Name, "info_test.txt")
	}
	if got.Slug != "info-test-txt" {
		t.Errorf("Read() = %v, want %v", got.Slug, "info-test-txt")
	}
	if got.Mime != "text/plain" {
		t.Errorf("Read() = %v, want %v", got.Mime, "text/plain")
	}
	if got.Utf8 != true {
		t.Errorf("Read() = %v, want %v", got.Utf8, true)
	}
	const want = "883643f5e9ed278732c92d9b6f834b96"
	if got.MD5 != want {
		t.Errorf("Read() = %v, want %v", got.MD5, want)
	}
	filesystem.Clean(tmp)
}

func Test_parse(t *testing.T) {
	tmp := sampleFile()
	f, err := os.Stat(tmp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	type args struct {
		data []byte
		stat os.FileInfo
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"empty", args{[]byte(""), f}, 0, false},
		{"string", args{[]byte("hello"), f}, 5, false},
		{"string", args{[]byte("世界你好"), f}, 8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Detail
			err := got.parse(tt.args.data, tt.args.stat, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.CharCount, tt.want) {
				t.Errorf("parse() = %v, want %v", got.CharCount, tt.want)
			}
		})
	}
	filesystem.Clean(tmp)
}

func Test_JSON(t *testing.T) {
	tests := []struct {
		name   string
		d      Detail
		indent bool
		want   bool
	}{
		{"no indent", Detail{}, false, true},
		{"indent", Detail{}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json.Valid(tt.d.JSON(tt.indent)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Text(t *testing.T) {
	tmp := sampleFile()
	var d Detail
	err := d.Read(tmp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	want := 604
	if got := len(d.Text(false)); got != want {
		t.Errorf("Text() = %v, want %v", got, want)
	}
	filesystem.Clean(tmp)
}

func millennia(name string) {
	mtime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(name, mtime, mtime); err != nil {
		log.Fatal(err)
	}
}
