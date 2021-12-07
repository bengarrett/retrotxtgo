package listcmd_test

import (
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/cmd/internal/listcmd"
	"github.com/gookit/color"
)

func ExamplePrintTable() {
	color.Enable = false
	names := []string{"cp437"}
	t := listcmd.PrintTable(names...)
	fmt.Printf("%d characters in the table", len(t))
	// Output: 1689 characters in the table
}

func ExamplePrintTables() {
	color.Enable = false
	t := listcmd.PrintTables()
	fmt.Printf("%d characters in the table", len(t))
	// Output: 75088 characters in the table
}

func TestPrintExamples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		if got := listcmd.PrintExamples(); got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}

func TestPrintable(t *testing.T) {
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
			if got := listcmd.Printable(tt.name); got != tt.want {
				t.Errorf("Printable() = %v, want %v", got, tt.want)
			}
		})
	}
}
