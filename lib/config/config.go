package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	// PermF is posix permission bits for files
	PermF os.FileMode = 0660
	// PermD is posix permission bits for directories
	PermD os.FileMode = 0700

	cmdPath   = "retrotxt config"
	namedFile = "config.yaml"
)

var scope = gap.NewScope(gap.User, "retrotxt")

// Formats choices for flags
type Formats struct {
	Info    []string
	Shell   []string
	Version []string
}

// Format flag choices for info, shell and version commands.
var Format = Formats{
	Info:    []string{"color", "json", "json.min", "text", "xml"},
	Shell:   []string{"bash", "powershell", "zsh"},
	Version: []string{"color", "json", "json.min", "text"},
}

func (f Formats) String(field string) string {
	switch field {
	case "info":
		return strings.Join(f.Info, ", ")
	case "shell":
		return strings.Join(f.Shell, ", ")
	case "version":
		return strings.Join(f.Version, ", ")
	}
	return ""
}

// Path is the absolute path and filename of the configuration file.
func Path() (dir string) {
	dir, err := scope.ConfigPath(namedFile)
	if err != nil {
		h, _ := os.UserHomeDir()
		return filepath.Join(h, namedFile)
	}
	return dir
}

// configMissing prints an error notice and exits.
func configMissing(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + " create"
	fmt.Printf("%s no config file is in use\n create it: %s\n",
		logs.Info(), logs.Cp(cmd))
	os.Exit(21)
}

// UpdateConfig saves all viper settings to the named file.
func UpdateConfig(name string, print bool) (err error) {
	if name == "" {
		name = viper.ConfigFileUsed()
	}
	out, err := yaml.Marshal(viper.AllSettings())
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(name, out, PermF)
	if err != nil {
		return err
	}
	if print {
		fmt.Println("The change is saved")
	}
	return err
}
