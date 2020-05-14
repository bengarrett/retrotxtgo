package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func Test_versionPrint(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"empty", "", false},
		{"invalid", "abcde", true},
		{"j", "j", false},
		{"jm", "jm", false},
		{"t", "t", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotErr := versionPrint(tt.format); (gotErr.Msg != nil) != tt.wantErr {
				t.Errorf("versionPrint() = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func Test_arch(t *testing.T) {
	tests := []struct {
		name string
		v    string
		want string
	}{
		{"empty", "", ""},
		{"invalid", "xxx", ""},
		{"386", "386", "32-bit Intel/AMD"},
		{"ppc64", "ppc64", "64-bit PowerPC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arch(tt.v); got != tt.want {
				t.Errorf("arch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_binary(t *testing.T) {
	tests := []struct {
		name     string
		dontWant string
	}{
		{"ok", "error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binary(); got[:5] == tt.dontWant {
				t.Errorf("binary() = %v, don't want %v", got, tt.dontWant)
			}
		})
	}
}

func Test_info(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"os and arch", fmt.Sprintf("%s/%s [%s CPU]", runtime.GOOS, runtime.GOARCH, arch(runtime.GOARCH))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := info()["os"]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_locBuildDate(t *testing.T) {
	d := time.Date(1980, 1, 31, 1, 34, 0, 0, time.UTC)
	tests := []struct {
		name string
		date string
		want string
	}{
		{"empty", "", ""},
		{"invalid", "abcde", "abcd"},
		{"ok", d.UTC().Format(time.RFC3339), "1980"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := locBuildDate(tt.date); len(got) > 4 && got[:4] != tt.want {
				t.Errorf("locBuildDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_goVer(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"ok", "1."}, // 1.14.1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := goVer(); len(got) > 2 && got[:2] != tt.want {
				t.Errorf("goVer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionJSON(t *testing.T) {
	type args struct {
		indent bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"no indent", args{false}, true},
		{"indent", args{true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json.Valid(versionJSON(tt.args.indent)); got != tt.want {
				t.Errorf("versionJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ExampleversionText() {
	fmt.Print(versionText(false)[:8])
	// Output: RetroTxt
}

// func capVersionText(c bool) (output string) {
// 	rescueStdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	color.Enable = true
// 	versionText(c)
// 	w.Close()
// 	bytes, _ := ioutil.ReadAll(r)
// 	os.Stdout = rescueStdout
// 	return strings.TrimSpace(string(bytes))
// }
