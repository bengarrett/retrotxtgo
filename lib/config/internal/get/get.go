package get

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

var (
	ErrBool   = errors.New("key is not a boolean (true/false) value")
	ErrString = errors.New("key is not a string (text) value")
	ErrUint   = errors.New("key is not a absolute number")
)

// Defaults for setting keys.
type Defaults map[string]interface{}

// Hints for configuration values.
type Hints map[string]string

// TextEditor returns the path of a configured or discovered text editor.
func TextEditor(w io.Writer) string {
	edit := viper.GetString("editor")
	_, err := exec.LookPath(edit)
	if err != nil {
		if edit != "" {
			fmt.Fprintf(w, "%s\nwill attempt to use the $EDITOR environment variable\n", err)
		}
		if err := viper.BindEnv("editor", "EDITOR"); err != nil {
			return DiscEditor()
		}
		edit = viper.GetString("editor")
		if _, err := exec.LookPath(edit); err != nil {
			return DiscEditor()
		}
	}
	return edit
}

// DiscEditor attempts to discover any known text editors on the host system.
func DiscEditor() string {
	editors := [5]string{"nano", "vim", "emacs"}
	if runtime.GOOS == "windows" {
		editors[3] = "notepad++.exe"
		editors[4] = "notepad.exe"
	}
	edit := ""
	for _, editor := range editors {
		if editor == "" {
			continue
		}
		if _, err := exec.LookPath(editor); err == nil {
			edit = editor
			break
		}
	}
	return edit
}
