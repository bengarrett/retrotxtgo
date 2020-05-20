package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Edit a configuration file.
func Edit() (err logs.IssueErr) {
	var nf = logs.IssueErr{
		Issue: "no suitable editor could be found",
		Err:   errors.New("set one by creating an $EDITOR environment variable in your shell configuration"),
	}
	PrintLocation()
	file := viper.ConfigFileUsed()
	if file == "" {
		configMissing(cmdPath, "edit")
	}
	var edit string
	if err := viper.BindEnv("editor", "EDITOR"); err != nil {
		edit = lookEdit()
		if edit == "" {
			return nf
		}
	} else {
		edit = viper.GetString("editor")
		if _, err := exec.LookPath(edit); err != nil {
			return logs.IssueErr{
				Issue: "editor command failed",
				Err:   err,
			}
		}
	}
	// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
	exe := exec.Command(edit, file)
	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	if err := exe.Run(); err != nil {
		return logs.IssueErr{
			Issue: "failed to run editor" + fmt.Sprintf(" %q", edit),
			Err:   err,
		}
	}
	return logs.IssueErr{}
}

func lookEdit() (edit string) {
	editors := []string{"nano", "micro", "vim", "emacs"}
	if runtime.GOOS == "windows" {
		editors = append(editors, "notepad++.exe", "notepad.exe")
	}
	for _, editor := range editors {
		if _, err := exec.LookPath(editor); err == nil {
			edit = editor
			break
		}
	}
	return edit
}
