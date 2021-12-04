package logs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

const status = "exit status 1"

var ErrTest = errors.New("error")

func init() {
	color.Enable = false
}

func ExampleSprint() {
	fmt.Println(logs.Sprint(ErrTest))
	fmt.Println(status)
	//Output: Problem:
	// error.
	// exit status 1
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
		{"text", fields{ErrTest, "hint"}, fmt.Sprintf("Problem:\nerror.\n run %s hint", meta.Bin)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logs.Hint(tt.fields.Hint, tt.fields.Error); got != tt.want {
				t.Errorf("Hint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"empty", nil, ""},
		{"test", ErrTest, "Problem:\nerror."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logs.Sprint(tt.err); got != tt.want {
				t.Errorf("Sprint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintCmd(t *testing.T) {
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
			if got := logs.SprintCmd(tt.args.name, tt.args.err); got != tt.want {
				t.Errorf("SprintCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintFlag(t *testing.T) {
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
			if got := logs.SprintFlag(tt.args.name, tt.args.flag, tt.args.err); got != tt.want {
				t.Errorf("SprintFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintMark(t *testing.T) {
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
			if got := logs.SprintMark(tt.args.value, tt.args.err, tt.args.errs); got != tt.want {
				t.Errorf("SprintMark() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintWrap(t *testing.T) {
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
			if got := logs.SprintWrap(tt.args.err, tt.args.errs); got != tt.want {
				t.Errorf("SprintWrap() = %v, want %v", got, tt.want)
			}
		})
	}
}
