package list_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/internal/list"
	"github.com/gookit/color"
)

func ExampleTable() {
	color.Enable = false
	names := []string{"cp437"}
	t, err := list.Table(names...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d characters in the table", len(t))
	// Output: 1690 characters in the table
}

func ExampleTables() {
	color.Enable = false
	t, err := list.Tables()
	if err != nil {
		log.Print(err)
	}
	fmt.Printf("%d characters in the table", len(t))
	// Output: 75117 characters in the table
}

func TestExamples(t *testing.T) {
	t.Run("example", func(t *testing.T) {
		got, err := list.Examples()
		if err != nil {
			t.Error(err)
			return
		}
		if got == nil {
			t.Errorf("examples() failed to return anything")
		}
	})
}

func TestTable(t *testing.T) {
	tests := []struct {
		name     string
		contains string
		wantErr  bool
	}{
		{"", "", true},
		{"ansi", "", true},
		{"ascii", "ANSI X3.4 1967/77/86 - Extended ASCII", false},
		{"CP437", "IBM Code Page 437 (DOS, OEM-US) - Extended ASCII", false},
		{"utf8", "UTF-8 - Unicode", false},
		{"shiftjis", "Shift JIS (Japanese) - Extended ASCII", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, err := list.Table(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("table returned error: %s, wanted %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(table, tt.contains) {
				t.Errorf("table does not contain the header: %s", tt.contains)
			}
		})
	}
}

func TestTables(t *testing.T) {
	tests := []struct {
		contains string
	}{
		//{"ANSI X3.4 1967/77/86 - Extended ASCII"},
		{"IBM Code Page 437 (DOS, OEM-US) - Extended ASCII"},
		{"UTF-8 - Unicode"},
		{"Shift JIS (Japanese) - Extended ASCII"},
	}
	tables, err := list.Tables()
	if err != nil {
		t.Error(err)
	}
	for _, tt := range tests {
		t.Run(tt.contains, func(t *testing.T) {
			if !strings.Contains(tables, tt.contains) {
				t.Errorf("tables does not contain the header: %s", tt.contains)
			}
		})
	}
}

func TestPrintable(t *testing.T) {
	tests := []struct {
		n    string
		name string
		want bool
	}{
		{"empty", "", false},
		{"utf16", "Utf-16le", false},
		{"utf32", "UTF-32Be", false},
		{"cp437", "CP-437", true},
	}
	for _, tt := range tests {
		t.Run(tt.n, func(t *testing.T) {
			if got := list.Printable(tt.name); got != tt.want {
				t.Errorf("Printable() = %v, want %v", got, tt.want)
			}
		})
	}
}
