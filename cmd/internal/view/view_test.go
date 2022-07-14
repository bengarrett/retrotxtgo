package view_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/bengarrett/retrotxtgo/lib/convert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestTransform(t *testing.T) {
	const a = "\x01\x02\x03\x7f"
	const s = "☺☻♥⌂"
	cp437 := charmap.CodePage437
	latin1 := charmap.ISO8859_1
	c := convert.Convert{}
	type args struct {
		in   encoding.Encoding
		out  encoding.Encoding
		conv *convert.Convert
		b    []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []rune
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"conv nil", args{
			nil, nil, nil, []byte(a),
		}, nil, true},
		{"no encoding", args{
			nil, nil, &c, []byte(a),
		}, nil, true},
		{"cp437", args{
			cp437, nil, &c, []byte(a),
		}, []rune(s), false},
		{"cp437->latin1", args{
			cp437, latin1, &c, []byte(a),
		}, []rune(s), false},
		{"latin1", args{
			latin1, nil, &c, []byte(a),
		}, []rune{' ', ' ', ' ', ' '}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := view.Transform(tt.args.in, tt.args.out, tt.args.conv, tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() = %q, want %q", got, tt.want)
			}
		})
	}
}
