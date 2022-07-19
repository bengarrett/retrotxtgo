package config_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/config"
)

func TestCreate(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testcreate-")
	if err != nil {
		log.Fatal("Cannot create the temporary test file", err)
	}
	defer os.Remove(tmpFile.Name())
	type args struct {
		name string
		ow   bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{}, true},
		{"temp", args{tmpFile.Name(), true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := config.Create(os.Stdout, tt.args.name, tt.args.ow); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
