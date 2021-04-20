package cmd

import (
	"testing"
)

func Test_examples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		if got := examples(); got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}

func Test_skipTable(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want bool
	}{
		{"empty", "", false},
		{"utf", "UTF-32Be", true},
	}
	for _, tt := range tests {
		t.Run(tt.n, func(t *testing.T) {
			if got := skipTable(tt.name); got != tt.want {
				t.Errorf("skipTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
