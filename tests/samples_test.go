package tests

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
	cleanFile(name)
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

func ExampleReadBytes() {
	const name = base + "cp437.ans"
	r, err := Base64Decode(LogoANSI)
	if err != nil {
		log.Fatal(err)
	}
	Save(r, name)
	t, err := ReadBytes(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(t))
	// Output: 26337
}

func ExampleSave() {
	const name = base + "save.txt"
	path, err := Save([]byte("hello world"), name)
	if err != nil {
		log.Fatal(err)
	}
	cleanFile(path)
	// Output:
}

func ExampleBase64Decode() {
	const name = base + "cp437.ans"
	r, err := Base64Decode(LogoANSI)
	if err != nil {
		log.Fatal(err)
	}
	Save(r, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(t))
	// Output: 26369
}

func ExampleBase64Encode() {
	const name = base + "newlines.txt"
	b := Base64Encode(newlines)
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
	cleanFile(name)
	fmt.Printf("\nread:\t%q", t) //string(d))
	// Output: source:	YQpiCmMuLi4K
	// result:	"a\nb\nc...\n"
	// read:	"a\nb\nc...\n"
}

func ExampleUTF8() {
	const name = base + "utf8.txt"
	result, _, err := UTF8(symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	cleanFile(name)
	fmt.Printf("%dB %s", len(t), t)
	// Output: 14B [☠|☮|♺]
}

func ExampleUTF16BE() {
	const name = base + "utf16be.txt"
	result, _, err := UTF16BE(symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	cleanFile(name)
	fmt.Print(len(t))
	// Output: 17
}

func ExampleUTF16LE() {
	const name = base + "utf16le.txt"
	result, _, err := UTF16LE(symbols)
	if err != nil {
		log.Fatal(err)
	}
	Save(result, name)
	t, err := ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	cleanFile(name)
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
	type args struct {
		text string
	}
	tests := []struct {
		name        string
		args        args
		wantEncoded string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotEncoded := Base64Encode(tt.args.text); gotEncoded != tt.wantEncoded {
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
		{"nl", newlines, newlines, false},
		{"utf", symbols, "[Γÿá|Γÿ«|ΓÖ║]", false},
		{"escapes", escapes, `bell:,back:,tab:	,form:,vertical:,quote:"`, false},
		{"digits", digits, "░░┼░┼░", false},
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
	cleanFile(name)
	match := (string(ansi) == raw)
	if !match {
		t.Errorf("Base64Decode() = %v, want %v", match, true)
	}
}
