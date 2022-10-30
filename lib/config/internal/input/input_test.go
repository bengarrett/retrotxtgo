package input_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/config/internal/input"
	"github.com/gookit/color"
)

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
		{"empty", args{}, "Set"},
		{"key", args{get.Keywords, "ooooh"}, "Replace"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP := input.PreviewPrompt(tt.args.name, tt.args.value)
			firstWord := strings.Split(strings.TrimSpace(gotP), " ")[0]
			if firstWord != tt.wantP {
				t.Errorf("PreviewPrompt() = %v, want %v", firstWord, tt.wantP)
			}
		})
	}
}

func TestColorScheme(t *testing.T) {
	color.Enable = false
	ue := input.Update{}
	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", true},
		{"valid name", input.Update{Name: "html.meta.color_scheme", Value: "abc"},
			"<meta name=\"color_scheme\" value=\"abc\">", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.ColorScheme(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("ColorScheme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("ColorScheme() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	t.Run("iterate", func(t *testing.T) {
		for key := range get.Reset() {
			s := input.Defaults(key)
			fmt.Println(s)
		}
	})
}

func TestEditor(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", false},
		{"valid name", input.Update{Name: "editor", Value: "abc"},
			"Set a text editor to launch", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.Editor(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("Editor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("Editor() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestLayout(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", true},
		{"valid name", input.Update{Name: "html.layout", Value: "abc"},
			"Choose a HTML template", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.Layout(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("Layout() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("Layout() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestServe(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", true},
		{"invalid value", input.Update{Name: "serve", Value: "abc"},
			"8086", false},
		{"valid value", input.Update{Name: "serve", Value: 8888},
			"8888", false},
		{"out of range value", input.Update{Name: "serve", Value: 999999},
			"8086", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.Serve(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("Serve() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("Serve() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestSaveDir(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""},
			"Choose a new", false},
		{"invalid value", input.Update{Name: "save_directory", Value: 012},
			"Choose a new directory", false},
		{"valid value", input.Update{Name: "save_directory", Value: "~"},
			"Choose a new directory", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.SaveDir(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("SaveDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("SaveDir() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestStyleHTML(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", true},
		{"invalid value", input.Update{Name: "style.html", Value: 012},
			" 0 <abap=\"abap\">", false},
		{"valid value", input.Update{Name: "style.html", Value: "nord"},
			" 0 <abap=\"abap\">", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.StyleHTML(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("StyleHTML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("StyleHTML() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestStyleInfo(t *testing.T) {
	color.Enable = false
	ue := input.Update{}

	tests := []struct {
		name    string
		u       input.Update
		wantW   string
		wantErr bool
	}{
		{"nil", ue, "", true},
		{"empty name", input.Update{Name: "", Value: nil}, "", true},
		{"invalid name", input.Update{Name: "xyz", Value: ""}, "", true},
		{"invalid value", input.Update{Name: "style.info", Value: 012},
			" 0 { \"abap\":\"abap\" }", false},
		{"valid value", input.Update{Name: "style.info", Value: "nord"},
			" 0 { \"abap\":\"abap\" }", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := input.StyleInfo(w, tt.u); (err != nil) != tt.wantErr {
				t.Errorf("StyleInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(w.String(), tt.wantW) {
				t.Errorf("StyleInfo() does not contain %v", tt.wantW)
				fmt.Printf("%s\n", w)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	const want = "Hello World."
	tests := []struct {
		s    string
		want string
	}{
		{"hello world.", want},
		{"HELLO world.", want},
		{"hEllO wOrLD.", want},
		{"ħEllO wOrLD.", "Ħello World."},
		{"123, hello world!?", "123, Hello World!?"},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			if got := input.Title(tt.s); got != tt.want {
				t.Errorf("Title() = %v, want %v", got, tt.want)
			}
		})
	}
}
