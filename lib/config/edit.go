package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

var ErrEd = errors.New("there is no configuration file to edit")

// Edit a configuration file.
func Edit(w io.Writer) error {
	fmt.Fprintf(w, "%s%s", str.Info(), Location())
	file := viper.ConfigFileUsed()
	if file == "" {
		return ErrEd
	}
	edit := get.TextEditor(w)
	if edit == "" {
		return fmt.Errorf("create an $EDITOR environment variable in your shell configuration: %w", ErrEditorNil)
	}
	// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
	exe := exec.Command(edit, file)
	exe.Stdin = os.Stdin
	exe.Stdout = w
	if err := exe.Run(); err != nil {
		e := fmt.Errorf("%s: %w", edit, ErrEditorRun)
		return fmt.Errorf("%s: %w", e, err)
	}
	return nil
}
