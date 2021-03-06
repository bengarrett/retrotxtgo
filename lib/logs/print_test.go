package logs

import (
	"errors"
	"testing"

	"github.com/gookit/color"
)

var ErrTest = errors.New("error")

func TestHint_String(t *testing.T) {
	color.Enable = false
	type fields struct {
		Error error
		Hint  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, ""},
		{"text", fields{nil, "hint"}, ""},
		{"text", fields{ErrTest, "hint"}, "problem: error\n         run retrotxt hint"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hint(tt.fields.Hint, tt.fields.Error); got != tt.want {
				t.Errorf("Hint() = %v, want %v", got, tt.want)
			}
		})
	}
}
