package key_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/key"
)

func Test_keys_validate(t *testing.T) {
	var opt key.Keys = []string{"hello", "world", "hi"}
	tests := []struct {
		name   string
		k      key.Keys
		key    string
		wantOk bool
	}{
		{"empty", key.Keys{}, "", false},
		{"no input", opt, "", false},
		{"wrong", opt, "someinput", false},
		{"ok 1", opt, "hi", true},
		{"ok 2", opt, "world", true},
		{"only one choice", opt, "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := tt.k.Validate(tt.key); gotOk != tt.wantOk {
				t.Errorf("Keys.Validate() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_keys_prompt(t *testing.T) {
	var opt key.Keys = []string{"hello", "world", "hi"}
	var stdin bytes.Buffer
	tests := []struct {
		name    string
		k       key.Keys
		input   string
		wantKey string
	}{
		{"empty", key.Keys{}, "", ""},
		{"no input", opt, "", ""},
		{"wrong", opt, "blahblah", ""},
		{"ok 1", opt, "hi", "hi"},
		{"ok 2", opt, "hello", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotKey := tt.k.Prompt(os.Stdout, &stdin, false); gotKey != tt.wantKey {
				t.Errorf("Keys.prompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func Test_keys_shortValidate(t *testing.T) {
	fruits := []string{"Apple", "Watermelon", "Orange", "Pear", "Cherry"}
	tests := []struct {
		name string
		k    key.Keys // should be unique
		key  string
		want string
	}{
		{"empty", key.Keys{""}, "", ""},
		{"match l1", key.Keys{"hi"}, "hi", "hi"},
		{"match s1", key.Keys{"hi"}, "h", "hi"},
		{"match fullname fruit", fruits, "Orange", "Orange"},
		{"match letter of fruit", fruits, "o", "Orange"},
		{"mismatch letter of fruit", fruits, "z", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.ShortValidate(tt.key); got != tt.want {
				t.Errorf("keys.ShortValidate(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func Test_keys_shortPrompt(t *testing.T) {
	fruits := key.Keys{"apple", "orange", "grape"}
	tests := []struct {
		name    string
		k       key.Keys
		r       io.Reader
		wantKey string
	}{
		{"no match", fruits, strings.NewReader("some input"), ""},
		{"match", fruits, strings.NewReader("orange"), "orange"},
		{"short", fruits, strings.NewReader("o"), "orange"},
		{"remove", fruits, strings.NewReader("-"), "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.k.ShortPrompt(os.Stdout, tt.r); gotKey != tt.wantKey {
				t.Errorf("key.shortPrompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
