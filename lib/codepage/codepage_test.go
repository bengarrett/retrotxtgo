package codepage

import (
	"reflect"
	"testing"
)

func TestUTF8(t *testing.T) {
	type args struct {
		b []byte
	}
	bom := append(BOM(), []byte("hello world")...)
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty string", args{}, false},
		{"ascii string", args{[]byte("hello world")}, false},
		{"arabic text", args{[]byte("مرحبا بالعالم")}, true},
		{"cyrillic text", args{[]byte("Привіт, народ")}, true},
		{"bom string", args{bom}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UTF8(tt.args.b); got != tt.want {
				t.Errorf("IsUTF8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToBOM(t *testing.T) {
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
			if got := ToBOM(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToBOM() = %v, want %v", got, tt.want)
			}
		})
	}
}
