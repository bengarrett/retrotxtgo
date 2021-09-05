package cmd

import (
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func Example_listTable() {
	color.Enable = false
	names := []string{"cp437"}
	t := listTable(names...)
	fmt.Printf("%d characters in the table", len(t))
	// Output: 1689 characters in the table
}
func Example_listTbls() {
	color.Enable = false
	t := listTbls()
	fmt.Printf("%d characters in the table", len(t))
	// Output: 75088 characters in the table
}

func Test_examples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		if got := examples(); got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}

func Test_usableTbl(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want bool
	}{
		{"empty", "", false},
		{"utf", "UTF-32Be", false},
	}
	for _, tt := range tests {
		t.Run(tt.n, func(t *testing.T) {
			if got := usableTbl(tt.name); got != tt.want {
				t.Errorf("usableTbl() = %v, want %v", got, tt.want)
			}
		})
	}
}
