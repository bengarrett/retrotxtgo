package config

import (
	"os"
	"reflect"
	"testing"

	"retrotxt.com/retrotxt/lib/str"
)

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"list", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := List(); (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		name string
		args string
	}{
		{"empty", ""},
		{"0", "0"},
		{"valid", "editor"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Set(tt.name)
		})
	}
}

func Test_copyKeys(t *testing.T) {
	tests := []struct {
		name     string
		keys     []string
		wantCopy []string
	}{
		{"empty", []string{}, []string{}},
		{"3 vals", []string{"a", "b", "c"}, []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCopy := copyKeys(tt.keys...); !reflect.DeepEqual(gotCopy, tt.wantCopy) {
				t.Errorf("copyKeys() = %v, want %v", gotCopy, tt.wantCopy)
			}
		})
	}
}

func Test_names_string(t *testing.T) {
	tests := []struct {
		name  string
		n     names
		theme bool
		want  string
	}{
		{"nil", nil, false, ""},
		{"empty", names{""}, false, ""},
		{"ok", names{"okay"}, false, "okay"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.string(tt.theme); got != tt.want {
				t.Errorf("names.string() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dirExpansion(t *testing.T) {
	u, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	w, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	tests := []struct {
		name    string
		wantDir string
	}{
		{"", ""},
		{"~", u},
		{".", w},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := dirExpansion(tt.name); gotDir != tt.wantDir {
				t.Errorf("dirExpansion() = %v, want %v", gotDir, tt.wantDir)
			}
		})
	}
}

func Test_colorElm(t *testing.T) {
	// set test mode for str.HighlightWriter()
	str.TestMode = true
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
			if got := colorElm(tt.elm, "html", "bw"); got != tt.want {
				t.Errorf("colorhtml() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColorCSS(t *testing.T) {
	tests := []struct {
		name string
		elm  string
		want string
	}{
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ColorCSS(tt.elm); got != tt.want {
				t.Errorf("ColorCSS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_previewPrompt(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name  string
		args  args
		wantP string
	}{
		{"empty", args{}, "Set a new value or leave blank to keep it unused:"},
		{"key", args{"html.meta.keywords", "ooooh"}, "Replace the current keywords, leave blank to keep as-is or use a dash [-] to remove:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotP := previewPrompt(tt.args.name, tt.args.value); gotP != tt.wantP {
				t.Errorf("previewPrompt() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}
