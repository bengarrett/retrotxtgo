package update

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/bengarrett/retrotxtgo/lib/config/internal/get"
	"github.com/bengarrett/retrotxtgo/lib/internal/save"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/bengarrett/retrotxtgo/meta"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// Bool the boolean value of the named setting.
func Bool(b bool, name string) string {
	switch b {
	case true:
		return fmt.Sprintf("\n  The %s is enabled.\n", str.ColFuz(name))
	default:
		return fmt.Sprintf("\n  The %s is not in use.\n", str.ColFuz(name))
	}
}

// String the string value of the named setting.
func String(w io.Writer, s, name, value string) {
	const sd = get.SaveDir
	switch s {
	case "":
		fmt.Fprintf(w, "\n  The empty %s setting is not in use.\n",
			str.ColFuz(name))
		if name == sd {
			fmt.Fprintf(w, "  Files created by %s will always be saved to the active directory.\n",
				meta.Name)
		}
	default:
		fmt.Fprintf(w, "\n  The %s is set to %q.", str.ColFuz(name), value)
		// print the operating system's ability to use the existing set values
		// does the 'editor' exist in the env path, does the save-directory exist?
		switch name {
		case "editor":
			_, err := exec.LookPath(value)
			fmt.Fprint(w, " ", str.Bool(err == nil))
		case sd:
			f := false
			if _, err := os.Stat(value); !os.IsNotExist(err) {
				f = true
			}
			fmt.Fprint(w, " ", str.Bool(f))
		default:
		}
		fmt.Fprintln(w)
	}
}

// Config saves all viper settings to the named file.
func Config(w io.Writer, name string, stdout bool) error {
	if name == "" {
		name = viper.ConfigFileUsed()
	}
	data, err := get.Marshal()
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
	if err = ioutil.WriteFile(name, out, save.FileMode); err != nil {
		return fmt.Errorf("config update saving data to the file failed: %q: %w", name, err)
	}
	if stdout {
		fmt.Fprintln(w, "The change is saved")
	}
	return nil
}
