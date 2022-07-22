package set_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

func TestDirExpansion(t *testing.T) {
	u, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"", ""},
		{"~", u},
		{".", w},
		{"..", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDir := set.DirExpansion(tt.name)
			if tt.name == ".." {
				if gotDir == "" {
					t.Errorf("DirExpansion() using the parent argument is empty")
				}
				if _, err := os.Stat(gotDir); err != nil {
					if errors.Is(err, os.ErrNotExist) {
						t.Errorf("DirExpansion() using the parent argument is not found: %s", gotDir)
					}
				}
				return
			}
			if gotDir != tt.wantDir {
				t.Errorf("DirExpansion() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func TestKeys(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"0 index", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := set.Keys()[0]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSkipWrite(t *testing.T) {
	const skipValue = 6060
	type args struct {
		name  string
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"nil", args{}, true},
		{"title err type", args{name: "html.title", value: 0}, true},
		{"title err", args{name: "title", value: "0"}, true},
		{"serve int", args{"serve", int(8080)}, false},
		{"serve uint", args{"serve", uint(8080)}, false},
		{"serve int skip", args{"serve", uint(skipValue)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := viper.AllKeys(); len(v) == 0 {
				fmt.Println("init serve example.")
				viper.Set("serve", skipValue)
			}
			err := set.SkipWrite(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SkipWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestRecommend(t *testing.T) {
	color.Enable = false
	value, want := "some command", " (suggestion: some command)"
	if got := set.Recommend(value); got != want {
		t.Errorf("Recommand() = %q, want %q", got, want)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		wantOk bool
	}{
		{"empty", "", false},
		{"editor", "editor", true},
		{"rt", "html.meta.retrotxt", true},
		{"typo", "html.meta.retrotx", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := set.Validate(tt.key); gotOk != tt.wantOk {
				t.Errorf("Validate() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestWrite(t *testing.T) {
	const s86 = "serve is set to \"8086\""
	color.Enable = false
	type args struct {
		name  string
		setup bool
		value interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"empty", args{}, "", true},
		{"nil int", args{name: "serve"}, "", true},
		{"nil string", args{name: "html.meta.color_scheme"}, "", true},
		{"uint", args{name: "serve", value: 8086}, s86, false},
		{"int", args{name: "serve", value: int(8086)}, s86, false},
		{"string", args{name: "html.meta.theme_color", value: "blue"},
			"html.meta.theme_color is set to \"blue\"", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.LoadTester(os.Stdout); err != nil {
				t.Error(err)
				return
			}
			w := &bytes.Buffer{}
			if err := set.Write(w, tt.args.name, tt.args.setup, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if err != nil {
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("Write() = does not contain %s", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestDirectory(t *testing.T) {
	color.Enable = false
	tmp := os.TempDir()
	type args struct {
		name  string
		setup bool
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{"nil", args{}, "", true},
		{"empty name", args{name: ""}, "", true},
		{"invalid name", args{name: "zxcvbnmasdfrj", setup: false}, "", true},
		{"temp", args{name: tmp, setup: false}, "skipped", false},
		{"rm", args{name: "-", setup: false}, "skipped", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cmd.LoadTester(os.Stdout); err != nil {
				t.Error(err)
				return
			}
			w := &bytes.Buffer{}
			err := set.Directory(w, tt.args.name, tt.args.setup)
			if (err != nil) != tt.wantErr && !errors.Is(err, set.ErrBreak) {
				t.Errorf("Directory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("Directory() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}
