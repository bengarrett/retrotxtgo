package info

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"hash/crc64"
	"io"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/bengarrett/bbs"
	"github.com/bengarrett/retrotxtgo/fsys"
	"github.com/bengarrett/retrotxtgo/nl"
	"github.com/bengarrett/sauce"
	"github.com/bengarrett/sauce/humanize"
	"github.com/charmbracelet/lipgloss"
	"github.com/mozillazg/go-slugify"
	"github.com/zRedShift/mimemagic"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Detail is the exported file details.
type Detail struct {
	XMLName    xml.Name     `json:"-"          xml:"file"`
	Name       string       `json:"filename"   xml:"name"`          // Name is the file name.
	Unicode    string       `json:"unicode"    xml:"unicode,attr"`  // Unicode is the file encoding if in Unicode.
	LineBreak  nl.LineBreak `json:"lineBreak"  xml:"line_break"`    // LineBreak is the line break used in the file.
	Count      Stats        `json:"counts"     xml:"counts"`        // Count is the file content statistics.
	Size       Sizes        `json:"size"       xml:"size"`          // Size is the file size in multiples.
	Lines      int          `json:"lines"      xml:"lines"`         // Lines is the number of lines in the file.
	Width      int          `json:"width"      xml:"width"`         // Width is the number of characters per line in the file, this may be inaccurate.
	Modified   ModDates     `json:"modified"   xml:"last_modified"` // Modified is the last modified date of the file.
	Sums       Checksums    `json:"checksums"  xml:"checksums"`     // Sums are the checksums of the file.
	Mime       Content      `json:"mime"       xml:"mime"`          // Mime is the file content metadata.
	Slug       string       `json:"slug"       xml:"id,attr"`       // Slug is the file name slugified.
	Sauce      sauce.Record `json:"sauce"      xml:"sauce"`         // Sauce is the SAUCE metadata.
	ZipComment string       `json:"zipComment" xml:"zip_comment"`   // ZipComment is the zip file comment.
	UTF8       bool         `json:"-"          xml:"-"`             // UTF8 is true if the file is UTF-8 encoded.
	LegacySums bool         `json:"-"          xml:"-"`             // LegacySums is true if the user requests legacy checksums.
	sauceIndex int          // sauceIndex is the index of the SAUCE record in the file.
}

// Checksums act as a fingerprint of the file for uniqueness and data corruption checks.
type Checksums struct {
	CRC32  string `json:"crc32"  xml:"crc32"`  // CRC32 is a cyclic redundancy check of the file.
	CRC64  string `json:"crc64"  xml:"crc64"`  // CRC64 is a cyclic redundancy check of the file.
	MD5    string `json:"md5"    xml:"md5"`    // MD5 is a weak cryptographic hash function.
	SHA256 string `json:"sha256" xml:"sha256"` // SHA256 is a strong cryptographic hash function.
}

// Content metadata from either MIME content type and magic file data.
type Content struct {
	Type  string `json:"-"        xml:"-"`
	Media string `json:"media"    xml:"media"`     // Media is the MIME media type.
	Sub   string `json:"subMedia" xml:"sub_media"` // Sub is the MIME sub type.
	Commt string `json:"comment"  xml:"comment"`   // Commt is the MIME comment.
}

// ModDates is the file last modified dates in multiple output formats.
type ModDates struct {
	Time  time.Time `json:"iso"   xml:"date"`       // Time is the last modified date of the file.
	Epoch int64     `json:"epoch" xml:"epoch,attr"` // Epoch is the last modified date of the file in seconds since the Unix epoch.
}

// Sizes of the file in multiples.
type Sizes struct {
	Bytes   int64  `json:"bytes"   xml:"bytes"`        // Bytes is the size of the file in bytes.
	Decimal string `json:"decimal" xml:"decimal,attr"` // Decimal is the size of the file with decimal units.
	Binary  string `json:"binary"  xml:"binary,attr"`  // Binary is the size of the file with binary units.
}

// Stats are the text file content statistics and counts.
type Stats struct {
	Chars    int `json:"characters"   xml:"characters"`    // Chars is the number of characters in the file.
	Controls int `json:"ansiControls" xml:"ansi_controls"` // Controls is the number of ANSI escape controls in the file.
	Words    int `json:"words"        xml:"words"`         // Words is the number of words in the file, this may be inaccurate.
}

// Format of the output text.
type Format int

const (
	ColorText Format = iota // ColorText is ANSI colored text.
	PlainText               // PlainText is standard text.
	JSON                    // JSON data-interchange format.
	JSONMin                 // JSONMin is JSON data minified.
	XML                     // XML markup data.
)

const (
	uc8         = "UTF-8"
	ans         = "ANSI controls"
	cmmt        = "comment"
	txt         = "text"
	zipComment  = "zip comment"
	octetStream = "application/octet-stream"
	zipType     = "application/zip"

	// Tree layout constants.
	treeCorner     = "├── "
	treeLastCorner = "└── "
	treeVertical   = "│   "
	treeSpace      = "    "

	// Header constants.
	headerPadding  = 2 // spaces on each side of header text
	paddingDivisor = 2 // divisor for centering header text
)

// lang returns the English Language tag used for numeric syntax formatting.
func lang() language.Tag {
	return language.English
}

// Ctrls counts the number of ANSI escape controls in the named file.
func (d *Detail) Ctrls(name string) error {
	r, err := os.Open(name)
	if err != nil {
		return err
	}
	defer r.Close()
	cnt, err := fsys.Controls(r)
	if err != nil {
		return err
	}
	d.Count.Controls = cnt
	return nil
}

// Marshal writes the Detail data in a given format syntax.
func (d *Detail) Marshal(w io.Writer, f Format) error {
	if w == nil {
		w = io.Discard
	}
	const jsTab = "    "
	const xmlTab = "\t"
	var err error
	switch f {
	case ColorText:
		d.marshalAsTree(w, true)
	case PlainText:
		d.marshalAsTree(w, false)
	case JSON:
		b, errj := json.MarshalIndent(d, "", jsTab)
		if errj != nil {
			return fmt.Errorf("detail json indent marshal: %w", errj)
		}
		_, err = w.Write(b)
	case JSONMin:
		b, errj := json.Marshal(d)
		if errj != nil {
			return fmt.Errorf("detail json marshal: %w", errj)
		}
		_, err = w.Write(b)
	case XML:
		b, errj := xml.MarshalIndent(d, "", xmlTab)
		if errj != nil {
			return fmt.Errorf("detail xml marshal: %w", errj)
		}
		_, err = w.Write(b)
	default:
		return fmt.Errorf("detail marshal %q: %w", f, ErrFmt)
	}
	if err != nil {
		return fmt.Errorf("detail marshal %q: %w", f, err)
	}
	return nil
}

// MimeUnknown detects non-Standard legacy data.
func (d *Detail) MimeUnknown() {
	if d.Mime.Commt != "unknown" {
		return
	}
	if d.Count.Controls > 0 {
		d.Mime.Commt = "Text document with ANSI controls"
		return
	}
	switch d.LineBreak.Decimal {
	case [2]rune{21}, [2]rune{133}:
		d.Mime.Commt = "EBCDIC encoded text document"
		return
	}
	if d.Mime.Type == octetStream {
		if !d.UTF8 && d.Count.Words > 0 {
			d.Mime.Commt = "US-ASCII encoded text document"
			return
		}
	}
}

// Parse the file and the raw data content.
func (d *Detail) Parse(name string, data ...byte) error {
	routines := 5
	if d.LegacySums {
		routines += 3
	}
	wg := sync.WaitGroup{}
	wg.Add(routines)
	go func() {
		defer wg.Done()
		d.sauceIndex = sauce.Index(data)
		if d.sauceIndex > 0 {
			d.Sauce = sauce.Decode(data)
		}
	}()
	go func() {
		defer wg.Done()
		d.mime(name, data...)
	}()
	go func() {
		defer wg.Done()
		stat, _ := os.Stat(name)
		d.input(len(data), stat)
	}()
	if d.LegacySums {
		go func() {
			defer wg.Done()
			crc32sum := crc32.ChecksumIEEE(data)
			d.Sums.CRC32 = strconv.FormatUint(uint64(crc32sum), 16)
		}()
		go func() {
			defer wg.Done()
			crc64sum := crc64.Checksum(data, crc64.MakeTable(crc64.ECMA))
			d.Sums.CRC64 = strconv.FormatUint(crc64sum, 16)
		}()
		go func() {
			defer wg.Done()
			md5sum := md5.Sum(data)
			d.Sums.MD5 = hex.EncodeToString(md5sum[:])
		}()
	}
	go func() {
		defer wg.Done()
		shasum := sha256.Sum256(data)
		d.Sums.SHA256 = hex.EncodeToString(shasum[:])
	}()
	go func() {
		defer wg.Done()
		d.UTF8 = utf8.Valid(data)
		d.Unicode = unicode(d.UTF8, data...)
	}()
	wg.Wait()
	return nil
}

func unicode(uni bool, b ...byte) string {
	UTF8Bom := []byte{0xEF, 0xBB, 0xBF}
	// little endianness, x86, ARM
	UTF16LEBom := []byte{0xFF, 0xFE}
	UTF32LEBom := []byte{0xFF, 0xFE, 0x00, 0x00}
	// big endianness, legacy mainframes, RISC
	UTF16BEBom := []byte{0xFE, 0xFF}
	UTF32BEBom := []byte{0x00, 0x00, 0xFE, 0xFF}
	switch {
	case bytes.HasPrefix(b, UTF8Bom):
		return uc8
	case bytes.HasPrefix(b, UTF16LEBom):
		return "UTF-16 LE"
	case bytes.HasPrefix(b, UTF16BEBom):
		return "UTF-16 BE"
	case bytes.HasPrefix(b, UTF32LEBom):
		return "UTF-32 LE"
	case bytes.HasPrefix(b, UTF32BEBom):
		return "UTF-32 BE"
	default:
		if uni {
			return "UTF-8 compatible"
		}
		return "no"
	}
}

func sauceDate(s string) string {
	t, err := time.Parse("20060102", s) // CCYYMMDD
	if err != nil {
		return ""
	}
	return humanize.DMY.Format(t.UTC())
}

// Read and parse the named file and content.
func (d *Detail) Read(name string) error {
	// Read file content
	p, err := fsys.ReadAllBytes(name)
	if err != nil {
		return err
	}
	return d.Parse(name, p...)
}

// ValidText reports whether the MIME content-type value is valid for text files.
func ValidText(mime string) bool {
	s := strings.Split(mime, "/")
	const req = 2
	if len(s) != req {
		return false
	}
	if s[0] == txt {
		return true
	}
	if mime == octetStream {
		return true
	}
	return false
}

// Len counts the number of characters used per line in the named file.
func (d *Detail) Len(name string) error {
	r, err := os.Open(name)
	if err != nil {
		return err
	}
	defer r.Close()
	w, err := fsys.Columns(r, d.LineBreak.Decimal)
	if err != nil {
		return err
	}
	if w < 0 {
		w = d.Count.Chars
	}
	d.Width = w
	return nil
}

// Words counts the number of words used in the named file.
func (d *Detail) Words(name string) error {
	r, err := os.Open(name)
	if err != nil {
		return err
	}
	defer r.Close()
	switch d.LineBreak.Decimal {
	case [2]rune{nl.NL}, [2]rune{nl.NEL}:
		if d.Count.Words, err = fsys.WordsEBCDIC(r); err != nil {
			return err
		}
	default:
		if d.Count.Words, err = fsys.Words(r); err != nil {
			return err
		}
	}
	return nil
}

// input parses simple statistical data on the file.
func (d *Detail) input(data int, stat fs.FileInfo) {
	if stat != nil {
		b := stat.Size()
		d.Size.Bytes = b
		d.Size.Binary = humanize.Binary(b, lang())
		d.Size.Decimal = humanize.Decimal(b, lang())
		d.Name = stat.Name()
		d.Modified.Time = stat.ModTime().UTC()
		d.Modified.Epoch = stat.ModTime().Unix()
		d.Slug = slugify.Slugify(stat.Name())
		return
	}
	b := int64(data)
	d.Size.Bytes = b
	d.Size.Binary = humanize.Binary(b, lang())
	d.Size.Decimal = humanize.Decimal(b, lang())
	d.Name = "n/a (stdin)"
	d.Slug = "n/a"
	d.Modified.Time = time.Now()
	d.Modified.Epoch = time.Now().Unix()
}

// marshalAsTree returns the marshaled detail data using a tree-like structure.
// This provides a more organized and visually appealing output format.
func (d *Detail) marshalAsTree(w io.Writer, useColors bool) { //nolint:cyclop,funlen,gocognit
	if w == nil {
		w = io.Discard
	}

	// Create styles based on whether we're using colors
	var (
		headerStyle func(string) string
		keyStyle    func(string) string
		valueStyle  func(string) string
		treeStyle   func(string) string
	)

	if useColors {
		// Color styles using lipgloss
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

		headerStyle = func(s string) string {
			return lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("231")).
				Padding(0, 1).
				Render(s)
		}

		keyStyle = func(s string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("39")).
				Render(s)
		}

		valueStyle = func(s string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("252")).
				Render(s)
		}

		treeStyle = func(s string) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				Render(s)
		}

		// Create header with lipgloss box
		header := headerStyle("File Information")
		border := borderStyle.Render(header)
		fmt.Fprintln(w, border)
	} else {
		// Plain text styles (no styling functions needed)
		keyStyle = func(s string) string { return s }
		valueStyle = func(s string) string { return s }
		treeStyle = func(s string) string { return s }

		// Create plain text box header
		const headerText = "File Information"
		const boxWidth = len(headerText) + 2*headerPadding
		fmt.Fprintln(w, "┌"+strings.Repeat("─", boxWidth)+"┐")
		padding := (boxWidth - len(headerText)) / paddingDivisor
		fmt.Fprintf(w, "│%s%s%s│\n",
			strings.Repeat(" ", padding),
			headerText,
			strings.Repeat(" ", padding))
		fmt.Fprintln(w, "└"+strings.Repeat("─", boxWidth)+"┘")
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, d.Name)

	data := d.marshalled()

	// Organize data into categories
	basicInfo := []struct{ k, v string }{}
	contentStats := []struct{ k, v string }{}
	fileMeta := []struct{ k, v string }{}
	checksums := []struct{ k, v string }{}
	sauceData := []struct{ k, v string }{}
	comments := []struct{ k, v string }{}

	for _, x := range data {
		if !d.validate(x) || d.skip(x) {
			continue
		}

		switch x.k {
		case "slug", "filename", "filetype", "Unicode", "line break":
			basicInfo = append(basicInfo, x)
		case "characters", "words", "size", "lines", "width", ans:
			contentStats = append(contentStats, x)
		case "modified", "media mime type":
			fileMeta = append(fileMeta, x)
		case "SHA256 checksum", "CRC64 ECMA", "CRC32", "MD5":
			checksums = append(checksums, x)
		case "title", "author", "group", "date", "original size", "file type", "data type",
			"description", "character width", "number of lines", "interpretation":
			sauceData = append(sauceData, x)
		case cmmt:
			comments = append(comments, x)
		case "\u00A0": // noBreakSpace
			comments = append(comments, x)
		case zipComment:
			// Handle separately
		}
	}

	// Track which sections we've displayed
	sections := []struct {
		name    string
		items   []struct{ k, v string }
		display bool
		isLast  bool
	}{
		{"Basic Information", basicInfo, len(basicInfo) > 0, false},
		{"Content Statistics", contentStats, len(contentStats) > 0, false},
		{"File Metadata", fileMeta, len(fileMeta) > 0, false},
		{"Checksums & Integrity", checksums, len(checksums) > 0, false},
		{"SAUCE Metadata", sauceData, len(sauceData) > 0, true},
	}

	// Display sections with tree structure
	for _, section := range sections {
		if !section.display {
			continue
		}

		// Determine connector for section header
		sectionConnector := treeCorner
		if section.isLast {
			sectionConnector = treeLastCorner
		}

		fmt.Fprintf(w, "%s%s\n", treeStyle(sectionConnector), treeStyle(section.name))

		// Display items in this section
		for j, item := range section.items {
			itemConnector := treeVertical
			if section.isLast {
				itemConnector = treeSpace
			}

			itemPrefix := treeCorner
			if j == len(section.items)-1 {
				itemPrefix = treeLastCorner
			}

			fmt.Fprintf(w, "%s%s%s: %s\n",
				treeStyle(itemConnector),
				treeStyle(itemPrefix),
				keyStyle(item.k),
				valueStyle(item.v))
		}

		// Display SAUCE comments if this is the SAUCE section and we have comments
		if section.name == "SAUCE Metadata" && len(comments) > 0 {
			commentConnector := "    "
			if section.isLast {
				commentConnector = "    "
			}

			fmt.Fprintf(w, "%s    └── Comments\n", treeStyle(commentConnector))
			for _, comment := range comments {
				fmt.Fprintf(w, "%s        %s\n", treeStyle(commentConnector), valueStyle(comment.v))
			}
		}
	}
}

// marshalled returns the data structure used for print marshaling.
func (d *Detail) marshalled() []struct{ k, v string } {
	const (
		noBreakSpace     = "\u00A0"
		symbolForNewline = "\u2424"
		// baseFieldCount represents the number of base fields in the detail structure
		baseFieldCount = 30
	)
	p := message.NewPrinter(lang())
	// Preallocate slice with capacity for all fields plus potential SAUCE comments
	data := make([]struct {
		k, v string
	}, 0, baseFieldCount+len(d.Sauce.Comnt.Comment))
	data = append(data,
		struct{ k, v string }{k: "slug", v: d.Slug},
		struct{ k, v string }{k: "filename", v: d.Name},
		struct{ k, v string }{k: "filetype", v: d.Mime.Commt},
		struct{ k, v string }{k: "Unicode", v: d.Unicode},
		struct{ k, v string }{k: "line break", v: fsys.LineBreak(d.LineBreak.Decimal, true)},
		struct{ k, v string }{k: "characters", v: p.Sprint(d.Count.Chars)},
		struct{ k, v string }{k: ans, v: p.Sprint(d.Count.Controls)},
		struct{ k, v string }{k: "words", v: p.Sprint(d.Count.Words)},
		struct{ k, v string }{k: "size", v: d.Size.Decimal},
		struct{ k, v string }{k: "lines", v: p.Sprint(d.Lines)},
		struct{ k, v string }{k: "width", v: p.Sprint(d.Width)},
		struct{ k, v string }{k: "modified", v: humanize.DMY.Format(d.Modified.Time.UTC())},
		struct{ k, v string }{k: "media mime type", v: d.Mime.Type},
		struct{ k, v string }{k: "SHA256 checksum", v: d.Sums.SHA256},
		struct{ k, v string }{k: "CRC64 ECMA", v: d.Sums.CRC64},
		struct{ k, v string }{k: "CRC32", v: d.Sums.CRC32},
		struct{ k, v string }{k: "MD5", v: d.Sums.MD5},
		struct{ k, v string }{k: zipComment, v: d.ZipComment},
	)
	// sauce data
	data = append(data,
		struct{ k, v string }{k: "title", v: d.Sauce.Title},
		struct{ k, v string }{k: "author", v: d.Sauce.Author},
		struct{ k, v string }{k: "group", v: d.Sauce.Group},
		struct{ k, v string }{k: "date", v: sauceDate(d.Sauce.Date.Value)},
		struct{ k, v string }{k: "original size", v: d.Sauce.FileSize.Decimal},
		struct{ k, v string }{k: "file type", v: d.Sauce.File.Name},
		struct{ k, v string }{k: "data type", v: d.Sauce.Data.Name},
		struct{ k, v string }{k: "description", v: d.Sauce.Desc},
		struct{ k, v string }{k: d.Sauce.Info.Info1.Info, v: strconv.FormatUint(uint64(d.Sauce.Info.Info1.Value), 10)},
		struct{ k, v string }{k: d.Sauce.Info.Info2.Info, v: strconv.FormatUint(uint64(d.Sauce.Info.Info2.Value), 10)},
		struct{ k, v string }{k: d.Sauce.Info.Info3.Info, v: strconv.FormatUint(uint64(d.Sauce.Info.Info3.Value), 10)},
		struct{ k, v string }{k: "interpretation", v: d.Sauce.Info.Flags.String()},
	)
	// sauce comment
	for i, line := range d.Sauce.Comnt.Comment {
		comment := struct{ k, v string }{
			k: noBreakSpace, v: line,
		}
		if i == 0 {
			comment.k = cmmt
		}
		data = append(data, comment)
	}
	return data
}

func (d *Detail) mime(name string, data ...byte) {
	mm := mimemagic.MatchMagic(data)
	d.Mime.Media = mm.Media
	d.Mime.Sub = mm.Subtype
	d.Mime.Type = fmt.Sprintf("%s/%s", mm.Media, mm.Subtype)
	d.Mime.Commt = mm.Comment
	if d.Mime.Commt == "plain text document" {
		reader := bytes.NewReader(data)
		if s := bbs.Find(reader).Name(); s != "" {
			d.Mime.Commt += fmt.Sprintf(" with %s BBS color codes", s)
		}
	}
	if ValidText(d.Mime.Type) {
		var err error
		b := bytes.NewBuffer(data)
		if d.Count.Chars, err = fsys.Runes(b); err != nil {
			fmt.Fprintf(os.Stdout, "mine sniffer failure, %s\n", err)
		}
		return
	}
	if d.Mime.Type == zipType {
		r, e := zip.OpenReader(name)
		if e != nil {
			fmt.Fprintf(os.Stdout, "open zip file failure: %s\n", e)
		}
		if r == nil {
			return
		}
		defer r.Close()
		d.ZipComment = r.Comment
	}
}

// skip reports whether the key and value data should be skipped.
func (d *Detail) skip(x struct{ k, v string }) bool {
	if !d.LegacySums {
		switch x.k {
		case "CRC32", "CRC64 ECMA", "MD5":
			return true
		}
	}
	return false
}

// validate reports whether the key and value data validate.
func (d *Detail) validate(x struct{ k, v string }) bool {
	if !ValidText(d.Mime.Type) {
		switch x.k {
		case uc8, "line break", "characters", ans, "words", "lines", "width":
			return false
		}
	} else if x.k == ans {
		if d.Count.Controls == 0 {
			return false
		}
	}
	if x.k == "description" && x.v == "" {
		return false
	}
	if x.k == d.Sauce.Info.Info1.Info && d.Sauce.Info.Info1.Value == 0 {
		return false
	}
	if x.k == d.Sauce.Info.Info2.Info && d.Sauce.Info.Info2.Value == 0 {
		return false
	}
	if x.k == d.Sauce.Info.Info3.Info && d.Sauce.Info.Info3.Value == 0 {
		return false
	}
	if x.k == "interpretation" && x.v == "" {
		return false
	}
	return true
}
