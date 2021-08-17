package cmd

import (
	"fmt"
	"testing"
)

const etag = `W/"3715383704fac6f3568e9039b347937a`

func ExamplecacheGet() {
	if err := cacheSet(etag, "0.0.1"); err != nil {
		fmt.Println(err)
	}
	e, v := cacheGet()
	fmt.Println("Etag", e)
	fmt.Println("Version", v)
	// Output: Etag W/"3715383704fac6f3568e9039b347937a
	// Version 0.0.1
}

func ExamplecacheSet() {
	if err := cacheSet(etag, "0.0.1"); err != nil {
		fmt.Println(err)
	}
	e, v := cacheGet()
	fmt.Println("Etag", e)
	fmt.Println("Version", v)
	// Output: Etag W/"3715383704fac6f3568e9039b347937a
	// Version 0.0.1
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
