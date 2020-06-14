package samples

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func ExampleCP437Decode() {
	const name = base + "cp437In.txt"
	result, err := CP437Decode(cp437hex)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(t)
	// Output: ═╣▓╠═
}

func ExampleCP437Encode() {
	const name = base + "cp437.txt"
	result, err := CP437Encode(utf)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(name)
	fmt.Print(len(t))
	// Output: 8
}

func ExampleHexDecode() {
	b, err := HexDecode("cdb9b0cccd")
	if err != nil {
		log.Fatal(err)
	}
	d, err := CP437Decode(string(b))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(d))
	// Output: ═╣░╠═
}

func ExampleHexEncode() {
	s, err := CP437Encode("═╣░╠═")
	if err != nil {
		log.Fatal(err)
	}
	h := HexEncode(string(s))
	fmt.Printf("%s %v", h, h)
	// Output: cdb9b0cccd [99 100 98 57 98 48 99 99 99 100]
}

func ExampleSave() {
	const name = base + "save.txt"
	path, err := Save([]byte("hello world"), name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(path)
	// Output:
}

func ExampleBase64Encode() {
	const name = base + "newlines.txt"
	b := Base64Encode(Newlines)
	r, err := Base64Decode(b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("source:\t%v\nresult:\t%q", b, r)
	Save(r, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(name)
	fmt.Printf("\nread:\t%q", t) //string(d))
	// Output: source:	YQpiCmMuLi4K
	// result:	"a\nb\nc...\n"
	// read:	"a\nb\nc...\n"
}

func ExampleUTF8() {
	const name = base + "utf8.txt"
	result, _, err := UTF8(Symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(name)
	fmt.Printf("%dB %s", len(t), t)
	// Output: 14B [☠|☮|♺]
}

func ExampleUTF16BE() {
	const name = base + "utf16be.txt"
	result, _, err := UTF16BE(Symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(name)
	fmt.Print(len(t))
	// Output: 17
}

func ExampleUTF16LE() {
	const name = base + "utf16le.txt"
	result, _, err := UTF16LE(Symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	Clean(name)
	fmt.Print(len(t))
	// Output: 17
}

func BenchmarkCP437Decode(b *testing.B) {
	const name = base + "cp437In.txt"
	result, err := CP437Decode(cp437hex)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(t)
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		wantEncoded string
	}{
		// can create test cases
		// https://www.base64encode.org/
		{"empty", "", ""},
		{"abc", "abc", "YWJj"},
		{"octets", `\0\x01\x02\x03\x04\x05\x06\x07\b\t\n\x0B\f\r\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1B\x1C\x1D\x1E\x1F !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_abcdefghijklmnopqrstuvwxyz{|}~\x7F`,
			"XDBceDAxXHgwMlx4MDNceDA0XHgwNVx4MDZceDA3XGJcdFxuXHgwQlxmXHJceDBFXHgwRlx4MTBceDExXHgxMlx4MTNceDE0XHgxNVx4MTZceDE3XHgxOFx4MTlceDFBXHgxQlx4MUNceDFEXHgxRVx4MUYgISIjJCUmXCcoKSorLC0uLzAxMjM0NTY3ODk6Ozw9Pj9AQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVpbXFxdXl9hYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ent8fX5ceDdG"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncoded := Base64Encode(tt.text); gotEncoded != tt.wantEncoded {
				t.Errorf("Base64Encode() = %v, want %v", gotEncoded, tt.wantEncoded)
			}
		})
	}
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
		{"nl", Newlines, Newlines, false},
		{"utf", Symbols, "[Γÿá|Γÿ«|ΓÖ║]", false},
		{"escapes", Escapes, `bell:,back:,tab:	,form:,vertical:,quote:"`, false},
		{"digits", Digits, "░░┼░┼░", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := CP437Decode(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("CP437Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotResult), tt.wantResult) {
				t.Errorf("CP437Decode() = %v, want %v", string(gotResult), tt.wantResult)
			}
		})
	}
}

func TestBase64Decode(t *testing.T) {
	const name = base + "sample.base64"
	raw, err := ReadLine(ansiFile(), "dos")
	if err != nil {
		log.Fatal(err)
	}
	ansi, err := Base64Decode(EncodeANSI())
	if err != nil {
		log.Fatal(err)
	}
	Save([]byte(EncodeANSI()), name)
	Clean(name)
	match := (string(ansi) == raw)
	if !match {
		t.Errorf("Base64Decode() = %v, want %v", match, true)
	}
}
