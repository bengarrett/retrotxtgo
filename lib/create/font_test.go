package create

import (
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

func ExampleFonts() {
	fmt.Print(Fonts()[2])
	// Output: vga
}

func TestAutoFont(t *testing.T) {
	tests := []struct {
		name string
		e    encoding.Encoding
		want Font
	}{
		{"empty", nil, VGA},
		{"jp", japanese.ShiftJIS, Mona},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AutoFont(tt.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AutoFont() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFont_File(t *testing.T) {
	tests := []struct {
		name string
		f    Font
		want string
	}{
		{"vga", VGA, "ibm-vga8.woff2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.File(); got != tt.want {
				t.Errorf("Font.File() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFamily(t *testing.T) {
	tests := []struct {
		name string
		want Font
	}{
		{"v", VGA},
		{"mona", Mona},
		{"a", Automatic},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Family(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Family() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFontCSS(t *testing.T) {
	type args struct {
		name  string
		embed bool
	}
	tests := []struct {
		name     string
		args     args
		wantSize int
		wantErr  bool
	}{
		{"empty", args{"", false}, 213, false}, // automatic returns vga
		{"vga embed", args{"vga", true}, 19740, false},
		{"vga", args{"vga", false}, 213, false},
		{"mona", args{"mona", false}, 214, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := FontCSS(tt.args.name, nil, tt.args.embed)
			gotSize := len(gotB)
			if (err != nil) != tt.wantErr {
				t.Errorf("FontCSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSize != tt.wantSize {
				t.Errorf("FontCSS() = %v, want %v", gotSize, tt.wantSize)
			}
		})
	}
}
