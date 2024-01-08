package fsys_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/internal/mock"
	"github.com/bengarrett/retrotxtgo/internal/tmp"
)

func ExampleSaveTemp() {
	file, _ := fsys.SaveTemp("example.txt", []byte("hello world")...)
	defer os.Remove(file)
	s, _ := os.Stat(file)
	fmt.Printf("%s, %d", s.Name(), s.Size())
	// Output:example.txt, 11
}

func ExampleTar() {
	name := tmp.File("tar_test.tar")
	file, _ := fsys.SaveTemp(name, []byte("x")...)
	defer os.Remove(file)
	_ = fsys.Tar(name, file)
	s, _ := os.Stat(file)
	fmt.Printf("%s, %d", s.Name(), s.Size())
	// Output:tar_test.tar, 1536
}

func ExampleTouch() {
	file, _ := fsys.Touch("example.txt")
	defer os.Remove(file)
	s, _ := os.Stat(file)
	fmt.Printf("%s, %d", s.Name(), s.Size())
	// Output:example.txt, 0
}

func ExampleWrite() {
	file, _ := fsys.Touch("example.txt")
	defer os.Remove(file)
	i, _, _ := fsys.Write(file, []byte("hello world")...)
	s, _ := os.Stat(file)
	fmt.Printf("%s, %d", s.Name(), i)
	// Output:example.txt, 10
}

func TestRead(t *testing.T) {
	t.Parallel()
	f := mock.FileExample("hello")
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		t.Cleanup(func() {
			os.Remove(f)
			os.Remove(large)
		})
		for _, tt := range tests {
			_, err := fsys.Read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read(%q) error = %v, wantErr %v", tt.args.name, err, tt.wantErr)
			}
		}
	})
}

func TestReadAllBytes(t *testing.T) {
	t.Parallel()
	f2 := mock.FileExample(mock.T()["Symbols"])
	f3 := mock.FileExample(mock.T()["Tabs"])
	f4 := mock.FileExample(mock.T()["Escapes"])
	f5 := mock.FileExample(mock.T()["Digits"])
	large := mock.LargeExample()
	t.Cleanup(func() {
		os.Remove(f2)
		os.Remove(f3)
		os.Remove(f4)
		os.Remove(f5)
		os.Remove(large)
	})
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotData, err := fsys.ReadAllBytes(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAllBytes(%q) error = %v, wantErr %v", tt.args.name, err, tt.wantErr)
				return
			}
			if tt.wantData != nil && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadAllBytes(%q) = %q, want %q", tt.args.name, string(gotData), string(tt.wantData))
			}
		}
	})
}

func TestReadChunk(t *testing.T) {
	t.Parallel()
	f1 := mock.FileExample(mock.T()["Newline"])
	f2 := mock.FileExample(mock.T()["Symbols"])
	f3 := mock.FileExample(mock.T()["Tabs"])
	f4 := mock.FileExample(mock.T()["Escapes"])
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		t.Cleanup(func() {
			os.Remove(f1)
			os.Remove(f2)
			os.Remove(f3)
			os.Remove(f4)
			os.Remove(large)
		})
		for _, tt := range tests {
			gotData, err := fsys.ReadChunk(tt.args.name, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadChunk(%q) error = %v, wantErr %v", tt.args.name, err, tt.wantErr)
				return
			}
			if tt.name == large && len(gotData) != 100 {
				t.Errorf("ReadChunk(%q) length = %v, want %v", tt.args.name, len(gotData), 100)
			}
			if tt.name != large && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadChunk(%q) = %v, want %v", tt.args.name, gotData, tt.wantData)
			}
		}
	})
}

func TestReadTail(t *testing.T) {
	t.Parallel()
	f1 := mock.FileExample(mock.T()["Newline"])
	f2 := mock.FileExample(mock.T()["Symbols"])
	f3 := mock.FileExample(mock.T()["Tabs"])
	f4 := mock.FileExample(mock.T()["Escapes"])
	large := mock.LargeExample()
	t.Cleanup(func() {
		os.Remove(f1)
		os.Remove(f2)
		os.Remove(f3)
		os.Remove(f4)
		os.Remove(large)
	})
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotData, err := fsys.ReadTail(tt.args.name, tt.args.offset)
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
		}
	})
}

func Test_word(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := fsys.Word(tt.s); got != tt.want {
				t.Errorf("Word() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestTouch(t *testing.T) {
	t.Parallel()
	tmp0 := filepath.Join(os.TempDir(), "testtouch")
	t.Cleanup(func() {
		os.Remove(tmp0)
	})
	type args struct {
		name string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
		wantErr  bool
	}{
		{"empty", args{}, "", true},
		{"tmp", args{tmp0}, tmp0, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			gotPath, err := fsys.Touch(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Touch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Touch() = %v, want %v", gotPath, tt.wantPath)
			}
		}
	})
}
