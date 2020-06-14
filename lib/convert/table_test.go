package convert

import (
	"testing"
)

func TestTable(t *testing.T) {
	tests := []struct {
		name    string
		wantNil bool
		wantErr bool
	}{
		{"IBM437", false, false},
		{"cp437", false, false},
		{"win", false, false},
		{"xxx", true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Table(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Table() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(got != nil) != tt.wantNil {
				t.Errorf("Table() = %v, want %v", got, tt.wantNil)
			}
		})
	}
}
