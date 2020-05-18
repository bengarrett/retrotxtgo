package filesystem

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/samples"
)

func exampleSave(s string, i int) string {
	var name = fmt.Sprintf("rt_fs_save%d.txt", i)
	if s == "" {
		s = "hello world"
	}
	path, err := samples.Save([]byte(s), name)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func TestRead(t *testing.T) {
	f := exampleSave("", 0)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Read() error = %v, wantErr %v", err, tt.wantErr)
				//return
			}
		})
	}
	samples.Clean(f)
}

func TestReadAllBytes(t *testing.T) {
	f2 := exampleSave(samples.Symbols, 2)
	f3 := exampleSave(samples.Tabs, 3)
	f4 := exampleSave(samples.Escapes, 4)
	f5 := exampleSave(samples.Digits, 5)
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
		{"utf8", args{f2}, []byte(samples.Symbols), false},
		{"tabs", args{f3}, []byte(samples.Tabs), false},
		{"escs", args{f4}, []byte(samples.Escapes), false},
		{"digs", args{f5}, []byte(samples.Digits), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadAllBytes(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAllBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadAllBytes() = %v, want %v", string(gotData), string(tt.wantData))
			}
		})
	}
	samples.Clean(f2)
	samples.Clean(f3)
	samples.Clean(f4)
	samples.Clean(f5)
}

func TestReadChunk(t *testing.T) {
	f1 := exampleSave(samples.Newlines, 1)
	f2 := exampleSave(samples.Symbols, 2)
	f3 := exampleSave(samples.Tabs, 3)
	f4 := exampleSave(samples.Escapes, 4)
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
		{"range +", args{f2, 20}, []byte(samples.Symbols), false},
		{"nl", args{f1, 4}, []byte("a\nb\n"), false},
		{"utf8", args{f2, 4}, []byte("[☠|☮"), false},
		{"tabs", args{f3, 7}, []byte("☠\tSkull"), false},
		{"escs", args{f4, 13}, []byte("bell:\a,back:\b"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadChunk(tt.args.name, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadChunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadChunk() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
	samples.Clean(f1)
	samples.Clean(f2)
	samples.Clean(f3)
	samples.Clean(f4)
}

func TestReadTail(t *testing.T) {
	f2 := exampleSave(samples.Symbols, 2)
	f3 := exampleSave(samples.Tabs, 3)
	f4 := exampleSave(samples.Escapes, 4)
	f5 := exampleSave(samples.Digits, 5)
	type args struct {
		name   string
		offset int64
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
		{"utf8", args{f2, 1}, []byte(samples.Symbols), false},
		{"tabs", args{f3, 1}, []byte(samples.Tabs), false},
		{"escs", args{f4, 1}, []byte(samples.Escapes), false},
		{"digs", args{f5, 1}, []byte(samples.Digits), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := ReadTail(tt.args.name, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("ReadTail() = %v, want %v", gotData, tt.wantData)
			}
		})
	}
	samples.Clean(f2)
	samples.Clean(f3)
	samples.Clean(f4)
	samples.Clean(f5)
}

func TestTailBytesEmpty(t *testing.T) {
	r, err := ReadTail("", 0)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %s but got %s", e, r)
	}
	if err == nil {
		t.Fatalf("Expected err but got nil")
	}
}

func TestTailBytes0(t *testing.T) {
	r, err := ReadTail("../textfiles/hi.txt", 0)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err != nil {
		t.Fatalf("Expected nil error, %q", err)
	}
}

func TestTailBytesPositive(t *testing.T) {
	r, err := ReadTail("../textfiles/hi.txt", 10)
	var e = []byte("")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err == nil {
		t.Fatalf("Expected an error, EOF")
	}
}

func TestTailBytesRange(t *testing.T) {
	_, err := ReadTail("../textfiles/hi.txt", -999999999)
	if err == nil {
		t.Fatalf("Expected an error, offset: value is too large")
	}
}

func TestTailBytesNegative(t *testing.T) {
	r, err := ReadTail("../textfiles/hi.txt", -9)
	var e = []byte("ልዑል")
	if bytes.Equal(r, e) == false {
		t.Fatalf("Expected %q but got %q", e, r)
	}
	if err != nil {
		t.Fatalf("Expected nil error, %q", err)
	}
}
