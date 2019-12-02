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
	"os"
	"reflect"
	"testing"
)

func Test_details(t *testing.T) {
	n := "../textfiles/hi.txt"
	type args struct {
		name string
	}
	got, err := details(n)
	if err != nil {
		t.Errorf("details() = %v, want %v", err, nil)
	}
	if got.Bytes != 39 {
		t.Errorf("details() = %v, want %v", got.Bytes, 39)
	}
	if got.Name != "hi.txt" {
		t.Errorf("details() = %v, want %v", got.Name, "hi.txt")
	}
	if got.Slug != "hi-txt" {
		t.Errorf("details() = %v, want %v", got.Slug, "hi-txt")
	}
	if got.Mime != "text/plain" {
		t.Errorf("details() = %v, want %v", got.Mime, "text/plain")
	}
	if got.Utf8 != true {
		t.Errorf("details() = %v, want %v", got.Utf8, true)
	}
	if got.MD5 != "a4098d2ae92f55f6aacc6865812a4291" {
		t.Errorf("details() = %v, want %v", got.MD5, "a4098d2ae92f55f6aacc6865812a4291")
	}
}

func Test_parse(t *testing.T) {
	f, err := os.Stat("../textfiles/hi.txt")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	type args struct {
		data []byte
		stat os.FileInfo
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{"empty", args{[]byte(""), f}, 0, false},
		{"string", args{[]byte("hello"), f}, 5, false},
		{"string", args{[]byte("世界你好"), f}, 8, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parse(tt.args.data, tt.args.stat)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.CharCount, tt.want) {
				t.Errorf("parse() = %v, want %v", got.CharCount, tt.want)
			}
		})
	}
}

func Test_infoJSON(t *testing.T) {
	var f Detail
	type args struct {
		indent bool
		f      Detail
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"no indent", args{false, f}, true},
		{"indent", args{true, f}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json.Valid(infoJSON(tt.args.indent, tt.args.f)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("infoJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_infoXML(t *testing.T) {
	n := "../textfiles/hi.txt"
	f, err := details(n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if got := len(infoXML(f)); got != 322 {
		t.Errorf("infoXML() = %v, want %v", got, 322)
	}
}
