package create

import "testing"

func TestPort(t *testing.T) {
	tests := []struct {
		name string
		port uint
		want bool
	}{
		{"empty", 0, true},
		{"www", 80, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Port(tt.port); got != tt.want {
				t.Errorf("Port() = %v, want %v", got, tt.want)
			}
		})
	}
}
