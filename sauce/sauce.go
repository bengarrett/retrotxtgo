//Package sauce to handle the opening and reading of text files
package sauce

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type data struct {
	id       []byte
	version  []byte
	title    []byte
	author   []byte
	group    []byte
	date     []byte
	filesize []byte
	datatype []byte
	filetype []byte
	tinfo1   []byte
	tinfo2   []byte
	tinfo3   []byte
	tinfo4   []byte
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
	FileSize uint8
	DataType string
	FileType string
	TypeInfo string
}

var datatypes = make(map[uint8]string)

//Exists ooof
func Exists(filename string) bool {
	file, _ := os.Open(filename)
	defer file.Close()

	buffer := make([]byte, 1024)

	_, err := file.Seek(-128, 2)
	if err != nil {
		panic(err)
	}

	_, err = file.Read(buffer)
	if err != nil {
		panic(err)
	}

	x := fmt.Sprintf("%X", buffer[90:94])
	y, _ := strconv.ParseInt(x, 16, 64)
	// fmt.Printf("%v", y)

	v := y
	v += y << 8
	v += y << 16
	v += (y << 24)

	fmt.Printf("\n[[%v]]\n", v)

	if string(buffer[0:7]) == "SAUCE00" {
		return true
	}
	return false
}

//Get sauce
func slice(b []byte) data {
	p := Scan(b)
	d := data{
		id:       b[p : p+5],
		version:  b[p+5 : p+7],
		title:    b[p+7 : p+42],
		author:   b[p+42 : p+62],
		group:    b[p+62 : p+82],
		date:     b[p+82 : p+90],
		filesize: b[p+90 : p+94],
		datatype: b[p+94 : p+95],
		filetype: b[p+95 : p+96],
		tinfo1:   b[p+96 : p+98],
		tinfo2:   b[p+98 : p+100],
		tinfo3:   b[p+100 : p+102],
		tinfo4:   b[p+102 : p+104], // unused, should always be 0
	}
	return d
}

func binconvert(b []byte) uint8 {
	// An unsigned binary value of 1 byte (0 to 255), 2 bytes (0 to 65535) or
	// 4 bytes (0 to 4294967295) stored in intel little-endian format.
	var f float64
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.LittleEndian, &f)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			fmt.Println("convert binary.Read failed:", err)
			return 0
		}

	}
	fmt.Printf("\t%d - %g\t", b, f)
	return uint8(f)
}

func id(d data) string {
	return fmt.Sprintf("%s", d.id)
}
func version(d data) string {
	return fmt.Sprintf("%s", d.version)
}
func title(d data) string {
	s := strings.TrimSpace(string(d.title))
	return fmt.Sprintf("%v", s)
}
func author(d data) string {
	s := strings.TrimSpace(string(d.author))
	return fmt.Sprintf("%v", s)
}
func group(d data) string {
	s := strings.TrimSpace(string(d.group))
	return fmt.Sprintf("%v", s)
}
func date(d data) string {
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
func filesize(d data) uint8 {
	fs := binconvert(d.filesize)
	fmt.Println(d.filesize, fs)
	return fs
}
func datatype(d data) string {
	types := [9]string{"None", "Character", "Graphics", "Vector", "Sound", "BinaryText", "XBin", "Archive", "Executable"}
	val := binconvert(d.datatype)
	return fmt.Sprintf("%v", types[val])
}
func filetype(d data) string {
	dt := binconvert(d.datatype) // todo move to d data
	ft := binconvert(d.filetype) // move to d data
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
	dt := binconvert(d.datatype) // todo move to d data
	ft := binconvert(d.filetype) // move to d data
	t1 := binconvert(d.tinfo1)
	t2 := binconvert(d.tinfo2)
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
func characterinfo(i uint8) string {
	return ""
}
func bitmapinfo(i uint8) string {
	return ""
}
func audioinfo(i uint8) string {
	return ""
}
func xbininfo(i uint8, t1 uint8, t2 uint8) string {
	if t1 == 0 {
		t1 = 80
	}
	return fmt.Sprintf("Width: %v Line height: %v", t1, t2)
}

//Scan blah blah
func Scan(b []byte) int {
	s := bytes.LastIndexAny(b, "SAUCE00")
	if s < 0 {
		return -1
	}
	//l := len(b) - 128
	// if s > -1 {
	// 	if len(b)-s < 128 {
	// 		return -1
	// 	}
	// }
	//fmt.Printf("\n>%v|%v=%v\t%v\n%q\n\n", l, s, s-l, 128-(s-l), b[len(b)-128:])
	if bytes.Compare(b[s:s+7], []byte("SAUCE00")) != 0 {
		s = s + 16 - 128
	}
	return s
}

//Get sauce oof
func Get(b []byte) Record {
	d := slice(b)
	// do checks
	r := Record{
		ID:       id(d),
		Version:  version(d),
		Title:    title(d),
		Author:   author(d),
		Group:    group(d),
		Date:     date(d),
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
	fmt.Printf("SAUCE metadata:\t\t%v\n", "Yes")
	s := Get(b)
	fmt.Printf("\nversion:\t\t%v\ntitle:\t\t%q\nauthor:\t\t%q\n", s.Version, s.Title, s.Author)
	fmt.Printf("group:\t\t%q\ndate:\t\t%s\nlsd:\t\t%s\n", s.Group, s.Date, s.LSDate)
	fmt.Printf("file size:\t%v\n", s.FileSize)
	fmt.Printf("data type:\t%q\n", s.DataType)
	fmt.Printf("file type:\t%q\n", s.FileType)
	fmt.Printf("type info:\t%v\n", s.TypeInfo)

	// se := sauce.Exists(path)
	// fmt.Printf("\n\nexists?:\t%v\n", se)
}