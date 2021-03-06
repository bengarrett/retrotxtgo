package config

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/viper"
)

// Edit a configuration file.
func Edit() error {
	PrintLocation()
	file := viper.ConfigFileUsed()
	if file == "" {
		configMissing(cmdPath, "edit")
		os.Exit(1)
	}
	edit := Editor()
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

// Editor returns the path of a configured or discovered text editor.
func Editor() (edit string) {
	edit = viper.GetString("editor")
	_, err := exec.LookPath(edit)
	if err != nil {
		if edit != "" {
			fmt.Printf("%s\nwill attempt to use the $EDITOR environment variable\n", err)
		}
		if err := viper.BindEnv("editor", "EDITOR"); err != nil {
			return lookEdit()
		}
		edit = viper.GetString("editor")
		if _, err := exec.LookPath(edit); err != nil {
			return lookEdit()
		}
	}
	return edit
}

func lookEdit() (edit string) {
	editors := [5]string{"nano", "vim", "emacs"}
	if runtime.GOOS == "windows" {
		editors[3] = "notepad++.exe"
		editors[4] = "notepad.exe"
	}
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
