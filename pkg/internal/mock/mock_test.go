package mock_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/internal/mock"
)

func Test_filler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		sizeMB     float64
		wantLength int
	}{
		{"0", 0, 0},
		{"0.1", 0.1, 100000},
		{"1.5", 1.5, 1500000},
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()
		for _, tt := range tests {
			s := mock.Filler(tt.sizeMB)
			if gotLength := len(s); gotLength != tt.wantLength {
				t.Errorf("Filler() = %v, want %v", gotLength, tt.wantLength)
			}
		}
	})
}
