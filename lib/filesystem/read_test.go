package filesystem

import (
	"os"
	"testing"
)

func TestReadColumns(t *testing.T) {
	tmp0 := fileExample("hello world\n", 0)
	tmp1 := fileExample("hello\x0aworld\x0a", 1)
	tmp2 := fileExample("hello 😄😄😄\n", 2)
	tmp3 := fileExample("hello\nworld\n", 3)
	tmp4 := fileExample("hello\x0d\x0aworld\x0d\x0a", 4)
	tmp5 := fileExample("", 5)
	tmp6 := fileExample("\x0d\x0a", 6)
	tmp7 := fileExample("let's\x0duse\x0dan old-skool\x0d8-bit microcomputer\x0dnewline\x0d", 7)
	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{"", -1, true},
		{tmp0, 11, false},
		{tmp1, 5, false},
		{tmp2, 18, false},
		{tmp3, 5, false},
		{tmp4, 5, false},
		{tmp5, 0, false},
		{tmp6, 0, false},
		{tmp7, 19, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := ReadColumns(tt.name)
			os.Remove(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadColumns() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ReadColumns() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReadControls(t *testing.T) {
	tmp0 := fileExample("\x1B\x5b0mhello world\n", 0)
	tmp1 := fileExample("\x1B\x5b1mhello world\x1B\x5b0m\n", 1)
	tmp2 := fileExample("hello \x1B\x5b1m😄😄😄\x1B\x5b0m\n", 2)
	tmp3 := fileExample("\x1B\x5b0m\x1B\x5b34mH\x1B\x5b1me\x1B\x5b32ml\x1B\x5b0;32ml\x1B\x5b1;36mo\x1B\x5b37m w\x1B\x5b0mo\x1B\x5b33mr\x1B\x5b1ml\x1B\x5b35md\x1B\x5b0;35m!\x1B\x5b37m\n", 3)
	tests := []struct {
		name      string
		wantCount int
		wantErr   bool
	}{
		{"", -1, true},
		{tmp0, 1, false},
		{tmp1, 2, false},
		{tmp2, 2, false},
		{tmp3, 13, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := ReadControls(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadControls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ReadControls() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
