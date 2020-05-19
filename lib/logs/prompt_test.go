package logs

import (
	"bytes"
	"testing"
)

func Test_promptstring(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name      string
		input     string
		wantWords string
	}{
		{"empty", "", ""},
		{"reset", "-", ""},
		{"hi", "hello", "hello"},
		{"hw", "hello world", "hello world"},
		{"emoji", "hi 😃", "hi 😃"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotWords := pstring(&stdin); gotWords != tt.wantWords {
				t.Errorf("pstring() = %v, want %v", gotWords, tt.wantWords)
			}
		})
	}
}

func Test_keys_validate(t *testing.T) {
	var opt keys = []string{"hello", "world", "hi"}
	tests := []struct {
		name   string
		k      keys
		key    string
		wantOk bool
	}{
		{"empty", keys{}, "", false},
		{"no input", opt, "", false},
		{"wrong", opt, "someinput", false},
		{"ok 1", opt, "hi", true},
		{"ok 2", opt, "world", true},
		{"only one choice", opt, "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := tt.k.validate(tt.key); gotOk != tt.wantOk {
				t.Errorf("keys.validate() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_keys_prompt(t *testing.T) {
	var opt keys = []string{"hello", "world", "hi"}
	var stdin bytes.Buffer
	tests := []struct {
		name    string
		k       keys
		input   string
		wantKey string
	}{
		{"empty", keys{}, "", ""},
		{"no input", opt, "", ""},
		{"wrong", opt, "blahblah", ""},
		{"ok 1", opt, "hi", "hi"},
		{"ok 2", opt, "hello", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotKey := tt.k.prompt(&stdin); gotKey != tt.wantKey {
				t.Errorf("keys.prompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
