package create_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/create"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
)

func ExampleFonts() {
	fmt.Print(create.Fonts()[2])
	// Output: vga
}

func TestAutoFont(t *testing.T) {
	tests := []struct {
		name string
		e    encoding.Encoding
		want create.Font
	}{
		{"empty", nil, create.VGA},
		{"jp", japanese.ShiftJIS, create.Mona},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := create.AutoFont(tt.e); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AutoFont() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFont_File(t *testing.T) {
	tests := []struct {
		name string
		f    create.Font
		want string
	}{
		{"vga", create.VGA, "ibm-vga8.woff2"},
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
		want create.Font
	}{
		{"v", create.VGA},
		{"mona", create.Mona},
		{"a", create.Automatic},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := create.Family(tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Family() = %v, want %v", got, tt.want)
			}
		})
	}
}
