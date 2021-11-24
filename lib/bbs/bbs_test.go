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
		{"pcboard with nulls", args{"hello\n\n@X01world"}, PCBoard},
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

func TestBBS_HTML(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		bbs     BBS
		args    args
		want    string
		wantErr bool
	}{
		{"empty", -1, args{}, "", true},
		{"plaintext", -1, args{"text"}, "", true},
		{"plaintext", ANSI, args{"\x27\x91text"}, "", true},
		{"celerity", Celerity, args{"|S|gHello|Rworld"},
			"<i class=\"PBg,PFw\">Hello</i><i class=\"PBR,PFw\">world</i>", false},
	}
	for _, tt := range tests {
		got, err := tt.bbs.HTML(tt.args.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("BBS.HTML() %v error = %v, wantErr %v", tt.name, err, tt.wantErr)
			return
		}
		if got.String() != tt.want {
			t.Errorf("BBS.HTML() = %v, want %v", got, tt.want)
		}
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
		{"false positive s", args{[]byte("Hello |sworld")}, -1},
		{"cel B", args{[]byte("Hello |Bworld")}, Celerity},
		{"cel W", args{[]byte("Hello world\n|WThis is a newline.")}, Celerity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindCelerity(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindCelerity() = %v, want %v", got, tt.want)
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
			if got := FindRenegade(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindRenegade() = %v, want %v", got, tt.want)
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
			if got := FindPCBoard(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindPCBoard() = %v, want %v", got, tt.want)
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
			if got := FindWildcat(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWildcat() = %v, want %v", got, tt.want)
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
			if got := FindWWIVHeart(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWWIVHeart() = %v, want %v", got, tt.want)
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
			if got := FindWWIVHash(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWWIVHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserBar(t *testing.T) {
	type args struct {
		s string
	}
	const black, white, red = "0", "7", "20"
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"|" + black + white + "Hello world"}, "<i class=\"P0,P7\">Hello world</i>", false},
		{"multi", args{"|" + black + white + "White |" + red + "Red Background"},
			"<i class=\"P0,P7\">White </i><i class=\"P20,P7\">Red Background</i>", false},
		{"newline", args{"|07White\n|20Red Background"},
			"<i class=\"P0,P7\">White\n</i><i class=\"P20,P7\">Red Background</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parserBar(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserBar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parserBar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parserCelerity(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"string", args{"the quick brown fox"}, "the quick brown fox", false},
		{"prefix", args{"|kHello world"}, "<i class=\"PBk,PFk\">Hello world</i>", false},
		{"background", args{"|S|bHello world"}, "<i class=\"PBb,PFw\">Hello world</i>", false},
		{"multi", args{"|S|gHello|Rworld"}, "<i class=\"PBg,PFw\">Hello</i><i class=\"PBR,PFw\">world</i>", false},
		{"newline", args{"|S|gHello\n|Rworld"}, "<i class=\"PBg,PFw\">Hello\n</i><i class=\"PBR,PFw\">world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parserCelerity(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserCelerity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parserCelerity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parsePCBoard(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"@X07Hello world"}, "<i class=\"PB0,PF7\">Hello world</i>", false},
		{"multi", args{"@X07Hello @X11world"}, "<i class=\"PB0,PF7\">Hello </i><i class=\"PB1,PF1\">world</i>", false},
		{"newline", args{"@X07Hello\n@X11world"}, "<i class=\"PB0,PF7\">Hello\n</i><i class=\"PB1,PF1\">world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePCBoard(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePCBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parsePCBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTelegard(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"`07Hello world"}, "<i class=\"PB0,PF7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTelegard(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTelegard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parseTelegard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseWHash(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"|#7Hello world"}, "<i class=\"P0,P7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWHash(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parseWHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseWHeart(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"\x037Hello world"}, "<i class=\"P0,P7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseWHeart(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWHeart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("parseWHeart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseWildcat(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"empty", args{}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"@0F@Hello world"}, "<i class=\"PB0,PFF\">Hello world</i>", false},
	}
	for _, tt := range tests {
		got, err := parseWildcat(tt.args.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("parseWildcat() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got.String() != tt.want {
			t.Errorf("parseWildcat() = %v, want %v", got, tt.want)
		}
	}
}

func Test_validateC(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"invalid Z", args{byte('Z')}, false},
		{"invalid 0", args{byte('0')}, false},
		{"normal black", args{byte('k')}, true},
		{"swap", args{byte('S')}, true},
		{"case test", args{byte('s')}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateC(tt.args.b); got != tt.want {
				t.Errorf("validateC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateP(t *testing.T) {
	type args struct {
		b byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"z", args{byte('z')}, false},
		{"Z", args{byte('Z')}, false},
		{"G", args{byte('Z')}, false},
		{"0", args{byte('0')}, true},
		{"9", args{byte('9')}, true},
		{"F", args{byte('F')}, true},
		{"f", args{byte('f')}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateP(tt.args.b); got != tt.want {
				t.Errorf("validateP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateR(t *testing.T) {
	type args struct {
		b [2]byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty", args{}, false},
		{"00", args{[2]byte{byte('0'), byte('0')}}, true},
		{"23", args{[2]byte{byte('2'), byte('3')}}, true},
		{"out of range a", args{[2]byte{byte('2'), byte('4')}}, false},
		{"out of range b", args{[2]byte{byte('3'), byte('0')}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateR(tt.args.b); got != tt.want {
				t.Errorf("validateR() = %v, want %v", got, tt.want)
			}
		})
	}
}
