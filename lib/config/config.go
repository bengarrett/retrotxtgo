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
	"gopkg.in/yaml.v2"
)

type files map[string]string

type ports struct {
	max uint
	min uint
	rec uint
}

var port = ports{
	max: logs.PortMax,
	min: logs.PortMin,
	rec: logs.PortRec,
}

// Formats blah
type Formats struct {
	Info    []string
	Shell   []string
	Version []string
}

// Format ...
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

// ConfigName is the default configuration filename.
const ConfigName string = "config.yaml"
const cmdPath = "retrotxt config"

// posix permissions for the configuration file and directory.
const perm, permDir os.FileMode = 0660, 0700

// operating system exit code.
const exit = 20

var (
	scope      = gap.NewScope(gap.User, "retrotxt")
	infoStyles string
)

var cfgNameFlag string // TO IMPLEMENT?
var configSetName string

// BuildVer retrotxt version
var BuildVer string // remove?

// Filepath is the absolute path and filename of the configuration file.
func Filepath() (dir string) {
	dir, err := scope.ConfigPath(ConfigName)
	if err != nil {
		h, _ := os.UserHomeDir()
		return filepath.Join(h, ConfigName)
	}
	return dir
}

// configExit prints an error notice and exits.
func configExit(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + "create"
	fmt.Printf("No config file is in use.\nTo create run: %s\n", logs.Cp(cmd))
	os.Exit(exit + 1)
}

// todo: writefile requires settings such as creating parent dirs
func writeConfig(update bool) {
	out, err := yaml.Marshal(viper.AllSettings())
	logs.ChkErr("could not create: settings", err)
	err = ioutil.WriteFile(Filepath(), out, perm)
	logs.ChkErr("could not write: settings", err)
	s := "Created a new"
	if update {
		s = "Updated the"
	}
	fmt.Printf("%s config file at: %s\n", s, logs.Cf(Filepath()))
}