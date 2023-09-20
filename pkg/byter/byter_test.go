package byter_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/byter"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"golang.org/x/text/encoding/charmap"
)

func ExampleDecode() {
	s := "\xCD\xB9\xB2\xCC\xCD"
	p, _ := byter.Decode(charmap.CodePage437, s)
	fmt.Printf("%s\n", p)
	p, _ = byter.Decode(charmap.ISO8859_1, s)
	fmt.Printf("%s\n", p)
	// Output: ═╣▓╠═
	// Í¹²ÌÍ
}

func ExampleEncode() {
	p, _ := byter.Encode(charmap.CodePage437, "═╣▓╠═")
	fmt.Printf("%X\n", p)
	p, _ = byter.Encode(charmap.ISO8859_1, "hello")
	fmt.Printf("%s (%X)\n", p, p)
	p, _ = byter.Encode(charmap.CodePage1140, "hello")
	fmt.Printf("%X\n", p) // not ascii compatible
	// Output: CDB9B2CCCD
	// hello (68656C6C6F)
	// 8885939396
}

func ExampleBOM() {
	fmt.Printf("#%x", byter.BOM())
	// Output: #efbbbf
}

func ExampleTrimEOF() {
	fmt.Printf("%q", byter.TrimEOF([]byte("hello world\x1a")))
	// Output: "hello world"
}

func ExampleMakeBytes() {
	fmt.Printf("%d", len(byter.MakeBytes()))
	// Output: 256
}

func ExampleHexDecode() {
	b, _ := byter.HexDecode("6F6B")
	fmt.Printf("%s\n", b)
	fmt.Printf("%d\n", b)
	// Output: ok
	// [111 107]
}

func ExampleHexEncode() {
	b := byter.HexEncode("ok")
	fmt.Printf("#%s\n", b)
	// Output: #6f6b
}

func TestTrimEOF(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := byter.TrimEOF(tt.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TrimEOF() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestMakeBytes(t *testing.T) {
	t.Parallel()
	if l := len(byter.MakeBytes()); l != 256 {
		t.Errorf("MakeBytes() = %v, want %v", l, 256)
	}
}

func TestMark(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := byter.Mark(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mark() = %v, want %v", got, tt.want)
			}
		}
	})
}

const c, h = "═╣░╠═", "cdb9b0cccd"

func TestHexDecode(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotResult, err := byter.HexDecode(tt.args.hexadecimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexDecode() = %v, want %v", gotResult, tt.wantResult)
			}
		}
	})
}

func TestHexEncode(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if gotResult := byter.HexEncode(tt.args.text); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexEncode() = %v, want %v", gotResult, tt.wantResult)
			}
		}
	})
}

const (
	cp437hex = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`
	utf      = "═╣ ░ ╠═"
	base     = "rt_sample-"
)

func TestCP437Decode(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotResult, err := byter.Decode(charmap.CodePage437, tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotResult), tt.wantResult) {
				t.Errorf("D437() = %v, want %v", string(gotResult), tt.wantResult)
			}
		}
	})
}

func TestDString(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c := tt.args.cp
			gotResult, err := byter.Decode(&c, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DString() = %v, want %v", gotResult, tt.wantResult)
			}
		}
	})
}

func TestEString(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c := tt.args.cp
			gotResult, err := byter.Encode(&c, tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Encode() = %v, want %v", gotResult, tt.wantResult)
			}
		}
	})
}

func TestD437(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotResult, err := byter.Decode(charmap.CodePage437, tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("D437() = %s, want %v", gotResult, tt.wantResult)
			}
		}
	})
}

func TestE437(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotResult, err := byter.Encode(charmap.CodePage437, tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("E437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("E437() = %v, want %v", gotResult, tt.wantResult)
			}
		}
	})
}
