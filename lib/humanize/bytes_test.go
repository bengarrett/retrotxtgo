package humanize

import (
	"testing"

	"golang.org/x/text/language"
)

func Test_binary(t *testing.T) {
	var us = language.AmericanEnglish
	type args struct {
		b int64
		t language.Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0", args{int64(0), us}, "0"},
		{"1", args{int64(1), us}, "1B"},
		{"1.5K", args{int64(1500), us}, "1.5 KiB"},
		{"2.5M", args{int64(2500000), us}, "2.38 MiB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Binary(tt.args.b, tt.args.t); got != tt.want {
				t.Errorf("Binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decimal(t *testing.T) {
	var us = language.AmericanEnglish
	type args struct {
		b int64
		t language.Tag
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"0", args{int64(0), us}, "0"},
		{"1", args{int64(1), us}, "1B"},
		{"1.5K", args{int64(1500), us}, "1.5 kB"},
		{"2.5M", args{int64(2500000), us}, "2.50 MB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decimal(tt.args.b, tt.args.t); got != tt.want {
				t.Errorf("decimal() = %v, want %v", got, tt.want)
			}
		})
	}
}
