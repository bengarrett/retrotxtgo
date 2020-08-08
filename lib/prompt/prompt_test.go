package prompt

import (
	"bytes"
	"testing"
)

func Test_pport(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name     string
		validate bool
		input    string
		wantPort uint
	}{
		{"empty", false, "", 0},
		{"empty validate", true, "", 0},
		{"no validation", false, "8000", 8000},
		{"validation", true, "8000", 8000},
		{"invalid 1", false, "90000", 90000},
		{"invalid 2", false, "abcde", 0},
		{"invalid 3", true, "90000", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotPort := pport(&stdin, tt.validate, false); gotPort != tt.wantPort {
				t.Errorf("pport() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestPortValid(t *testing.T) {
	tests := []struct {
		name   string
		port   uint
		wantOk bool
	}{
		{"0", 0, true},
		{"80", 80, true},
		{"8080", 8080, true},
		{"8888", 8888, true},
		{"88888", 88888, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := PortValid(tt.port); gotOk != tt.wantOk {
				t.Errorf("PortValid() = %v, want %v", gotOk, tt.wantOk)
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
			if gotKey := tt.k.prompt(&stdin, false); gotKey != tt.wantKey {
				t.Errorf("keys.prompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func Test_promptRead(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name      string
		input     string
		wantInput string
		wantErr   bool
	}{
		{"empty", "", "", false},
		{"nl", "\n", "", false},
		{"tab", "\t", "", false},
		{"hi", "hello", "hello", false},
		{"hw", "hello world", "hello world", false},
		{"emoji", "hi 😃", "hi 😃", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			gotInput, err := promptRead(&stdin)
			if (err != nil) != tt.wantErr {
				t.Errorf("promptRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInput != tt.wantInput {
				t.Errorf("promptRead() = %v, want %v", gotInput, tt.wantInput)
			}
		})
	}
}

func Test_parseYN(t *testing.T) {
	type args struct {
		input      string
		yesDefault bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"yes 1", args{"", true}, true},
		{"yes 1", args{"y", false}, true},
		{"yes 2", args{"yes", false}, true},
		{"no 1", args{"no", false}, false},
		{"no 2", args{"no", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseYN(tt.args.input, tt.args.yesDefault); got != tt.want {
				t.Errorf("parseYN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keys_shortValidate(t *testing.T) {
	fruits := []string{"Apple", "Watermelon", "Orange", "Pear", "Cherry"}
	tests := []struct {
		name string
		k    keys // should be unique
		key  string
		want string
	}{
		{"empty", keys{""}, "", ""},
		{"match l1", keys{"hi"}, "hi", "hi"},
		{"match s1", keys{"hi"}, "h", "hi"},
		{"match fullname fruit", fruits, "Orange", "Orange"},
		{"match letter of fruit", fruits, "o", "Orange"},
		{"mismatch letter of fruit", fruits, "z", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.shortValidate(tt.key); got != tt.want {
				t.Errorf("keys.shortValidate(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}
