package config_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
	"github.com/gookit/color"
)

func TestPath(t *testing.T) {
	tests := []struct {
		name    string
		wantDir string
	}{
		{"def", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDir := config.Path(); gotDir == tt.wantDir {
				t.Errorf("Path() = \"\"")
			}
		})
	}
}

func TestSetConfig(t *testing.T) {
	color.Enable = false
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testsetconfig-")
	if err != nil {
		log.Fatal("Cannot create the temporary test file", err)
	}
	defer os.Remove(tmpFile.Name())
	tests := []struct {
		name    string
		flag    string
		wantErr bool
	}{
		{"default", "", false},
		{"invalid", "this-file-doesnt-exist", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.SetConfig(tt.flag); (err != nil) != tt.wantErr {
				t.Errorf("SetConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
