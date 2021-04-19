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
		Error Generic
		Hint  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "\n         run retrotxt "},
		{"text", fields{Generic{}, "hint"}, "\n         run retrotxt hint"},
		{"text", fields{Generic{"issue", "arg", nil}, "hint"}, "\n         run retrotxt hint"},
		{"text", fields{Generic{"issue", "arg", ErrTest}, "hint"}, "problem: issue arg, error\n         run retrotxt hint"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Hint{
				Error: tt.fields.Error,
				Hint:  tt.fields.Hint,
			}
			if got := h.String(); got != tt.want {
				t.Errorf("Hint.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
