// Package sample opens and encodes the example textfiles embedded into the program.
package sample

import (
	"fmt"
	"log"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"retrotxt.com/retrotxt/lib/convert"
)

func ExampleOpen() {
	b, err := Open("037")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(len(b))
	// Output:130
}

func TestFlags_Open(t *testing.T) {
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
		{"no convert", fields{nil, charmap.CodePage1047}, args{&file, "037"}, false, false},
		{"convert 1252", fields{nil, charmap.CodePage437}, args{&file, "1252"}, false, false},
		{"convert to cp437", fields{charmap.Windows1252, charmap.CodePage437}, args{&file, "1252"}, false, false},
		{"ansi", fields{}, args{&file, "ansi"}, true, false},
		{"cp437 dump", fields{}, args{&file, "437.cr"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Flags{
				From: tt.fields.From,
				To:   tt.fields.To,
			}
			gotS, err := f.Open(tt.args.name, tt.args.conv)
			if (err != nil) != tt.wantErr {
				t.Errorf("Flags.Open() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			r := bool(len(gotS.Runes) > 0)
			if r != tt.wantRunes {
				t.Errorf("Flags.Open() = %v, want %v", r, tt.wantRunes)
			}
		})
	}
}

func TestValid(t *testing.T) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Valid(tt.args.name); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
