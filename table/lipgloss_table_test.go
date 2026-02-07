package table

import (
	"bytes"
	"testing"
)

func TestLipglossTable(t *testing.T) {
	testRows := []Row{
		{Name: "IBM Code Page 437", Value: "cp437", Numeric: "437", Alias: "msdos"},
		{Name: "* IBM Code Page 037", Value: "cp037", Numeric: "37", Alias: "ibm037"},
		{Name: "† Big5", Value: "big5", Numeric: "", Alias: "big-5"},
		{Name: "⁑ ASA X3.4 1963", Value: "ascii-63", Numeric: "1963", Alias: ""},
	}

	buf := new(bytes.Buffer)
	err := LipglossTable(buf, testRows)
	if err != nil {
		t.Fatalf("LipglossTable failed: %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("LipglossTable produced empty output")
	}

	// Check that the output contains expected content
	if !contains(output, "IBM Code Page 437") {
		t.Error("Output should contain 'IBM Code Page 437'")
	}
	if !contains(output, "cp437") {
		t.Error("Output should contain 'cp437'")
	}
}

func TestListWithStyle(t *testing.T) {
	tests := []struct {
		name     string
		style    string
		wantErr  bool
		expectFn func(string) bool
	}{
		{
			name:    "auto style",
			style:   "auto",
			wantErr: false,
			expectFn: func(s string) bool {
				return s != "" // Just check it produces output
			},
		},
		{
			name:    "lipgloss style",
			style:   "lipgloss",
			wantErr: false,
			expectFn: func(s string) bool {
				return contains(s, "┌") || contains(s, "│") // Check for lipgloss borders
			},
		},
		{
			name:    "tabwriter style",
			style:   "tabwriter",
			wantErr: false,
			expectFn: func(s string) bool {
				return contains(s, "Formal name") && contains(s, "Named value")
			},
		},
		{
			name:    "invalid style",
			style:   "invalid",
			wantErr: false, // Should default to one of the valid styles
			expectFn: func(s string) bool {
				return s != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := ListWithStyle(buf, tt.style)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("ListWithStyle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			if !tt.expectFn(output) {
				t.Errorf("ListWithStyle() output did not meet expectations for style %s", tt.style)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && 
		   (len(s) >= len(substr)) && 
		   (s == substr || 
		    len(s) > len(substr) && (s[:len(substr)] == substr || 
		                           s[len(s)-len(substr):] == substr ||
		                           containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}