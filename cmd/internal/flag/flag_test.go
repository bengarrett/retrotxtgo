package flag_test

import (
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/flag"
	"github.com/bengarrett/retrotxtgo/lib/sample"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func Test_configInfo(t *testing.T) {
	tests := []struct {
		name     string
		wantExit bool
	}{
		{"output", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExit := flag.ConfigInfo(); gotExit != tt.wantExit {
				t.Errorf("ConfigInfo() = %v, want %v", gotExit, tt.wantExit)
			}
		})
	}
}

func Test_dfaultInput(t *testing.T) {
	tests := []struct {
		name string
		want encoding.Encoding
	}{
		{"nil", charmap.CodePage437},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flag.DfaultInput(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dfaultInput() = %v, want %v", got, tt.want)
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
			got, err := flag.InitEncodings(tt.args.cmd, tt.args.dfault)
			if (err != nil) != tt.wantErr {
				t.Errorf("initEncodings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initEncodings() = %v, want %v", got, tt.want)
			}
		})
	}
}
