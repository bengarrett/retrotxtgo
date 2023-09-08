package mock_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
)

func Test_filler(t *testing.T) {
	tests := []struct {
		name       string
		sizeMB     float64
		wantLength int
	}{
		{"0", 0, 0},
		{"0.1", 0.1, 100000},
		{"1.5", 1.5, 1500000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLength, _ := mock.Filler(tt.sizeMB); gotLength != tt.wantLength {
				t.Errorf("Filler() = %v, want %v", gotLength, tt.wantLength)
			}
		})
	}
}
