package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var hi = filepath.Clean("../../textfiles/hi.txt")

func TestDetail_infoSwitch(t *testing.T) {
	tests := []struct {
		name    string
		f       Detail
		format  string
		wantErr bool
	}{
		{"empty", Detail{}, "", true},
		{"invalid", Detail{}, "invalid", true},
		{"color", Detail{}, "color", false},
		{"j", Detail{}, "j", false},
		{"jm", Detail{}, "jm", false},
		{"text", Detail{}, "text", false},
		{"x", Detail{}, "x", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := tt.f.infoSwitch(tt.format); (gotErr.Msg != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_details(t *testing.T) {
	got, err := details(hi)
	if err != nil {
		t.Errorf("details() = %v, want %v", err, nil)
	}
	if got.Bytes != 40 {
		t.Errorf("details() = %v, want %v", got.Bytes, 40)
	}
	if got.Name != "hi.txt" {
		t.Errorf("details() = %v, want %v", got.Name, "hi.txt")
	}
	if got.Slug != "hi-txt" {
		t.Errorf("details() = %v, want %v", got.Slug, "hi-txt")
	}
	if got.Mime != "text/plain" {
		t.Errorf("details() = %v, want %v", got.Mime, "text/plain")
	}
	if got.Utf8 != true {
		t.Errorf("details() = %v, want %v", got.Utf8, true)
	}
	const want = "1b466b6448d7ff10e2f8f7160d936987"
	if got.MD5 != want {
		t.Errorf("details() = %v, want %v", got.MD5, want)
	}
}

func Test_parse(t *testing.T) {
	f, err := os.Stat(hi)
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
			got, err := parse(tt.args.data, tt.args.stat)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.CharCount, tt.want) {
				t.Errorf("parse() = %v, want %v", got.CharCount, tt.want)
			}
		})
	}
}

func TestDetail_infoJSON(t *testing.T) {
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
			if got := json.Valid(tt.d.infoJSON(tt.indent)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Detail.infoJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_infoText(t *testing.T) {
	d, err := details(hi)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	want := 491
	if got := len(d.infoText(false)); got != want {
		t.Errorf("infoText() = %v, want %v", got, want)
	}
}

func Test_infoXML(t *testing.T) {
	d, err := details(hi)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	want := 430
	if got := len(d.infoXML()); got != want {
		t.Errorf("infoXML() = %v, want %v", got, want)
	}
}
