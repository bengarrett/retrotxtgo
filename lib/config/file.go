package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type configErr struct {
	FileUsed string
	Err      error
}

func (e configErr) String() string {
	return (logs.Err{
		Issue: "config file",
		Arg:   e.FileUsed,
		Msg:   e.Err}).String()
}

// InitDefaults initialises flag and configuration defaults.
func InitDefaults() {
	for key, val := range Defaults {
		viper.SetDefault(key, val)
		viper.Set(key, val)
	}
}

// configMissing prints an error notice and exits.
func configMissing(name string, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + " create"
	fmt.Printf("%s no config file is in use\n create it: %s\n",
		str.Info(), str.Cp(cmd))
	os.Exit(21)
}

func exit(e configErr) {
	if e.Err != nil {
		fmt.Println(e.String())
		os.Exit(1)
	}
	os.Exit(0)
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

// PrintLocation prints the location of the current configuration file.
func PrintLocation() {
	s := str.Cb(fmt.Sprintf("Config file: %s",
		viper.ConfigFileUsed()))
	if diff := len(Missing()); diff > 0 {
		if diff == 1 {
			s += str.Cb(fmt.Sprint(" (1 setting is missing)"))
		} else {
			s += str.Cb(fmt.Sprintf(" (%d settings are missing)", diff))
		}
	}
	fmt.Println(s)
}

// SetConfig reads and loads a configuration file.
func SetConfig(flag string) {
	viper.SetConfigType("yaml")
	var path = flag
	if flag == "" {
		path = Path()
	}
	viper.SetConfigFile(path)
	var e = configErr{
		FileUsed: viper.ConfigFileUsed(),
		Err:      nil,
	}
	if e.Err = viper.ReadInConfig(); e.Err != nil {
		if flag == "" {
			if errors.Is(e.Err, os.ErrNotExist) {
				// auto-generate new config except when the --config flag is used
				e.Err = Create(viper.ConfigFileUsed(), false)
				exit(e)
			}
			// internal config fail
			exit(e)
		}
		if errors.Is(e.Err, os.ErrNotExist) {
			// initialise a new, default config file if conditions are met
			if len(os.Args) > 2 {
				switch strings.Join(os.Args[1:3], ".") {
				case "config.create", "config.delete":
					// never auto-generate a config when these arguments are given
					return
				}
			}
			e.Err = errors.New("does not exist\n\t use the command: retrotxt config create --config=" + flag)
			exit(e)
		} else {
			// user config file fail
			exit(e)
		}
	} else if flag != "" {
		if len(viper.AllKeys()) == 0 {
			e.Err = errors.New("is not a retrotxt config file")
			exit(e)
		}
		// always print the config location when the --config flag is used
		if len(os.Args) > 0 && os.Args[1] == "config" {
			// except when the config command is in use
			return
		}
		PrintLocation()
	}
	// otherwise settings are loaded from default config
}

// UpdateConfig saves all viper settings to the named file.
func UpdateConfig(name string, print bool) (err error) {
	if name == "" {
		name = viper.ConfigFileUsed()
	}
	data, err := Marshal()
	logs.Check("config.update:", err)
	out, err := yaml.Marshal(&data)
	if err != nil {
		return err
	}
	// prepend comments
	cmt := []byte("# RetroTxt configuration file")
	out = bytes.Join([][]byte{cmt, out}, []byte("\n"))
	err = ioutil.WriteFile(name, out, PermF)
	if err != nil {
		return err
	}
	if print {
		fmt.Println("The change is saved")
	}
	return err
}
