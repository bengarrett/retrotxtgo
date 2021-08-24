package logs

import (
	"errors"
	"testing"
)

var (
	ErrLongMsg  = errors.New("this error message has at least three words")
	ErrBF       = errors.New("bad flag test")
	ErrType     = errors.New("invalid argument type")
	ErrBoolType = errors.New("invalid argument test flag strconv.ParseBool")
)

func Test_execute(t *testing.T) {
	type args struct {
		err  error
		args []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{}, ""},
		{"no strings", args{ErrTest, []string{}}, "cmd error args: word count is too short, less than 3"},
		{"default", args{ErrLongMsg, []string{""}}, "Problem:\n this error message has at least three words\n run retrotxt  --help"},
		{"strings", args{ErrBF, []string{"some", "command"}}, "Problem:\n with the some --test flag, bad flag test"},
		{"type", args{ErrBoolType, []string{"some", "command"}},
			"Problem:\n with the some --strconv.ParseBool flag, the value must be either true or false"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := execute(tt.args.err, true, tt.args.args...); got != tt.want {
				t.Errorf("execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
