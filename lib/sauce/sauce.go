// Package sauce to handle the opening and reading of text files.
package sauce

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"retrotxt.com/retrotxt/lib/humanize"
)

// TODO: handle comments!!

const (
	sauceID = "SAUCE00"
	invalid = "invalid value"
	noPref  = "no preference"
)

// Record layout for the SAUCE metadata.
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

// DataTypes includes both the SAUCE DataType value and name.
type DataTypes struct {
	Type DataType `json:"type" xml:"type"`
	Name string   `json:"name" xml:"name"`
}

// DataType is the type of data.
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

// FileTypes includes both the SAUCE FileType value and name.
type FileTypes struct {
	Type FileType `json:"type" xml:"type"`
	Name string   `json:"name" xml:"name"`
}

// FileType is the type of file.
type FileType uint

// TypeInfos includes the SAUCE fields dependant on DataType and FileType.
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

func (a *ANSIFlags) String() (s string) {
	if a.Decimal == 0 {
		return s
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
	return strings.Join(l, ", ")
}

// Flags is the SAUCE Flags field.
type Flags uint8

func (f Flags) parse() ANSIFlags {
	const binary5Bits = "%05b"
	bin := fmt.Sprintf(binary5Bits, f)
	r := []rune(bin)
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
	switch b {
	case "0":
		return "blink mode"
	case "1":
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
	switch ls {
	case "00":
		return noPref
	case "01":
		return "select 8 pixel font"
	case "10":
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
	switch ar {
	case "00":
		return noPref
	case "01":
		return "stretch pixels"
	case "10":
		return "square pixels"
	default:
		return invalid
	}
}

// Character based files.
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
		"ScreamTracker module",
		"ScreamTracker 3 module",
		"MultiTracker module",
		"Farandole Composer module",
		"Ultra Tracker module",
		"Dual Module Player module",
		"X-Tracker module",
		"Oktalyzer module",
		"AdLib Visual Composer FM audio",
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

type (
	record   []byte
	id       [5]byte
	version  [2]byte
	title    [35]byte
	author   [20]byte
	group    [20]byte
	date     [8]byte
	fileSize [4]byte
	dataType [1]byte
	fileType [1]byte
	tInfo1   [2]byte
	tInfo2   [2]byte
	tInfo3   [2]byte
	tInfo4   [2]byte
	comments [1]byte
	comment  [64]byte
	tFlags   [1]byte
	tInfoS   [22]byte
)

func (t tInfoS) String() string {
	const nul = 0
	s := ""
	for _, b := range t {
		if b == nul {
			continue
		}
		s += string(b)
	}
	return s
}

// this sauce data struct intentionally shares the key names with the type key names.
// so the `data.version` item uses the type named `version` which is a [2]byte value.
type data struct {
	id       id
	version  version
	title    title
	author   author
	group    group
	date     date
	filesize fileSize
	datatype dataType
	filetype fileType
	tinfo1   tInfo1
	tinfo2   tInfo2
	tinfo3   tInfo3
	tinfo4   tInfo4
	comments comments
	tFlags   tFlags
	tInfoS   tInfoS
}

func (r record) date(i int) date {
	var d date
	const (
		start = 82
		end   = start + len(d)
	)
	for j, c := range r[start+i : end+i] {
		d[j] = c
	}
	return d
}

func (d *data) dates() Dates {
	t := d.parseDate()
	u := t.Unix()
	return Dates{
		Value: fmt.Sprintf("%s", d.date),
		Time:  t,
		Epoch: u,
	}
}

func (d *data) dataType() DataTypes {
	dt := DataType(unsignedBinary1(d.datatype))
	return DataTypes{
		Type: dt,
		Name: fmt.Sprintf("%v", dt.String()),
	}
}

func (d *data) description() (s string) {
	dt, ft := unsignedBinary1(d.datatype), unsignedBinary1(d.filetype)
	c := Character(ft)
	if DataType(dt) != character {
		return s
	}
	switch c {
	case ascii, ansi, ansiMation, ripScript, pcBoard, avatar, html, source, tundraDraw:
		return c.Desc()
	}
	return s
}

func (d *data) fileType() FileTypes {
	data, file := unsignedBinary1(d.datatype), unsignedBinary1(d.filetype)
	switch DataType(data) {
	case none:
		return FileTypes{FileType(none), none.String()}
	case character:
		c := Character(file)
		return FileTypes{FileType(c), c.String()}
	case bitmap:
		b := Bitmap(file)
		return FileTypes{FileType(b), b.String()}
	case vector:
		v := Vector(file)
		return FileTypes{FileType(v), v.String()}
	case audio:
		a := Audio(file)
		return FileTypes{FileType(a), a.String()}
	case binaryText:
		return FileTypes{FileType(binaryText), binaryText.String()}
	case xBin:
		return FileTypes{FileType(xBin), xBin.String()}
	case archive:
		a := Archive(file)
		return FileTypes{FileType(a), a.String()}
	case executable:
		return FileTypes{FileType(executable), executable.String()}
	default:
		return FileTypes{FileType(0), "error"}
	}
}

func (d *data) parseDate() (t time.Time) {
	da := d.date
	dy, err := strconv.Atoi(string(da[0:4]))
	if err != nil {
		fmt.Println("day conversion failed:", err)
		return t
	}
	dm, err := strconv.Atoi(string(da[4:6]))
	if err != nil {
		fmt.Println("month conversion failed:", err)
		return t
	}
	dd, err := strconv.Atoi(string(da[6:8]))
	if err != nil {
		fmt.Println("year conversion failed:", err)
		return t
	}
	return time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC)
}

func (d *data) sizes() Sizes {
	value := unsignedBinary4(d.filesize)
	en := language.English
	return Sizes{
		Bytes:   value,
		Decimal: humanize.Decimal(int64(value), en),
		Binary:  humanize.Binary(int64(value), en),
	}
}

func (d *data) typeInfo() TypeInfos {
	dt, ft := unsignedBinary1(d.datatype), unsignedBinary1(d.filetype)
	t1, t2, t3 := unsignedBinary2(d.tinfo1), unsignedBinary2(d.tinfo2), unsignedBinary2(d.tinfo3)
	flag := Flags(unsignedBinary1(d.tFlags))
	ti := TypeInfos{
		TypeInfo{t1, ""},
		TypeInfo{t2, ""},
		TypeInfo{t3, ""},
		flag.parse(),
		d.tInfoS.String(),
	}
	switch DataType(dt) {
	case none:
		return ti // golangci-lint deadcode placeholder
	case character:
		switch Character(ft) {
		case ascii, ansi, ansiMation, pcBoard, avatar, tundraDraw:
			ti.Info1.Info = "character width"
			ti.Info2.Info = "number of lines"
		case ripScript:
			ti.Info1.Info = "pixel width"
			ti.Info2.Info = "character screen height"
			ti.Info3.Info = "number of colors"
		case html, source:
			return ti
		}
	case bitmap:
		switch Bitmap(ft) {
		case gif, pcx, lbm, tga, fli, flc, bmp, gl, dl, wpg, png, jpg, mpg, avi:
			ti.Info1.Info = "pixel width"
			ti.Info2.Info = "pixel height"
			ti.Info3.Info = "pixel depth"
		}
	case vector:
		switch Vector(ft) {
		case dxf, dwg, wpvg, kinetix:
			return ti
		}
	case audio:
		switch Audio(ft) {
		case smp8, smp8s, smp16, smp16s:
			ti.Info1.Info = "sample rate"
		case mod, composer669, stm, s3m, mtm, far, ult, amf, dmf, okt, rol, cmf, midi,
			sadt, voc, wave, patch8, patch16, xm, hsc, it:
			return ti
		}
	case binaryText:
		return ti
	case xBin:
		ti.Info1.Info = "character width"
		ti.Info2.Info = "number of lines"
	case archive:
		switch Archive(ft) {
		case zip, arj, lzh, arc, tar, zoo, rar, uc2, pak, sqz:
			return ti
		}
	case executable:
		return ti
	}
	return ti
}

func (r record) author(i int) author {
	var a author
	const (
		start = 42
		end   = start + len(a)
	)
	for j, c := range r[start+i : end+i] {
		a[j] = c
	}
	return a
}

func (r record) comments(i int) comments {
	return comments{r[i+104]}
}

func (r record) dataType(i int) dataType {
	return dataType{r[i+94]}
}

func (r record) extract() data {
	i := Scan(r...)
	if i == -1 {
		return data{}
	}
	return data{
		id:       r.id(i),
		version:  r.version(i),
		title:    r.title(i),
		author:   r.author(i),
		group:    r.group(i),
		date:     r.date(i),
		filesize: r.fileSize(i),
		datatype: r.dataType(i),
		filetype: r.fileType(i),
		tinfo1:   r.tInfo1(i),
		tinfo2:   r.tInfo2(i),
		tinfo3:   r.tInfo3(i),
		tinfo4:   r.tInfo4(i),
		comments: r.comments(i),
		tFlags:   r.tFlags(i),
		tInfoS:   r.tInfoS(i),
	}
}

func (r record) fileSize(i int) fileSize {
	return fileSize{r[i+90], r[i+91], r[i+92], r[i+93]}
}

func (r record) fileType(i int) fileType {
	return fileType{r[i+95]}
}

func (r record) group(i int) group {
	var g group
	const (
		start = 62
		end   = start + len(g)
	)
	for j, c := range r[start+i : end+i] {
		g[j] = c
	}
	return g
}

func (r record) id(i int) id {
	return id{r[i+0], r[i+1], r[i+2], r[i+3], r[i+4]}
}

func (r record) tFlags(i int) tFlags {
	return tFlags{r[i+105]}
}

func (r record) title(i int) title {
	var t title
	const (
		start = 7
		end   = start + len(t)
	)
	for j, c := range r[start+i : end+i] {
		t[j] = c
	}
	return t
}

func (r record) tInfo1(i int) tInfo1 {
	return tInfo1{r[i+96], r[i+97]}
}

func (r record) tInfo2(i int) tInfo2 {
	return tInfo2{r[i+98], r[i+99]}
}

func (r record) tInfo3(i int) tInfo3 {
	return tInfo3{r[i+100], r[i+101]}
}

func (r record) tInfo4(i int) tInfo4 {
	return tInfo4{r[i+102], r[i+103]}
}

func (r record) tInfoS(i int) tInfoS {
	var s tInfoS
	const (
		start = 106
		end   = start + len(s)
	)
	for j, c := range r[start+i : end+i] {
		if c == 0 {
			continue
		}
		s[j] = c
	}
	return s
}

func (r record) version(i int) version {
	return version{r[i+5], r[i+6]}
}

// Scan returns the position of the SAUCE00 ID or -1 if no ID exists.
func Scan(b ...byte) (index int) {
	const sauceSize, maximum = 128, 512
	id, l := []byte(sauceID), len(b)
	var backwardsLoop = func(i int) int {
		return l - 1 - i
	}
	// search for the id sequence in b
	for i := range b {
		if i > maximum {
			break
		}
		i = backwardsLoop(i)
		if i < sauceSize {
			break
		}
		// do matching in reverse
		if b[i] != id[6] {
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
		if b[i-6] != id[0] {
			continue // S
		}
		return i - 6
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
		ID:       fmt.Sprintf("%s", d.id),
		Version:  fmt.Sprintf("%s", d.version),
		Title:    strings.TrimSpace(fmt.Sprintf("%s", d.title)),
		Author:   strings.TrimSpace(fmt.Sprintf("%s", d.author)),
		Group:    strings.TrimSpace(fmt.Sprintf("%s", d.group)),
		Date:     d.dates(),
		FileSize: d.sizes(),
		Data:     d.dataType(),
		File:     d.fileType(),
		Info:     d.typeInfo(),
		Desc:     d.description(),
	}
}

func unsignedBinary1(b [1]byte) (value uint8) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsignedBinary1 failed:", err)
	}
	return value
}

func unsignedBinary2(b [2]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsignedBinary2 failed:", err)
	}
	return value
}

func unsignedBinary4(b [4]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("unsignedBinary4 failed:", err)
	}
	return value
}
