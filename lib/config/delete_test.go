package config

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/logs"
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
		wantErr logs.IssueErr
	}{
		{"ok", logs.IssueErr{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "ok" {
				viper.SetConfigFile(tmpFile.Name())
			}
			if gotErr := Delete(); gotErr != tt.wantErr {
				t.Errorf("Delete() = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}