package filesystem

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// fileExample the string to a text file.
func fileExample(s string, i int) (path string) {
	var name = fmt.Sprintf("rt_fs_save%d.txt", i)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// largeExample generates and saves a 800k file of random us-ascii text.
func largeExample() (path string) {
	const name = "rs_mega_example_save.txt"
	_, s := filler(0.8)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// megaExample generates and saves a 1.5MB file of random us-ascii text.
func megaExample() (path string) {
	const name = "rs_giga_mega_save.txt"
	_, s := filler(1.5)
	path, err := SaveTemp(name, []byte(s)...)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

// filler generates random us-ascii text.
func filler(sizeMB float64) (length int, random string) {
	if sizeMB <= 0 {
		return length, random
	}
	// make characters to randomize
	const (
		// ascii code points (rune codes)
		start    = 33  // "!"
		end      = 122 // "z"
		charsLen = end - start + 1
	)
	chars := make([]rune, charsLen)
	for c, i := 0, start; i <= end; i++ {
		chars[c] = rune(i)
		c++
	}
	// initialize rune slice
	f := (math.Pow(1000, 2) * sizeMB)
	s := make([]rune, int(f))
	// generate random string
	for i := range s {
		s[i] = chars[rand.Intn(charsLen)]
	}
	return len(s), string(s)
}

func BenchmarkReadLarge(b *testing.B) {
	large := largeExample()
	defer Clean(large)
	_, err := Read(large)
	if err != nil {
		log.Fatal(err)
	}
}

func BenchmarkReadMega(b *testing.B) {
	mega := megaExample()
	defer Clean(mega)
	_, err := Read(mega)
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleSave() {
	path, err := SaveTemp("examplesave.txt", []byte("hello world")...)
	if err != nil {
		log.Fatal(err)
	}
	defer Clean(path)
	// Output:
}

func Test_filler(t *testing.T) {
	tests := []struct {
		name       string
		sizeMB     float64
		wantLength int
	}{
		{"0", 0, 0},
		{"0.1", 0.1, 100000},
		{"1.5", 1.5, 1500000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLength, _ := filler(tt.sizeMB); gotLength != tt.wantLength {
				t.Errorf("filler() = %v, want %v", gotLength, tt.wantLength)
			}
		})
	}
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
	type dirTests []struct {
		name    string
		wantDir string
	}
	var tests dirTests
	if runtime.GOOS == "windows" {
		// WINDOWS
		tests = dirTests{
			{fmt.Sprintf("C:%shome%suser", s, s), fmt.Sprintf("C:%shome%suser", s, s)},
			{"~", h},
			{filepath.Join("~", "foo"), filepath.Join(h, "foo")},
			{".", w},
			{fmt.Sprintf(".%sfoo", s), filepath.Join(w, "foo")},
			{fmt.Sprintf("..%sfoo", s), filepath.Join(wp, "foo")},
			{fmt.Sprintf("~%s..%sfoo", s, s), filepath.Join(hp, "foo")},
			{fmt.Sprintf("d:%sroot%sfoo%s..%sblah", s, s, s, s), fmt.Sprintf("D:%sroot%sblah", s, s)},
			{fmt.Sprintf("z:%sroot%sfoo%s.%sblah", s, s, s, s), fmt.Sprintf("Z:%sroot%sfoo%sblah", s, s, s)},
		}
	} else {
		// LINUX, UNIX
		tests = dirTests{
			{fmt.Sprintf("%shome%suser", s, s), fmt.Sprintf("%shome%suser", s, s)},
			{"~", h},
			{filepath.Join("~", "foo"), filepath.Join(h, "foo")},
			{".", w},
			{fmt.Sprintf(".%sfoo", s), filepath.Join(w, "foo")},
			{fmt.Sprintf("..%sfoo", s), filepath.Join(wp, "foo")},
			{fmt.Sprintf("~%s..%sfoo", s, s), filepath.Join(hp, "foo")},
			{fmt.Sprintf("%sroot%sfoo%s..%sblah", s, s, s, s), fmt.Sprintf("%sroot%sblah", s, s)},
			{fmt.Sprintf("%sroot%sfoo%s.%sblah", s, s, s, s), fmt.Sprintf("%sroot%sfoo%sblah", s, s, s)},
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := DirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("DirExpansion(%v) = %v, want %v", tt.name, gotDir, tt.wantDir)
			}
		})
	}
}

func TestRead(t *testing.T) {
	f := fileExample("hello", 0)
	large := largeExample()
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
			_, err := Read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	Clean(f)
	Clean(large)
}

func TestReadAllBytes(t *testing.T) {
	f2 := fileExample(T()["Symbols"], 2)
	f3 := fileExample(T()["Tabs"], 3)
	f4 := fileExample(T()["Escapes"], 4)
	f5 := fileExample(T()["Digits"], 5)
	large := largeExample()
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
		{"utf8", args{f2}, []byte(T()["Symbols"]), false},
		{"tabs", args{f3}, []byte(T()["Tabs"]), false},
		{"escs", args{f4}, []byte(T()["Escapes"]), false},
		{"digs", args{f5}, []byte(T()["Digits"]), false},
		{"1.5MB", args{large}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadAllBytes(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAllBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantData != nil && !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadAllBytes() = %q, want %q", string(gotData), string(tt.wantData))
			}
		})
	}
	Clean(f2)
	Clean(f3)
	Clean(f4)
	Clean(f5)
	Clean(large)
}

func TestReadChunk(t *testing.T) {
	f1 := fileExample(T()["Newline"], 1)
	f2 := fileExample(T()["Symbols"], 2)
	f3 := fileExample(T()["Tabs"], 3)
	f4 := fileExample(T()["Escapes"], 4)
	large := largeExample()
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
		{"range -", args{f2, -20}, nil, false},
		{"range +", args{f2, 20}, []byte(T()["Symbols"]), false},
		{"nl", args{f1, 4}, []byte("a\nb\n"), false},
		{"utf8", args{f2, 4}, []byte("[☠|☮"), false},
		{"tabs", args{f3, 7}, []byte("☠\tSkull"), false},
		{"escs", args{f4, 13}, []byte("bell:\a,back:\b"), false},
		{large, args{large, 100}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadChunk(tt.args.name, tt.args.size)
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
	//Clean(f1)
	Clean(f2)
	Clean(f3)
	Clean(f4)
	Clean(large)
}

func TestReadTail(t *testing.T) {
	f1 := fileExample(T()["Newline"], 1)
	f2 := fileExample(T()["Symbols"], 2)
	f3 := fileExample(T()["Tabs"], 3)
	f4 := fileExample(T()["Escapes"], 4)
	large := largeExample()
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
			gotData, err := ReadTail(tt.args.name, tt.args.offset)
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
	Clean(f1)
	Clean(f2)
	Clean(f3)
	Clean(f4)
	Clean(large)
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
			if got := word(tt.s); got != tt.want {
				t.Errorf("word() = %v, want %v", got, tt.want)
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
			gotPath, err := Touch(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Touch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotPath != tt.wantPath {
				t.Errorf("Touch() = %v, want %v", gotPath, tt.wantPath)
			}
		})
	}
	Clean(tmpFile)
}
