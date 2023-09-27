// Package sample opens and encodes the example textfiles embedded into the program.
package sample_test

import (
	"fmt"
	"reflect"
	"testing"
	"unicode/utf8"

	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/sample"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleMap() {
	s := sample.Map()["037"]
	fmt.Printf("%s - %s, %s", s.Name, s.Description, s.Encoding)
	// Output: text/cp037.txt - EBCDIC 037 IBM mainframe test, IBM Code Page 037
}

func ExampleOpen() {
	b, _ := sample.Open("037")
	fmt.Print(len(b))
	// Output:130
}

func ExampleFlags_Open() {
	c := convert.Convert{}
	f := sample.Flags{Input: charmap.CodePage037}
	r, _ := f.Open(&c, "037")
	fmt.Print(string(r[0:15]))
	// Output: RetroTxt EBCDIC
}

func TestFlags_Open(t *testing.T) {
	t.Parallel()
	var file convert.Convert
	type fields struct {
		From encoding.Encoding
		To   encoding.Encoding
	}
	type args struct {
		conv *convert.Convert
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantRunes bool
		wantErr   bool
	}{
		{"empty", fields{}, args{nil, ""}, false, true},
		{"missing", fields{}, args{nil, "abcde"}, false, true},
		{"no conv", fields{}, args{nil, "037"}, false, true},
		{"okay 037", fields{}, args{&file, "037"}, true, false},
		{"okay 1252", fields{nil, nil}, args{&file, "1252"}, true, false},
		{"no convert", fields{nil, charmap.CodePage1047}, args{&file, "037"}, true, false},
		{"convert 1252", fields{nil, charmap.CodePage437}, args{&file, "1252"}, true, false},
		{"convert to cp437", fields{charmap.Windows1252, charmap.CodePage437}, args{&file, "1252"}, true, false},
		{"ansi", fields{}, args{&file, "ansi"}, true, false},
		{"cp437 dump", fields{}, args{&file, "437.cr"}, true, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			f := sample.Flags{
				Input: tt.fields.From,
				// Output: tt.fields.To,
			}
			gotR, err := f.Open(tt.args.conv, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Flags.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r := bool(len(gotR) > 0)
			if r != tt.wantRunes {
				t.Errorf("Flags.Open() runes = %v, want %v", r, tt.wantRunes)
			}
		}
	})
}

func TestValid(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"invalid", args{"invalid filename"}, false},
		{"okay", args{"ansi.setmodes"}, true},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := sample.Valid(tt.args.name); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestTransform(t *testing.T) {
	t.Parallel()
	hi := []byte("hello world")
	dos := []byte("▒▓█")
	type args struct {
		e encoding.Encoding
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"empty", args{}, nil, true},
		{"no encoding", args{nil, hi}, nil, true},
		{"no text", args{charmap.CodePage437, nil}, []byte{}, false},
		{"ascii text", args{charmap.CodePage437, hi}, hi, false},
		{"dos text", args{charmap.CodePage437, dos}, []byte{177, 178, 219}, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := sample.Transform(tt.args.e, tt.args.b...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Transform() = %v, want %v", got, tt.want)
			}
		}
	})
}

func TestOpen(t *testing.T) {
	t.Parallel()
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantLen int
		wantErr bool
	}{
		{"empty", args{}, 0, true},
		{"name error", args{"name that doesn't exist"}, 0, true},
		{"865", args{"865"}, 117, false},
		{"utf8", args{"utf8"}, 128, false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			got, err := sample.Open(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLen := utf8.RuneCount(got); gotLen != tt.wantLen {
				t.Errorf("Open() = %v, want %v", gotLen, tt.wantLen)
			}
		}
	})
}
