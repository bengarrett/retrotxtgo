package transform

import (
	"reflect"
	"testing"
)

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
			if got := AddBOM(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddBOM() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValid(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{"", false},
		{"437", true},
		{"cp437", true},
		{"CP437", true},
		{"CP 437", true},
		{"CP-437", true},
		{"IBM437", true},
		{"IBM-437", true},
		{"IBM 437", true},
		{"ISO 8859-1", true},
		{"ISO8859-1", true},
		{"ISO88591", true},
		{"isolatin1", true},
		{"latin1", true},
		{"88591", false},
		{"windows1254", true},
		{"win1254", true},
		{"cp1254", true},
		{"1254", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.name); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
