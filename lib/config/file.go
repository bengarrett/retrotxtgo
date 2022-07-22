package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/logs"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	gap "github.com/muesli/go-app-paths"
	"github.com/spf13/viper"
)

const NameFile = "config.yaml"

// InitDefaults initializes flag and configuration defaults.
func InitDefaults() {
	for key, val := range get.Reset() {
		viper.SetDefault(key, val)
		viper.Set(key, val)
	}
}

// Path returns the absolute path of the configuration file.
func Path() string {
	dir, err := gap.NewScope(gap.User, meta.Dir).ConfigPath(NameFile)
	if err != nil {
		var h string
		if h, err = os.UserHomeDir(); err != nil {
			return ""
		}
		return filepath.Join(h, NameFile)
	}
	return dir
}

// Location returns the absolute path of the current configuration file
// and the status of any missing settings.
func Location() string {
	s := fmt.Sprintf("Config file: %s", str.Path(viper.ConfigFileUsed()))
	if diff := len(Missing()); diff > 0 {
		if diff == 1 {
			s += str.ColSec(" (1 setting is missing)")
		} else {
			s += str.ColSec(fmt.Sprintf(" (%d settings are missing)", diff))
		}
	}
	return fmt.Sprintln(s)
}

// SetConfig reads and loads a configuration file.
func SetConfig(w io.Writer, flag string) error {
	if w == nil {
		return ErrWriter
	}
	viper.SetConfigType("yaml")
	path := flag
	if flag == "" {
		path = Path()
	}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return readInCfgErr(w, flag, err)
	}
	if flag != "" {
		if len(viper.AllKeys()) == 0 {
			return fmt.Errorf("set config: %w", logs.ErrConfigRead)
		}
		// always print the config location when the --config flag is used
		if len(os.Args) > 0 && os.Args[1] == "config" {
			// except when the config command is in use
			return nil
		}
		fmt.Fprint(w, Location())
	}
	// otherwise settings are loaded from default config
	return nil
}

func readInCfgErr(w io.Writer, flag string, err error) error {
	if w == nil {
		return ErrWriter
	}
	if flag == "" {
		if errors.Is(err, os.ErrNotExist) {
			// auto-generate new config except when the --config flag is used
			return Create(w, viper.ConfigFileUsed(), false)
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
		fmt.Fprintln(w, logs.Hint(fmt.Sprintf("config create --config=%s", flag), err))
		return err
	}
	// user given config file fail
	if strings.Contains(err.Error(), "found character that cannot start any token") {
		return logs.ErrConfigRead
	}
	return err
}

// ConfigMissing prints an config file error notice.
func ConfigMissing(w io.Writer, name, suffix string) {
	cmd := strings.TrimSuffix(name, suffix) + " create"
	if used := viper.ConfigFileUsed(); used != "" {
		fmt.Fprintf(w, "%s %q config file is missing\ncreate it: %s\n",
			str.Info(), used, str.ColPri(cmd+" --config="+used))
		return
	}
	fmt.Fprintf(w, "%s no config file is in use\ncreate it: %s\n",
		str.Info(), str.ColPri(cmd))
}

// Save all viper settings to the named file.
func Save(w io.Writer, name string) error {
	if name != "" {
		viper.SetConfigName(name)
	}
	if viper.ConfigFileUsed() == "" {
		viper.SetConfigFile(Path())
	}
	for key, val := range get.Reset() {
		if !viper.IsSet(key) {
			viper.Set(key, val)
		}
	}
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	if w != nil {
		fmt.Fprintln(w, "The change is saved")
	}
	return nil
}
