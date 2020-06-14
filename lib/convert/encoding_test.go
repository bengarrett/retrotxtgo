package convert

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestSet_Transform(t *testing.T) {
	tests := []struct {
		name     string
		codepage string
		str      string
		want     string
		wantErr  bool
	}{
		{"null", "ascii", "\x00", "␀", false},
		{"CP037", "cp037", "\xc8\x51\xba\x93\xcf", "Hé[lõ", false},
		{"bell", "cp037", "ring a \x07", "ring a ␇", false},
		{"CP437", "cp437", "H\x82ll\x93 \x9d\xa7\xf4\x9c\xbe", "Héllô ¥º⌠£╛", false},
		{"⌂", "cp437", "Home sweat \x7f", "Home sweat ⌂", false},
		{"mac", "macintosh", "\x11 command + \x12 shift.", "⌘ command + ⇧ shift.", false},
		{"latin1", "latin1", "abcde", "abcde", false},
		{"6e", "iso-8859-6-e", "ring a \x07", "ring a ␇", false},
		{"koi8", "koi", "\xf5\xf2\xf3\xf3", "УРСС", false},
		{"jp", "shiftjis", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := Set{}
			set.B = []byte(tt.str)
			_, err := set.Transform(tt.codepage)
			set.Swap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(set.R) != tt.want {
				t.Errorf("Set.Transform() = %v, want %v", string(set.R), tt.want)
			}
		})
	}
}

func TestSet_SwapANSI(t *testing.T) {
	tests := []struct {
		name     string
		codepage string
		str      string
		want     []rune
		wantErr  bool
	}{
		{"null", "ascii", "\x00", []rune("␀"), false},
		{"CP037", "cp037", "\xc8\x51\xba\x93\xcf", []rune("Hé[lõ"), false},
		{"ansi dos", "cp437", "\x1b\x5b0m", []rune{27, 91, 48, 109}, false},
		{"ansi win", "cp1252", "\x1b\x5b0m", []rune{27, 91, 48, 109}, false},
		{"panic", "cp1252", "\x1b", []rune{9243}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set := Set{}
			set.B = []byte(tt.str)
			_, err := set.Transform(tt.codepage)
			set.Swap()
			set.SwapANSI()
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(set.R, tt.want) {
				t.Errorf("Set.Transform() = %v %q, want %v", set.R, set.R, tt.want)
			}
		})
	}
}

func TestToBOM(t *testing.T) {
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
		{"existing bom string", args{BOM()}, []byte{239, 187, 191}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddBOM(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddBOM() = %v, want %v", got, tt.want)
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
		{"IBM Codepage 437", nil, true},
		{"CP-437", charmap.CodePage437, false},
		{"oem-us", charmap.CodePage437, false},
		{"ibm-37", charmap.CodePage037, false},
		{"858", charmap.CodePage858, false},
		{"mac", charmap.Macintosh, false},
		{"cp1004", charmap.Windows1252, false},
		{"latin1", charmap.ISO8859_1, false},
		{"ISO-8859-1", charmap.ISO8859_1, false},
		{"latin1", charmap.ISO8859_1, false},
		{"ansi_x3.4-1968", charmap.Windows1252, false},
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

func TestMakeBytes(t *testing.T) {
	if l := len(MakeBytes()); l != 256 {
		t.Errorf("MakeBytes() = %v, want %v", l, 256)
	}
}

func TestCutEOF(t *testing.T) {
	b := []byte("hello\x1aworld")
	s := Set{B: b}
	s.CutEOF()
	if string(s.B) != "hello" {
		t.Errorf("TestCutEOF() = %v, want %v", string(s.B), "hello")
	}
}

func TestRunesControls(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		//{"empty", "", ""},
		{"hi", "hello world", "hello world"},
		{"nul", "\x00", "␀"},
		{"us", "\x1f", "␟"},
		{"device controls", "\x10\x11", "␐␑"},
	}
	for _, tt := range tests {
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("windows-1252"); err != nil {
			t.Error(err)
		}
		s.Swap()
		t.Run(tt.name, func(t *testing.T) {
			if got, w := s.R, []rune(tt.want); !reflect.DeepEqual(got, w) {
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
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("koi8-r"); err != nil {
			t.Error(err)
		}
		s.RunesKOI8()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
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
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("iso-8859-1"); err != nil {
			t.Error(err)
		}
		s.RunesLatin()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
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
		{"dos pipes", "|!\x7c", "¦!¦"},
	}
	for _, tt := range tests {
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("cp437"); err != nil {
			t.Error(err)
		}
		s.RunesDOS()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
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
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("mac"); err != nil {
			t.Error(err)
		}
		s.RunesMacintosh()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
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
		{"ansi control example", "\x1b[0m;", "[0m;"},
		{"DEL", "\x7f", "␡"},
		{"invalid", "\x90", " "},
	}
	for _, tt := range tests {
		s := Set{B: []byte(tt.text)}
		if _, err := s.Transform("Windows-1252"); err != nil {
			t.Error(err)
		}
		s.RunesWindows()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
				t.Errorf("TestRunesWindows() = %s, want %s", got, tt.want)
			}
		})
	}
}
func TestRunesEBCDIC(t *testing.T) {
	// EBCDIC codepages are not compatible with ISO/IEC 646 (ASCII)
	// so a number of these tests either convert input UTF-8 text into CP037
	tx, _ := charmap.CodePage037.NewEncoder().Bytes([]byte("ring my "))
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
		{"Æ", []byte{158}, "Æ", false},
		{"ring", []byte("ring my bell"), "ring my bell", true},
		{"ring symbol", append(tx, []byte{7}...), "ring my ␡", false},
		{"c ben", []byte{180, 64, 194, 133, 149}, "© Ben", false},
		{"c ben", []byte("© Ben"), "© Ben", true},
	}
	for _, tt := range tests {
		c := tt.text
		// if tt.encode {
		// 	c, _ = charmap.CodePage037.NewEncoder().Bytes(tt.text)
		// }
		s := Set{
			B: c,
		}
		if _, err := s.Transform("cp037"); err != nil {
			t.Error(err)
		}
		s.RunesEBCDIC()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
				fmt.Println(c)
				t.Errorf("RunesEBCDIC() = '%v' (0x%X), want '%v' (0x%X)", got, got, tt.want, tt.want)
			}
		})
	}
}

func TestSet_Newlines(t *testing.T) {
	tests := []struct {
		name string
		text []byte
		want [2]rune
	}{
		{"unix", []byte("hello\x0aworld\x0a"), [2]rune{10}},
		{"win", []byte("hello\x0d\x0aworld\x0d\x0a\x1a"), [2]rune{13, 10}},
		{"c64", []byte("hello\x0dworld\x0d"), [2]rune{13}},
		{"ibm", []byte("hello\x15world\x15"), [2]rune{21}},
		{"mix", []byte("\x15Windows newline: \x0d\x0a\x15Unix newline: \x0a\x15"), [2]rune{21}},
	}
	for _, tt := range tests {
		s := Set{
			B: tt.text,
		}
		if _, err := s.Transform(""); err != nil {
			t.Error(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := s.Newlines(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Set.Newlines() = %v, want %v", got, tt.want)
			}
		})
	}
}
