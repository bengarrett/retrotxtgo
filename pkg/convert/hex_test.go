package convert_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
)

const c, h = "═╣░╠═", "cdb9b0cccd"

func TestHexDecode(t *testing.T) {
	samp, err := convert.E437(c)
	if err != nil {
		t.Errorf("HexDecode() E437() error = %v", err)
	}
	type args struct {
		hexadecimal string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		{"═╣░╠═", args{h}, samp, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := convert.HexDecode(tt.args.hexadecimal)
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
	samp, err := convert.E437(c)
	if err != nil {
		t.Errorf("HexDecode() E437() error = %v", err)
	}
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
			if gotResult := convert.HexEncode(tt.args.text); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("HexEncode() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
