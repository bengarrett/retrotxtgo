package config

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
)

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"0 index", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Keys()[0]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getBool(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want bool
	}{
		{"bool", genr, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getBool(tt.key); got != tt.want {
				t.Errorf("getBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUint(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want uint
	}{
		{"uint", serve, meta.WebPort},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUint(tt.key); got != tt.want {
				t.Errorf("getUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getString(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"string", layoutTmpl, "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getString(tt.key); got != tt.want {
				t.Errorf("getString() = %v, want %v", got, tt.want)
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
			_, err := Marshal()
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMissing(t *testing.T) {
	tests := []struct {
		name      string
		wantCount int
	}{
		{"want no values", len(Keys())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotList := len(Missing()); gotList != tt.wantCount {
				t.Errorf("Missing() = %v, want %v", gotList, tt.wantCount)
			}
		})
	}
}

func TestKeySort(t *testing.T) {
	got, keys := KeySort(), Keys()
	lenG, lenK := len(got), len(keys)
	if lenG == 0 {
		t.Errorf("KeySort() is empty")
		return
	}
	if lenG != lenK {
		t.Errorf("KeySort() length = %d, want %d", lenG, lenK)
	}
	want := fontFamily
	if s := got[0]; s != want {
		t.Errorf("KeySort()[0] = %s, want %s", s, want)
	}
}
