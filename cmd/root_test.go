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

	"github.com/spf13/cobra"
)

func TestLayoutDefault(t *testing.T) {
	tests := []struct {
		name string
		want PageData
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LayoutDefault(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LayoutDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Execute()
		})
	}
}

func Test_initConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initConfig()
		})
	}
}

func TestErrorFmt_ErrorPrint(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorFmt
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ErrorPrint(); got != tt.want {
				t.Errorf("ErrorFmt.ErrorPrint() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorFmt_GoErr(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorFmt
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.GoErr()
		})
	}
}

func TestErrorFmt_FlagErr(t *testing.T) {
	tests := []struct {
		name string
		e    *ErrorFmt
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.FlagErr()
		})
	}
}

func TestErrorFmt_UsageErr(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		e    *ErrorFmt
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.e.UsageErr(tt.args.cmd)
		})
	}
}

func TestFileMissingErr(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			FileMissingErr()
		})
	}
}
