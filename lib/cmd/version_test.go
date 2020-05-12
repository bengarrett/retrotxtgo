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
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func Test_arch(t *testing.T) {
	type args struct {
		v string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty", args{""}, ""},
		{"invalid", args{"xxx"}, ""},
		{"386", args{"386"}, "32-bit Intel/AMD"},
		{"ppc64", args{"ppc64"}, "64-bit PowerPC"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arch(tt.args.v); got != tt.want {
				t.Errorf("arch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_binary(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binary(); got != tt.want {
				t.Errorf("binary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_info(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"os and arch", fmt.Sprintf("%s/%s [%s CPU]", runtime.GOOS, runtime.GOARCH, arch(runtime.GOARCH))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := info()["os"]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionJSON(t *testing.T) {
	type args struct {
		indent bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"no indent", args{false}, true},
		{"indent", args{true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json.Valid(versionJSON(tt.args.indent)); got != tt.want {
				t.Errorf("versionJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_versionText(t *testing.T) {
	type args struct {
		c bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionText(tt.args.c)
		})
	}
}
