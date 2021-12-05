package util_test

import (
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/filesystem/internal/util"
)

func TestWindows(t *testing.T) {
	type args struct {
		i   int
		p   string
		os  string
		dir string
	}
	tests := []struct {
		name     string
		args     args
		wantS    string
		wantCont bool
	}{
		{"home", args{1, "", "linux", "/home/retro"}, "/home/retro", false},
		{"c drive", args{0, "c:", "windows", ""}, "C:\\", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, gotCont := util.Windows(tt.args.i, tt.args.p, tt.args.os, tt.args.dir)
			if gotS != tt.wantS {
				t.Errorf("Windows() gotS = %v, want %v", gotS, tt.wantS)
			}
			if gotCont != tt.wantCont {
				t.Errorf("Windows() gotCont = %v, want %v", gotCont, tt.wantCont)
			}
		})
	}
}
