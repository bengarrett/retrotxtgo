package info_test

import (
	"errors"
	"log"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/info"
	"github.com/bengarrett/retrotxtgo/pkg/info/internal/detail"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"github.com/bengarrett/retrotxtgo/static"
)

func rawData() []byte {
	b, err := static.Text.ReadFile("text/sauce.txt")
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func sampleFile() string {
	b := []byte(mock.T()["Tabs"]) // Tabs and Unicode glyphs
	path, err := fsys.SaveTemp("info_test.txt", b...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func TestMarshal(t *testing.T) {
	tmp := sampleFile()
	type args struct {
		filename string
		format   detail.Format
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", detail.PlainText}, true},
		{"file not exist", args{"notexistingfile", detail.JSON}, true},
		{"color", args{tmp, detail.ColorText}, false},
		{"json", args{tmp, detail.JSON}, false},
		{"json.min", args{tmp, detail.JSONMin}, false},
		{"text", args{tmp, detail.PlainText}, false},
		{"xml", args{tmp, detail.XML}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := info.Marshal(tt.args.filename, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	fsys.Clean(tmp)
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
			_, err := info.Stdin(tt.args.format, tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Stdin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNames_Info(t *testing.T) {
	const fileToTest = "internal/detail/detail.go"
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
		{"empty", fields{}, args{}, nil},
		{"bad dir", fields{}, args{name: "some invalid filename"}, nil},
		{"temp file", fields{}, args{name: fileToTest, format: "json.min"}, nil},
		{"temp dir", fields{}, args{name: ".", format: "json.min"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := info.Names{
				Index:  tt.fields.Index,
				Length: tt.fields.Length,
			}
			if _, got := n.Info(tt.args.name, tt.args.format); !errors.Is(got, tt.wantErr) {
				t.Errorf("Names.Info() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
