package read_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/read"
)

func Test_promptRead(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name      string
		input     string
		wantInput string
		wantErr   bool
	}{
		{"empty", "", "", false},
		{"nl", "\n", "", false},
		{"tab", "\t", "", false},
		{"hi", "hello", "hello", false},
		{"hw", "hello world", "hello world", false},
		{"emoji", "hi 😃", "hi 😃", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			gotInput, err := read.Read(&stdin)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInput != tt.wantInput {
				t.Errorf("Read() = %v, want %v", gotInput, tt.wantInput)
			}
		})
	}
}

func Test_parseYN(t *testing.T) {
	type args struct {
		input      string
		yesDefault bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"yes 1", args{"", true}, true},
		{"yes 1", args{"y", false}, true},
		{"yes 2", args{"yes", false}, true},
		{"no 1", args{"no", false}, false},
		{"no 2", args{"no", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := read.ParseYN(tt.args.input, tt.args.yesDefault); got != tt.want {
				t.Errorf("ParseYN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pstring(t *testing.T) {
	type args struct {
		r io.Reader
	}
	a := strings.NewReader("my request")
	b := strings.NewReader("-")
	c := strings.NewReader("\x0D")
	tests := []struct {
		name      string
		args      args
		wantWords string
	}{
		{"input", args{a}, "my request"},
		{"remove", args{b}, "-"},
		{"enter", args{c}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWords, _ := read.Parse(tt.args.r); gotWords != tt.wantWords {
				t.Errorf("Parse() = %v, want %v", gotWords, tt.wantWords)
			}
		})
	}
}
