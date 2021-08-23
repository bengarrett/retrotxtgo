package prompt

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/gookit/color"
)

func init() {
	color.Enable = false
}

func ExampleYesNo() {
	yn := YesNo("Say hello", true)
	fmt.Print(yn)
	// Output:Say hello? [Yes/no] true
}

func Test_pport(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name     string
		validate bool
		input    string
		wantPort uint
	}{
		{"empty", false, "", 0},
		{"empty validate", true, "", 0},
		{"no validation", false, "8000", 8000},
		{"validation", true, "8000", 8000},
		{"invalid 1", false, "90000", 90000},
		{"invalid 2", false, "abcde", 0},
		{"invalid 3", true, "90000", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotPort := pport(&stdin, tt.validate, false); gotPort != tt.wantPort {
				t.Errorf("pport() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestPortValid(t *testing.T) {
	port := strconv.Itoa(int(meta.WebPort))
	tests := []struct {
		name   string
		port   uint
		wantOk bool
	}{
		{port, meta.WebPort, true},
		{"0", 0, false},
		{"80", 80, true},
		{"8888", 8888, true},
		{"88888", 88888, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := PortValid(tt.port); gotOk != tt.wantOk {
				t.Errorf("PortValid() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_keys_validate(t *testing.T) {
	var opt keys = []string{"hello", "world", "hi"}
	tests := []struct {
		name   string
		k      keys
		key    string
		wantOk bool
	}{
		{"empty", keys{}, "", false},
		{"no input", opt, "", false},
		{"wrong", opt, "someinput", false},
		{"ok 1", opt, "hi", true},
		{"ok 2", opt, "world", true},
		{"only one choice", opt, "hello world", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOk := tt.k.validate(tt.key); gotOk != tt.wantOk {
				t.Errorf("keys.validate() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_keys_prompt(t *testing.T) {
	var opt keys = []string{"hello", "world", "hi"}
	var stdin bytes.Buffer
	tests := []struct {
		name    string
		k       keys
		input   string
		wantKey string
	}{
		{"empty", keys{}, "", ""},
		{"no input", opt, "", ""},
		{"wrong", opt, "blahblah", ""},
		{"ok 1", opt, "hi", "hi"},
		{"ok 2", opt, "hello", "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			if gotKey := tt.k.prompt(&stdin, false); gotKey != tt.wantKey {
				t.Errorf("keys.prompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func Test_promptRead(t *testing.T) {
	var stdin bytes.Buffer
	tests := []struct {
		name      string
		input     string
		wantInput string
		wantErr   bool
	}{
		{"empty", "", "", false},
		{"nl", "\n", "", false},
		{"tab", "\t", "", false},
		{"hi", "hello", "hello", false},
		{"hw", "hello world", "hello world", false},
		{"emoji", "hi 😃", "hi 😃", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdin.Write([]byte(tt.input + "\n")) // \n is a requirement
			gotInput, err := promptRead(&stdin)
			if (err != nil) != tt.wantErr {
				t.Errorf("promptRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotInput != tt.wantInput {
				t.Errorf("promptRead() = %v, want %v", gotInput, tt.wantInput)
			}
		})
	}
}

func Test_parseYN(t *testing.T) {
	type args struct {
		input      string
		yesDefault bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"yes 1", args{"", true}, true},
		{"yes 1", args{"y", false}, true},
		{"yes 2", args{"yes", false}, true},
		{"no 1", args{"no", false}, false},
		{"no 2", args{"no", true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseYN(tt.args.input, tt.args.yesDefault); got != tt.want {
				t.Errorf("parseYN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_keys_shortValidate(t *testing.T) {
	fruits := []string{"Apple", "Watermelon", "Orange", "Pear", "Cherry"}
	tests := []struct {
		name string
		k    keys // should be unique
		key  string
		want string
	}{
		{"empty", keys{""}, "", ""},
		{"match l1", keys{"hi"}, "hi", "hi"},
		{"match s1", keys{"hi"}, "h", "hi"},
		{"match fullname fruit", fruits, "Orange", "Orange"},
		{"match letter of fruit", fruits, "o", "Orange"},
		{"mismatch letter of fruit", fruits, "z", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.k.shortValidate(tt.key); got != tt.want {
				t.Errorf("keys.shortValidate(%s) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func Test_pstring(t *testing.T) {
	type args struct {
		r io.Reader
	}
	a := strings.NewReader("my request")
	b := strings.NewReader("-")
	c := strings.NewReader("\x0D")
	tests := []struct {
		name      string
		args      args
		wantWords string
	}{
		{"input", args{a}, "my request"},
		{"remove", args{b}, "-"},
		{"enter", args{c}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWords := pstring(tt.args.r); gotWords != tt.wantWords {
				t.Errorf("pstring() = %v, want %v", gotWords, tt.wantWords)
			}
		})
	}
}

func Test_keys_shortPrompt(t *testing.T) {
	var fruits = keys{"apple", "orange", "grape"}
	tests := []struct {
		name    string
		k       keys
		r       io.Reader
		wantKey string
	}{
		{"no match", fruits, strings.NewReader("some input"), ""},
		{"match", fruits, strings.NewReader("orange"), "orange"},
		{"short", fruits, strings.NewReader("o"), "orange"},
		{"remove", fruits, strings.NewReader("-"), "-"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotKey := tt.k.shortPrompt(tt.r); gotKey != tt.wantKey {
				t.Errorf("keys.shortPrompt() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestSkipSet(t *testing.T) {
	if ss := SkipSet(false); ss != "" {
		t.Errorf("SkipSet(true) = %s, want empty", ss)
	}
	const want = "  skipped setting"
	if ss := SkipSet(true); ss != want {
		t.Errorf("SkipSet(true) = %s, want %s", ss, want)
	}
}

// mockInput uses the os pipe to mock the user input.
// os.Pipe() https://stackoverflow.com/questions/46365221/fill-os-stdin-for-function-that-reads-from-it
func mockInput(input string) (*os.File, error) {
	s := []byte(input)
	r, w, err := os.Pipe()
	if err != nil {
		return r, err
	}
	_, err = w.Write(s)
	if err != nil {
		return r, err
	}
	w.Close()
	return r, nil
}

func TestPort(t *testing.T) {
	type args struct {
		validate bool
		setup    bool
		input    string
	}
	tests := []struct {
		name     string
		args     args
		wantPort uint
	}{
		{"empty", args{}, 0},
		{"bad input", args{true, true, "abc"}, 0},
		{"valid", args{true, false, "80"}, 80},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mockInput(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotPort := Port(tt.args.validate, tt.args.setup); gotPort != tt.wantPort {
				t.Errorf("Port() = %v, want %v", gotPort, tt.wantPort)
			}
		})
	}
}

func TestIndexStrings(t *testing.T) {
	options := []string{"number1", "number2", "number3"}
	type args struct {
		options *[]string
		setup   bool
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, false, "xyz"}, ""},
		{"ok input", args{&options, false, "number2"}, "number2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mockInput(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := IndexStrings(tt.args.options, tt.args.setup); gotKey != tt.wantKey {
				t.Errorf("IndexStrings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestShortStrings(t *testing.T) {
	options := []string{"alpha", "beta", "gamma"}
	type args struct {
		options *[]string
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, "xyz"}, ""},
		{"ok input", args{&options, "b"}, "beta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mockInput(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := ShortStrings(tt.args.options); gotKey != tt.wantKey {
				t.Errorf("ShortStrings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"mixed input", "123 abc", "123 abc"},
		{"ok input", "hello world", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mockInput(tt.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := String(); gotKey != tt.want {
				t.Errorf("String() = %v, want %v", gotKey, tt.want)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	options := []string{"alpha", "beta", "gamma"}
	type args struct {
		options *[]string
		input   string
	}
	tests := []struct {
		name    string
		args    args
		wantKey string
	}{
		{"empty", args{}, ""},
		{"bad input", args{&options, "xyz"}, ""},
		{"ok input", args{&options, "beta"}, "beta"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := mockInput(tt.args.input)
			if err != nil {
				t.Error(err)
			}
			stdin := os.Stdin
			defer func() {
				os.Stdin = stdin
			}()
			os.Stdin = r
			if gotKey := Strings(tt.args.options, false); gotKey != tt.wantKey {
				t.Errorf("Strings() = %v, want %v", gotKey, tt.wantKey)
			}
		})
	}
}
