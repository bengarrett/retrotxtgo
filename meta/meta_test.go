// Package meta handles the metadata generated through the go builder using ldflags.
package meta_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
)

func Example_digits() {
	fmt.Println(meta.Digits("v1.0 (init release)"))
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
			if got := meta.Digits(tt.s); got != tt.want {
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
		wantVersion meta.Version
	}{
		{"empty", "", false, meta.Version{-1, -1, -1}},
		{"text", "hello world", false, meta.Version{-1, -1, -1}},
		{"zero", meta.GoBuild, true, meta.Version{0, 0, 0}},
		{"vzero", "v0.0.0", true, meta.Version{0, 0, 0}},
		{"ver str", "v1.2.3 (super-release)", true, meta.Version{1, 2, 3}},
		{"short", "V1", true, meta.Version{1, 0, 0}},
		{"short.", "V1.1", true, meta.Version{1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion := meta.Semantic(tt.ver)
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
		v    meta.Version
		want string
	}{
		{"empty", meta.Version{}, "x0.0.0"},
		{"alpha", meta.Version{0, 0, 1}, "α0.0.1"},
		{"beta", meta.Version{0, 5, 11}, "β0.5.11"},
		{"release", meta.Version{10, 0, 9}, "10.0.9"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
