package sgr_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/sgr"
	"github.com/bengarrett/retrotxtgo/pkg/ansi/internal/xterm"
)

// func TestDataStream(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		s    string
// 		want []Ps
// 	}{
// 		{"empty", "", nil},
// 		// {"reset", "0m", []Ps{0}},
// 		// {"invalid", "1234m", nil},
// 		// {"bold bg/fg", "37;40;1m", []Ps{37, 40, 1}},
// 		// {"reset and style", "0;42;31;3m", []Ps{0, 42, 31, 3}},
// 		//{"xterm fg", "38;5;0m", []Ps{}},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := DataStream([]byte(tt.s)); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("DataStream() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestPs_Valid(t *testing.T) {
	tests := []struct {
		name string
		p    sgr.Ps
		want bool
	}{
		{"empty", -1, false},
		{"bold", sgr.Bold, true},
		{"bright white", sgr.BrightWhite, true},
		{"noexist middle", 60, false},
		{"too big", 999, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Valid(); got != tt.want {
				t.Errorf("Ps.Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRGBDecimal(t *testing.T) {
	type args struct {
		r uint8
		g uint8
		b uint8
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"black", args{0, 0, 0}, "000000"},
		{"white", args{255, 255, 255}, "ffffff"},
		{"red", args{255, 0, 0}, "ff0000"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sgr.RGBDecimal(tt.args.r, tt.args.g, tt.args.b); fmt.Sprintf("%06x", got) != tt.want {
				t.Errorf("RGBDecimal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecimalRGB(t *testing.T) {
	tests := []struct {
		name    string
		i       float64
		wantR   uint8
		wantG   uint8
		wantB   uint8
		wantErr bool
	}{
		{"neg", -1, 0, 0, 0, true},
		{"max neg", -16777215, 0, 0, 0, true},
		{"min", 0, 0, 0, 0, false},
		{"max", 16777215, 255, 255, 255, false},
		{"half", 8388607, 127, 255, 255, false},
		{"neg", 16777216, 0, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotG, gotB, err := sgr.DecimalRGB(tt.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecimalRGB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotR != tt.wantR {
				t.Errorf("DecimalRGB() gotR = %v, want %v", gotR, tt.wantR)
			}
			if gotG != tt.wantG {
				t.Errorf("DecimalRGB() gotG = %v, want %v", gotG, tt.wantG)
			}
			if gotB != tt.wantB {
				t.Errorf("DecimalRGB() gotB = %v, want %v", gotB, tt.wantB)
			}
		})
	}
}

func TestXterm(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		wantFG xterm.Color
		wantBG xterm.Color
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sgr.DataStream([]byte(tt.s))
			if !reflect.DeepEqual(got.Foreground, tt.wantFG) {
				t.Errorf("DataStream() = \n%v, want \n%v", got, tt.wantFG)
			}
		})
	}
}

// func TestDataStream(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		s    string
// 		want Attributes
// 	}{
// 		// {"empty", "", Toggle{}},
// 		// {"reset", "0m", Toggle{Foreground: fg.White, Background: bg.Black, Font: Primary}},
// 		// {"invalid", "1234m",
// 		// 	Toggle{PrintChars: []byte("1234m")}},
// 		// {"normal bold", "0;1mHi",
// 		// 	Toggle{Bold: true, Foreground: fg.White, Background: bg.Black, Font: Primary,
// 		// 		PrintChars: []byte("Hi")}},
// 		// {"bold then reset", "1;0mHi",
// 		// 	Toggle{Bold: false, Foreground: fg.White, Background: bg.Black, Font: Primary,
// 		// 		PrintChars: []byte("Hi")}},
// 		// {"bold bg/fg", "37;40;1mABC",
// 		// 	Toggle{Bold: true, Foreground: fg.White, Background: bg.Black,
// 		// 		PrintChars: []byte("ABC")}},
// 		// {"ega 8px", "37;40;15mEGA",
// 		// 	Toggle{Font: IbmEGA8, Foreground: fg.White, Background: bg.Black,
// 		// 		PrintChars: []byte("EGA")}},
// 		{"xterm", "48;5;233;38;5;243mXterm colors!",
// 			Attributes{}},
// 		// 		// {"reset and style", "0;42;31;3m", []Ps{0, 42, 31, 3}},
// 		// 		//{"xterm fg", "38;5;0m", []Ps{}},

// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := DataStream([]byte(tt.s)); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("DataStream() = \n%v, want \n%v", got, tt.want)
// 			}
// 		})
// 	}
// }
