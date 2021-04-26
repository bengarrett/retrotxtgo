package info

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"retrotxt.com/retrotxt/lib/filesystem"
	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/static"
)

func rawData() []byte {
	b, err := static.Text.ReadFile("text/sauce.txt")
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func millennia(name string) {
	mtime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(name, mtime, mtime); err != nil {
		log.Fatal(err)
	}
}

func sampleFile() string {
	b := []byte(filesystem.T()["Tabs"]) // Tabs and Unicode glyphs
	path, err := filesystem.SaveTemp("info_test.txt", b...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func ExampleText() {
	var file Detail
	tmp := sampleFile()
	millennia(tmp)
	if err := file.read(tmp); err != nil {
		log.Fatal(err)
	}
	data, _ := file.marshal(XML)
	filesystem.Clean(tmp)
	s := strings.ReplaceAll(string(data), "\t", "")
	ln := strings.Split(s, "\n")
	fmt.Println(ln[0])
	// Output: <file utf8="true" id="info-test-txt">
}

func TestValidText(t *testing.T) {
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
	var d Detail
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d.Mime.Type = tt.contentType
			if got := d.validText(); got != tt.want {
				t.Errorf("IsText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Marshal(t *testing.T) {
	tmp := sampleFile()
	type args struct {
		filename string
		format   Format
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", PlainText}, true},
		{"file not exist", args{"notexistingfile", JSON}, true},
		{"color", args{tmp, ColorText}, false},
		{"json", args{tmp, JSON}, false},
		{"json.min", args{tmp, JSONMin}, false},
		{"text", args{tmp, PlainText}, false},
		{"xml", args{tmp, XML}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Marshal(tt.args.filename, tt.args.format, 0, 0); (err != nil) != tt.wantErr {
				t.Errorf("Print() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	filesystem.Clean(tmp)
}

func Test_read(t *testing.T) {
	tmp := sampleFile()
	fmt.Println("path:", tmp)
	var got Detail
	err := got.read(tmp)
	if err != nil {
		t.Errorf("read() = %v, want %v", err, nil)
	}
	if got.Size.Bytes != 57 {
		t.Errorf("read() = %v, want %v", got.Size.Bytes, 57)
	}
	if got.Name != "info_test.txt" {
		t.Errorf("read() = %v, want %v", got.Name, "info_test.txt")
	}
	if got.Slug != "info-test-txt" {
		t.Errorf("read() = %v, want %v", got.Slug, "info-test-txt")
	}
	if got.Mime.Type != "text/plain" {
		t.Errorf("read() = %v, want %v", got.Mime, "text/plain")
	}
	if got.Utf8 != true {
		t.Errorf("read() = %v, want %v", got.Utf8, true)
	}
	const want = "883643f5e9ed278732c92d9b6f834b96"
	if got.Sums.MD5 != want {
		t.Errorf("read() = %v, want %v", got.Sums.MD5, want)
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
		{"string", args{[]byte("世界你好"), f}, 4, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Detail
			err := got.parse("", tt.args.stat, tt.args.data...)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Count.Chars, tt.want) {
				t.Errorf("parse() = %v, want %v", got.Count.Chars, tt.want)
			}
		})
	}
	filesystem.Clean(tmp)
}

func Test_marshal_json(t *testing.T) {
	tests := []struct {
		name   string
		d      Detail
		format Format
		want   bool
	}{
		{"no indent", Detail{}, JSONMin, true},
		{"indent", Detail{}, JSON, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j, _ := tt.d.marshal(tt.format)
			if got := json.Valid(j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_marshal_text(t *testing.T) {
	const want = 727
	var d Detail
	tmp := sampleFile()
	err := d.read(tmp)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	b, _ := d.marshal(PlainText)
	if got := len(b); got != want {
		t.Errorf("marshal() = %v, want %v", got, want)
	}
	filesystem.Clean(tmp)
}

func TestStdin(t *testing.T) {
	type args struct {
		format string
		b      []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, false},
		{"empty xml", args{format: "xml"}, false},
		{"color", args{format: "c", b: rawData()}, false},
		{"text", args{format: "text", b: rawData()}, false},
		{"json", args{format: "json", b: rawData()}, false},
		{"json.min", args{format: "jm", b: rawData()}, false},
		{"xml", args{format: "x", b: rawData()}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Stdin(tt.args.format, tt.args.b...); (err != nil) != tt.wantErr {
				t.Errorf("Stdin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNames_Info(t *testing.T) {
	type fields struct {
		Index  int
		Length int
	}
	type args struct {
		name   string
		format string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"empty", fields{}, args{}, logs.ErrNameNil},
		{"bad dir", fields{}, args{name: "some invalid filename"}, logs.ErrFileNil},
		{"temp dir", fields{}, args{name: os.TempDir(), format: "json.min"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Names{
				Index:  tt.fields.Index,
				Length: tt.fields.Length,
			}
			if got := n.Info(tt.args.name, tt.args.format); !errors.Is(got, tt.wantErr) {
				t.Errorf("Names.Info() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
