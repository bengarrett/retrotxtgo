package version

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func ExampleJSON() {
	fmt.Printf("%s", JSON(true))
	fmt.Printf("%s", JSON(false))
}

func Test_JSON(t *testing.T) {
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
			if got := json.Valid(JSON(tt.args.indent)); got != tt.want {
				t.Errorf("JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Print(t *testing.T) {
	tests := []struct {
		name   string
		format string
		wantOk bool
	}{
		{"empty", "", true},
		{"invalid", "abcde", false},
		{"j", "j", true},
		{"jm", "jm", true},
		{"t", "t", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := Print(tt.format); gotOk != tt.wantOk {
				t.Errorf("Print() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func ExampleSprint() {
	fmt.Print(Sprint(false))
}

func TestSprint(t *testing.T) {
	tests := []struct {
		name  string
		color bool
		empty bool
	}{
		{"text", false, false},
		{"color", true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotText := Sprint(tt.color); (len(gotText) == 0) != tt.empty {
				t.Errorf("Sprint() = %v, want %v", gotText, tt.empty)
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
			if got := information()["os"]; !reflect.DeepEqual(got, tt.want) {
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
			if got := localBuild(tt.date); len(got) > 4 && got[:4] != tt.want {
				t.Errorf("locBuildDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_semantic(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"ok", "1"}, // 1.14.1
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := semantic(); len(got) > 2 && strings.Split(got, ".")[0] != tt.want {
				t.Errorf("semantic() = %v, want %v", got, tt.want)
			}
		})
	}
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
