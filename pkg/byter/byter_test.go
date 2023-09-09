package byter_test

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/byter"
	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"golang.org/x/text/encoding/charmap"
)

func ExampleBOM() {
	fmt.Fprintf(os.Stdout, "%X", byter.BOM())
	// Output: EFBBBF
}

func TestTrimEOF(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []byte
	}{
		{"empty", nil, nil},
		{"none", []byte("hello world"), []byte("hello world")},
		{"one", []byte("hello\x1aworld"), []byte("hello")},
		{"two", []byte("hello\x1aworld\x1athis should be hidden"), []byte("hello")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byter.TrimEOF(tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimEOF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeBytes(t *testing.T) {
	if l := len(byter.MakeBytes()); l != 256 {
		t.Errorf("MakeBytes() = %v, want %v", l, 256)
	}
}

func TestMark(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"empty string", args{}, []byte{239, 187, 191}},
		{"ascii string", args{[]byte("hi")}, []byte{239, 187, 191, 104, 105}},
		{"existing bom string", args{byter.BOM()}, []byte{239, 187, 191}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byter.Mark(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mark() = %v, want %v", got, tt.want)
			}
		})
	}
}

const c, h = "═╣░╠═", "cdb9b0cccd"

func TestHexDecode(t *testing.T) {
	samp, err := byter.Encode(charmap.CodePage437, c)
	if err != nil {
		t.Errorf("HexDecode() E437() error = %v", err)
	}
	type args struct {
		hexadecimal string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"═╣░╠═", args{h}, samp, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := byter.HexDecode(tt.args.hexadecimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexDecode() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestHexEncode(t *testing.T) {
	samp, err := byter.Encode(charmap.CodePage437, c)
	if err != nil {
		t.Errorf("HexDecode() E437() error = %v", err)
	}
	type args struct {
		text string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
	}{
		{"═╣░╠═", args{string(samp)}, []byte{99, 100, 98, 57, 98, 48, 99, 99, 99, 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := byter.HexEncode(tt.args.text); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexEncode() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

const (
	cp437hex = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`
	utf      = "═╣ ░ ╠═"
	base     = "rt_sample-"
)

func ExampleD437() {
	const name = base + "cp437In.txt"
	result, err := byter.Decode(charmap.CodePage437, cp437hex)
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
			gotResult, err := byter.Decode(charmap.CodePage437, tt.s)
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
	result, err := byter.Encode(charmap.CodePage437, utf)
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
			gotResult, err := byter.Decode(&c, tt.args.s)
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
			gotResult, err := byter.Encode(&c, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Encode() = %v, want %v", gotResult, tt.wantResult)
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
			gotResult, err := byter.Decode(charmap.CodePage437, tt.s)
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
			gotResult, err := byter.Encode(charmap.CodePage437, tt.s)
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
