// Package sauce to handle the opening and reading of text files.
package sauce

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"retrotxt.com/retrotxt/lib/humanize"
)

// TODO: handle comments!!

const sauceID = "SAUCE00"

var (
	// ErrFmt format error.
	ErrFmt = errors.New("format is not known")
	// ErrNoName name cannot be empty.
	ErrNoName = errors.New("name cannot be empty")
	// ErrNoDir directories not usable with command.
	ErrNoDir = errors.New("directories are not usable with this command")
	// ErrNoFile file does not exist.
	ErrNoFile = errors.New("file does not exist")
)

// Record layout for printing.
type Record struct {
	ID       string    `json:"id"`
	Version  string    `json:"version"`
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Group    string    `json:"group"`
	Date     string    `json:"date"`
	LSDate   string    `json:"lsdate"`
	FileSize string    `json:"filesize"`
	Data     DataTypes `json:"dataType"`
	File     FileTypes `json:"fileType"`
	Info     TypeInfos `json:"typeInfo"`
}

type DataTypes struct {
	Type DataType `json:"type"`
	Name string   `json:"name"`
}

type FileTypes struct {
	Type FileType `json:"type"`
	Name string   `json:"name"`
}

type TypeInfos struct {
	Info1 TypeInfo  `json:"1"`
	Info2 TypeInfo  `json:"2"`
	Info3 TypeInfo  `json:"3"`
	Flags ANSIFlags `json:"flags"`
	Font  string    `json:"fontName"`
}

type TypeInfo struct {
	Value uint16 `json:"value"`
	Info  string `json:"info"`
}

type ANSIFlags struct {
	Decimal Flags    `json:"decimal"`
	Binary  string   `json:"binary"`
	B       ANSIFlag `json:"nonBlinkMode"`
	LS      ANSIFlag `json:"letterSpacing"`
	AR      ANSIFlag `json:"aspectRatio"`
}

type Flags uint8

type ANSIFlag struct {
	Flag string `json:"flag"`
	Info string `json:"interpretation"`
}

type (
	// DataType is the type of data.
	DataType uint
	// FileType is the type of file.
	FileType uint
	// Character based files.
	Character uint
	// Bitmap graphic and animation files.
	Bitmap uint
	// Vector graphic files.
	Vector uint
	// Audio or music files.
	Audio uint
	// BinaryText is a raw memory copy of a text mode screen.
	BinaryText uint
	// XBin or eXtended BinaryText files.
	XBin uint
	// Archive and compressed files.
	Archive uint
	// Executable program files.
	Executable uint
)

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
		"undefined", "text", "bitmap graphic or animation", "vector graphic",
		"audio or music", "binary text", "extended binary text", "archive", "executable",
	}[d]
	return fmt.Sprintf("%s file", s)
}

const (
	// ASCII text file with no formatting codes or color codes.
	ASCII Character = iota
	// ANSI text file with coloring codes and cursor positioning.
	ANSI
	// ANSIMation are ANSI text files that rely on fixed screen sizes.
	ANSIMation
	// RIPScript are Remote Imaging Protocol graphics.
	RIPScript
	// PCBoard color codes and macros, and ANSI codes.
	PCBoard
	// Avatar color codes, and ANSi codes.
	Avatar
	// HTML markup files.
	HTML
	// Source code for a programming language.
	Source
	// TundraDraw files, like ANSI, but with a custom palette.
	TundraDraw
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

func (c Character) info(t1, t2, t3 uint16, x string) string {
	switch c {
	case ASCII, ANSI, ANSIMation:
		if t1 == 0 && t2 == 0 {
			return ""
		}
		s := fmt.Sprintf("character width: %d, lines: %d", t1, t2)
		if x != "" {
			return fmt.Sprintf("%s; %s", s, x)
		}
		return s
	case RIPScript:
		if t1 == 0 && t2 == 0 && t3 == 0 {
			return ""
		}
		return fmt.Sprintf("pixel width: %d, height: %d, colors: %d", t1, t2, t3)
	case PCBoard, Avatar, TundraDraw:
		if t1 == 0 && t2 == 0 {
			return ""
		}
		return fmt.Sprintf("character width: %d, lines: %d", t1, t2)
	case HTML, Source:
		return ""
	}
	return "charcter error"
}

const (
	// GIF Graphics Interchange Format.
	GIF Bitmap = iota
	// PCX ZSoft Paintbrush.
	PCX
	// LBM DeluxePaint LBM/IFF.
	LBM
	// TGA Targa true color.
	TGA
	// FLI Autodesk animation.
	FLI
	// FLC Autodesk animation.
	FLC
	// BMP Windows or OS/2 Bitmap.
	BMP
	// GL Grasp GL animation.
	GL
	// DL animation.
	DL
	// WPG WordPerfect Bitmap.
	WPG
	// PNG Portable Network Graphics.
	PNG
	// JPG JPEG File Interchange Format.
	JPG
	// MPG Moving Picture Experts Group.
	MPG
	// AVI Audio Video Interleave.
	AVI
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

const (
	// DXF Drawing Exchange Format for AutoCAD and AutoDRAW CAD.
	DXF Vector = iota
	// DWG AutoCAD Drawing is the native binary format for AutoDesk CAD products.
	DWG
	// WPVG WordPerfect Graphics vector graphics (WPG).
	WPVG
	// Kinetix 3D Studio and 3D Studio MAX product line (3DS).
	Kinetix
)

func (v Vector) String() string {
	return [...]string{
		"AutoDesk CAD vector graphic",
		"WordPerfect vector graphic",
		"3D Studio vector graphic",
	}[v]
}

const (
	// MOD NoiseTracker 4, 6 or 8 channels.
	MOD Audio = iota
	// Composer669 an 8 channel module by Renaissance (669).
	Composer669
	// STM Future Crew 4 channel ScreamTracker.
	STM
	// S3M Future Crew variable channel ScreamTracker 3.
	S3M
	// MTM Renaissance variable channel MultiTracker.
	MTM
	// FAR Farandole composer.
	FAR
	// ULT Ultra Tracker.
	ULT
	// AMF DMP/DSMI Advanced Module Format.
	AMF
	// DMF Delusion Digital Music Format (X-Tracker).
	DMF
	// OKT Oktalyzer.
	OKT
	// ROL AdLib ROL file (FM audio).
	ROL
	// CMF Creative Music File (FM audio).
	CMF
	// MID aka MIDI (Musical Instrument Digital Interface).
	MID
	// SADT SAdT composer (FM audio).
	SADT
	// VOC Creative Voice file.
	VOC
	// WAV Waveform Audio file format.
	WAV
	// SMP8 Raw, single channel 8-bit sample.
	SMP8
	// SMP8S Raw, stereo 8-bit sample.
	SMP8S
	// SMP16 Raw, single channel 16-bit sample.
	SMP16
	// SMP16S Raw, stereo 16-bit sample.
	SMP16S
	// PATCH8 8-bit patch file.
	PATCH8
	// PATCH16 16-bit patch file.
	PATCH16
	// XM FastTracker ][ module.
	XM
	// HSC Tracker (FM audio).
	HSC
	// IT Impulse Tracker.
	IT
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

const (
	// ZIP originally from PKWare but now an open format.
	ZIP Archive = iota
	// ARJ Archive by Robert Jung.
	ARJ
	// LZH by Yoshizaki Haruyasu, also known as LHA.
	LZH
	// ARC by System Enhancement Associates.
	ARC
	// TAR or a tarball is an open archive format.
	TAR
	// ZOO format using LZW compression by Rahul Dhesi.
	ZOO
	// RAR Roshal Archive by Eugene Roshal.
	RAR
	// UC2 UltraCompressor II.
	UC2
	// PAK format is an extension of ARC also known as GSARC.
	PAK
	// SQZ Squeeze It by Jonas Hammarberg.
	SQZ
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
	// ansiflags
	aspectRatio   [2]string
	letterSpacing [2]string
	nonBlinkMode  [1]string
)

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

func (d *data) dataType() DataTypes {
	dt := DataType(unsignedBinary1(d.datatype))
	return DataTypes{
		Type: dt,
		Name: fmt.Sprintf("%v", dt.String()),
	}
}

func (d *data) fileType() FileTypes {
	data, file := unsignedBinary1(d.datatype), unsignedBinary1(d.filetype)
	switch DataType(data) {
	case none:
		return FileTypes{FileType(none), none.String()}
	case character:
		c := Character(file)
		// if s := strings.TrimSpace(fmt.Sprintf("%s", d.tInfoS)); s != "" {
		// 	return fmt.Sprintf("%s, %s", s, c.String())
		// }
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
		// if s := strings.TrimSpace(fmt.Sprintf("%s", d.tInfoS)); s != "" {
		// 	return fmt.Sprintf("%s, %s", s, binaryText.String())
		// }
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

func (d *data) fileSize() string {
	const kB = 1000
	value, p := unsignedBinary4(d.filesize), message.NewPrinter(language.English)
	if value < kB {
		return p.Sprintf("%d bytes", value)
	}
	h := humanize.Bytes(int64(value), language.AmericanEnglish)
	return p.Sprintf("%s (%d bytes)", h, value)
}

func (d *data) lsDate() string {
	da := d.date
	dy, err := strconv.Atoi(string(da[0:4]))
	if err != nil {
		fmt.Println("day conversion failed:", err)
		return fmt.Sprintf("%s", da)
	}
	dm, err := strconv.Atoi(string(da[4:6]))
	if err != nil {
		fmt.Println("month conversion failed:", err)
		return fmt.Sprintf("%s", da)
	}
	dd, err := strconv.Atoi(string(da[6:8]))
	if err != nil {
		fmt.Println("year conversion failed:", err)
		return fmt.Sprintf("%s", da)
	}
	t := time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC)
	return fmt.Sprintf("%s", t.Format("2 Jan 2006"))
}

func (d *data) typeInfo() TypeInfos {
	dt, ft := unsignedBinary1(d.datatype), unsignedBinary1(d.filetype)
	t1, t2, t3 := unsignedBinary2(d.tinfo1), unsignedBinary2(d.tinfo2), unsignedBinary2(d.tinfo3)
	flag := unsignedBinary1(d.tFlags)
	font := fmt.Sprintf("%v", d.tInfoS)
	ti := TypeInfos{
		TypeInfo{t1, ""},
		TypeInfo{t2, ""},
		TypeInfo{t3, ""},
		ansiFlags(Flags(flag)),
		font,
	}
	switch DataType(dt) {
	case none:
	case character:
		switch Character(ft) {
		case ASCII, ANSI, ANSIMation, PCBoard, Avatar, TundraDraw:
			ti.Info1.Info = "character width"
			ti.Info2.Info = "number of lines"
		case RIPScript:
			ti.Info1.Info = "pixel width"
			ti.Info2.Info = "character screen height"
			ti.Info3.Info = "number of colors"
		}
	case bitmap:
		ti.Info1.Info = "pixel width"
		ti.Info2.Info = "pixel height"
		ti.Info3.Info = "pixel depth"
	case vector:
	case audio:
		switch Audio(ft) {
		case SMP8, SMP8S, SMP16, SMP16S:
			ti.Info1.Info = "sample rate"
		case MOD, Composer669, STM, S3M, MTM, FAR, ULT, AMF, DMF, OKT, ROL, CMF, MID,
			SADT, VOC, WAV, PATCH8, PATCH16, XM, HSC, IT:
		}
	case binaryText:
	case xBin:
		ti.Info1.Info = "character width"
		ti.Info2.Info = "number of lines"
	case archive, executable:
	}
	return ti
}

func ansiFlags(f Flags) ANSIFlags {
	bin := fmt.Sprintf("%05b", f)
	r := []rune(bin)
	b, ls, ar := string(r[0]), string(r[1:3]), string(r[3:5])
	a := ANSIFlags{
		Decimal: f,
		Binary:  bin,
		B: ANSIFlag{
			Flag: b,
			Info: ansiB(b),
		},
		LS: ANSIFlag{
			Flag: ls,
			Info: ansiLS(ls),
		},
		AR: ANSIFlag{
			Flag: ar,
			Info: ansiAR(ar),
		},
	}
	return a
}

func ansiB(b string) string {
	switch b {
	case "0":
		return "blink mode"
	case "1":
		return "non-blink mode"
	default:
		return "invalid value"
	}
}

func ansiLS(ls string) string {
	switch ls {
	case "00":
		return "no preference"
	case "01":
		return "select 8 pixel font"
	case "10":
		return "select 9 pixel font"
	default:
		return "invalid value"
	}
}

func ansiAR(ar string) string {
	switch ar {
	case "00":
		return "no preference"
	case "01":
		return "stretch pixels"
	case "10":
		return "square pixels"
	default:
		return "invalid value"
	}
}

// extract sauce record.
func (r record) extract() data {
	i := Scan(r)
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

func (r record) id(i int) id {
	return id{r[i+0], r[i+1], r[i+2], r[i+3], r[i+4]}
}

func (r record) version(i int) version {
	return version{r[i+5], r[i+6]}
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

func (r record) fileSize(i int) fileSize {
	return fileSize{r[i+90], r[i+91], r[i+92], r[i+93]}
}

func (r record) dataType(i int) dataType {
	return dataType{r[i+94]}
}

func (r record) fileType(i int) fileType {
	return fileType{r[i+95]}
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

func (r record) comments(i int) comments {
	return comments{r[i+104]}
}

func (r record) tFlags(i int) tFlags {
	return tFlags{r[i+105]}
}

func (r record) tInfoS(i int) tInfoS {
	var s tInfoS
	const (
		start = 106
		end   = start + len(s)
	)
	for j, c := range r[start+i : end+i] {
		fmt.Printf("%v\n", c)
		if c == 0 {
			continue
		}
		s[j] = c
	}
	return s
}

// Text format the record data.
// func Text(b []byte) {
// 	var info = func(t string) string {
// 		return str.Cinf(fmt.Sprintf("%s\t", t))
// 	}
// 	s := Parse(b...)
// 	if s.ID != "SAUCE" {
// 		return
// 	}
// 	var data = []struct {
// 		k, v string
// 	}{
// 		{k: "title", v: s.Title},
// 		{k: "author", v: s.Author},
// 		{k: "group", v: s.Group},
// 		{k: "date", v: s.LSDate},
// 		{k: "filesize", v: s.FileSize},
// 		{k: "type", v: s.DataType},
// 		{k: "file", v: s.FileType},
// 		{k: "info", v: s.TypeInfo},
// 	}
// 	var buf bytes.Buffer
// 	w := new(tabwriter.Writer)
// 	w.Init(&buf, 0, 8, 0, '\t', 0)
// 	for _, x := range data {
// 		if x.k == "filesize" && s.FileSize == "0" {
// 			continue
// 		}
// 		if x.k == "info" && s.TypeInfo == "" {
// 			continue
// 		}
// 		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
// 	}
// 	if err := w.Flush(); err != nil {
// 		logs.Fatal("flush of tab writer failed", "", err)
// 	}
// 	fmt.Print(buf.String())
// }

// Scan returns the position of the SAUCE00 ID or -1 if no ID exists.
func Scan(b []byte) (index int) {
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
		Date:     fmt.Sprintf("%s", d.date),
		LSDate:   d.lsDate(), // todo: change to time type
		FileSize: d.fileSize(),
		Data:     d.dataType(),
		File:     d.fileType(),
		Info:     d.typeInfo(),
		//
	}
}

func unsignedBinary1(b [1]byte) (value uint8) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return value
}

func unsignedBinary2(b [2]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return value
}

func unsignedBinary4(b [4]byte) (value uint16) {
	buf := bytes.NewReader(b[:])
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return value
}
