package config_test

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/spf13/viper"
)

func TestDelete(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdelete-")
	if err != nil {
		log.Fatal("Cannot create the temporary test file", err)
	}
	defer os.Remove(tmpFile.Name())
	tests := []struct {
		name    string
		wantErr error
	}{
		{"ok", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "ok" {
				viper.SetConfigFile(tmpFile.Name())
			}
			if _, gotErr := config.Delete(true); !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("Delete() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
