package convert

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/filesystem"
	"golang.org/x/text/encoding/charmap"
)

const (
	cp437hex = "\xCD\xB9\xB2\xCC\xCD" // `═╣░╠═`
	utf      = "═╣ ░ ╠═"

	base = "rt_sample-"

	permf os.FileMode = 0644
)

func ExampleD437() {
	const name = base + "cp437In.txt"
	result, err := D437(cp437hex)
	if err != nil {
		log.Fatal(err)
	}
	filesystem.SaveTemp(name, result)
	t, err := filesystem.ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(t)
	// Output: ═╣▓╠═
}
func TestCP437Decode(t *testing.T) {
	tests := []struct {
		name       string
		s          string
		wantResult string
		wantErr    bool
	}{
		{"empty", "", "", false},
		{"hex", cp437hex, "═╣▓╠═", false},
		{"nl", filesystem.Newlines, filesystem.Newlines, false},
		{"utf", filesystem.Symbols, "[Γÿá|Γÿ«|ΓÖ║]", false},
		{"escapes", filesystem.Escapes, `bell:,back:,tab:	,form:,vertical:,quote:"`, false},
		{"digits", filesystem.Digits, "░░┼░┼░", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := D437(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(string(gotResult), tt.wantResult) {
				t.Errorf("D437() = %v, want %v", string(gotResult), tt.wantResult)
			}
		})
	}
}

func ExampleE437() {
	const name = base + "cp437.txt"
	result, err := E437(utf)
	if err != nil {
		log.Fatal(err)
	}
	filesystem.SaveTemp(name, result)
	t, err := filesystem.ReadText(name)
	if err != nil {
		log.Fatal(err)
	}
	filesystem.Clean(name)
	fmt.Print(len(t))
	// Output: 8
}

func TestDString(t *testing.T) {
	type args struct {
		s  string
		cp charmap.Charmap
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := DString(tt.args.s, tt.args.cp)
			if (err != nil) != tt.wantErr {
				t.Errorf("DString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("DString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestEString(t *testing.T) {
	type args struct {
		s  string
		cp charmap.Charmap
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := EString(tt.args.s, tt.args.cp)
			if (err != nil) != tt.wantErr {
				t.Errorf("EString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("EString() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestD437(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := D437(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("D437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("D437() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestE437(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantResult []byte
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := E437(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("E437() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("E437() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}