package config

import (
	"os"
	"reflect"
	"testing"
)

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"list", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := List(); (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"empty", ""},
		{"0", "0"},
		{"valid", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.name)
		})
	}
}

func Test_copyKeys(t *testing.T) {
	tests := []struct {
		name     string
		keys     []string
		wantCopy []string
	}{
		{"empty", []string{}, []string{}},
		{"3 vals", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCopy := copyKeys(tt.keys); !reflect.DeepEqual(gotCopy, tt.wantCopy) {
				t.Errorf("copyKeys() = %v, want %v", gotCopy, tt.wantCopy)
			}
		})
	}
}

func Test_names_string(t *testing.T) {
	tests := []struct {
		name  string
		n     names
		theme bool
		want  string
	}{
		{"nil", nil, false, ""},
		{"empty", names{""}, false, ""},
		{"ok", names{"okay"}, false, "okay"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.string(tt.theme); got != tt.want {
				t.Errorf("names.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dirExpansion(t *testing.T) {
	u, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"", ""},
		{"~", u},
		{".", w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := dirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("dirExpansion() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}
