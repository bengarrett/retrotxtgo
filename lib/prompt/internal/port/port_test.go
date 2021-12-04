package port_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/port"
	"github.com/bengarrett/retrotxtgo/meta"
)

func TestValid(t *testing.T) {
	p := strconv.Itoa(int(meta.WebPort))
	tests := []struct {
		name   string
		p      uint
		wantOk bool
	}{
		{p, meta.WebPort, true},
		{"0", 0, false},
		{"80", 80, true},
		{"8888", 8888, true},
		{"88888", 88888, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := port.Valid(tt.p); gotOk != tt.wantOk {
				t.Errorf("Valid() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_Port(t *testing.T) {
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
			if gotPort := port.Port(&stdin, tt.validate, false); gotPort != tt.wantPort {
				t.Errorf("Port() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}
