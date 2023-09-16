package info_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/info"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"github.com/bengarrett/retrotxtgo/static"
)

func ExampleInfo() {
	s := strings.Builder{}
	_ = info.Info(&s, "testdata/example.txt", "text")
	x := strings.Split(s.String(), "\n")
	for _, v := range x {
		if strings.Contains(v, "SHA256 checksum") {
			fmt.Print(v)
		}
	}
	// Output: SHA256 checksum  4b187b0e6bc12541659eed5845d9dbe0914d4fc026f849bd03c255775a97d878
}

func ExampleMarshal() {
	s := strings.Builder{}
	_ = info.Marshal(&s, "testdata/example.txt", info.JSON)
	fmt.Printf("%d bytes and json? %t", len(s.String()), json.Valid([]byte(s.String())))
	// Output: 2363 bytes and json? true
}

func ExampleStream() {
	s := strings.Builder{}
	file, _ := os.Open("testdata/example.txt")
	b := make([]byte, 25)
	file.Read(b)
	info.Stream(&s, "text", b...)
	stdin := strings.Contains(s.String(), "n/a (stdin)")
	fmt.Print("stdin? ", stdin)
	// Output: stdin? true
}

func rawData() []byte {
	b, err := static.Text.ReadFile("text/sauce.txt")
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func sampleFileB() string {
	b := []byte(mock.T()["Tabs"]) // Tabs and Unicode glyphs
	path, err := fsys.SaveTemp("info_test.txt", b...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func TestMarshal(t *testing.T) {
	tmp := sampleFileB()
	type args struct {
		filename string
		format   info.Format
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", info.PlainText}, true},
		{"file not exist", args{"notexistingfile", info.JSON}, true},
		{"color", args{tmp, info.ColorText}, false},
		{"json", args{tmp, info.JSON}, false},
		{"json.min", args{tmp, info.JSONMin}, false},
		{"text", args{tmp, info.PlainText}, false},
		{"xml", args{tmp, info.XML}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := info.Marshal(nil, tt.args.filename, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fsys.Clean(tmp)
}

func TestStream(t *testing.T) {
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
			err := info.Stream(nil, tt.args.format, tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stream() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
