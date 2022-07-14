package example_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/example"
	"github.com/gookit/color"
)

func TestExample(t *testing.T) {
	color.Enable = false
	tests := []struct {
		contains string
	}{
		{"# save the text files as webpages"},
		{"config setup"},
		{"config info"},
		{"# print a HTML file created from file.txt"},
		{"list codepages"},
		{"list examples"},
		{"table cp437"},
		{"info text.asc logo.jpg"},
		{"config set --list"},
		{"view file.txt"},
		{""},
	}
	val := -1
	for _, tt := range tests {
		val++
		t.Run(fmt.Sprintf("example_%d", val), func(t *testing.T) {
			if !strings.Contains(example.Example(val).String(), tt.contains) {
				t.Errorf("example %v does not contain the expected string: %q", val, tt.contains)
				fmt.Printf("%s", example.Example(val))
			}
		})
	}
}
