package version_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/version"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

func TestTemplate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want string
	}{
		{"©", "©"},
		{"© date", meta.Copyright},
		{"path", "path:"},
	}
	color.Enable = false
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			s := strings.Builder{}
			_ = version.Template(&s)
			if !strings.Contains(s.String(), tt.want) {
				t.Errorf("Template() does not contain %v", tt.want)
				t.Error(s.String())
				return
			}
		}
	})
}
