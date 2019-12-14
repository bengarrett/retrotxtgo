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
	"strings"
	"testing"
)

func Test_createLayouts(t *testing.T) {
	l := strings.Split(createLayouts(), ",")
	if got := len(l); got != 5 {
		t.Errorf("createTemplates() = %v, want %v", got, 5)
	}
	if got := createTemplates()["body"]; got != "body-content" {
		t.Errorf("createTemplates() = %v, want %v", got, "body-content")
	}
}

func Test_createTemplates(t *testing.T) {
	tests := []struct {
		name string
		key  string
		want string
	}{
		{"empty", "", ""},
		{"body", "body", "body-content"},
		{"standard", "standard", "standard"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTemplates()[tt.key]; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTemplates() = %v, want %v", got, tt.want)
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
	w = "RetroTxt | example"
	got = pagedata(d).PageTitle
	if got != w {
		t.Errorf("pagedata().PageTitle = %v, want %v", got, w)
	}
	w = ""
	got = pagedata(d).MetaDesc
	if got != w {
		t.Errorf("pagedata().MetaDesc = %v, want %v", got, w)
	}
	htmlLayout = "standard"
	w = ""
	got = pagedata(d).MetaAuthor
	if got != w {
		t.Errorf("pagedata().MetaAuthor = %v, want %v", got, w)
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
