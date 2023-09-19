package flag_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/pkg/convert"
	"github.com/bengarrett/retrotxtgo/pkg/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func ExampleEndOfFile() {
	var f convert.Flag
	f.Controls = []string{"eof"}
	fmt.Fprint(os.Stdout, flag.EndOfFile(f))
	// Output: true
}

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

func TestInputOriginal(t *testing.T) {
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
			got, err := flag.InputOriginal(tt.args.cmd, tt.args.dfault)
			if (err != nil) != tt.wantErr {
				t.Errorf("InputOriginal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InputOriginal() = %v, want %v", got, tt.want)
			}
		})
	}
}
