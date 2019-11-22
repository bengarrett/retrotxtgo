//Package sauce to handle the opening and reading of text files
package sauce

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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
	FileSize int
	DataType string
	FileType string
}

var datatypes = make(map[float64]string)
var filetypes = make(map[float64]string)

//Get sauce
func slice(b []byte) data {
	p := Scan(b)
	d := data{
		id:      b[p : p+5],
		version: b[p+5 : p+7],
		title:   b[p+7 : p+42],
		author:  b[p+42 : p+62],
		group:   b[p+62 : p+82],
		date:    b[p+82 : p+90],
		// An unsigned binary value of 1 byte (0 to 255), 2 bytes (0 to 65535) or
		// 4 bytes (0 to 4294967295) stored in intel little-endian format.
		filesize: b[p+90 : p+94],
		datatype: b[p+94 : p+95],
		filetype: b[p+95 : p+96],
	}
	return d
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
func filesize(d data) int {
	var fs float64
	buf := bytes.NewReader(d.filesize)
	err := binary.Read(buf, binary.LittleEndian, &fs)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			fmt.Println("filesize binary.Read failed:", err)
		}

	}
	return int(fs)
}
func datatype(d data) string {
	var val float64
	buf := bytes.NewReader(d.datatype)
	fmt.Print(buf)
	err := binary.Read(buf, binary.LittleEndian, &val)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			fmt.Println("datatype binary.Read failed:", err)
		}

	}
	// move to func
	datatypes[0] = "None"
	datatypes[1] = "Character"
	datatypes[2] = "Bitmap"
	datatypes[3] = "Vector"
	datatypes[4] = "Audio"
	datatypes[5] = "Binary text"
	datatypes[6] = "XBin"
	datatypes[7] = "Archive"
	datatypes[8] = "Executable"

	return fmt.Sprintf("%v", datatypes[val])
}
func filetype(d data) string {
	var pi float64
	buf := bytes.NewReader(d.datatype)
	fmt.Print(buf)
	binary.Read(buf, binary.LittleEndian, &pi)

	return fmt.Sprintf("%v", filetypes[pi]) // todo create map and return data types
}

//Scan blah blah
func Scan(b []byte) int {
	s := bytes.Index(bytes.ToUpper(b), []byte("SAUCE00"))
	if s > -1 {
		if len(b)-s < 128 {
			return -1
		}
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
	}
	return r
}

//Print sauce
func Print(r Record) {
	fmt.Println(r)
}
