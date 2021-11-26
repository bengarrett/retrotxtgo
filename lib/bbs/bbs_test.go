package bbs_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/bengarrett/retrotxtgo/lib/bbs"
)

const (
	ansiEsc = "\x1B\x5B"
)

func TestBBS_String(t *testing.T) {
	tests := []struct {
		name string
		b    bbs.BBS
		want string
	}{
		{"too small", -1, ""},
		{"too big", 111, ""},
		{"first", bbs.Celerity, "Celerity |"},
		{"last", bbs.WWIVHeart, "WWIV ♥"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("BBS.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBBS_Name(t *testing.T) {
	tests := []struct {
		name string
		b    bbs.BBS
		want string
	}{
		{"too small", -1, ""},
		{"too big", 111, ""},
		{"first", bbs.Celerity, "Celerity"},
		{"last", bbs.WWIVHeart, "WWIV ♥"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Name(); got != tt.want {
				t.Errorf("BBS.Name() = %v, want %v", got, tt.want)
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
		want bbs.BBS
	}{
		{"empty", args{""}, -1},
		{"ansi", args{ansiEsc + "0;"}, bbs.ANSI},
		{"cls", args{"@CLS@"}, -1},
		{"pcb+ans", args{"@CLS@" + ansiEsc + "0;"}, bbs.ANSI},
		{"pcb+ans", args{"@CLS@Hello world\nThis is some text." + ansiEsc + "0;"}, bbs.ANSI},
		{"celerity", args{"Hello world\n|WThis is a newline."}, bbs.Celerity},
		{"renegade", args{"Hello world\n|09This is a newline."}, bbs.Renegade},
		{"pcboard", args{"Hello world\n@X01This is a newline."}, bbs.PCBoard},
		{"telegard", args{"Hello world\n`09This is a newline."}, bbs.Telegard},
		{"wildcat", args{"Hello world\n@01@This is a newline."}, bbs.Wildcat},
		{"wwiv ♥", args{"Hello world\n\x031This is a newline."}, bbs.WWIVHeart},
		{"pcboard with nulls", args{"hello\n\n@X01world"}, bbs.PCBoard},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.args.s)
			if got := bbs.Find(r); !reflect.DeepEqual(got, tt.want) {
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
		bbs     bbs.BBS
		args    args
		want    string
		wantErr bool
	}{
		{"empty", -1, args{}, "", true},
		{"plaintext", -1, args{"text"}, "", true},
		{"plaintext", bbs.ANSI, args{"\x27\x91text"}, "", true},
		{"celerity", bbs.Celerity, args{"|S|gHello|Rworld"},
			"<i class=\"PBg PFw\">Hello</i><i class=\"PBR PFw\">world</i>", false},
		{"xss", bbs.Celerity, args{"|S|gABC<script>alert('xss');</script>D|REF"},
			"<i class=\"PBg PFw\">ABC&lt;script&gt;alert(&#39;xss&#39;);&lt;/script&gt;D</i><i class=\"PBR PFw\">EF</i>", false},
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
		want bbs.BBS
	}{
		{"empty", args{[]byte{}}, -1},
		{"ansi", args{[]byte(ansiEsc + "0;")}, -1},
		{"false positive z", args{[]byte("Hello |Zworld")}, -1},
		{"false positive s", args{[]byte("Hello |sworld")}, -1},
		{"cel B", args{[]byte("Hello |Bworld")}, bbs.Celerity},
		{"cel W", args{[]byte("Hello world\n|WThis is a newline.")}, bbs.Celerity},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindCelerity(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		want bbs.BBS
	}{
		{"empty", args{nil}, -1},
		{"celerity", args{[]byte("Hello |Bworld")}, -1},
		{"first", args{[]byte("|00")}, bbs.Renegade},
		{"end", args{[]byte("|23")}, bbs.Renegade},
		{"newline", args{[]byte("Hello world\n|15This is a newline.")}, bbs.Renegade},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindRenegade(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		want bbs.BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("@XHello world")}, -1},
		{"incomplete", args{[]byte("@X0Hello world")}, -1},
		{"out of range", args{[]byte("@X0GHello world")}, -1},
		{"first", args{[]byte("@X00Hello world")}, bbs.PCBoard},
		{"end", args{[]byte("@XFFHello world")}, bbs.PCBoard},
		{"newline", args{[]byte("Hello world\n@X00This is a newline.")}, bbs.PCBoard},
		{"false pos", args{[]byte("PCBoard @X code")}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindPCBoard(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		want bbs.BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("@Hello world")}, -1},
		{"incomplete", args{[]byte("@0Hello world")}, -1},
		{"out of range", args{[]byte("@0@GHello world")}, -1},
		{"first", args{[]byte("@00@Hello world")}, bbs.Wildcat},
		{"end", args{[]byte("@FF@Hello world")}, bbs.Wildcat},
		{"newline", args{[]byte("Hello world\n@00@This is a newline.")}, bbs.Wildcat},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindWildcat(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		want bbs.BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed", args{[]byte("\x03Hello world")}, -1},
		{"first", args{[]byte("\x030Hello world")}, bbs.WWIVHeart},
		{"last", args{[]byte("\x039Hello world")}, bbs.WWIVHeart},
		{"lots of numbers", args{[]byte("\x0398765 Hello world")}, bbs.WWIVHeart},
		{"newline", args{[]byte("Hello world\n\x031This is a newline.")}, bbs.WWIVHeart},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindWWIVHeart(tt.args.b); !reflect.DeepEqual(got, tt.want) {
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
		want bbs.BBS
	}{
		{"empty", args{nil}, -1},
		{"malformed |#", args{[]byte("|#Hello world")}, -1},
		{"malformed |0", args{[]byte("|0Hello world")}, -1},
		{"malformed #0", args{[]byte("#0Hello world")}, -1},
		{"first", args{[]byte("|#0Hello world")}, bbs.WWIVHash},
		{"last", args{[]byte("|#9Hello world")}, bbs.WWIVHash},
		{"lots of numbers", args{[]byte("|#98765 Hello world")}, bbs.WWIVHash},
		{"newline", args{[]byte("Hello world\n|#1This is a newline.")}, bbs.WWIVHash},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bbs.FindWWIVHash(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWWIVHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SplitBars(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty", args{""}, 0},
		{"first", args{"|00"}, 1},
		{"last", args{"|23"}, 1},
		{"out of range", args{"|24"}, 0},
		{"incomplete", args{"|2"}, 0},
		{"multiples", args{"|01Hello|00 |10world"}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(bbs.SplitBars(tt.args.s)); got != tt.want {
				fmt.Println(bbs.SplitBars(tt.args.s))
				t.Errorf("SplitBars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParserBars(t *testing.T) {
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
		{"false pos", args{"hello|world"}, "hello|world", false},
		{"false pos double", args{"| hello world |"}, "| hello world |", false},
		{"prefix", args{"|" + black + white + "Hello world"}, "<i class=\"P0 P7\">Hello world</i>", false},
		{"multi", args{"|" + black + white + "White |" + red + "Red Background"},
			"<i class=\"P0 P7\">White </i><i class=\"P20 P7\">Red Background</i>", false},
		{"newline", args{"|07White\n|20Red Background"},
			"<i class=\"P0 P7\">White\n</i><i class=\"P20 P7\">Red Background</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParserBars(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParserBars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParserBars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParserCelerity(t *testing.T) {
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
		{"prefix", args{"|kHello world"}, "<i class=\"PBk PFk\">Hello world</i>", false},
		{"background", args{"|S|bHello world"},
			"<i class=\"PBb PFw\">Hello world</i>", false},
		{"multi", args{"|S|gHello|Rworld"},
			"<i class=\"PBg PFw\">Hello</i><i class=\"PBR PFw\">world</i>", false},
		{"newline", args{"|S|gHello\n|Rworld"},
			"<i class=\"PBg PFw\">Hello\n</i><i class=\"PBR PFw\">world</i>", false},
		{"false positive", args{"| Hello world |"}, "| Hello world |", false},
		{"double bar", args{"||pipes"}, "||pipes", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParserCelerity(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParserCelerity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParserCelerity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SplitPCBoard(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty", args{""}, 0},
		{"first", args{"@X00"}, 1},
		{"last", args{"@XFF"}, 1},
		{"out of range", args{"@XFG"}, 0},
		{"incomplete", args{"@X0"}, 0},
		{"multiples", args{"@X01Hello@X00 @X10world"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := len(bbs.SplitPCBoard(tt.args.s)); got != tt.want {
				fmt.Println(bbs.SplitPCBoard(tt.args.s))
				t.Errorf("SplitPCBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParsePCBoard(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args

		want    string
		wantErr bool
	}{
		{"empty", args{""}, "", false},
		{"string", args{"hello world"}, "hello world", false},
		{"prefix", args{"@X07Hello world"}, "<i class=\"PB0 PF7\">Hello world</i>", false},
		{"casing", args{"@xaBHello world"}, "<i class=\"PBA PFB\">Hello world</i>", false},
		{"multi", args{"@X07Hello @X11world"}, "<i class=\"PB0 PF7\">Hello </i><i class=\"PB1 PF1\">world</i>", false},
		{"newline", args{"@X07Hello\n@X11world"}, "<i class=\"PB0 PF7\">Hello\n</i><i class=\"PB1 PF1\">world</i>", false},
		{"false pos 0", args{"@X code for PCBoard"}, "@X code for PCBoard", false},
		{"false pos 1", args{"PCBoard @X code"}, "PCBoard @X code", false},
		{"false pos 2", args{"PCBoard @Xcode"}, "PCBoard @Xcode", false},
		{"false pos 3", args{"Does PCBoard @X code offer a red @X?"}, "Does PCBoard @X code offer a red @X?", false},
		{"combo", args{"@X07@Xcodes combo"}, "<i class=\"PB0 PF7\">@Xcodes combo</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParsePCBoard(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePCBoard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParsePCBoard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseTelegard(t *testing.T) {
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
		{"prefix", args{"`07Hello world"}, "<i class=\"PB0 PF7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParseTelegard(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTelegard() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParseTelegard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseWHash(t *testing.T) {
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
		{"prefix", args{"|#7Hello world"}, "<i class=\"P0 P7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParseWHash(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParseWHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseWHeart(t *testing.T) {
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
		{"prefix", args{"\x037Hello world"}, "<i class=\"P0 P7\">Hello world</i>", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bbs.ParseWHeart(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWHeart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want {
				t.Errorf("ParseWHeart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseWildcat(t *testing.T) {
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
		{"prefix", args{"@0F@Hello world"}, "<i class=\"PB0 PFF\">Hello world</i>", false},
	}
	for _, tt := range tests {
		got, err := bbs.ParseWildcat(tt.args.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseWildcat() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		if got.String() != tt.want {
			t.Errorf("ParseWildcat() = %v, want %v", got, tt.want)
		}
	}
}
