package cmd

import (
	"path/filepath"
	"testing"
)

// move to cli
// func TestDetail_infoSwitch(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		f       Detail
// 		format  string
// 		wantErr bool
// 	}{
// 		{"empty", Detail{}, "", true},
// 		{"invalid", Detail{}, "invalid", true},
// 		{"color", Detail{}, "color", false},
// 		{"j", Detail{}, "j", false},
// 		{"jm", Detail{}, "jm", false},
// 		{"text", Detail{}, "text", false},
// 		{"x", Detail{}, "x", false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if gotErr := tt.f.infoSwitch(tt.format); (gotErr.Msg != nil) != tt.wantErr {
// 				t.Errorf("parse() error = %v, wantErr %v", gotErr, tt.wantErr)
// 			}
// 		})
// 	}
// }

// Todo: move this to a test library that generates this file in a temp directory.
// temp file will hardcode os.stat info and also remove itself after use.
var hi = filepath.Clean("../../textfiles/hi.txt")

func Test_infoPrint(t *testing.T) {
	type args struct {
		filename string
		format   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty", args{"", ""}, true},
		{"file not exist", args{"notexistingfile", "json"}, true},
		{"invalid fmt", args{hi, "yaml"}, true},
		{"color", args{hi, ""}, false},
		{"json", args{hi, "json"}, false},
		{"json.min", args{hi, "jm"}, false},
		{"text", args{hi, "t"}, false},
		{"xml", args{hi, "xml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := infoPrint(tt.args.filename, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("infoPrint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
