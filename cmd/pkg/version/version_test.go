package version_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/version"
	"github.com/gookit/color"
)

func TestTemplate(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"©", "©"},
		{"© date", version.Copyright},
		{"path", "path:"},
	}
	color.Enable = false
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := version.Template()
			if !strings.Contains(got, tt.want) {
				t.Errorf("Template() does not contain %v", tt.want)
				t.Error(got)
				return
			}
		})
	}
}