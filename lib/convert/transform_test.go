package convert

import (
	"testing"

	"golang.org/x/text/encoding/unicode"
)

func Test_skipCtrlCodes(t *testing.T) {
	tests := []struct {
		name string
		ctrl []string
		want []rune
	}{
		{"nil", []string{}, []rune{}},
		{"bs", []string{"bs"}, []rune{8}},
		{"v,del", []string{"v", "del"}, []rune{11, 127}},
		{"invalid", []string{"xxx"}, []rune{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c = Convert{}
			c.lineBreaks = true
			c.Flags.Controls = tt.ctrl
			c.skipCtrlCodes()
			if got := c.ignores; string(got) != string(tt.want) {
				t.Errorf("Convert.skipCtrlCodes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArgs_Dump(t *testing.T) {
	tests := []struct {
		name    string
		b       []byte
		want    string
		wantErr bool
	}{
		{"empty", []byte(""), "", true},
		{"hi", []byte("hello\nworld"), "hello\nworld", false},
	}
	for _, tt := range tests {
		var a = Convert{}
		if len(tt.b) > 0 {
			a.Input.Encoding = unicode.UTF8
		}
		t.Run(tt.name, func(t *testing.T) {
			gotUtf8, err := a.Dump(tt.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotUtf8) != tt.want {
				t.Errorf("Args.Dump() gotUtf8 = %v, want %v", string(gotUtf8), tt.want)
			}
		})
	}
}
