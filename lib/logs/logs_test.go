package logs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	//"github.com/bengarrett/retrotxtgo/lib/filesystem"

	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

func TestErr_String(t *testing.T) {
	color.Disable()
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
		// this will cause the test to exit to the os
		//{"err", args{"some issue", errors.New("some error")}, false},
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
	// set test mode for str.HighlightWriter()
	str.TestMode = true
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
		{"basic", "<h1>hello</h1>", "\n<\x1b[1mh1\x1b[0m>hello</\x1b[1mh1\x1b[0m>\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := colorhtml(tt.elm, "bw"); got != tt.want {
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
	if err := os.RemoveAll(file); err != nil {
		fmt.Fprintln(os.Stderr, "removing path:", err)
	}
}

func TestPath(t *testing.T) {
	if got := Path(); !filepath.IsAbs(got) {
		t.Errorf("Path() is empty or not an absolute path, it should return a directory")
	}
}

func Test_unknownFlag(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{"empty", "", false},
		{"random", "the quick brown fox", false},
		{"ok long", "unknown shorthand flag --test", true},
		{"ok short", "unknown flag --test", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// if got := unknownFlag(tt.s); got != tt.want {
			// 	t.Errorf("unknownFlag() = %v, want %v", got, tt.want)
			// }
		})
	}
}
