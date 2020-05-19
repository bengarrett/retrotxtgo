package logs

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/bengarrett/retrotxtgo/samples"
)

func TestErr_String(t *testing.T) {
	tests := []struct {
		name string
		e    Err
		want string
	}{
		{"empty", Err{}, ""},
		{"abc", Err{Issue: "A", Arg: "B", Msg: errors.New("C")}, "problem: A B, C"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Err.String() = %q, want %v", got, tt.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	type args struct {
		issue string
		err   error
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		{"empty", args{"", nil}, true},
		{"ok", args{"some issue", nil}, true},
		{"err", args{"some issue", errors.New("some error")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := Check(tt.args.issue, tt.args.err); gotOk != tt.wantOk {
				t.Errorf("Check() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_check(t *testing.T) {
	type args struct {
		issue string
		err   error
	}
	tests := []struct {
		name     string
		args     args
		wantMsg  string
		wantCode int
	}{
		{"empty", args{}, "\n", 1},
		{"valid", args{issue: "abc", err: errors.New("some error")}, "abc some error\n", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotCode := check(tt.args.issue, tt.args.err)
			if gotMsg != tt.wantMsg {
				t.Errorf("check() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotCode != tt.wantCode {
				t.Errorf("check() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_checkArgument(t *testing.T) {
	type args struct {
		arg  string
		args []string
	}
	tests := []struct {
		name     string
		args     args
		wantMsg  string
		wantCode int
	}{
		{"empty", args{}, "problem: invalid argument \"\"\n", 10},
		{"ok", args{"???", []string{"json", "text", "xml"}},
			"problem: invalid argument \"???\" choices: json, text, xml\nplease use one of the argument choices shown above\n", 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotCode := checkArgument(tt.args.arg, tt.args.args)
			if gotMsg != tt.wantMsg {
				t.Errorf("checkArgument() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotCode != tt.wantCode {
				t.Errorf("checkArgument() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func TestErr_check(t *testing.T) {
	var e Err
	tests := []struct {
		name     string
		e        Err
		wantMsg  string
		wantCode int
	}{
		{"empty", e, "", 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMsg, gotCode := tt.e.check()
			if gotMsg != tt.wantMsg {
				t.Errorf("Err.check() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
			if gotCode != tt.wantCode {
				t.Errorf("Err.check() gotCode = %v, want %v", gotCode, tt.wantCode)
			}
		})
	}
}

func Test_colorhtml(t *testing.T) {
	type args struct {
		elm string
	}
	tests := []struct {
		name string
		elm  string
		want string
	}{
		{"empty", "", ""},
		{"str", "hello", "\nhello\n"},
		{"basic", "<h1>hello</h1>", "\n[38;5;102m<[0m[38;5;25mh1[0m[38;5;102m>[0mhello[38;5;102m</[0m[38;5;25mh1[0m[38;5;102m>[0m\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorhtml(&tt.elm); got != tt.want {
				t.Errorf("colorhtml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func testSave() string {
	return filepath.Join(os.TempDir(), "rt_log_savetest")
}

func Test_save(t *testing.T) {
	file := testSave()
	type args struct {
		err  error
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"empty+file", args{nil, file}, true},
		{"ok", args{errors.New("some problem"), file}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := save(tt.args.err, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	samples.Clean(file)
}

func TestFilepath(t *testing.T) {
	if got := Filepath(); !filepath.IsAbs(got) {
		t.Errorf("Filepath() is empty or not an absolute path, it should return a directory")
	}
}
