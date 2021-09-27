package bbs

import (
	"reflect"
	"strings"
	"testing"
)

const (
	ansiEsc = "\x1B\x5B"
)

func TestBBS_String(t *testing.T) {
	tests := []struct {
		name string
		b    BBS
		want string
	}{
		{"too small", -1, ""},
		{"too big", 111, ""},
		{"first", Celerity, "Celerity |"},
		{"last", WWIVHeart, "WWIV ♥"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BBS.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{""}, -1},
		{"ansi", args{ansiEsc + "0;"}, ANSI},
		{"cls", args{"@CLS@"}, -1},
		{"pcb+ans", args{"@CLS@" + ansiEsc + "0;"}, ANSI},
		{"pcb+ans", args{"@CLS@Hello world\nThis is some text." + ansiEsc + "0;"}, ANSI},
		{"celerity", args{"Hello world\n|WThis is a newline."}, Celerity},
		{"renegade", args{"Hello world\n|09This is a newline."}, Renegade},
		{"pcboard", args{"Hello world\n@X01This is a newline."}, PCBoard},
		{"telegard", args{"Hello world\n`09This is a newline."}, Telegard},
		{"wildcat", args{"Hello world\n@01@This is a newline."}, Wildcat},
		{"wwiv ♥", args{"Hello world\n\x031This is a newline."}, WWIVHeart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.args.s)
			if got := Find(r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() = %d, want %d", got, tt.want)
			}
		})
	}
}

func Test_findCelerity(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{[]byte{}}, -1},
		{"ansi", args{[]byte(ansiEsc + "0;")}, -1},
		{"false positive z", args{[]byte("Hello |Zworld")}, -1},
		{"false positive b", args{[]byte("Hello |bworld")}, -1},
		{"cel B", args{[]byte("Hello |Bworld")}, Celerity},
		{"cel W", args{[]byte("Hello world\n|WThis is a newline.")}, Celerity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findCelerity(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findCelerity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findRenegade(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{nil}, -1},
		{"celerity", args{[]byte("Hello |Bworld")}, -1},
		{"first", args{[]byte("|00")}, Renegade},
		{"end", args{[]byte("|23")}, Renegade},
		{"newline", args{[]byte("Hello world\n|15This is a newline.")}, Renegade},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findRenegade(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findRenegade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findPCBoard(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("@XHello world")}, -1},
		{"incomplete", args{[]byte("@X0Hello world")}, -1},
		{"out of range", args{[]byte("@X0GHello world")}, -1},
		{"first", args{[]byte("@X00Hello world")}, PCBoard},
		{"end", args{[]byte("@XFFHello world")}, PCBoard},
		{"newline", args{[]byte("Hello world\n@X00This is a newline.")}, PCBoard},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findPCBoard(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findPCBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findWildcat(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("@Hello world")}, -1},
		{"incomplete", args{[]byte("@0Hello world")}, -1},
		{"out of range", args{[]byte("@0@GHello world")}, -1},
		{"first", args{[]byte("@00@Hello world")}, Wildcat},
		{"end", args{[]byte("@FF@Hello world")}, Wildcat},
		{"newline", args{[]byte("Hello world\n@00@This is a newline.")}, Wildcat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findWildcat(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findWildcat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findWWIVHeart(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("\x03Hello world")}, -1},
		{"first", args{[]byte("\x030Hello world")}, WWIVHeart},
		{"last", args{[]byte("\x039Hello world")}, WWIVHeart},
		{"lots of numbers", args{[]byte("\x0398765 Hello world")}, WWIVHeart},
		{"newline", args{[]byte("Hello world\n\x031This is a newline.")}, WWIVHeart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findWWIVHeart(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findWWIVHeart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findWWIVHash(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name string
		args args
		want BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed |#", args{[]byte("|#Hello world")}, -1},
		{"malformed |0", args{[]byte("|0Hello world")}, -1},
		{"malformed #0", args{[]byte("#0Hello world")}, -1},
		{"first", args{[]byte("|#0Hello world")}, WWIVHash},
		{"last", args{[]byte("|#9Hello world")}, WWIVHash},
		{"lots of numbers", args{[]byte("|#98765 Hello world")}, WWIVHash},
		{"newline", args{[]byte("Hello world\n|#1This is a newline.")}, WWIVHash},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findWWIVHash(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findWWIVHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
