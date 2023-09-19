package update_test

import (
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/update"
	"github.com/bengarrett/retrotxtgo/meta"
)

const (
	etag  = `W/"3715383704fac6f3568e9039b347937a`
	alpha = `0.0.1`
)

func ExampleString() {
	update.Notice(os.Stdout, alpha, "1.0.0")
	// Output:┌───────────────────────────────────────────┐
	// │ A newer edition of Retrotxt is available! │
	// │   Learn more at https://retrotxt.com/go   │
	// │              α0.0.1 → 1.0.0               │
	// └───────────────────────────────────────────┘
}

func TestCacheSet(t *testing.T) {
	t.Run("cache set", func(t *testing.T) {
		if err := update.CacheSet(etag, alpha); err != nil {
			t.Error(err)
		}
	})
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
