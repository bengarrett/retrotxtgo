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
	"testing"
)

func TestErrorFmt_errorPrint(t *testing.T) {
	err := errors.New("err-text")
	h := ErrorFmt{"error", "test", err}
	tests := []struct {
		name string
		e    *ErrorFmt
		want string
	}{
		{"default", &h, "\nERROR: error test err-text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.errorPrint(); got != tt.want {
				t.Errorf("ErrorFmt.errorPrint() = %v, want %v", got, tt.want)
			}
		})
	}
}
