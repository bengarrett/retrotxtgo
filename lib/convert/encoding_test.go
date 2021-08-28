// nolint:gochecknoglobals,dupl
package convert

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

func ExampleCharacters() {
	fmt.Print(string(Characters()[DEL]))
	fmt.Print(string(Characters()[SquareRoot]))
	// Output: Δ✓
}

var (
	ccp037 = charmap.CodePage037
	cp437  = charmap.CodePage437
	cp865  = charmap.CodePage865
	cp1252 = charmap.Windows1252
	koi    = charmap.KOI8R
	iso1   = charmap.ISO8859_1
	iso6e  = charmap.ISO8859_6E
	iso15  = charmap.ISO8859_15
	jis    = japanese.ShiftJIS
	mac    = charmap.Macintosh
	u8     = unicode.UTF8
	u8bom  = unicode.UTF8BOM
	u16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	u16le  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
)

func toEncoding(e encoding.Encoding, s string) []byte {
	b, err := e.NewEncoder().Bytes([]byte(s))
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func TestSet_Transform(t *testing.T) {
	tests := []struct {
		name     string
		codepage encoding.Encoding
		text     []byte
		want     string
		wantErr  bool
	}{
		{"u8", u8, toEncoding(u8, "⌚ Watch"), "⌚ Watch", false},
		{"u8bom", u8bom, toEncoding(u8, "⌚ Watch"), "⌚ Watch", false},
		{"u16le", u16le, toEncoding(u16le, "⌚ Watch"), "⌚ Watch", false},
		{"u16be", u16be, toEncoding(u16be, "⌚ Watch"), "⌚ Watch", false},
		{"null", u8, []byte("\x00"), "␀", false},
		{"CP037", ccp037, []byte("\xc8\x51\xba\x93\xcf"), "Hé[lõ", false},
		{"bell", ccp037, []byte("ring a \x07"), "ring a ␇", false},
		{"CP437", cp437, []byte("H\x82ll\x93 \x9d\xa7\xf4\x9c\xbe"), "Héllô ¥º⌠£╛", false},
		{"⌂", cp437, []byte("Home sweat \x7f"), "Home sweat ⌂", false},
		{"mac", mac, []byte("\x11 command + \x12 shift."), "⌘ command + ⇧ shift.", false},
		{"latin1", iso1, toEncoding(iso1, "currency sign ¤"), "currency sign ¤", false},
		{"latin15", iso15, toEncoding(iso15, "euro sign €"), "euro sign €", false},
		{"6e", iso6e, []byte("ring a \x07"), "ring a ␇", false},
		{"koi8", koi, []byte("\xf5\xf2\xf3\xf3"), "УРСС", false},
		{"jp", jis, toEncoding(jis, "abc"), "abc", false},
		{"865", cp865, toEncoding(cp865, "currency sign ¤"), "currency sign ¤", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Convert{}
			data.Input.Bytes = tt.text
			data.Input.Encoding = tt.codepage
			err := data.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data.Swap()
			if string(data.Output) != tt.want {
				t.Errorf("Convert.Transform() = %q, want %q", data.Output, tt.want)
			}
		})
	}
}

func TestANSI(t *testing.T) {
	tests := []struct {
		name     string
		codepage encoding.Encoding
		str      string
		want     []rune
		wantErr  bool
	}{
		{"null", u8, "\x00", []rune("␀"), false},
		{"CP037", ccp037, "\xc8\x51\xba\x93\xcf", []rune("Hé[lõ"), false},
		{"ansi dos", cp437, "\x1b\x5b0m", []rune{27, 91, 48, 109}, false},
		{"ansi win", cp1252, "\x1b\x5b0m", []rune{27, 91, 48, 109}, false},
		{"panic", cp1252, "\x1b", []rune{9243}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := Convert{}
			data.Input.Bytes = []byte(tt.str)
			data.Input.Encoding = tt.codepage
			err := data.Transform()
			data.Swap().ANSIControls()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(data.Output, tt.want) {
				t.Errorf("Convert.Transform() = %v %q, want %v", data.Output, data.Output, tt.want)
			}
		})
	}
}

func TestEncoding(t *testing.T) {
	tests := []struct {
		name    string
		want    encoding.Encoding
		wantErr bool
	}{
		{"ibm437", charmap.CodePage437, false},
		{"CP437", charmap.CodePage437, false},
		{"IBM Code Page 437", charmap.CodePage437, false},
		{"CP-437", charmap.CodePage437, false},
		// {"oem-us", charmap.CodePage437, false}, // US-ASCII
		{"ibm-37", charmap.CodePage037, false},
		{"858", charmap.CodePage858, false},
		{"mac", charmap.Macintosh, false},
		{"cp1004", charmap.Windows1252, false},
		{"latin1", charmap.ISO8859_1, false},
		{"ISO-8859-1", charmap.ISO8859_1, false},
		{"latin1", charmap.ISO8859_1, false},
		// {"ansi_x3.4-1968", charmap.Windows1252, false}, // US-ASCII
		{"oem-850", charmap.CodePage850, false},
		{"oem-852", charmap.CodePage852, false},
		{"oem-855", charmap.CodePage855, false},
		{"oem-858", charmap.CodePage858, false},
		{"oem-860", charmap.CodePage860, false},
		{"oem-862", charmap.CodePage862, false},
		{"oem-863", charmap.CodePage863, false},
		{"oem-865", charmap.CodePage865, false},
		{"oem-866", charmap.CodePage866, false},
		{"oem-1047", charmap.CodePage1047, false},
		{"oem-1140", charmap.CodePage1140, false},
		{"cp28591", charmap.ISO8859_1, false},
		{"windows-28592", charmap.ISO8859_2, false},
		{"cp28593", charmap.ISO8859_3, false},
		{"cp28594", charmap.ISO8859_4, false},
		{"cp28595", charmap.ISO8859_5, false},
		{"cp28596", charmap.ISO8859_6, false},
		{"cp28597", charmap.ISO8859_7, false},
		{"cp28598", charmap.ISO8859_8, false},
		{"cp28599", charmap.ISO8859_9, false},
		{"cp919", charmap.ISO8859_10, false},
		{"cp874", charmap.Windows874, false},
		{"cp921", charmap.ISO8859_13, false},
		{"cp28604", charmap.ISO8859_14, false},
		{"cp923", charmap.ISO8859_15, false},
		{"cp28606", charmap.ISO8859_16, false},
		{"cp878", charmap.KOI8R, false},
		{"cp1168", charmap.KOI8U, false},
		{"cp10000", charmap.Macintosh, false},
		{"oem-1250", charmap.Windows1250, false},
		{"oem-1251", charmap.Windows1251, false},
		{"oem-1252", charmap.Windows1252, false},
		{"oem-1253", charmap.Windows1253, false},
		{"oem-1254", charmap.Windows1254, false},
		{"oem-1255", charmap.Windows1255, false},
		{"oem-1256", charmap.Windows1256, false},
		{"oem-1257", charmap.Windows1257, false},
		{"oem-1258", charmap.Windows1258, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encoding(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoding() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfFileX(t *testing.T) {
	b := []byte("hello\x1aworld")
	if got := EndOfFile(b...); string(got) != "hello" {
		t.Errorf("TestEndOfFile() = %v, want %v", string(got), "hello")
	}
}

func TestRunesControls(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		// {"empty", "", ""}, // returns two different empty values?
		{"hi", "hello world", "hello world"},
		{"nul", "\x00", "␀"},
		{"us", "\x1f", "␟"},
		{"device controls", "\x10\x11", "␐␑"},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = cp1252
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.Swap()
		t.Run(tt.name, func(t *testing.T) {
			if got, w := d.Output, []rune(tt.want); !reflect.DeepEqual(got, w) {
				t.Errorf("TestRunesControls() = %v (%X) [%s], want %v (%X) [%s]",
					got, got, string(got), w, w, string(w))
			}
		})
	}
}

func TestRunesKOI8(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"hi", "hello world", "hello world"},
		{"lines", "\x82\x80\x80hi\x80\x80\x83", "┌──hi──┐"},
		{"invalid", "\x00=NULL & \x1f=?", " =NULL &  =?"},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = koi
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesKOI8()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesKOI8() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRunesLatin(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"hi", "hello world", "hello world"},
		{"high", "\xbd of 5 is 2.5", "½ of 5 is 2.5"},
		{"invalid", "\x00=NULL & \x9f=?", " =NULL &  =?"},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = iso1
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesLatin()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesLatin() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRunesDOS(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"hi", "\x01hello world\x7f", "☺hello world⌂"},
		{"dos pipes", "|!\x7c", "|!|"},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = cp437
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesDOS()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesDOS() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRunesMacintosh(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"hi", "hello world", "hello world"},
		{"controls", "\x11+\x12+Z", "⌘+⇧+Z"},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = mac
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesMacintosh()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesMacintosh() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestRunesWindows(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{"empty", "", ""},
		{"hi", "hello world", "hello world"},
		{"ansi control example", "\x1b[0m;", "\x1b[0m;"},
		{"DEL", "\x7f", "␡"},
		{"invalid", "\x90", " "},
	}
	for _, tt := range tests {
		d := Convert{}
		d.Input.Bytes = []byte(tt.text)
		d.Input.Encoding = cp1252
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesWindows()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesWindows() = %s, want %s", got, tt.want)
			}
		})
	}
}
func TestRunesEBCDIC(t *testing.T) {
	// EBCDIC codepages are not compatible with ISO/IEC 646 (ASCII)
	// so a number of these tests either convert input UTF-8 text into CP037
	tx, err := charmap.CodePage037.NewEncoder().Bytes([]byte("ring my "))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name   string
		text   []byte
		want   string
		encode bool
	}{
		{"empty", []byte{}, "", true},
		{"nul", []byte{0}, "\u2400", true},
		{"ht", []byte{9}, "\u2409", true},
		{"invalid", []byte{4}, "\u0020", false},
		{"bell", []byte{7}, "␇", false},
		{"Æ", []byte{158}, " ", false},
		{"ring", []byte("ring my bell"), "ring my bell", true},
		{"ring symbol", append(tx, []byte{7}...), "ring my ␡", false},
		{"c ben", []byte{180, 64, 194, 133, 149}, "  Ben", false},
		{"c ben", []byte("© Ben"), "  Ben", true},
	}
	for _, tt := range tests {
		c := tt.text
		d := Convert{}
		d.Input.Bytes = c
		d.Input.Encoding = charmap.CodePage037
		if err := d.Transform(); err != nil {
			t.Error(err)
		}
		d.RunesEBCDIC()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(d.Output); got != tt.want {
				fmt.Println(c)
				t.Errorf("RunesEBCDIC() = '%v' (0x%X), want '%v' (0x%X)", got, got, tt.want, tt.want)
			}
		})
	}
}

func Test_equalLB(t *testing.T) {
	var o, n, lf, win = [2]rune{0, 0}, [2]rune{-1, -1}, [2]rune{10, 0}, [2]rune{13, 0}
	type args struct {
		r  [2]rune
		nl [2]rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"negative", args{n, lf}, false},
		{"nil lf", args{o, lf}, false},
		{"nil win", args{o, win}, false},
		{"lf", args{lf, lf}, true},
		{"win", args{win, win}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalLB(tt.args.r, tt.args.nl); got != tt.want {
				t.Errorf("equalLB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encode32(t *testing.T) {
	tests := []struct {
		name string
		a    string
		want encoding.Encoding
	}{
		{"empty", "", nil},
		{"le", "UTF-32", utf32.UTF32(utf32.LittleEndian, utf32.UseBOM)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encode32(tt.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encode32() = %v, want %v", got, tt.want)
			}
		})
	}
}
