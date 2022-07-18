package cmd_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/gookit/color"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

// TODO: implement https://github.com/spf13/afero for temp dirs.

var ErrTmpCfg = errors.New("cannot create the temporary test file")

const missing = "These settings are missing and should be configured"

func init() {
}

func createConfig(name string) (string, error) {
	appFS := afero.NewMemMapFs()
	if name == "" {
		name = "test_retrotxt_cfg-"
	}
	//tmpFile, err := ioutil.TempFile(os.TempDir(), name)
	tmpFile, err := afero.TempFile(appFS, "", "ioutil-test")
	if err != nil {
		return "", fmt.Errorf("%w, temp dir: %s", ErrTmpCfg, os.TempDir())
	}
	viper.SetConfigFile(tmpFile.Name())
	return viper.ConfigFileUsed(), nil
}

func Test_createConfig(t *testing.T) {
	t.Run("noname", func(t *testing.T) {
		tmp, err := createConfig("")
		if err != nil {
			t.Errorf("createConfig returned an error: %s", err)
			return
		}
		defer func() {
			if err := os.Remove(tmp); err != nil {
				t.Error(err)
			}
		}()
	})
}

func Test_ConfigErr(t *testing.T) {
	t.Run("config invalid", func(t *testing.T) {
		const invalid = "zxcvbnnm"
		gotB, err := infoT.tester([]string{"--test", invalid})
		if err == nil {
			t.Errorf("using this invalid config command did not return an error: %s", invalid)
			t.Error(gotB)
		}
	})
}

func Test_ConfigCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantFormal string
	}{
		// help
		{"settings", nil,
			"RetroTxtGo settings"},
		{"config", nil,
			"retrotxt config"},
		{"reset", nil,
			"create      Create or reset the config file"},
		// create
		{"create", []string{"create"}, "A config file already exists"},
		{"create ow", []string{"create", "--overwrite"}, "A new file was created"},
		// delete
		{"delete", []string{"delete"}, "hello"},
		// {"edit", []string{"edit"}, "hello"},
		// {"info", []string{"info"}, "hello"},
		// {"set", []string{"set"}, "hello"},
		// {"set -c", []string{"set", "-c"}, "hello"},
		// {"setup", []string{"setup"}, "hello"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			color.Enable = false
			args := tt.args
			if args != nil {
				args = append(args, "--tester")
			}
			gotB, err := confT.tester(args)
			if err != nil {
				t.Error(err)
				return
			}
			if !bytes.Contains(gotB, []byte(tt.wantFormal)) {
				fmt.Println(string(gotB))
				t.Errorf("could not find %q text in config usage", tt.wantFormal)
			}
			if bytes.Contains(gotB, []byte(missing)) {
				t.Error("config file is missing some expected settings!")
			}
		})
	}
}
