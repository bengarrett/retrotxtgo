package cmd

import (
	"testing"
)

func Test_codepages(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"ok", "\nSupported legacy codepages and encodings"},
	}
	const l = 41
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := codepages(); len(got) > l && got[:l] != tt.want {
				t.Errorf("codepages() = %q, want %q", got[:l], tt.want)
			}
		})
	}
}
