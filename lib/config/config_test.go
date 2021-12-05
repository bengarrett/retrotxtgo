package config_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
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
			if got := config.Keys()[0]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
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
			_, err := config.Marshal()
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
		{"want no values", len(config.Keys())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotList := len(config.Missing()); gotList != tt.wantCount {
				t.Errorf("Missing() = %v, want %v", gotList, tt.wantCount)
			}
		})
	}
}

func TestKeySort(t *testing.T) {
	got, keys := config.KeySort(), config.Keys()
	lenG, lenK := len(got), len(keys)
	if lenG == 0 {
		t.Errorf("KeySort() is empty")
		return
	}
	if lenG != lenK {
		t.Errorf("KeySort() length = %d, want %d", lenG, lenK)
	}
	want := get.FontFamily
	if s := got[0]; s != want {
		t.Errorf("KeySort()[0] = %s, want %s", s, want)
	}
}
