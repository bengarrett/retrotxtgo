package convert

import (
	"testing"
)

func TestConvert_controls(t *testing.T) {
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
			c.useBreaks = true
			c.Flags.Controls = tt.ctrl
			c.unicodeControls()
			if got := c.Output.ignores; string(got) != string(tt.want) {
				t.Errorf("Convert.unicodeControls() got = %v, want %v", got, tt.want)
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
		{"empty", []byte(""), "", false},
		{"hi", []byte("hello\nworld"), "hello\nworld", false},
	}
	for _, tt := range tests {
		var a = Convert{}
		t.Run(tt.name, func(t *testing.T) {
			gotUtf8, err := a.Dump(&tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.Dump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotUtf8) != tt.want {
				t.Errorf("Args.Dump() = %v, want %v", string(gotUtf8), tt.want)
			}
		})
	}
}

func TestArgs_Text(t *testing.T) {
	tests := []struct {
		name    string
		b       []byte
		want    string
		wantErr bool
	}{
		{"empty", []byte(""), "", false},
		{"hi", []byte("hello\nworld"), "hello\nworld", false},
	}
	for _, tt := range tests {
		var a = Convert{}
		t.Run(tt.name, func(t *testing.T) {
			gotUtf8, err := a.Text(&tt.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("Args.Text() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotUtf8) != tt.want {
				t.Errorf("Args.Text() = %v, want %v", gotUtf8, tt.want)
			}
		})
	}
}
