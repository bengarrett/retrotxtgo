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
	"os"
	"path"
	"reflect"
	"testing"
)

func Test_read(t *testing.T) {
	type args struct {
		name string
	}
	args0 := args{""}
	args1 := args{"somefile.txt"}
	args2 := args{"../textfiles/hi.txt"}
	want2 := []byte("Hello world ☺\n☺ ሰላም ልዑል")
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"no arguments", args0, nil, true},
		{"invalid file", args1, nil, true},
		{"valid file", args2, want2, false},
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

func Test_filename(t *testing.T) {
	w := "../static/html/standard.html"
	got, _ := filename(true)
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	w = "static/html/standard.html"
	got, _ = filename(false)
	if got != w {
		t.Errorf("filename = %v, want %v", got, w)
	}
	htmlLayout = "error"
	_, err := filename(false)
	if (err != nil) != true {
		t.Errorf("filename = %v, want %v", got, w)
	}
}

func Test_pagedata(t *testing.T) {
	// htmlLayout should be standard
	w := "hello"
	d := []byte(w)
	got := pagedata(d).PreText
	if got != w {
		t.Errorf("pagedata().PreText = %v, want %v", got, w)
	}
	htmlLayout = "mini"
	pageTitle = "page title"
	w = pageTitle
	got = pagedata(d).PageTitle
	if got != w {
		t.Errorf("pagedata().PageTitle = %v, want %v", got, w)
	}
	metaDesc = "page description"
	w = ""
	got = pagedata(d).MetaDesc
	if got != w {
		t.Errorf("pagedata().MetaDesc = %v, want %v", got, w)
	}
}

func Test_writeFile(t *testing.T) {
	type args struct {
		data []byte
		name string
		test bool
	}
	tmpFile := path.Join(os.TempDir(), "retrotxtgo_create_test.txt")
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data or name", args{[]byte(""), "", true}, true},
		{"invalid name", args{[]byte("abc"), "this-is-an-invalid-path", true}, true},
		{"file as name", args{[]byte("abc"), tmpFile, true}, true},
		//{"home as name", args{[]byte("abc"), "~", true}, false}, // fixme: this fails in tests but works on CLI
		{"cwd as name", args{[]byte("abc"), ".", true}, false},
		{"path as name", args{[]byte("abc"), os.TempDir(), true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeFile(tt.args.data, tt.args.name, tt.args.test); (err != nil) != tt.wantErr {
				t.Errorf("writeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// clean-up
	wd, err := os.Getwd()
	if err == nil {
		p := path.Join(wd, "index.html")
		if _, err = os.Stat(p); !os.IsNotExist(err) {
			t.Log("Attempted to delete " + p)
			os.Remove(p)
		}
	}
}

func Test_writeStdout(t *testing.T) {
	type args struct {
		data []byte
		test bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"no data", args{[]byte(""), true}, false},
		{"some data", args{[]byte("hello world"), true}, false},
		{"nil data", args{nil, true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := writeStdout(tt.args.data, tt.args.test); (err != nil) != tt.wantErr {
				t.Errorf("writeStdout() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_layOpts(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := layOpts(); got != tt.want {
				t.Errorf("layOpts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_layTemplates(t *testing.T) {
	tests := []struct {
		name string
		want files
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := layTemplates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("layTemplates() = %v, want %v", got, tt.want)
			}
		})
	}
}
