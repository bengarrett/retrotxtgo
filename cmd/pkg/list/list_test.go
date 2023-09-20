package list_test

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/cmd/pkg/list"
	"github.com/gookit/color"
)

func ExampleTable() {
	color.Enable = false
	names := []string{"cp437"}
	s := strings.Builder{}
	if err := list.Table(&s, names...); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%d characters in the table", len(s.String()))
	// Output: 1690 characters in the table
}

func ExampleTables() {
	color.Enable = false
	s := strings.Builder{}
	if err := list.Tables(&s); err != nil {
		log.Print(err)
	}
	const val = 70000
	l := s.Len()
	fmt.Fprintf(os.Stdout, "characters > 70000: %v", l > val)
	// Output: characters > 70000: true
}

func TestExamples(t *testing.T) {
	t.Parallel()
	t.Run("example", func(t *testing.T) {
		t.Parallel()
		s := strings.Builder{}
		err := list.Examples(&s)
		if err != nil {
			t.Error(err)
			return
		}
		if s.String() == "" {
			t.Errorf("examples() failed to return anything")
		}
	})
}

func TestTable(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		contains string
		wantErr  bool
	}{
		{"", "", true},
		{"aix", "", true},
		{"ascii", "ANSI X3.4 1967/77/86 - Extended ASCII", false},
		{"CP437", "IBM Code Page 437 (DOS, OEM-US) - Extended ASCII", false},
		{"utf8", "UTF-8 - Unicode", false},
		{"shiftjis", "Shift JIS (Japanese) - Extended ASCII", false},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			s := strings.Builder{}
			err := list.Table(&s, tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("table returned error: %s, wanted %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(s.String(), tt.contains) {
				t.Errorf("table does not contain the header: %s", tt.contains)
			}
		}
	})
}

func TestTables(t *testing.T) {
	t.Parallel()
	tests := []struct {
		contains string
	}{
		// {"ANSI X3.4 1967/77/86 - Extended ASCII"},
		{"IBM Code Page 437 (DOS, OEM-US) - Extended ASCII"},
		{"UTF-8 - Unicode"},
		{"Shift JIS (Japanese) - Extended ASCII"},
	}
	s := &strings.Builder{}
	if err := list.Tables(s); err != nil {
		t.Error(err)
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if !strings.Contains(s.String(), tt.contains) {
				t.Errorf("tables does not contain the header: %s", tt.contains)
			}
		}
	})
}

func TestPrintable(t *testing.T) {
	t.Parallel()
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
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			if got := list.Printable(tt.name); got != tt.want {
				t.Errorf("Printable() = %v, want %v", got, tt.want)
			}
		}
	})
}
