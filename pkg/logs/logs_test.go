package logs_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/bengarrett/retrotxtgo/pkg/logs"
	"github.com/gookit/color"
)

var ErrTest = errors.New("something went wrong")

func ExampleHint() {
	err := errors.New("oops")
	fmt.Print(logs.Hint(err, "helpme"))
	// Output: Problem:
	// oops.
	//  run retrotxt helpme
}

func ExampleSprint() {
	err := errors.New("oops")
	fmt.Print(logs.Sprint(err))
	// Output: Problem:
	// oops.
}

func ExampleSprintCmd() {
	err := errors.New("oops")
	fmt.Print(logs.SprintCmd(err, "helpme"))
	// Output: Problem:
	//  the command helpme does not exist, oops
}

func ExampleSprintFlag() {
	err := errors.New("oops")
	fmt.Print(logs.SprintFlag(err, "error", "err"))
	// Output: Problem:
	//  with the error --err flag, oops
}

func ExampleSprintS() {
	err := errors.New("oops")
	wrap := errors.New("uh-oh")
	fmt.Print(logs.SprintS(err, wrap, "we have some errors"))
	// Output: Problem:
	//  oops "we have some errors": uh-oh
}

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
		{"text", fields{ErrTest, "hint"}, fmt.Sprintf("Problem:\nerror.\n run %s hint", meta.Bin)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := logs.Hint(tt.fields.Error, tt.fields.Hint); got != tt.want {
				t.Errorf("Hint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprint(t *testing.T) {
	color.Enable = false
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
	color.Enable = false
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
			if got := logs.SprintCmd(tt.args.err, tt.args.name); got != tt.want {
				t.Errorf("SprintCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintFlag(t *testing.T) {
	color.Enable = false
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
			if got := logs.SprintFlag(tt.args.err, tt.args.name, tt.args.flag); got != tt.want {
				t.Errorf("SprintFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSprintS(t *testing.T) {
	color.Enable = false
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
			if got := logs.SprintS(tt.args.err, tt.args.errs, tt.args.value); got != tt.want {
				t.Errorf("SprintS() = %v, want %v", got, tt.want)
			}
		})
	}
}
