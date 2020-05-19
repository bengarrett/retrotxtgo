package info

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Todo: move this to a test library that generates this file in a temp directory.
// temp file will hardcode os.stat info and also remove itself after use.
var hi = filepath.Clean("../../textfiles/hi.txt")

func Test_File(t *testing.T) {
	var got Detail
	err := got.Read(hi)
	if err != nil {
		t.Errorf("Read() = %v, want %v", err, nil)
	}
	if got.Bytes != 40 {
		t.Errorf("Read() = %v, want %v", got.Bytes, 40)
	}
	if got.Name != "hi.txt" {
		t.Errorf("Read() = %v, want %v", got.Name, "hi.txt")
	}
	if got.Slug != "hi-txt" {
		t.Errorf("Read() = %v, want %v", got.Slug, "hi-txt")
	}
	if got.Mime != "text/plain" {
		t.Errorf("Read() = %v, want %v", got.Mime, "text/plain")
	}
	if got.Utf8 != true {
		t.Errorf("Read() = %v, want %v", got.Utf8, true)
	}
	const want = "1b466b6448d7ff10e2f8f7160d936987"
	if got.MD5 != want {
		t.Errorf("Read() = %v, want %v", got.MD5, want)
	}
}

func Test_parse(t *testing.T) {
	f, err := os.Stat(hi)
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
			var got Detail
			err := got.parse(tt.args.data, tt.args.stat)
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

func Test_JSON(t *testing.T) {
	tests := []struct {
		name   string
		d      Detail
		indent bool
		want   bool
	}{
		{"no indent", Detail{}, false, true},
		{"indent", Detail{}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := json.Valid(tt.d.JSON(tt.indent)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Text(t *testing.T) {
	var d Detail
	err := d.Read(hi)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	want := 491
	if got := len(d.Text(false)); got != want {
		t.Errorf("Text() = %v, want %v", got, want)
	}
}

func ExampleDetail_XML() {
	var file Detail
	file.Read(hi)
	data, _ := file.XML()
	fmt.Printf("%s", data)
	// Output:
	// <file id="hi-txt">
	// 	<name>hi.txt</name>
	// 	<content>
	// 		<mime>text/plain</mime>
	// 		<utf8>true</utf8>
	// 	</content>
	// 	<size>
	// 		<bytes>40</bytes>
	// 		<value>40 bytes</value>
	// 		<character-count>22</character-count>
	// 	</size>
	// 	<checksum>
	// 		<md5>1b466b6448d7ff10e2f8f7160d936987</md5>
	// 		<sha256>3ec92fce657848240c9f9eb6887dbf49a6331ac071759440c41396248ca501fb</sha256>
	// 	</checksum>
	// 	<modified>2019-12-03T23:04:45.1715585+11:00</modified>
	// </file>
}
