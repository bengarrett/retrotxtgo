package config

import "testing"

func Test_hr(t *testing.T) {
	tests := []struct {
		name string
		w    uint
		want string
	}{
		{"empty", 0, ""},
		{"5", 5, "-----"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hr(tt.w); got != tt.want {
				t.Errorf("hr() = %v, want %v", got, tt.want)
			}
		})
	}
}
