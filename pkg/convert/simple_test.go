package convert_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"golang.org/x/text/encoding/charmap"
)

const (
	cp437hex = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`
	utf      = "═╣ ░ ╠═"
	base     = "rt_sample-"
)

func ExampleD437() {
	const name = base + "cp437In.txt"
	result, err := convert.D437(cp437hex)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fsys.SaveTemp(name, result...)
	if err != nil {
		log.Fatal(err)
	}
	t, err := fsys.ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(os.Stdout, t)
	// Output: ═╣▓╠═
}

func TestCP437Decode(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantResult string
		wantErr    bool
	}{
		{"empty", "", "", false},
		{"hex", cp437hex, "═╣▓╠═", false},
		{"nl", mock.T()["Newline"], mock.T()["Newline"], false},
		{"utf", mock.T()["Symbols"], "[Γÿá|Γÿ«|ΓÖ║]", false},
		{"escapes", mock.T()["Escapes"], "bell:\a,back:\b,tab:\t,form:\f,vertical:\v,quote:\"", false},
		{"digits", mock.T()["Digits"], "░░┼░┼░", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := convert.D437(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotResult), tt.wantResult) {
				t.Errorf("D437() = %v, want %v", string(gotResult), tt.wantResult)
			}
		})
	}
}

func ExampleE437() {
	const name = base + "cp437.txt"
	result, err := convert.E437(utf)
	if err != nil {
		log.Fatal(err)
	}
	_, err = fsys.SaveTemp(name, result...)
	if err != nil {
		log.Fatal(err)
	}
	t, err := fsys.ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fsys.Clean(name)
	fmt.Fprint(os.Stdout, len(t))
	// Output: 8
}

func TestDString(t *testing.T) {
	type args struct {
		s  string
		cp charmap.Charmap
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"empty", args{}, []byte{}, false},
		{"hello", args{"hello world", *charmap.CodePage437}, []byte("hello world"), false},
		{"mainframe", args{
			string([]byte{136, 133, 147, 147, 150, 64, 166, 150, 153, 147, 132}),
			*charmap.CodePage037,
		}, []byte("hello world"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.cp
			gotResult, err := convert.DString(tt.args.s, &c)
			if (err != nil) != tt.wantErr {
				t.Errorf("DString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestEString(t *testing.T) {
	type args struct {
		s  string
		cp charmap.Charmap
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"empty", args{}, []byte{}, false},
		{
			"hello",
			args{"hello world", *charmap.CodePage437},
			[]byte("hello world"), false,
		},
		{
			"mainframe",
			args{"hello world", *charmap.CodePage037},
			[]byte{136, 133, 147, 147, 150, 64, 166, 150, 153, 147, 132},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.cp
			gotResult, err := convert.EString(tt.args.s, &c)
			if (err != nil) != tt.wantErr {
				t.Errorf("EString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("EString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestD437(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantResult []byte
		wantErr    bool
	}{
		{"empty", "", []byte{}, false},
		{"hello", "hello world", []byte("hello world"), false},
		{"hex", "\xe0 alpha \xe1 beta", []byte("α alpha ß beta"), false},
		{"octal", "\253 half \254 quarter", []byte("½ half ¼ quarter"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := convert.D437(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("D437() = %s, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestE437(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantResult []byte
		wantErr    bool
	}{
		{"empty", "", []byte{}, false},
		{"hello", "hello world", []byte("hello world"), false},
		{"hex", "α alpha ß beta", []byte("\xe0 alpha \xe1 beta"), false},
		{"octal", "½ half ¼ quarter", []byte("\253 half \254 quarter"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := convert.E437(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("E437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("E437() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
