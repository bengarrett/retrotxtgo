package input_test

import (
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
)

func Test_previewPrompt(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name  string
		args  args
		wantP string
	}{
		{"empty", args{}, "Set"},
		{"key", args{get.Keywords, "ooooh"}, "Replace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP := input.PreviewPrompt(tt.args.name, tt.args.value)
			firstWord := strings.Split(strings.TrimSpace(gotP), " ")[0]
			if firstWord != tt.wantP {
				t.Errorf("PreviewPrompt() = %v, want %v", firstWord, tt.wantP)
			}
		})
	}
}
