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
	"github.com/bengarrett/retrotxtgo/meta"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// InitDefaults initializes flag and configuration defaults.
func InitDefaults() {
	for key, val := range Reset() {
		viper.SetDefault(key, val)
		viper.Set(key, val)
	}
}

// Path is the absolute path and filename of the configuration file.
func Path() string {
	dir, err := gap.NewScope(gap.User, meta.Dir).ConfigPath(namedFile)
	if err != nil {
		h := ""
		if h, err = os.UserHomeDir(); err != nil {
			return ""
		}
		return filepath.Join(h, namedFile)
	}
	return dir
}

// Location returns the absolute path of the current configuration file
// and the status of any missing settings.
func Location() string {
	s := fmt.Sprintf("Config file: %s", str.Path(viper.ConfigFileUsed()))
	if diff := len(Missing()); diff > 0 {
		if diff == 1 {
			s += str.Cb(" (1 setting is missing)")
		} else {
			s += str.Cb(fmt.Sprintf(" (%d settings are missing)", diff))
		}
	}
	return fmt.Sprintln(s)
}

// SetConfig reads and loads a configuration file.
func SetConfig(flag string) (err error) {
	viper.SetConfigType("yaml")
	path := flag
	if flag == "" {
		path = Path()
	}
	viper.SetConfigFile(path)
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
			// initialize a new, default config file if conditions are met
			const min = 2
			if len(os.Args) > min {
				switch strings.Join(os.Args[1:3], ".") {
				case "config.create", "config.delete":
					// never auto-generate a config when these arguments are given
					return nil
				}
			}
			fmt.Println(logs.Hint(fmt.Sprintf("config create --config=%s", flag), err))
			return err
		}
		// user given config file fail
		if strings.Contains(err.Error(), "found character that cannot start any token") {
			return logs.ErrCfg
		}
		return err
	}
	if flag != "" {
		if len(viper.AllKeys()) == 0 {
			return fmt.Errorf("set config: %w", logs.ErrCfg)
		}
		// always print the config location when the --config flag is used
		if len(os.Args) > 0 && os.Args[1] == "config" {
			// except when the config command is in use
			return nil
		}
		fmt.Print(Location())
	}
	// otherwise settings are loaded from default config
	return nil
}

// UpdateConfig saves all viper settings to the named file.
func UpdateConfig(name string, stdout bool) (err error) {
	if name == "" {
		name = viper.ConfigFileUsed()
	}
	data, err := Marshal()
	if err != nil {
		return fmt.Errorf("config update marshal failed: %w", err)
	}
	out, err := yaml.Marshal(&data)
	if err != nil {
		return fmt.Errorf("config update marshal data failed: %w", err)
	}
	// prepend comments
	cmt := []byte(fmt.Sprintf("# %s configuration file", meta.Name))
	out = bytes.Join([][]byte{cmt, out}, []byte("\n"))
	if err = ioutil.WriteFile(name, out, filemode); err != nil {
		return fmt.Errorf("config update saving data to the file failed: %q: %w", name, err)
	}
	if stdout {
		fmt.Println("The change is saved")
	}
	return nil
}

// configMissing prints an error notice and exits.
func configMissing(name, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + " create"
	used := viper.ConfigFileUsed()
	if used != "" {
		fmt.Printf("%s %q config file is missing\ncreate it: %s\n",
			str.Info(), used, str.Cp(cmd+" --config="+used))
		return
	}
	fmt.Printf("%s no config file is in use\ncreate it: %s\n",
		str.Info(), str.Cp(cmd))
}
