// Package convert is extends Go's x/text/encoding capability
// to convert legacy encoded text to a modern UTF-8 encoding.
package convert

import (
	"fmt"
	"reflect"
	"testing"
)

func ExampleBOM() {
	fmt.Printf("%X", BOM())
	// Output: EFBBBF
}

func TestEndOfFile(t *testing.T) {
	tests := []struct {
		name string
		b    []byte
		want []byte
	}{
		{"empty", nil, nil},
		{"none", []byte("hello world"), []byte("hello world")},
		{"one", []byte("hello\x1aworld"), []byte("hello")},
		{"two", []byte("hello\x1aworld\x1athis should be hidden"), []byte("hello")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EndOfFile(tt.b...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EndOfFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeBytes(t *testing.T) {
	if l := len(MakeBytes()); l != 256 {
		t.Errorf("MakeBytes() = %v, want %v", l, 256)
	}
}

func TestMark(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"empty string", args{}, []byte{239, 187, 191}},
		{"ascii string", args{[]byte("hi")}, []byte{239, 187, 191, 104, 105}},
		{"existing bom string", args{BOM()}, []byte{239, 187, 191}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Mark(tt.args.b...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Mark() = %v, want %v", got, tt.want)
			}
		})
	}
}
