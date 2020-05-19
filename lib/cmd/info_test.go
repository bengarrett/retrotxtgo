package cmd

import (
	"log"
	"testing"

	"github.com/bengarrett/retrotxtgo/samples"
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

var sampleFile = func() string {
	path, err := samples.Save([]byte(samples.Tabs), "info_test.txt")
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func Test_infoPrint(t *testing.T) {
	tmp := sampleFile()
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
		{"invalid fmt", args{tmp, "yaml"}, true},
		{"color", args{tmp, ""}, false},
		{"json", args{tmp, "json"}, false},
		{"json.min", args{tmp, "jm"}, false},
		{"text", args{tmp, "t"}, false},
		{"xml", args{tmp, "xml"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := infoPrint(tt.args.filename, tt.args.format); (err != nil) != tt.wantErr {
				t.Errorf("infoPrint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	samples.Clean(tmp)
}
