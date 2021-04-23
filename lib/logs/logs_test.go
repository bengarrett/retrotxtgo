package logs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gookit/color"
)

var (
	ErrC       = errors.New("c")
	ErrRandom  = errors.New("some problem")
	ErrLogTest = errors.New("log test")
)

func ExampleLog() {
	t := fmt.Sprintf("%s", ErrLogTest)
	Log(ErrLogTest)
	last, _ := LastEntry()
	i := len(last) - len(t) - 1
	fmt.Print(last[i:])
	// Output:log test
}

func TestErr_String(t *testing.T) {
	color.Disable()
	tests := []struct {
		name string
		e    Argument
		want string
	}{
		{"empty", Argument{}, ""},
		{"abc", Argument{Issue: "a", Arg: "b", Err: ErrC}, "problem: a b, c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.String(); got != tt.want {
				t.Errorf("Err.String() = %q, want %v", got, tt.want)
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
		{"ok", args{ErrRandom, file}, false},
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
