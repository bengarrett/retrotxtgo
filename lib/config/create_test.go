package config

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
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
		//{"temp no overwrite", args{tmpFile.Name(), false}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Create(tt.args.name, tt.args.ow); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
