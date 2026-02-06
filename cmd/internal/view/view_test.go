package view_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/view"
	"github.com/bengarrett/retrotxtgo/convert"
	"github.com/nalgeon/be"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestTransform(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := view.Transform(tt.args.conv, tt.args.in, tt.args.out, tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() = %q, want %q", got, tt.want)
			}
		}
	})
}

// Test error handling in Transform.
func TestTransformErrorHandling(t *testing.T) {
	t.Parallel()

	// Test with empty bytes
	c := &convert.Convert{}
	r, err := view.Transform(c, nil, nil)
	be.Err(t, err, nil)
	be.Equal(t, len(r), 0)

	// Test with nil converter
	//nolint:ineffassign,staticcheck,wastedassign // r and err are used in assertions below
	r, err = view.Transform(nil, nil, nil)
	be.True(t, err != nil)
	be.Equal(t, err, view.ErrConv)
}

// Test different encoding scenarios.
func TestEncodingScenarios(t *testing.T) {
	t.Parallel()

	c := &convert.Convert{}
	cp437 := charmap.CodePage437
	latin1 := charmap.ISO8859_1

	// Test CP437 to UTF-8
	cp437Bytes := []byte{0x01, 0x02, 0x03} // CP437 special chars
	r, err := view.Transform(c, cp437, nil, cp437Bytes...)
	be.Err(t, err, nil)
	be.True(t, len(r) > 0)

	// Test Latin1 to UTF-8
	latin1Bytes := []byte{0xE9, 0xE8, 0xE0} // Latin1 accented chars
	r, err = view.Transform(c, latin1, nil, latin1Bytes...)
	be.Err(t, err, nil)
	be.True(t, len(r) > 0)

	// Test UTF-8 to UTF-8 (identity)
	utf8Bytes := []byte("Hello World")
	r, err = view.Transform(c, nil, nil, utf8Bytes...)
	be.Err(t, err, nil)
	be.Equal(t, string(r), "Hello World")
}

// Test edge cases in Transform.
func TestTransformEdgeCases(t *testing.T) {
	t.Parallel()

	c := &convert.Convert{}
	cp437 := charmap.CodePage437

	// Test with empty input
	r, err := view.Transform(c, cp437, nil)
	be.Err(t, err, nil)
	be.Equal(t, len(r), 0)

	// Test with special characters
	special := []byte{0x01, 0x02, 0x03, 0x7f} // CP437 special chars
	r, err = view.Transform(c, cp437, nil, special...)
	be.Err(t, err, nil)
	be.True(t, len(r) > 0)

	// Test with Unicode characters
	//nolint:gosmopolitan // Intentional use of Han script to test Unicode handling
	unicodeBytes := []byte("Hello 世界")
	r, err = view.Transform(c, nil, nil, unicodeBytes...)
	be.Err(t, err, nil)
	be.True(t, len(r) > 0)
}
