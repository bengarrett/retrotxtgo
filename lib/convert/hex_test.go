package convert

import (
	"reflect"
	"testing"
)

const c, h = "═╣░╠═", "cdb9b0cccd"

var samp, _ = E437(c)

func TestHexDecode(t *testing.T) {
	type args struct {
		hexadecimal string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"═╣░╠═", args{h}, []byte(samp), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := HexDecode(tt.args.hexadecimal)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexDecode() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestHexEncode(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
	}{
		{"═╣░╠═", args{string(samp)}, []byte{99, 100, 98, 57, 98, 48, 99, 99, 99, 100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := HexEncode(tt.args.text); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexEncode() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
