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
	const (
		setEdit = "Set a text editor to launch when using config edit:"
		hFE     = "Embed the font as Base64 text within the HTML"
		hFF     = "Choose a font, automatic, mona, vga"
	)
	inf := []string{"info", "--style", "none"}
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
		{"delete", []string{"delete"}, "File is deleted"},
		{"delete rm", []string{"rm"}, "File is deleted"},
		// edit
		// {"edit", []string{"edit"}, "hello"},
		// info
		{"info", inf, "\"editor\": \"\","},
		{"info", inf, "\"html.font.embed\": false,"},
		{"info", inf, "\"html.font.family\": \"vga\","},
		{"info", inf, "\"html.layout\": \"standard\","},
		{"info", inf, "\"html.meta.author\": \"\","},
		{"info", inf, "\"html.meta.color-scheme\": \"\","},
		{"info", inf, "\"html.meta.description\": \"\","},
		{"info", inf, "\"html.meta.generator\": true,"},
		{"info", inf, "\"html.meta.keywords\": \"\","},
		{"info", inf, "\"html.meta.notranslate\": false,"},
		{"info", inf, "\"html.meta.referrer\": \"\","},
		{"info", inf, "\"html.meta.retrotxt\": true,"},
		{"info", inf, "\"html.meta.robots\": \"\","},
		{"info", inf, "\"html.meta.theme-color\": \"\","},
		{"info", inf, "\"html.title\": \"RetroTxtGo\","},
		{"info", inf, "\"save-directory\": \"\","},
		{"info", inf, "\"serve\": 8086,"},
		{"info", inf, "\"style.html\": \"lovelace\","},
		{"info", inf, "\"style.info\": \"dracula\""},
		// info styles
		{"info configs", []string{"info", "--configs"},
			"Syntax highlighter for the config info output."},
		{"info styles", []string{"info", "--styles"},
			"retrotxt info --style=\"xcode-dark\""},
		// set styles
		{"set --list", []string{"set", "--list"},
			"Available RetroTxtGo configurations and settings"},
		{"set --list", []string{"set", "--list"},
			"-l, --list   list all the available setting names"},
		{"set help", []string{"set"}, "Edit a RetroTxtGo setting."},
		// set config items
		//{"set", []string{"set", "editor"}, setEdit},
		{"set", []string{"set", "html.font.embed"}, hFE},
		// {"set", []string{"set", "html.font.family"}, hFF},
		// {"set", []string{"set", "html.layout"}, hFF},
		// {"set", []string{"set", "html.meta.author"}, hFF},
		// {"set", []string{"set", "html.meta.color-scheme"}, hFF},
		// {"set", []string{"set", "html.meta.description"}, hFF},
		// {"set", []string{"set", "html.meta.generator"}, hFF},
		// {"set", []string{"set", "html.meta.keywords"}, hFF},
		// {"set", []string{"set", "html.meta.notranslate"}, hFF},
		// {"set", []string{"set", "html.meta.referrer"}, hFF},
		// {"set", []string{"set", "html.meta.retrotxt"}, hFF},
		// {"set", []string{"set", "html.meta.robots"}, hFF},
		// {"set", []string{"set", "html.meta.theme-color"}, hFF},
		// {"set", []string{"set", "html.title"}, hFF},
		// {"set", []string{"set", "save-directory"}, hFF},
		// {"set", []string{"set", "serve"}, hFF},
		// {"set", []string{"set", "style.html"}, hFF},
		// {"set", []string{"set", "style.info"}, hFF},
		// setup items
		// when --tester is enabled, setup skips all prompts and enters defaults
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
				t.Errorf("could not find %q text in config usage", tt.wantFormal)
				fmt.Println("<<", string(gotB), ">>")
			}
			if bytes.Contains(gotB, []byte(missing)) {
				t.Error("config file is missing some expected settings!")
			}
		})
	}
}
