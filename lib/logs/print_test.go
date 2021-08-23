package logs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

var ErrTest = errors.New("error")

func init() {
	color.Enable = false
}

func TestHint_String(t *testing.T) {
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
		{"text", fields{ErrTest, "hint"}, fmt.Sprintf("Problem:\n error\n run %s hint", meta.Bin)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Hint(tt.fields.Hint, tt.fields.Error); got != tt.want {
				t.Errorf("Hint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"empty", nil, ""},
		{"test", ErrTest, "Problem:\n error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Errorf(tt.err); got != tt.want {
				t.Errorf("Errorf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemCmd(t *testing.T) {
	type args struct {
		name string
		err  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{"", nil}, ""},
		{"no name", args{"", ErrTest}, ""},
		{"error", args{"test", ErrTest}, "Problem:\n the command test does not exist, error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProblemCmd(tt.args.name, tt.args.err); got != tt.want {
				t.Errorf("ProblemCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemFlag(t *testing.T) {
	type args struct {
		name string
		flag string
		err  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"no name", args{"", "", ErrTest}, ""},
		{"error", args{"error", "err", ErrTest}, "Problem:\n with the error --err flag, error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProblemFlag(tt.args.name, tt.args.flag, tt.args.err); got != tt.want {
				t.Errorf("ProblemFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemMark(t *testing.T) {
	type args struct {
		value string
		err   error
		errs  error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"no name", args{"", ErrTest, ErrTest}, ""},
		{"errors", args{"error", ErrTest, ErrTest}, "Problem:\n error \"error\": error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProblemMark(tt.args.value, tt.args.err, tt.args.errs); got != tt.want {
				t.Errorf("ProblemMark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProblemf(t *testing.T) {
	type args struct {
		err  error
		errs error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"partial", args{ErrTest, nil}, ""},
		{"errors", args{ErrTest, ErrTest}, "Problem:\nerror: error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Problemf(tt.args.err, tt.args.errs); got != tt.want {
				t.Errorf("Problemf() = %v, want %v", got, tt.want)
			}
		})
	}
}
