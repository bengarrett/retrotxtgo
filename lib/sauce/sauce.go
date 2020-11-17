//Package sauce to handle the opening and reading of text files
package sauce

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"retrotxt.com/retrotxt/lib/logs"
	"retrotxt.com/retrotxt/lib/str"
)

const sauceID = "SAUCE00"

type (
	// None undefined filetype.
	None uint
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

const (
	// GIF CompuServe Graphics Interchange Format.
	GIF Bitmap = iota
	// PCX ZSoft Paintbrush.
	PCX
	// LBM DeluxePaint LBM/IFF.
	LBM
	// TGA Targa truecolor.
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
	// WPG Wordperfect Bitmap.
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
	// ULT UltraTracker.
	ULT
	// AMF DMP/DSMI Advanced Module Format.
	AMF
	// DMF Delusion Digital Music Format (XTracker).
	DMF
	// OKT Oktalyser.
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
}

//Record blah
type Record struct {
	ID       string
	Version  string
	Title    string
	Author   string
	Group    string
	Date     string
	LSDate   string
	FileSize string
	DataType string
	FileType string
	TypeInfo string
}

//var datatypes = make(map[uint8]string)

//Get sauce
func slice(b record) data {
	i := Scan(b)
	if i == -1 {
		return data{}
	}
	return data{
		id:       b.id(i),
		version:  b.version(i),
		title:    b.title(i),
		author:   b.author(i),
		group:    b.group(i),
		date:     b.date(i),
		filesize: b.fileSize(i),
		datatype: b.dataType(i),
		filetype: b.fileType(i),
		tinfo1:   b.tInfo1(i),
		tinfo2:   b.tInfo2(i),
		tinfo3:   b.tInfo3(i),
		tinfo4:   b.tInfo4(i),
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
	return tInfo4{r[i+102], r[i+102]}
}

func getID(d data) string {
	return fmt.Sprintf("%s", d.id)
}
func getVersion(d data) string {
	return fmt.Sprintf("%s", d.version)
}
func getTitle(d data) string {
	s := strings.TrimSpace(fmt.Sprintf("%s", d.title))
	return fmt.Sprintf("%v", s)
}
func getAuthor(d data) string {
	s := strings.TrimSpace(fmt.Sprintf("%s", d.author))
	return fmt.Sprintf("%v", s)
}
func getGroup(d data) string {
	s := strings.TrimSpace(fmt.Sprintf("%s", d.group))
	return fmt.Sprintf("%v", s)
}
func getDate(d data) string {
	return fmt.Sprintf("%s", d.date)
}
func lsdate(d data) string {
	da := d.date
	dy, _ := strconv.Atoi(string(da[0:4]))
	dm, _ := strconv.Atoi(string(da[4:6]))
	dd, _ := strconv.Atoi(string(da[6:8]))
	t := time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC)
	year, month, day := t.Date()
	return fmt.Sprintf("%v-%v-%v", year, month, day)
}

func unsignedBinary(b []byte) (value uint16) {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return value
}

func unsignedBinary8(b []byte) (value uint8) {
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &value)
	if err != nil {
		fmt.Println("binary.Read failed:", err)
	}
	return value
}

func filesize(d data) string {
	value := unsignedBinary(d.filesize[:])
	return fmt.Sprintf("%d", value)
}
func datatype(d data) string {
	types := [9]string{"None", "Character", "Graphics", "Vector", "Sound", "BinaryText", "XBin", "Archive", "Executable"}

	fmt.Println(d.datatype[:])
	val := unsignedBinary8(d.datatype[:])

	fmt.Println("ti1", d.tinfo1, unsignedBinary(d.tinfo1[:]))
	fmt.Println("ti2", d.tinfo2, unsignedBinary(d.tinfo2[:]))

	return fmt.Sprintf("%v", types[val])
}
func filetype(d data) string {
	dt := unsignedBinary8(d.datatype[:])
	ft := unsignedBinary8(d.filetype[:])
	fmt.Println("dt", dt, "ft", ft)
	type filetypes struct {
		data  uint8
		types []string
	}
	// these values were copied from tehmaze/sauce
	// https://github.com/tehmaze/sauce
	var fts = filetypes{0, []string{"-"}}
	fts = filetypes{1, []string{"ASCII", "ANSi", "ANSiMation", "RIP", "PCBoard", "Avatar", "HTML", "Source"}}
	fts = filetypes{2, []string{"GIF", "PCX", "LBM/IFF", "TGA", "FLI", "FLC", "BMP", "GL", "DL", "WPG", "PNG", "JPG", "MPG", "AVI"}}
	fts = filetypes{3, []string{"DX", "DWG", "WPG", "3DS"}}
	fts = filetypes{4, []string{"MOD", "669", "STM", "S3M", "MTM", "FAR", "ULT", "AMF", "DMF", "OKT", "ROL", "CMF", "MIDI", "SADT", "VOC", "WAV", "SMP8", "SMP8S", "SMP16", "SMP16S", "PATCH8", "PATCH16", "XM", "HSC", "IT"}}
	fts = filetypes{5, []string{"-"}}
	fts = filetypes{6, []string{"-"}}
	fts = filetypes{7, []string{"ZIP", "ARJ", "LZH", "ARC", "TAR", "ZOO", "RAR", "UC2", "PAK", "SQZ"}}
	fts = filetypes{8, []string{"-"}}
	fts.data = dt
	return fmt.Sprintf("%v", fts.types[ft])
}
func typeinfo(d data) string {
	dt := unsignedBinary(d.filesize[:])
	ft := unsignedBinary(d.datatype[:])
	t1 := unsignedBinary(d.tinfo1[:])
	t2 := unsignedBinary(d.tinfo2[:])
	switch dt {
	case 0:
	case 3:
	case 5:
	case 7:
	case 8:
		return ""
	case 1:
		return characterinfo(ft)
	case 2:
		return bitmapinfo(ft)
	case 4:
		return audioinfo(ft)
	case 6:
		return xbininfo(ft, t1, t2)
	}
	return ""
}
func characterinfo(i uint16) string {
	return ""
}
func bitmapinfo(i uint16) string {
	return ""
}
func audioinfo(i uint16) string {
	return ""
}
func xbininfo(i uint16, t1 uint16, t2 uint16) string {
	if t1 == 0 {
		t1 = 80
	}
	return fmt.Sprintf("Width: %v Line height: %v", t1, t2)
}

//Scan blah blah
func Scan(b []byte) (index int) {
	const sauceSize = 128
	index = bytes.LastIndexAny(b, sauceID)
	if index < 0 {
		return -1
	}
	if !bytes.Equal(b[index:index+len(sauceID)], []byte(sauceID)) {
		index = index + 16 - sauceSize
	}
	if (len(b) - index - sauceSize) < 0 {
		// sauce data is expected to be at least 128 bytes
		return -1
	}
	return index
}

//Get sauce oof
func Get(b []byte) Record {
	d := slice(b)
	// do checks
	r := Record{
		ID:       getID(d),
		Version:  getVersion(d),
		Title:    getTitle(d),
		Author:   getAuthor(d),
		Group:    getGroup(d),
		Date:     getDate(d),
		LSDate:   lsdate(d),
		FileSize: filesize(d),
		DataType: datatype(d),
		FileType: filetype(d),
		TypeInfo: typeinfo(d),
	}
	return r
}

//Print sauce
func Print(b []byte) {
	var info = func(t string) string {
		return str.Cinf(fmt.Sprintf("%s\t", t))
	}

	s := Get(b)

	var data = []struct {
		k, v string
	}{
		{k: "title", v: s.Title},
		{k: "author", v: s.Author},
		{k: "group", v: s.Group},
		{k: "date", v: s.LSDate},
		{k: "filesize", v: s.FileSize},
		{k: "type", v: s.DataType},
		{k: "file", v: s.FileType},
		{k: "info", v: s.TypeInfo},
	}
	var buf bytes.Buffer
	w := new(tabwriter.Writer)
	w.Init(&buf, 0, 8, 0, '\t', 0)
	for _, x := range data {
		if x.k == "filesize" && s.FileSize == "0" {
			continue
		}
		fmt.Fprintf(w, "\t %s\t  %s\n", x.k, info(x.v))
	}
	if err := w.Flush(); err != nil {
		logs.Fatal("flush of tab writer failed", "", err)
	}
	fmt.Println(buf.String())

	// fmt.Printf("\nversion:\t\t%v\ntitle:\t\t%q\nauthor:\t\t%q\n", s.Version, s.Title, s.Author)
	// fmt.Printf("group:\t\t%q\ndate:\t\t%s\nlsd:\t\t%s\n", s.Group, s.Date, s.LSDate)
	// fmt.Printf("file size:\t%v\n", s.FileSize)
	// fmt.Printf("data type:\t%q\n", s.DataType)
	// fmt.Printf("file type:\t%q\n", s.FileType)
	// fmt.Printf("type info:\t%v\n", s.TypeInfo)
	// se := sauce.Exists(path)
	// fmt.Printf("\n\nexists?:\t%v\n", se)
}
