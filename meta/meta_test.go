// Package meta handles the metadata generated through the go builder using ldflags.
package meta

import (
	"fmt"
	"reflect"
	"testing"
)

func Example_digits() {
	fmt.Println(digits("v1.0 (init release)"))
	// Output: 1.0
}

func Test_digits(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"digits", "01234567890", "01234567890"},
		{"symbols", "~!@#$%^&*()_+", ""},
		{"mixed", "A0B1C2D3E4F5G6H7I8J9K0L", "01234567890"},
		{"semantic", "v1.0.0 (FINAL)", "1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := digits(tt.s); got != tt.want {
				t.Errorf("digits() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSemantic(t *testing.T) {
	tests := []struct {
		name        string
		ver         string
		wantOk      bool
		wantVersion Version
	}{
		{"empty", "", false, Version{-1, -1, -1}},
		{"text", "hello world", false, Version{-1, -1, -1}},
		{"zero", "0.0.0", true, Version{0, 0, 0}},
		{"vzero", "v0.0.0", true, Version{0, 0, 0}},
		{"ver str", "v1.2.3 (super-release)", true, Version{1, 2, 3}},
		{"short", "V1", true, Version{1, 0, 0}},
		{"short.", "V1.1", true, Version{1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion := Semantic(tt.ver)
			gotOk := gotVersion.Valid()
			if gotOk != tt.wantOk {
				t.Errorf("Semantic() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotVersion, tt.wantVersion) {
				t.Errorf("Semantic() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name string
		v    Version
		want string
	}{
		{"empty", Version{}, "α0.0.0"},
		{"alpha", Version{0, 0, 1}, "α0.0.1"},
		{"beta", Version{0, 5, 11}, "β0.5.11"},
		{"release", Version{10, 0, 9}, "10.0.9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
