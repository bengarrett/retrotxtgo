package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"retrotxt.com/retrotxt/lib/str"
)

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

// Path is the absolute path and filename of the configuration file.
func Path() (dir string) {
	dir, err := scope.ConfigPath(namedFile)
	if err != nil {
		h := ""
		if h, err = os.UserHomeDir(); err != nil {
			return ""
		}
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
func SetConfig(flag string) (err error) {
	viper.SetConfigType("yaml")
	var path = flag
	if flag == "" {
		path = Path()
	}
	viper.SetConfigFile(path)
	// viper.ConfigFileUsed()
	if err = viper.ReadInConfig(); err != nil {
		if flag == "" {
			if errors.Is(err, os.ErrNotExist) {
				// auto-generate new config except when the --config flag is used
				return Create(viper.ConfigFileUsed(), false)
			}
			// internal config fail
			return err
		}
		if errors.Is(err, os.ErrNotExist) {
			// initialise a new, default config file if conditions are met
			if len(os.Args) > 2 {
				switch strings.Join(os.Args[1:3], ".") {
				case "config.create", "config.delete":
					// never auto-generate a config when these arguments are given
					return nil
				}
			}
			return fmt.Errorf("does not exist\n\t use the command:retrotxt config create --config=%s", flag)
		}
		// user config file fail
		return nil
	} else if flag != "" {
		if len(viper.AllKeys()) == 0 {
			return errors.New("is not a retrotxt config file")
		}
		// always print the config location when the --config flag is used
		if len(os.Args) > 0 && os.Args[1] == "config" {
			// except when the config command is in use
			return nil
		}
		PrintLocation()
	}
	// otherwise settings are loaded from default config
	return nil
}

// UpdateConfig saves all viper settings to the named file.
func UpdateConfig(name string, print bool) (err error) {
	if name == "" {
		name = viper.ConfigFileUsed()
	}
	data, err := Marshal()
	if err != nil {
		return fmt.Errorf("config update marshal failed: %s", err)
	}
	out, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("config update marshal data failed: %s", err)
	}
	// prepend comments
	cmt := []byte("# RetroTxt configuration file")
	out = bytes.Join([][]byte{cmt, out}, []byte("\n"))
	if err = ioutil.WriteFile(name, out, filemode); err != nil {
		return fmt.Errorf("config update saving data to the file failed: %q: %s", name, err)
	}
	if print {
		fmt.Println("The change is saved")
	}
	return nil
}
