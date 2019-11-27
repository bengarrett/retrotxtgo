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
	"reflect"
	"testing"
)

func Test_read(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := read(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tmpl(t *testing.T) {
	type args struct {
		args    []string
		testing bool
	}
	args0 := args{[]string{""}, true}
	args1 := args{[]string{"somefile.txt"}, true}
	args2 := args{[]string{"../textfiles/hi.txt"}, true}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no arguments", args0, true},
		{"invalid file", args1, true},
		{"valid file", args2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tmpl(tt.args.args, tt.args.testing); (err != nil) != tt.wantErr {
				t.Errorf("tmpl() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
