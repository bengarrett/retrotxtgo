// Package sauce to handle the reading and parsing of embedded SAUCE metadata.

package sauce

import (
	"testing"
)

func TestDataType_String(t *testing.T) {
	tests := []struct {
		name string
		d    DataType
		want string
	}{
		{"none", 0, "undefined"},
		{"executable", 8, "executable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("DataType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lsBit_String(t *testing.T) {
	tests := []struct {
		name string
		ls   lsBit
		want string
	}{
		{"empty", "", invalid},
		{"00", "00", noPref},
		{"8px", "01", "select 8 pixel font"},
		{"9px", "10", "select 9 pixel font"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ls.String(); got != tt.want {
				t.Errorf("lsBit.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Character
		want string
	}{
		{"first", ascii, "ASCII text"},
		{"last", tundraDraw, "TundraDraw color text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Character.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_Desc(t *testing.T) {
	tests := []struct {
		name string
		c    Character
	}{
		{"first", ascii},
		{"last", tundraDraw},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Desc(); got == "" {
				t.Errorf("Character.Desc() = %q", got)
			}
		})
	}
}

func TestBitmap_String(t *testing.T) {
	tests := []struct {
		name string
		b    Bitmap
		want string
	}{
		{"first", gif, "GIF image"},
		{"last", avi, "AVI video"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Bitmap.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_String(t *testing.T) {
	tests := []struct {
		name string
		v    Vector
		want string
	}{
		{"first", dxf, "AutoDesk CAD vector graphic"},
		{"last", kinetix, "3D Studio vector graphic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Vector.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudio_String(t *testing.T) {
	tests := []struct {
		name string
		a    Audio
		want string
	}{
		{"first", mod, "NoiseTracker module"},
		{"midi", midi, "MIDI audio"},
		{"okt", okt, "Oktalyzer module"},
		{"last", it, "Impulse Tracker module"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Audio.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_String(t *testing.T) {
	tests := []struct {
		name string
		a    Archive
		want string
	}{
		{"zip", zip, "ZIP compressed archive"},
		{"squeeze", sqz, "Squeeze It compressed archive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Archive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
