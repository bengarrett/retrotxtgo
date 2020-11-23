// Package sauce to handle the reading and parsing of embedded SAUCE metadata.

package sauce

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"retrotxt.com/retrotxt/internal/pack"
)

const commentResult = "Any comments go here.                                           "

var (
	rawData     = pack.Get("text/sauce.txt")
	packData    = record(rawData)
	exampleData = packData.extract()
	sauceIndex  = Scan(rawData...)
)

func TestDataType_String(t *testing.T) {
	tests := []struct {
		name string
		d    DataType
		want string
	}{
		{"none", 0, "undefined"},
		{"executable", 8, "executable"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.String(); got != tt.want {
				t.Errorf("DataType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lsBit_String(t *testing.T) {
	tests := []struct {
		name string
		ls   lsBit
		want string
	}{
		{"empty", "", invalid},
		{"00", "00", noPref},
		{"8px", "01", "select 8 pixel font"},
		{"9px", "10", "select 9 pixel font"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ls.String(); got != tt.want {
				t.Errorf("lsBit.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_String(t *testing.T) {
	tests := []struct {
		name string
		c    Character
		want string
	}{
		{"first", ascii, "ASCII text"},
		{"last", tundraDraw, "TundraDraw color text"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.String(); got != tt.want {
				t.Errorf("Character.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCharacter_Desc(t *testing.T) {
	tests := []struct {
		name string
		c    Character
	}{
		{"first", ascii},
		{"last", tundraDraw},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.Desc(); got == "" {
				t.Errorf("Character.Desc() = %q", got)
			}
		})
	}
}

func TestBitmap_String(t *testing.T) {
	tests := []struct {
		name string
		b    Bitmap
		want string
	}{
		{"first", gif, "GIF image"},
		{"last", avi, "AVI video"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.String(); got != tt.want {
				t.Errorf("Bitmap.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVector_String(t *testing.T) {
	tests := []struct {
		name string
		v    Vector
		want string
	}{
		{"first", dxf, "AutoDesk CAD vector graphic"},
		{"last", kinetix, "3D Studio vector graphic"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.v.String(); got != tt.want {
				t.Errorf("Vector.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAudio_String(t *testing.T) {
	tests := []struct {
		name string
		a    Audio
		want string
	}{
		{"first", mod, "NoiseTracker module"},
		{"midi", midi, "MIDI audio"},
		{"okt", okt, "Oktalyzer module"},
		{"last", it, "Impulse Tracker module"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Audio.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_String(t *testing.T) {
	tests := []struct {
		name string
		a    Archive
		want string
	}{
		{"zip", zip, "ZIP compressed archive"},
		{"squeeze", sqz, "Squeeze It compressed archive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.String(); got != tt.want {
				t.Errorf("Archive.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_comment(t *testing.T) {
	tests := []struct {
		name string
		data data
		want []string
	}{
		{"empty", data{}, nil},
		{"example", exampleData, []string{commentResult}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &comnt{
				lines: tt.data.comnt.lines,
			}
			d := &data{
				comments: tt.data.comments,
				comnt:    *c,
			}
			if gotC := d.commentBlock(); !reflect.DeepEqual(gotC.Comment, tt.want) {
				t.Errorf("data.commentBlock() = %v, want %v", gotC.Comment, tt.want)
			}
		})
	}
}

func Test_commentByBreak(t *testing.T) {
	tests := []struct {
		name      string
		b         []byte
		wantLines []string
	}{
		{"empty", []byte{}, nil},
		{"example", exampleData.comnt.lines, []string{commentResult}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := commentByBreak(tt.b); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("commentByBreak() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func Test_commentByLine(t *testing.T) {
	tests := []struct {
		name      string
		b         []byte
		wantLines []string
	}{
		{"empty", []byte{}, nil},
		{"example", exampleData.comnt.lines, []string{commentResult}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLines := commentByLine(tt.b); !reflect.DeepEqual(gotLines, tt.wantLines) {
				t.Errorf("commentByLine() = %v, want %v", gotLines, tt.wantLines)
			}
		})
	}
}

func Test_data_dates(t *testing.T) {
	tests := []struct {
		name string
		date date
		want Dates
	}{
		{"example", exampleData.date, Dates{
			Value: "20161126",
			Time:  time.Date(2016, 11, 26, 0, 0, 0, 0, time.UTC),
			Epoch: 1480118400,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				date: tt.date,
			}
			if got := d.dates(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data.dates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_dataType(t *testing.T) {
	tests := []struct {
		name     string
		datatype dataType
		want     DataTypes
	}{
		{"none", [1]byte{0}, DataTypes{Type: none, Name: none.String()}},
		{"archive", [1]byte{7}, DataTypes{Type: archive, Name: archive.String()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				datatype: tt.datatype,
			}
			if got := d.dataType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data.dataType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_description(t *testing.T) {
	type fields struct {
		datatype dataType
		filetype fileType
	}
	tests := []struct {
		name   string
		fields fields
		wantS  string
	}{
		{"none", fields{[1]byte{0}, [1]byte{0}}, ""},
		{"pc board", fields{[1]byte{1}, [1]byte{4}}, pcBoard.Desc()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				datatype: tt.fields.datatype,
				filetype: tt.fields.filetype,
			}
			if gotS := d.description(); gotS != tt.wantS {
				t.Errorf("data.description() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func Test_data_fileType(t *testing.T) {
	type fields struct {
		datatype dataType
		filetype fileType
	}
	tests := []struct {
		name   string
		fields fields
		want   FileTypes
	}{
		{"none", fields{[1]byte{0}, [1]byte{0}}, FileTypes{FileType(none), none.String()}},
		{"audio", fields{[1]byte{4}, [1]byte{0}}, FileTypes{FileType(mod), mod.String()}},
		{"executable", fields{[1]byte{8}, [1]byte{0}}, FileTypes{FileType(executable), executable.String()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				datatype: tt.fields.datatype,
				filetype: tt.fields.filetype,
			}
			if got := d.fileType(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data.fileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_sizes(t *testing.T) {
	tests := []struct {
		name     string
		filesize fileSize
		want     Sizes
	}{
		{"none", fileSize([4]byte{}), Sizes{0, "0", "0"}},
		{"1 byte", fileSize([4]byte{1}), Sizes{1, "1B", "1B"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				filesize: tt.filesize,
			}
			if got := d.sizes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("data.sizes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_data_typeInfo(t *testing.T) {
	type fields struct {
		datatype dataType
		filetype fileType
		tinfo1   tInfo1
		tinfo2   tInfo2
		tinfo3   tInfo3
		tFlags   tFlags
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, ""},
		{"ascii", fields{datatype: [1]byte{1}, filetype: [1]byte{0}}, "character width"},
		{"rip script", fields{datatype: [1]byte{1}, filetype: [1]byte{3}}, "pixel width"},
		{"smp16s", fields{datatype: [1]byte{4}, filetype: [1]byte{19}}, "sample rate"},
		{"binary text", fields{datatype: [1]byte{5}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &data{
				datatype: tt.fields.datatype,
				filetype: tt.fields.filetype,
				tinfo1:   tt.fields.tinfo1,
				tinfo2:   tt.fields.tinfo2,
				tinfo3:   tt.fields.tinfo3,
				tFlags:   tt.fields.tFlags,
			}
			if got := d.typeInfo(); !reflect.DeepEqual(got.Info1.Info, tt.want) {
				t.Errorf("data.typeInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_record_author(t *testing.T) {
	tests := []struct {
		name string
		r    record
		i    int
		want string
	}{
		{"example", packData, sauceIndex, "Sauce author        "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.author(tt.i); !reflect.DeepEqual(fmt.Sprintf("%s", got), tt.want) {
				t.Errorf("record.author() = %v, want %v", fmt.Sprintf("%q", got), tt.want)
			}
		})
	}
}

func Test_record_comnt(t *testing.T) {
	type args struct {
		count      comments
		sauceIndex int
	}
	tests := []struct {
		name       string
		r          record
		args       args
		wantLength int
	}{
		{"example", packData, args{count: [1]byte{1}, sauceIndex: sauceIndex}, 1 * comntLineSize},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotBlock := tt.r.comnt(tt.args.count, tt.args.sauceIndex); !reflect.DeepEqual(gotBlock.length, tt.wantLength) {
				t.Errorf("record.comnt() = %v, want %v", gotBlock.length, tt.wantLength)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want string
	}{
		{"empty", []byte{}, ""},
		{"example", rawData, "Sauce title"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.data...); !reflect.DeepEqual(got.Title, tt.want) {
				t.Errorf("Parse() = %v, want %v", got.Title, tt.want)
			}
		})
	}
}
