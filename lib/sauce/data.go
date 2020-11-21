package sauce

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"retrotxt.com/retrotxt/lib/humanize"
)

// ErrDate to integer failure.
var ErrDate = errors.New("parse date to integer conversion")

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
	comnt    comnt
}

type comnt struct {
	index  int
	length int
	count  comments
	lines  []byte
}

func (d *data) commentBlock() (c Comments) {
	newlineCount := len(strings.Split(string(d.comnt.lines), "\n"))
	c.ID = comntID
	c.Count = int(unsignedBinary1(d.comnt.count))
	if newlineCount > 0 {
		// comments with newlines are technically invalid but they exist in the wild.
		// https://github.com/16colo-rs/16c/issues/67
		c.Comment = commentByNewline(d.comnt.lines)
		return c
	}
	c.Comment = commentByLine(d.comnt.lines)
	return c
}

func commentByNewline(b []byte) (lines []string) {
	r := bytes.NewReader(b)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func commentByLine(b []byte) (lines []string) {
	s, l := "", 0
	var resetLine = func() {
		s, l = "", 0
	}
	for _, c := range b {
		l++
		s += string(c)
		if l == comntLineSize {
			lines = append(lines, s)
			resetLine()
		}
	}
	return lines
}

func (d *data) dates() Dates {
	t, err := d.parseDate()
	if err != nil {
		fmt.Printf("sauce date error: %s\n", err)
	}
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

func (d *data) parseDate() (t time.Time, err error) {
	da := d.date
	dy, err := strconv.Atoi(string(da[0:4]))
	if err != nil {
		return t, fmt.Errorf("year failed: %v: %w", dy, ErrDate)
	}
	dm, err := strconv.Atoi(string(da[4:6]))
	if err != nil {
		return t, fmt.Errorf("month failed: %v: %w", dm, ErrDate)
	}
	dd, err := strconv.Atoi(string(da[6:8]))
	if err != nil {
		return t, fmt.Errorf("day failed: %v: %w", dd, ErrDate)
	}
	return time.Date(dy, time.Month(dm), dd, 0, 0, 0, 0, time.UTC), nil
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
