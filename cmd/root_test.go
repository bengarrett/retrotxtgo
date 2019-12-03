/*
Copyright © 2019 Ben Garrett <code.by.ben@gmail.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/viper"
)

func TestInitDefaults(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"layout", "create.layout", "standard"},
		{"save dir", "create.save-directory", ""},
	}
	InitDefaults()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := viper.GetString(tt.key); got != tt.want {
				t.Errorf("InitDefaults() %v = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

// func Test_initConfig(t *testing.T) {
// 	cfgFile = ""
// 	initConfig()
// 	want := "%HOME/.retrotxtgo.yaml"
// 	if got := viper.ConfigFileUsed(); got != want {
// 		t.Errorf("initConfig() ConfigFileUsed() = %v, want %v", got, want)
// 	}
// }

func TestCheckFlag(t *testing.T) {
	type args struct {
		e ErrorFmt
	}
	a := ErrorFmt{"", "", nil}
	b := ErrorFmt{"x", "y", errors.New("z")}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"blank", args{}, false},
		{"empty", args{a}, false},
		{"xyz", args{b}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go
			if os.Getenv("BE_CRASHER") == "1" {
				CheckFlag(tt.args.e)
				return
			}
			cmd := exec.Command(os.Args[0], "-test.run=TestCrasher")
			cmd.Env = append(os.Environ(), "BE_CRASHER=1")
			err := cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}
			t.Errorf("process ran with err %v, want exit status 1", err)
		})
	}
}
