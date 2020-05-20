package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/spf13/viper"
)

// Edit a configuration file.
func Edit() {
	cfg := viper.ConfigFileUsed()
	if cfg == "" {
		configMissing(cmdPath, "edit")
	}
	var edit string
	if err := viper.BindEnv("editor", "EDITOR"); err != nil {
		editors := []string{"nano", "vim", "emacs"}
		if runtime.GOOS == "windows" {
			editors = append(editors, "notepad++.exe", "notepad.exe")
		}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				edit = editor
				break
			}
		}
		if edit != "" {
			fmt.Printf("There is no %s environment variable set so using: %s\n",
				logs.Ci("EDITOR"), logs.Cp(edit))
		} else {
			editNotFound()
		}
	} else {
		edit = viper.GetString("editor")
		if _, err := exec.LookPath(edit); err != nil {
			logs.Check("edit command not found", exec.ErrNotFound)
			os.Exit(exit + 4)
		} else {
			editNotFound()
		}
	}
	// credit: https://stackoverflow.com/questions/21513321/how-to-start-vim-from-go
	exe := exec.Command(edit, cfg)
	exe.Stdin = os.Stdin
	exe.Stdout = os.Stdout
	if err := exe.Run(); err != nil {
		fmt.Printf("%s\n", err)
	}
}

func editNotFound() {
	log.Println("no suitable editor could be found\nplease set one by creating a $EDITOR environment variable in your shell configuration")
	os.Exit(exit + 3)
}
