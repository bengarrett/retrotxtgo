package update_test

import (
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/update"
	"github.com/bengarrett/retrotxtgo/meta"
)

const (
	etag  = `W/"3715383704fac6f3568e9039b347937a`
	alpha = `0.0.1`
)

func ExampleString() {
	s := update.String(alpha, "1.0.0")
	fmt.Println(s)
	// Output:┌─────────────────────────────────────────────┐
	// │ A newer edition of RetroTxtGo is available! │
	// │    Learn more at https://retrotxt.com/go    │
	// │               α0.0.1 → 1.0.0                │
	// └─────────────────────────────────────────────┘
}

func ExampleCacheSet() {
	if err := update.CacheSet(etag, alpha); err != nil {
		fmt.Println(err)
	}
	e, v := update.CacheGet()
	fmt.Println("Etag", e)
	fmt.Println("Version", v)
	// Output: Etag W/"3715383704fac6f3568e9039b347937a
	// Version 0.0.1
}

func TestCheck(t *testing.T) {
	t.Run("check as sourcecode", func(t *testing.T) {
		meta.App.Version = alpha
		s, err := update.Check()
		if err != nil {
			t.Error(err)
		}
		if s == "" {
			t.Error("expected a tag from github")
		}
	})
}

func TestCompare(t *testing.T) {
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
			if got := update.Compare(tt.args.current, tt.args.fetched); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}
