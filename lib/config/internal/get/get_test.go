package get_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/meta"
)

func TestBool(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    bool
		wantErr bool
	}{
		{"empty", "", false, true},
		{"bad key", "xyz", false, true},
		{"valid", "html.font.embed", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := get.Bool(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"bad key", "xyz", "", true},
		{"valid", get.LayoutTmpl, "standard", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := get.String(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Bool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUInt(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		want    uint
		wantErr bool
	}{
		{"empty", "", 0, true},
		{"bad key", "xyz", 0, true},
		{"valid", "serve", meta.WebPort, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := get.UInt(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("UInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"marshal", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := get.Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
