package config

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
)

// Edit a configuration file.
func Edit() error {
	fmt.Printf("%s%s",
		str.Info(), Location())
	file := viper.ConfigFileUsed()
	if file == "" {
		configMissing(cmdPath(), "edit")
		os.Exit(1)
	}
	edit := get.TextEditor()
	if edit == "" {
		return fmt.Errorf("create an $EDITOR environment variable in your shell configuration: %w", ErrEditorNil)
	}
	// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
	exe := exec.Command(edit, file)
	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	if err := exe.Run(); err != nil {
		e := fmt.Errorf("%s: %w", edit, ErrEditorRun)
		return fmt.Errorf("%s: %w", e, err)
	}
	return nil
}
