package flag_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/create"
	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want encoding.Encoding
	}{
		{"nil", charmap.CodePage437},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flag.Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initEncodings(t *testing.T) {
	type args struct {
		cmd    *cobra.Command
		dfault string
	}
	tests := []struct {
		name    string
		args    args
		want    sample.Flags
		wantErr bool
	}{
		{"empty", args{}, sample.Flags{}, false},
		{"default", args{nil, "CP437"}, sample.Flags{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flag.EncodeAndHide(tt.args.cmd, tt.args.dfault)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeAndHide() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeAndHide() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSAUCE(t *testing.T) {
	name, err := filepath.Abs("../../../static/text/sauce.txt")
	if err != nil {
		t.Error(err)
		return
	}
	f, err := filesystem.ReadAllBytes(name)
	if err != nil {
		t.Error(err)
		return
	}
	got := flag.SAUCE(&f)
	if reflect.DeepEqual(got, create.SAUCE{}) {
		t.Error("SAUCE result is empty")
		return
	}
	if !got.Use {
		t.Error("SAUCE.Use result is false")
	}
	const wantTitle = "Sauce title"
	if got.Title != wantTitle {
		t.Errorf("SAUCE.Title = %q, want %q", got.Title, wantTitle)
	}
	const wantAuthor = "Sauce author"
	if got.Author != wantAuthor {
		t.Errorf("SAUCE.Author = %q, want %q", got.Title, wantAuthor)
	}
	const wantGroup = "Sauce group"
	if got.Group != wantGroup {
		t.Errorf("SAUCE.Group = %q, want %q", got.Group, wantGroup)
	}
	const wantDesc = "ASCII text file with no formatting codes or color codes."
	if got.Description != wantDesc {
		t.Errorf("SAUCE.Description = %q, want %q", got.Description, wantDesc)
	}
	const wantWidth = 977
	if got.Width != wantWidth {
		t.Errorf("SAUCE.Width = %d, want %d", got.Width, wantWidth)
	}
	const wantLines = 9
	if got.Lines != wantLines {
		t.Errorf("SAUCE.Lines = %d, want %d", got.Lines, wantLines)
	}
}
