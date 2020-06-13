package version

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

const etag = `W/"3715383704fac6f3568e9039b347937a`

func ExampleCacheGet() {
	if err := CacheSet(etag, "0.0.1"); err != nil {
		fmt.Println(err)
	}
	e, v := CacheGet()
	fmt.Println("Etag", e)
	fmt.Println("Version", v)
	// Output: Etag W/"3715383704fac6f3568e9039b347937a
	// Version 0.0.1
}

func ExampleCacheSet() {
	if err := CacheSet(etag, "0.0.1"); err != nil {
		fmt.Println(err)
	}
	e, v := CacheGet()
	fmt.Println("Etag", e)
	fmt.Println("Version", v)
	// Output: Etag W/"3715383704fac6f3568e9039b347937a
	// Version 0.0.1
}

func ExampleDigits() {
	fmt.Println(Digits("v1.0 (init release)"))
	// Output: 1.0
}

func ExampleJSON() {
	fmt.Print(json.Valid(JSON(true)), json.Valid(JSON(false)))
	// Output: true true
}

func ExampleSprint() {
	fmt.Print(Sprint(false)[:8])
	// Output: RetroTxt
}

func TestDigits(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{"empty", "", ""},
		{"digits", "01234567890", "01234567890"},
		{"symbols", "~!@#$%^&*()_+", ""},
		{"mixed", "A0B1C2D3E4F5G6H7I8J9K0L", "01234567890"},
		{"semantic", "v1.0.0 (FINAL)", "1.0.0"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Digits(tt.s); got != tt.want {
				t.Errorf("Digits() = %v, want %v", got, tt.want)
			}
		})
	}
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

func TestSemantic(t *testing.T) {
	tests := []struct {
		name        string
		ver         string
		wantOk      bool
		wantVersion Version
	}{
		{"empty", "", false, Version{0, 0, 0}},
		{"text", "hello world", false, Version{0, 0, 0}},
		{"zero", "0.0.0", true, Version{0, 0, 0}},
		{"vzero", "v0.0.0", true, Version{0, 0, 0}},
		{"ver str", "v1.2.3 (super-release)", true, Version{1, 2, 3}},
		{"short", "V1", true, Version{1, 0, 0}},
		{"short.", "V1.1", true, Version{1, 1, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotVersion := Semantic(tt.ver)
			if gotOk != tt.wantOk {
				t.Errorf("Semantic() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotVersion, tt.wantVersion) {
				t.Errorf("Semantic() gotVersion = %v, want %v", gotVersion, tt.wantVersion)
			}
		})
	}
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

func Test_compare(t *testing.T) {
	type args struct {
		current string
		fetched string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{"", ""}, false},
		{"v1", args{"v1", ""}, false},
		{"v1.0", args{"v1.0", ""}, false},
		{"v1.0.0", args{"v1.0.0", ""}, false},
		{"v1.0.0", args{"v1.0.0", "v1.0.0"}, false},
		{"v1.0.0", args{"v1.0.1", "v1.0.0"}, false},
		{"v1.0.1", args{"v1.0.0", "v1.0.1"}, true},
		{"v1.1.1", args{"v1.0.0", "v1.1.1"}, true},
		{"v2.0.1", args{"v1.0.0", "v2.0.1"}, true},
		{"v2", args{"v1", "v2.0.1"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compare(tt.args.current, tt.args.fetched); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
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

func Test_locBuild(t *testing.T) {
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
				t.Errorf("localBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
