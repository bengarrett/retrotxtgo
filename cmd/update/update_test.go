package update_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/update"
)

const (
	etag  = `W/"3715383704fac6f3568e9039b347937a`
	alpha = `0.0.1`
)

func TestCacheSet(t *testing.T) {
	t.Parallel()
	t.Run("cache set", func(t *testing.T) {
		t.Parallel()
		if err := update.CacheSet(etag, alpha); err != nil {
			t.Error(err)
		}
	})
}

func TestCompare(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := update.Compare(tt.args.current, tt.args.fetched); got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		}
	})
}
