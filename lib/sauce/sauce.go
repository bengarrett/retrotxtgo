// Package sauce parses SAUCE (Standard Architecture for Universal Comment Extensions) metadata.
// http://www.acid.org/info/sauce/sauce.htm
package sauce

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

const (
	invalid       = "invalid value"
	noPref        = "no preference"
	sauceID       = "SAUCE00"
	comntID       = "COMNT"
	comntLineSize = 64
	comntMaxLines = 255
)

// Record is the container for SAUCE data.
type Record struct {
	ID       string    `json:"id" xml:"id,attr"`
	Version  string    `json:"version" xml:"version,attr"`
	Title    string    `json:"title" xml:"title"`
	Author   string    `json:"author" xml:"author"`
	Group    string    `json:"group" xml:"group"`
	Date     Dates     `json:"date" xml:"date"`
	FileSize Sizes     `json:"filesize" xml:"filesize"`
	Data     DataTypes `json:"dataType"  xml:"data_type"`
	File     FileTypes `json:"fileType" xml:"file_type"`
	Info     TypeInfos `json:"typeInfo"  xml:"type_info"`
	Desc     string    `json:"-" xml:"-"`
	Comnt    Comments  `json:"comments" xml:"comments"`
}

// Dates in multiple output formats.
type Dates struct {
	Value string    `json:"value" xml:"value"`
	Time  time.Time `json:"iso" xml:"date"`
	Epoch int64     `json:"epoch" xml:"epoch,attr"`
}

// Sizes of the file data in multiples.
type Sizes struct {
	Bytes   uint16 `json:"bytes" xml:"bytes"`
	Decimal string `json:"decimal" xml:"decimal,attr"`
	Binary  string `json:"binary" xml:"binary,attr"`
}

// DataTypes both the SAUCE DataType value and name.
type DataTypes struct {
	Type DataType `json:"type" xml:"type"`
	Name string   `json:"name" xml:"name"`
}

// DataType is the data type (SAUCE DataType).
type DataType uint

const (
	none DataType = iota
	character
	bitmap
	vector
	audio
	binaryText
	xBin
	archive
	executable
)

func (d DataType) String() string {
	s := [...]string{
		"undefined", "text or character stream", "bitmap graphic or animation", "vector graphic",
		"audio or music", "binary text", "extended binary text", "archive", "executable",
	}[d]
	return s
}

// FileTypes, both the SAUCE FileType value and name.
type FileTypes struct {
	Type FileType `json:"type" xml:"type"`
	Name string   `json:"name" xml:"name"`
}

// FileType is the type of file (SAUCE FileType).
type FileType uint

// TypeInfos includes the SAUCE fields dependant on both DataType and FileType.
type TypeInfos struct {
	Info1 TypeInfo  `json:"1" xml:"1"`
	Info2 TypeInfo  `json:"2" xml:"2"`
	Info3 TypeInfo  `json:"3" xml:"3"`
	Flags ANSIFlags `json:"flags" xml:"flags"`
	Font  string    `json:"fontName" xml:"fontname"`
}

// TypeInfo includes the SAUCE TInfo value and meaning.
type TypeInfo struct {
	Value uint16 `json:"value" xml:"value"`
	Info  string `json:"info" xml:"info,attr"`
}

// ANSIFlags are the interpretation of the SAUCE Flags field.
type ANSIFlags struct {
	Decimal         Flags      `json:"decimal" xml:"decimal,attr"`
	Binary          string     `json:"binary" xml:"binary,attr"`
	B               ANSIFlagB  `json:"nonBlinkMode" xml:"non_blink_mode"`
	LS              ANSIFlagLS `json:"letterSpacing" xml:"letter_spacing"`
	AR              ANSIFlagAR `json:"aspectRatio" xml:"aspect_ratio"`
	Interpretations string     `json:"-" xml:"-"`
}

func (a *ANSIFlags) String() string {
	if a.Decimal == 0 {
		return ""
	}
	b, ls, ar := a.B.Info, a.LS.Info, a.AR.Info
	l := []string{}
	if b != noPref {
		l = append(l, b)
	}
	if ls != noPref {
		l = append(l, ls)
	}
	if ar != noPref {
		l = append(l, ar)
	}
	if strings.TrimSpace(strings.Join(l, "")) == "" {
		return ""
	}
	return strings.Join(l, ", ")
}

// Flags is the SAUCE Flags field.
type Flags uint8

func (f Flags) parse() ANSIFlags {
	const binary5Bits, minLen = "%05b", 6
	bin := fmt.Sprintf(binary5Bits, f)
	r := []rune(bin)
	if len(r) < minLen {
		return ANSIFlags{
			Decimal: f,
			Binary:  bin,
		}
	}
	b, ls, ar := string(r[0]), string(r[1:3]), string(r[3:5])
	return ANSIFlags{
		Decimal: f,
		Binary:  bin,
		B:       ANSIFlagB{Flag: bBit(b), Info: bBit(b).String()},
		LS:      ANSIFlagLS{Flag: lsBit(ls), Info: lsBit(ls).String()},
		AR:      ANSIFlagAR{Flag: arBit(ar), Info: arBit(ar).String()},
	}
}

// ANSIFlagB is the interpretation of the SAUCE Flags non-blink mode binary bit.
type ANSIFlagB struct {
	Flag bBit   `json:"flag" xml:"flag"`
	Info string `json:"interpretation" xml:"interpretation,attr"`
}

type bBit string

func (b bBit) String() string {
	const blink, non = "0", "1"
	switch b {
	case blink:
		return "blink mode"
	case non:
		return "non-blink mode"
	default:
		return invalid
	}
}

// ANSIFlagLS is the interpretation of the SAUCE Flags letter spacing binary bits.
type ANSIFlagLS struct {
	Flag lsBit  `json:"flag" xml:"flag"`
	Info string `json:"interpretation" xml:"interpretation,attr"`
}

type lsBit string

func (ls lsBit) String() string {
	const none, eight, nine = "00", "01", "10"
	switch ls {
	case none:
		return noPref
	case eight:
		return "select 8 pixel font"
	case nine:
		return "select 9 pixel font"
	default:
		return invalid
	}
}

// ANSIFlagAR is the interpretation of the SAUCE Flags aspect ratio binary bits.
type ANSIFlagAR struct {
	Flag arBit  `json:"flag" xml:"flag"`
	Info string `json:"interpretation" xml:"interpretation,attr"`
}

type arBit string

func (ar arBit) String() string {
	const none, strect, square = "00", "01", "10"
	switch ar {
	case none:
		return noPref
	case strect:
		return "stretch pixels"
	case square:
		return "square pixels"
	default:
		return invalid
	}
}

// Comments contain the optional SAUCE comment block.
type Comments struct {
	ID      string   `json:"id" xml:"id,attr"`
	Count   int      `json:"count" xml:"count,attr"`
	Comment []string `json:"lines" xml:"line"`
}

// Character based files more commonly referred as text files.
type Character uint

const (
	ascii Character = iota
	ansi
	ansiMation
	ripScript
	pcBoard
	avatar
	html
	source
	tundraDraw
)

func (c Character) String() string {
	return [...]string{
		"ASCII text",
		"ANSI color text",
		"ANSIMation",
		"RIPScript",
		"PCBoard color text",
		"Avatar color text",
		"HTML markup",
		"Programming source code",
		"TundraDraw color text",
	}[c]
}

// Desc is the character description.
func (c Character) Desc() string {
	return [...]string{
		"ASCII text file with no formatting codes or color codes.",
		"ANSI text file with coloring codes and cursor positioning.",
		"ANSIMation are ANSI text files that rely on fixed screen sizes.",
		"RIPScript are Remote Imaging Protocol graphics.",
		"PCBoard color codes and macros, and ANSI codes.",
		"Avatar color codes, and ANSi codes.",
		"HTML markup files.",
		"Source code for a programming language.",
		"TundraDraw files, like ANSI, but with a custom palette.",
	}[c]
}

// Bitmap graphic and animation files.
type Bitmap uint

const (
	gif Bitmap = iota
	pcx
	lbm
	tga
	fli
	flc
	bmp
	gl
	dl
	wpg
	png
	jpg
	mpg
	avi
)

func (b Bitmap) String() string {
	return [...]string{
		"GIF image",
		"ZSoft Paintbrush image",
		"DeluxePaint image",
		"Targa true color image",
		"Autodesk Animator animation",
		"Autodesk Animator animation",
		"BMP Windows/OS2 bitmap",
		"Grasp GL animation",
		"DL animation",
		"WordPerfect graphic",
		"PNG image",
		"Jpeg photo",
		"MPEG video",
		"AVI video",
	}[b]
}

// Vector graphic files.
type Vector uint

const (
	dxf Vector = iota
	dwg
	wpvg
	kinetix
)

func (v Vector) String() string {
	return [...]string{
		"AutoDesk CAD vector graphic",
		"AutoDesk CAD vector graphic",
		"WordPerfect vector graphic",
		"3D Studio vector graphic",
	}[v]
}

// Audio or music files.
type Audio uint

const (
	mod Audio = iota
	composer669
	stm
	s3m
	mtm
	far
	ult
	amf
	dmf
	okt
	rol
	cmf
	midi
	sadt
	voc
	wave
	smp8
	smp8s
	smp16
	smp16s
	patch8
	patch16
	xm
	hsc
	it
)

func (a Audio) String() string {
	return [...]string{
		"NoiseTracker module",
		"Composer 669 module",
		"ScreamTracker module",
		"ScreamTracker 3 module",
		"MultiTracker module",
		"Farandole Composer module",
		"Ultra Tracker module",
		"Dual Module Player module",
		"X-Tracker module",
		"Oktalyzer module",
		"AdLib Visual Composer FM audio",
		"Creative Music FM audio",
		"MIDI audio",
		"SAdT composer FM audio",
		"Creative Voice File",
		"Waveform audio",
		"single channel 8-bit sample",
		"stereo 8-bit sample",
		"single channel 16-bit sample",
		"stereo 16-bit sample",
		"8-bit patch file",
		"16-bit patch file",
		"Extended Module",
		"Hannes Seifert Composition FM audio",
		"Impulse Tracker module",
	}[a]
}

// BinaryText is a raw memory copy of a text mode screen.
type BinaryText uint

func (b BinaryText) String() string {
	return "Binary text or a .BIN file"
}

// XBin or eXtended BinaryText files.
type XBin uint

func (x XBin) String() string {
	return "Extended binary text or a XBin file"
}

// Archive and compressed files.
type Archive uint

const (
	zip Archive = iota
	arj
	lzh
	arc
	tar
	zoo
	rar
	uc2
	pak
	sqz
)

func (a Archive) String() string {
	return [...]string{
		"ZIP compressed archive",
		"ARJ compressed archive",
		"LHA compressed archive",
		"ARC compressed archive",
		"Tarball tape archive",
		"ZOO compressed archive",
		"RAR compressed archive",
		"UltraCompressor II compressed archive",
		"PAK ARC compressed archive",
		"Squeeze It compressed archive",
	}[a]
}

// Executable program files.
type Executable uint

func (e Executable) String() string {
	return "Executable program file"
}

// Scan returns the position of the SAUCE00 ID or -1 if no ID exists.
func Scan(b ...byte) (index int) {
	const sauceSize, maximum = 128, 512
	id, l := []byte(sauceID), len(b)
	backwardsLoop := func(i int) int {
		return l - 1 - i
	}
	// search for the id sequence in b
	const indexEnd = 6
	for i := range b {
		if i > maximum {
			break
		}
		i = backwardsLoop(i)
		if i < sauceSize {
			break
		}
		// do matching in reverse
		if b[i] != id[indexEnd] {
			continue // 0
		}
		if b[i-1] != id[5] {
			continue // 0
		}
		if b[i-2] != id[4] {
			continue // E
		}
		if b[i-3] != id[3] {
			continue // C
		}
		if b[i-4] != id[2] {
			continue // U
		}
		if b[i-5] != id[1] {
			continue // A
		}
		if b[i-indexEnd] != id[0] {
			continue // S
		}
		return i - indexEnd
	}
	return -1
}

// Parse and extract the record data.
func Parse(data ...byte) Record {
	const empty = "\x00\x00"
	r := record(data)
	d := r.extract()
	if string(d.version[:]) == empty {
		return Record{}
	}
	return Record{
		ID:       d.id.String(),
		Version:  d.version.String(),
		Title:    strings.TrimSpace(d.title.String()),
		Author:   strings.TrimSpace(d.author.String()),
		Group:    strings.TrimSpace(d.group.String()),
		Date:     d.dates(),
		FileSize: d.sizes(),
		Data:     d.dataType(),
		File:     d.fileType(),
		Info:     d.typeInfo(),
		Desc:     d.description(),
		Comnt:    d.commentBlock(),
	}
}

func unsignedBinary1(b [1]byte) (value uint8) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsigned 1 byte, LE binary failed:", err)
	}
	return value
}

func unsignedBinary2(b [2]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsigned 2 bytes, LE binary failed:", err)
	}
	return value
}

func unsignedBinary4(b [4]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsigned 4 byte, LE binary failed:", err)
	}
	return value
}
