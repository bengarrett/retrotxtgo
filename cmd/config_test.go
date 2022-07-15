package cmd_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
)

var ErrTmpCfg = errors.New("cannot create the temporary test file")

func createConfig(name string) (string, error) {
	if name == "" {
		name = "test_retrotxt_cfg-"
	}
	tmpFile, err := ioutil.TempFile(os.TempDir(), name)
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
