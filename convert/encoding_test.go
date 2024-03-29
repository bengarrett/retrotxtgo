//nolint:gochecknoglobals,dupl
package convert_test

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/convert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
)

func ExampleEncodeAlias() {
	s := convert.EncodeAlias("437")
	fmt.Println(s)
	s = convert.EncodeAlias("koi8u")
	fmt.Println(s)
	// Output: IBM437
	// KOI8-U
}

func ExampleEncodeUTF32() {
	s := convert.EncodeUTF32("utf-32")
	fmt.Println(s)
	// Output: UTF-32LE (Use BOM)
}

func ExampleEncoder() {
	e, _ := convert.Encoder("cp437")
	fmt.Println(e)
	e, _ = convert.Encoder("1252")
	fmt.Println(e)
	// Output: IBM Code Page 437
	// Windows 1252
}

func ExampleHumanize() {
	s := convert.Humanize("cp437")
	fmt.Println(s)
	s = convert.Humanize("cp1252")
	fmt.Println(s)
	// Output: IBM437
	// Windows-1252
}

func ExampleShorten() {
	s := convert.Shorten("cp437")
	fmt.Println(s)
	s = convert.Shorten("IBM-437")
	fmt.Println(s)
	// Output: 437
	// 437
}

func ExampleSwap() {
	fmt.Println(string(convert.Swap(convert.DEL)))
	fmt.Println(string(convert.Swap(convert.SquareRoot)))
	// Output: Δ
	// ✓
}

func ExamplePicture() {
	const lf = byte(138)
	r := convert.Picture(lf)
	fmt.Printf("%U %s\n", r, string(r))
	const esc = byte(155)
	r = convert.Picture(esc)
	fmt.Printf("%U %s\n", r, string(r))
	// Output: U+240A ␊
	// U+241B ␛
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
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			data := &convert.Convert{}
			data.Input.Input = tt.text
			data.Input.Encoding = tt.codepage
			err := data.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data, err = data.Swap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Swap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(data.Output) != tt.want {
				t.Errorf("Convert.Transform() = %q, want %q", data.Output, tt.want)
			}
		}
	})
}

func TestANSI(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			data := &convert.Convert{}
			data.Input.Input = []byte(tt.str)
			data.Input.Encoding = tt.codepage
			err := data.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data, err = data.Swap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Swap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			data.ANSIControls()
			if !reflect.DeepEqual(data.Output, tt.want) {
				t.Errorf("Convert.Transform() = %v %q, want %v", data.Output, data.Output, tt.want)
			}
		}
	})
}

func TestEncoder(t *testing.T) { //nolint:funlen
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := convert.Encoder(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encoder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encoder() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestRunesKOI8(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"hi", "hello world", "hello world", false},
		{"lines", "\x82\x80\x80hi\x80\x80\x83", "┌──hi──┐", false},
		{"invalid", "\x00=NULL & \x1f=?", " =NULL &  =?", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			d := convert.Convert{}
			d.Input.Input = []byte(tt.text)
			d.Input.Encoding = koi
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesKOI8()
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesKOI8() = %s, want %s", got, tt.want)
			}
		}
	})
}

func TestRunesLatin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"hi", "hello world", "hello world", false},
		{"high", "\xbd of 5 is 2.5", "½ of 5 is 2.5", false},
		{"invalid", "\x00=NULL & \x9f=?", " =NULL &  =?", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			d := convert.Convert{}
			d.Input.Input = []byte(tt.text)
			d.Input.Encoding = iso1
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesLatin()
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesLatin() = %s, want %s", got, tt.want)
			}
		}
	})
}

func TestRunesDOS(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"hi", "\x01hello world\x7f", "☺hello world⌂", false},
		{"dos pipes", "|!\x7c", "|!|", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			d := convert.Convert{}
			d.Input.Input = []byte(tt.text)
			d.Input.Encoding = cp437
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesDOS()
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesDOS() = %s, want %s", got, tt.want)
			}
		}
	})
}

func TestRunesMacintosh(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"hi", "hello world", "hello world", false},
		{"controls", "\x11+\x12+Z", "⌘+⇧+Z", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			d := convert.Convert{}
			d.Input.Input = []byte(tt.text)
			d.Input.Encoding = mac
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesMacintosh()
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesMacintosh() = %s, want %s", got, tt.want)
			}
		}
	})
}

func TestRunesWindows(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		text    string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"hi", "hello world", "hello world", false},
		{"ansi control example", "\x1b[0m;", "\x1b[0m;", false},
		{"DEL", "\x7f", "␡", false},
		{"invalid", "\x90", " ", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			d := convert.Convert{}
			d.Input.Input = []byte(tt.text)
			d.Input.Encoding = cp1252
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesWindows()
			if got := string(d.Output); got != tt.want {
				t.Errorf("TestRunesWindows() = %s, want %s", got, tt.want)
			}
		}
	})
}

func TestRunesEBCDIC(t *testing.T) {
	t.Parallel()
	// EBCDIC codepages are not compatible with ISO/IEC 646 (ASCII)
	// so a number of these tests either convert input UTF-8 text into CP037
	tx, err := charmap.CodePage037.NewEncoder().Bytes([]byte("ring my "))
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		text    []byte
		want    string
		encode  bool
		wantErr bool
	}{
		{"empty", []byte{}, "", true, false},
		{"nul", []byte{0}, "\u2400", true, false},
		{"ht", []byte{9}, "\u2409", true, false},
		{"invalid", []byte{4}, "\u0020", false, false},
		{"bell", []byte{7}, "␇", false, false},
		{"Æ", []byte{158}, " ", false, false},
		{"ring", []byte("ring my bell"), "ring my bell", true, false},
		{"ring symbol", append(tx, []byte{7}...), "ring my ␡", false, false},
		{"c ben", []byte{180, 64, 194, 133, 149}, "  Ben", false, false},
		{"c ben", []byte("© Ben"), "  Ben", true, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c := tt.text
			d := convert.Convert{}
			d.Input.Input = c
			d.Input.Encoding = charmap.CodePage037
			err := d.Transform()
			if (err != nil) != tt.wantErr {
				t.Errorf("Convert.Transform() error = %v, wantErr %v", err, tt.wantErr)
			}
			d.RunesEBCDIC()
			if got := string(d.Output); got != tt.want {
				t.Errorf("RunesEBCDIC() = '%v' (0x%X), want '%v' (0x%X)", got, got, tt.want, tt.want)
			}
		}
	})
}

func Test_EqualLB(t *testing.T) {
	t.Parallel()
	o, n, lf, win := [2]rune{0, 0}, [2]rune{-1, -1}, [2]rune{10, 0}, [2]rune{13, 0}
	type args struct {
		r  [2]rune
		nl [2]rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, true},
		{"negative", args{n, lf}, false},
		{"nil lf", args{o, lf}, false},
		{"nil win", args{o, win}, false},
		{"lf", args{lf, lf}, true},
		{"win", args{win, win}, true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.EqualLB(tt.args.r, tt.args.nl); got != tt.want {
				t.Errorf("EqualLB() = %v, want %v", got, tt.want)
			}
		}
	})
}

func Test_EncodeUTF32(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		a    string
		want encoding.Encoding
	}{
		{"empty", "", nil},
		{"le", "UTF-32", utf32.UTF32(utf32.LittleEndian, utf32.UseBOM)},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.EncodeUTF32(tt.a); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeUTF32() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestHumanize(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"cp437", args{"CP-437"}, "IBM437"},
		{"win", args{"windows-1252"}, "Windows-1252"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.Humanize(tt.args.name); got != tt.want {
				t.Errorf("Humanize() = %v, want %v", got, tt.want)
			}
		}
	})
}

func Test_Shorten(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"cp437", args{"CP-437"}, "437"},
		{"win", args{"windows-1252"}, "1252"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.Shorten(tt.args.name); got != tt.want {
				t.Errorf("shorten() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestConvert_swaps(t *testing.T) {
	t.Parallel()
	hi := []rune("hello world 😃 ⌂")
	tests := []struct {
		name   string
		output []rune
		swaps  []string
		want   []rune
	}{
		{"empty", nil, nil, nil},
		{"nil swap", hi, nil, hi},
		{"empty swap", hi, []string{}, hi},
		{"replace DEL", hi, []string{"root"}, hi},
		{"replace Home", hi, []string{"house"}, []rune("hello world 😃 Δ")},
		{"replace Home", hi, strings.Split("n,b,h,p,r", ","), []rune("hello world 😃 Δ")},
		{"replace NULLS", []rune("hello\u0000world"), []string{"null"}, []rune("hello world")},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c := convert.Convert{
				Output: tt.output,
			}
			c.Args.SwapChars = tt.swaps
			_, _ = c.Swaps()
			if got := c.Output; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Convert.Swaps() = %v, want %v", string(got), string(tt.want))
			}
		}
	})
}

func Test_picture(t *testing.T) {
	t.Parallel()
	const err = 65533
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want rune
	}{
		{"empty", args{}, err},
		{"Q", args{byte(51)}, err},
		{"NULL", args{byte(0 + convert.Row8)}, 9216},
		{"SOH", args{byte(1 + convert.Row8)}, 9217},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.Picture(tt.args.b); got != tt.want {
				t.Errorf("convert.Picture() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestConvert_skipIgnores(t *testing.T) {
	t.Parallel()
	type args struct {
		i      int
		output []rune
		ignore []rune
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"ign h", args{0, []rune("hello"), []rune("h")}, true},
		{"ign H", args{1, []rune("hello"), []rune("h")}, false},
		{"ign H", args{0, []rune("hello"), []rune("H")}, false},
		{"ign 📙", args{1, []rune(" 📙 "), []rune("abcde📙")}, true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			c := convert.Convert{
				Output: tt.args.output,
			}
			c.Input.Ignore = tt.args.ignore
			if got := c.SkipIgnore(tt.args.i); got != tt.want {
				t.Errorf("Convert.SkipIgnore() = %v, want %v", got, tt.want)
			}
		}
	})
}

func Test_EncodeAlias(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
	}{
		{"", ""},
		{"37", "IBM037"},
		{"10000", "Macintosh"},
		{"win", "Windows-1252"},
		{"8", "ISO-8859-8"},
		{"11", "ISO-8859-11"},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := convert.EncodeAlias(tt.name); got != tt.want {
				t.Errorf("EncodeAlias() = %v, want %v", got, tt.want)
			}
		}
	})
}
