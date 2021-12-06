package prompt_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/internal/mock"
	"github.com/bengarrett/retrotxtgo/lib/prompt"
	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleYesNo() {
	yn := prompt.YesNo("Say hello", true)
	fmt.Print(yn)
	// Output:Say hello? [Yes/no] true
}

func TestSkipSet(t *testing.T) {
	if ss := prompt.SkipSet(false); ss != "" {
		t.Errorf("SkipSet(true) = %s, want empty", ss)
	}
	const want = "  skipped setting"
	if ss := prompt.SkipSet(true); ss != want {
		t.Errorf("SkipSet(true) = %s, want %s", ss, want)
	}
}

//nolint:dupl
func TestPort(t *testing.T) {
	type args struct {
		validate bool
		setup    bool
		input    string
	}
	tests := []struct {
		name     string
		args     args
		wantPort uint
	}{
		{"empty", args{}, 0},
		{"bad input", args{true, true, "abc"}, 0},
		{"valid", args{true, false, "80"}, 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotPort := prompt.Port(tt.args.validate, tt.args.setup); gotPort != tt.wantPort {
				t.Errorf("Port() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

//nolint:dupl
func TestIndexStrings(t *testing.T) {
	options := []string{"number1", "number2", "number3"}
	type args struct {
		options *[]string
		setup   bool
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, false, "xyz"}, ""},
		{"ok input", args{&options, false, "number2"}, "number2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := prompt.IndexStrings(tt.args.options, tt.args.setup); gotKey != tt.wantKey {
				t.Errorf("IndexStrings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestShortStrings(t *testing.T) {
	options := []string{"alpha", "beta", "gamma"}
	type args struct {
		options *[]string
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, "xyz"}, ""},
		{"ok input", args{&options, "b"}, "beta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := prompt.ShortStrings(tt.args.options); gotKey != tt.wantKey {
				t.Errorf("ShortStrings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"mixed input", "123 abc", "123 abc"},
		{"ok input", "hello world", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := prompt.String(); gotKey != tt.want {
				t.Errorf("String() = %v, want %v", gotKey, tt.want)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	options := []string{"alpha", "beta", "gamma"}
	type args struct {
		options *[]string
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, "xyz"}, ""},
		{"ok input", args{&options, "beta"}, "beta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mock.Input(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := prompt.Strings(tt.args.options, false); gotKey != tt.wantKey {
				t.Errorf("Strings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
