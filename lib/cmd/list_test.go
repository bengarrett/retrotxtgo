package cmd

import (
	"fmt"
	"testing"

	"github.com/gookit/color"
)

func Example_listTable() {
	color.Enable = false
	args := []string{"cp437"}
	t := listTable(args)
	fmt.Printf("%d characters in the table", len(t))
	// Output: 1690 characters in the table
}
func Example_listTables() {
	color.Enable = false
	t := listAllTables()
	fmt.Printf("%d characters in the table", len(t))
	// Output: 74975 characters in the table
}

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
