package fsys_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/fsys"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"github.com/bengarrett/retrotxtgo/pkg/nl"
)

func TestReadColumns(t *testing.T) {
	tmp0 := mock.FileExample("hello world\n", 0)
	tmp1 := mock.FileExample("hello\x0aworld\x0a", 1)
	tmp2 := mock.FileExample("hello 😄😄😄\n", 2)
	tmp3 := mock.FileExample("hello\nworld\n", 3)
	tmp4 := mock.FileExample("hello\x0d\x0aworld\x0d\x0a", 4)
	tmp5 := mock.FileExample("", 5)
	tmp6 := mock.FileExample("\x0d\x0a", 6)
	tmp7 := mock.FileExample("let's\x0duse\x0dan old-skool\x0d8-bit microcomputer\x0dnewline\x0d", 7)
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
			gotCount, err := fsys.ReadColumns(tt.name)
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
	tmp0 := mock.FileExample("\x1B\x5b0mhello world\n", 0)
	tmp1 := mock.FileExample("\x1B\x5b1mhello world\x1B\x5b0m\n", 1)
	tmp2 := mock.FileExample("hello \x1B\x5b1m😄😄😄\x1B\x5b0m\n", 2)
	tmp3 := mock.FileExample("\x1B\x5b0m\x1B\x5b34mH\x1B\x5b1me\x1B\x5b32ml\x1B\x5b0;32ml\x1B\x5b1;36mo\x1B\x5b37m "+
		"w\x1B\x5b0mo\x1B\x5b33mr\x1B\x5b1ml\x1B\x5b35md\x1B\x5b0;35m!\x1B\x5b37m\n", 3)
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
			gotCount, err := fsys.ReadControls(tt.name)
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

func TestIsPipe(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"none", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fsys.IsPipe(); got != tt.want {
				t.Errorf("IsPipe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadLine(t *testing.T) {
	tmp0 := mock.FileExample("hello\nworld\n", 0)
	type args struct {
		name      string
		linebreak nl.LineBreaks
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{"none", args{"", nl.NL}, "", true},
		{"tmp0", args{tmp0, nl.NL}, "hello\nworld\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, err := fsys.ReadLine(tt.args.name, tt.args.linebreak)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotText != tt.wantText {
				t.Errorf("ReadLine() = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestReadLines(t *testing.T) {
	tmp0 := mock.FileExample("hello\nworld\n", 0)
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"none", args{""}, -1, true},
		{"tmp0", args{tmp0}, 2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := fsys.ReadLines(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ReadLines() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReadText(t *testing.T) {
	tmp0 := mock.FileExample("hello\nworld\n", 0)
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantText string
		wantErr  bool
	}{
		{"empty", args{}, "", true},
		{"invalid", args{"this_file_doesnt_exist"}, "", true},
		{"tmp0", args{tmp0}, "hello\nworld\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotText, err := fsys.ReadText(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotText != tt.wantText {
				t.Errorf("ReadText() = %v, want %v", gotText, tt.wantText)
			}
		})
	}
}

func TestReadWords(t *testing.T) {
	tmp0 := mock.FileExample("hello\nworld,\nmy name is Ben\n", 0)
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		args      args
		wantCount int
		wantErr   bool
	}{
		{"empty", args{}, -1, true},
		{"invalid", args{"this_file_doesnt_exist"}, -1, true},
		{"tmp0", args{tmp0}, 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCount, err := fsys.ReadWords(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("ReadWords() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func TestReadPipe(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"hi", args{"hello world"}, []byte("hello world\n"), false},
		{"nl", args{"hello\nworld"}, []byte("hello\nworld\n"), false},
		{"utf8", args{"hello 😄"}, []byte("hello 😄\n"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			gotB, err := fsys.ReadPipe()
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadPipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !bytes.Equal(gotB, tt.want) {
				t.Errorf("ReadPipe() = %v, want %v", gotB, tt.want)
			}
		})
	}
}
