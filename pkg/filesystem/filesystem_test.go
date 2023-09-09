package filesystem_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/filesystem"
	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
	"github.com/bengarrett/retrotxtgo/pkg/internal/tmp"
)

const windows = "windows"

func ExampleTar() {
	tmpTar := tmp.File("tar_test.tar")
	tmpFile, err := filesystem.SaveTemp(tmpTar, []byte("x")...)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile)
	if err := filesystem.Tar(tmpTar, tmpFile); err != nil {
		log.Print(err)
		return
	}
	f, err := os.Stat(tmpFile)
	if err != nil {
		log.Print(err)
		return
	}
	fmt.Printf("%s, %d", f.Name(), f.Size())
	// Output:tar_test.tar, 1536
}

func BenchmarkReadLarge(b *testing.B) {
	large := mock.LargeExample()
	if _, err := filesystem.Read(large); err != nil {
		filesystem.Clean(large)
		log.Fatal(err)
	}
	filesystem.Clean(large)
}

func BenchmarkReadMega(b *testing.B) {
	mega := mock.MegaExample()
	if _, err := filesystem.Read(mega); err != nil {
		filesystem.Clean(mega)
		log.Fatal(err)
	}
	filesystem.Clean(mega)
}

func ExampleClean() {
	path, err := filesystem.SaveTemp("examplesave.txt", []byte("hello world")...)
	if err != nil {
		filesystem.Clean(path)
		log.Fatal(err)
	}
	filesystem.Clean(path)
	// Output:
}

func Test_DirExpansion(t *testing.T) {
	h, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	hp := filepath.Dir(h)
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	wp := filepath.Dir(w)
	s := string(os.PathSeparator)

	var tests mock.DirTests
	if runtime.GOOS == windows {
		tests = mock.WindowsTests(h, hp, s, w, wp)
	} else {
		tests = mock.NixTests(h, hp, s, w, wp)
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if gotDir := filesystem.DirExpansion(tt.Name); gotDir != tt.WantDir {
				t.Errorf("DirExpansion(%v) = %v, want %v", tt.Name, gotDir, tt.WantDir)
			}
		})
	}
}

func TestRead(t *testing.T) {
	f := mock.FileExample("hello", 0)
	large := mock.LargeExample()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{""}, true},
		{"invalid", args{"/invalid-file"}, true},
		{"dir", args{os.TempDir()}, true},
		{"valid", args{f}, false},
		{"1.5MB", args{large}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := filesystem.Read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	filesystem.Clean(f)
	filesystem.Clean(large)
}

func TestReadAllBytes(t *testing.T) {
	f2 := mock.FileExample(mock.T()["Symbols"], 2)
	f3 := mock.FileExample(mock.T()["Tabs"], 3)
	f4 := mock.FileExample(mock.T()["Escapes"], 4)
	f5 := mock.FileExample(mock.T()["Digits"], 5)
	large := mock.LargeExample()
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"empty", args{""}, nil, true},
		{"invalid", args{"/invalid-file"}, nil, true},
		{"dir", args{os.TempDir()}, nil, true},
		{"utf8", args{f2}, []byte(mock.T()["Symbols"]), false},
		{"tabs", args{f3}, []byte(mock.T()["Tabs"]), false},
		{"escs", args{f4}, []byte(mock.T()["Escapes"]), false},
		{"digs", args{f5}, []byte(mock.T()["Digits"]), false},
		{"1.5MB", args{large}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := filesystem.ReadAllBytes(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAllBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantData != nil && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadAllBytes() = %q, want %q", string(gotData), string(tt.wantData))
			}
		})
	}
	filesystem.Clean(f2)
	filesystem.Clean(f3)
	filesystem.Clean(f4)
	filesystem.Clean(f5)
	filesystem.Clean(large)
}

func TestReadChunk(t *testing.T) {
	f1 := mock.FileExample(mock.T()["Newline"], 1)
	f2 := mock.FileExample(mock.T()["Symbols"], 2)
	f3 := mock.FileExample(mock.T()["Tabs"], 3)
	f4 := mock.FileExample(mock.T()["Escapes"], 4)
	large := mock.LargeExample()
	type args struct {
		name string
		size int
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"empty", args{"", 0}, nil, true},
		{"invalid", args{"/invalid-file", 0}, nil, true},
		{"dir", args{os.TempDir(), 0}, nil, true},
		{"range 0", args{"", 10}, nil, true},
		{"range -", args{f2, -20}, []byte{}, false},
		{"range +", args{f2, 20}, []byte(mock.T()["Symbols"]), false},
		{"nl", args{f1, 4}, []byte("a\nb\n"), false},
		{"utf8", args{f2, 4}, []byte("[☠|☮"), false},
		{"tabs", args{f3, 7}, []byte("☠\tSkull"), false},
		{"escs", args{f4, 13}, []byte("bell:\a,back:\b"), false},
		{large, args{large, 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := filesystem.ReadChunk(tt.args.name, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadChunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == large && len(gotData) != 100 {
				t.Errorf("ReadChunk() length = %v, want %v", len(gotData), 100)
			}
			if tt.name != large && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadChunk() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
	filesystem.Clean(f1)
	filesystem.Clean(f2)
	filesystem.Clean(f3)
	filesystem.Clean(f4)
	filesystem.Clean(large)
}

func TestReadTail(t *testing.T) {
	f1 := mock.FileExample(mock.T()["Newline"], 1)
	f2 := mock.FileExample(mock.T()["Symbols"], 2)
	f3 := mock.FileExample(mock.T()["Tabs"], 3)
	f4 := mock.FileExample(mock.T()["Escapes"], 4)
	large := mock.LargeExample()
	type args struct {
		name   string
		offset int
	}
	tests := []struct {
		name     string
		args     args
		wantData []byte
		wantErr  bool
	}{
		{"empty", args{"", 0}, nil, true},
		{"invalid", args{"/invalid-file", 0}, nil, true},
		{"dir", args{os.TempDir(), 0}, nil, true},
		{"range", args{"", 10}, nil, true},
		{"utf8", args{f2, 4}, []byte("☮|♺]"), false},
		{"tabs", args{f3, 11}, []byte("♺\tRecycling"), false},
		{"escs", args{f4, 9}, []byte("\v,quote:\""), false},
		{large, args{large, 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := filesystem.ReadTail(tt.args.name, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.name == large && len(gotData) != 100 {
				t.Errorf("ReadChunk() length = %v, want %v", len(gotData), 100)
			}
			if tt.name != large && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadTail() = %q, want %q", string(gotData), string(tt.wantData))
			}
		})
	}
	filesystem.Clean(f1)
	filesystem.Clean(f2)
	filesystem.Clean(f3)
	filesystem.Clean(f4)
	filesystem.Clean(large)
}

func Test_word(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", false},
		{"1", "something", true},
		{"2", "some things", true},
		{"!@#", "!@#", true},
		{"1234.5", "1234.5", true},
		{"你好世界", "你好世界", true},
		{"😀", "😀", false},
		{"😀smiley", "😀smiley", false},
		{"▃▃▃▃▃", "▃▃▃▃▃", false},
		{"nl", "hello\nworld", true},
		{"nl😀", "hello\n😀", true},
		{"😀nl", "😀\nsmiley", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filesystem.Word(tt.s); got != tt.want {
				t.Errorf("Word() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTouch(t *testing.T) {
	type args struct {
		name string
	}
	tmpFile := filepath.Join(os.TempDir(), "testtouch")
	tests := []struct {
		name     string
		args     args
		wantPath string
		wantErr  bool
	}{
		{"empty", args{}, "", true},
		{"tmp", args{tmpFile}, tmpFile, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, err := filesystem.Touch(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Touch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Touch() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
	filesystem.Clean(tmpFile)
}
