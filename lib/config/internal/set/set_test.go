package set_test

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/set"
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := set.DirExpansion(tt.name); gotDir != tt.wantDir {
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
		want    bool
		wantErr bool
	}{
		{"nil", args{}, false, true},
		{"title err type", args{name: "html.title", value: 0}, false, true},
		{"title err", args{name: "title", value: "0"}, false, true},
		{"serve int", args{"serve", uint(8080)}, false, true},
		{"serve int skip", args{"serve", uint(skipValue)}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if v := viper.AllKeys(); len(v) == 0 {
				fmt.Println("init serve example.")
				viper.Set("serve", skipValue)
			}
			got, err := set.SkipWrite(tt.args.name, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SkipWrite() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SkipWrite() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecommend(t *testing.T) {
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
