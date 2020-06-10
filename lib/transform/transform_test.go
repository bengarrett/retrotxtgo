package transform

import (
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestSet_Transform(t *testing.T) {
	set := Set{}
	tests := []struct {
		name     string
		codepage string
		str      string
		want     string
		wantErr  bool
	}{
		//{"CP037", "cp037", "\xc8\x51\xba\x93\xcf", "Hé[lõ", true},
		//{"CP437", "cp437", "H\x82ll\x93 \x9d\xa7\xf4\x9c\xbe", "Héllô ¥º⌠£╛", true},
		{"⌂", "cp437", "Home sweat \x7f", "Home seat ⌂", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			set.B = []byte(tt.str)
			_, err := set.Transform(tt.codepage)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set.Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(set.B) != tt.want {
				t.Errorf("Set.Transform() = %v, want %v", string(set.B), tt.want)
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
		{"CP-437", nil, true},
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

func TestReplace(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"ibm437", "ibm437"},
		{"CP437", "cp437"},
		{"ISO8859-1", "iso-8859-1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Replace(tt.name); got != tt.want {
				t.Errorf("Replace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
		{"437", true},
		{"cp437", true},
		{"CP437", true},
		{"CP 437", true},
		{"CP-437", true},
		{"IBM437", true},
		{"IBM-437", true},
		{"IBM 437", true},
		{"ISO 8859-1", true},
		{"ISO8859-1", true},
		{"ISO88591", true},
		{"isolatin1", true},
		{"latin1", true},
		{"88591", false},
		{"windows1254", true},
		{"win1254", true},
		{"cp1254", true},
		{"1254", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.name); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeMap(t *testing.T) {
	if l := len(MakeMap()); l != 256 {
		t.Errorf("MakeMap() = %v, want %v", l, 256)
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
		s.RunesControls()
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
		if _, err := s.Transform("windows-1252"); err != nil {
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
	// NOTE: EBCDIC codepages are not compatible with ISO/IEC 646 (ASCII)
	// so a number of these tests either convert input UTF-8 ttext into CP037
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
		{"bell", []byte{47}, "␇", false},
		{"Æ", []byte{158}, "Æ", false},
		{"ring", []byte("ring my bell"), "ring my bell", true},
		{"ring symbol", append(tx, []byte{47}...), "ring my ␡", false},
		{"c ben", []byte{180, 64, 194, 133, 149}, "© Ben", false},
		{"c ben", []byte("© Ben"), "© Ben", true},
	}
	for _, tt := range tests {
		c := tt.text
		if tt.encode {
			c, _ = charmap.CodePage037.NewEncoder().Bytes(tt.text)
		}
		s := Set{
			B: c,
		}
		if _, err := s.Transform("cp037"); err != nil {
			t.Error(err)
		}
		s.RunesEBCDIC()
		t.Run(tt.name, func(t *testing.T) {
			if got := string(s.R); got != tt.want {
				t.Errorf("RunesEBCDIC() = '%v' (0x%X), want '%v' (0x%X)", got, got, tt.want, tt.want)
			}
		})
	}
}
