// internal header
// https://go.dev/blog/examples
// https://pkg.go.dev/github.com/fluhus/godoc-tricks#section-documentation
// http://wiki.synchro.net/custom:colors
//
// future funcs
// BBS.ReplaceAll || BBS.Remove || BBS.Discard
// Marshal(interface) => []byte
// Encode(dst, src []byte)
// EncodeToString(src []byte) string
//
// HTMLEscape(dst *bytes.Buffer, src []byte)
// func main() {
// var out bytes.Buffer
// json.HTMLEscape(&out, []byte(`{"Name":"<b>HTML content</b>"}`))
// out.WriteTo(os.Stdout)
// }
// dst.Write() // dst.WriteByte // dst.WriteString()
//
// func Valid(data []byte) bool
// Valid reports whether data is a valid JSON encoding.
//
// html.HTML type html.CSS type

// Package bbs is a blah.
// 1) what is this package and its purpose
// 2) history of bbs text vs ansi for colouring
// 3) similarities to html
package bbs

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrColorCodes = errors.New("no bbs color codes found")
	ErrANSI       = errors.New("ansi escape code found")

	//go:embed static/*
	static embed.FS
)

// Bulletin Board System color code format.
// Other than for Find, the ANSI type is not supported by this library.
type BBS int

const (
	ANSI      BBS = iota // ANSI escape sequences.
	Celerity             // Celerity BBS pipe codes.
	PCBoard              // PCBoard BBS @ codes.
	Renegade             // Renegade BBS pipe codes.
	Telegard             // Telegard BBS grave accent codes.
	Wildcat              // Wildcat! BBS @ codes.
	WWIVHash             // WWIV BBS # codes.
	WWIVHeart            // WWIV BBS ♥ codes.
)

// colorInt template data for integer based color codes.
type colorInt struct {
	Background int
	Foreground int
	Content    string
}

// colorStr template data for string based color codes.
type colorStr struct {
	Background string
	Foreground string
	Content    string
}

const (
	// ClearCmd is a PCBoard specific control to clear the screen that's occasionally found in ANSI text.
	ClearCmd string = "@CLS@"

	// CelerityMatch is a regular expression to match Celerity BBS color codes.
	CelerityMatch string = `\|(k|b|g|c|r|m|y|w|d|B|G|C|R|M|Y|W|S)`

	// PCBoardMatch is a case-insensitive, regular expression to match PCBoard BBS color codes.
	PCBoardMatch string = "(?i)@X([0-9A-F][0-9A-F])"

	// RenegadeMatch is a regular expression to match Renegade BBS color codes.
	RenegadeMatch string = `\|(0[0-9]|1[1-9]|2[0-3])`

	// TelegardMatch is a case-insensitive, regular expression to match Telegard BBS color codes.
	TelegardMatch string = "(?i)`([0-9|A-F])([0-9|A-F])"

	// WildcatMatch is a case-insensitive, regular expression to match Wildcat! BBS color codes.
	WildcatMatch string = `(?i)@([0-9|A-F])([0-9|A-F])@`

	// WWIVHashMatch is a regular expression to match WWIV BBS # color codes.
	WWIVHashMatch string = `\|#(\d)`

	// WWIVHeartMatch is a regular expression to match WWIV BBS ♥ color codes.
	WWIVHeartMatch string = `\x03(\d)`

	celerityCodes = "kbgcrmywdBGCRMYWS"
)

// FieldsBars slices a string into substrings separated by "|" vertical bar codes.
// The first two bytes of each substring will contain a colour value.
// Vertical bar codes are used by Renegade, WWIV hash and WWIV heart formats.
// An empty slice is returned when no valid bar code values exists.
func FieldsBars(src string) []string {
	const sep rune = 65535
	m := regexp.MustCompile(RenegadeMatch)
	repl := fmt.Sprintf("%s$1", string(sep))
	res := m.ReplaceAllString(src, repl)
	if !strings.ContainsRune(res, sep) {
		return nil
	}

	spl := strings.Split(res, string(sep))
	app := []string{}
	for _, val := range spl {
		if val != "" {
			app = append(app, val)
		}
	}
	return app
}

// ParserBars parses the string for BBS color codes that use
// vertical bar prefixes to apply a HTML template.
func parserBars(dst *bytes.Buffer, src string) error {
	const idiomaticTpl = `<i class="P{{.Background}} P{{.Foreground}}">{{.Content}}</i>`
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return err
	}

	d := colorInt{}
	bars := FieldsBars(src)
	if len(bars) == 0 {
		_, err := dst.WriteString(src)
		return err
	}

	for _, color := range bars {
		n, err := strconv.Atoi(color[0:2])
		if err != nil {
			continue
		}
		if barForeground(n) {
			d.Foreground = n
		}
		if barBackground(n) {
			d.Background = n
		}
		d.Content = color[2:]
		if err := tmpl.Execute(dst, d); err != nil {
			return err
		}
	}
	return nil
}

func barBackground(n int) bool {
	if n < 16 {
		return false
	}
	if n > 23 {
		return false
	}
	return true
}

func barForeground(n int) bool {
	if n < 0 {
		return false
	}
	if n > 15 {
		return false
	}
	return true
}

// FieldsCelerity slices a string into substrings separated by "|" vertical bar codes.
// The first byte of each substring will contain a Celerity colour value,
// that are comprised of a single, alphabetic character.
// An empty slice is returned when no valid Celerity code values exists.
func FieldsCelerity(src string) []string {
	// The format uses the vertical bar "|" followed by a case sensitive single alphabetic character.
	const sep rune = 65535
	m := regexp.MustCompile(CelerityMatch)
	repl := fmt.Sprintf("%s$1", string(sep))
	res := m.ReplaceAllString(src, repl)
	if !strings.ContainsRune(res, sep) {
		return []string{}
	}

	spl := strings.Split(res, string(sep))
	clean := []string{}
	for _, val := range spl {
		if val != "" {
			clean = append(clean, val)
		}
	}
	return clean
}

// ParserCelerity parses the string for the unique Celerity BBS color codes
// to apply a HTML template.
func parserCelerity(dst *bytes.Buffer, src string) error {
	const idiomaticTpl, swapCmd = `<i class="PB{{.Background}} PF{{.Foreground}}">{{.Content}}</i>`, "S"
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return err
	}

	//buf, background := bytes.Buffer{}, false
	background := false
	d := colorStr{
		Foreground: "w",
		Background: "k",
	}

	bars := FieldsCelerity(src)
	if len(bars) == 0 {
		_, err := dst.WriteString(src)
		return err
	}
	for _, color := range bars {
		if color == swapCmd {
			background = !background
			continue
		}
		if !background {
			d.Foreground = string(color[0])
		}
		if background {
			d.Background = string(color[0])
		}
		d.Content = color[1:]
		if err := tmpl.Execute(dst, d); err != nil {
			return err
		}
	}
	return nil
}

// FieldsPCBoard slices a string into substrings separated by PCBoard @X codes.
// The first two bytes of each substring will contain background
// and foreground hex colour values.
// An empty slice is returned when no valid @X code values exists.
func FieldsPCBoard(s string) []string {
	const sep rune = 65535
	m := regexp.MustCompile(PCBoardMatch)
	repl := fmt.Sprintf("%s$1", string(sep))
	res := m.ReplaceAllString(s, repl)
	if !strings.ContainsRune(res, sep) {
		return []string{}
	}

	spl := strings.Split(res, string(sep))
	clean := []string{}
	for _, val := range spl {
		if val != "" {
			clean = append(clean, val)
		}
	}
	return clean
}

// parserPCBoard parses the string for the common PCBoard BBS color codes
// to apply a HTML template.
func parserPCBoard(dst *bytes.Buffer, src string) error {
	const idiomaticTpl = `<i class="PB{{.Background}} PF{{.Foreground}}">{{.Content}}</i>`
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return err
	}

	d := colorStr{}
	xcodes := FieldsPCBoard(src)
	if len(xcodes) == 0 {
		_, err := dst.WriteString(src)
		return err
	}
	for _, color := range xcodes {
		d.Background = strings.ToUpper(string(color[0]))
		d.Foreground = strings.ToUpper(string(color[1]))
		d.Content = color[2:]
		if err := tmpl.Execute(dst, d); err != nil {
			return err
		}
	}
	return nil
}

// HTML writes to dst the HTML equivalent of BBS color codes with matching CSS color classes.
// The first found color code format is used for the remainder of the Reader.
func HTML(dst *bytes.Buffer, src io.Reader) error {
	var r1 bytes.Buffer

	r2 := io.TeeReader(src, &r1)

	find := Find(r2)
	b, err := io.ReadAll(&r1)
	if err != nil {
		return err
	}
	return find.HTML(dst, b)
}

// HTMLCelerity writes to dst the HTML equivalent of Celerity BBS color codes with
// matching CSS color classes.
func HTMLCelerity(dst *bytes.Buffer, src string) error {
	return parserCelerity(dst, src)
}

// HTMLRenegade writes to dst the HTML equivalent of Renegade BBS color codes with
// matching CSS color classes.
func HTMLRenegade(dst *bytes.Buffer, src string) error {
	return parserBars(dst, src)
}

// HTMLPCBoard writes to dst the HTML equivalent of PCBoard BBS color codes with
// matching CSS color classes.
func HTMLPCBoard(dst *bytes.Buffer, src string) error {
	return parserPCBoard(dst, src)
}

// HTMLTelegard writes to dst the HTML equivalent of Telegard BBS color codes with
// matching CSS color classes.
func HTMLTelegard(dst *bytes.Buffer, src string) error {
	r := regexp.MustCompile(TelegardMatch)
	x := r.ReplaceAllString(src, `@X$1$2`)
	return parserPCBoard(dst, x)
}

// HTMLWildcat writes to dst the HTML equivalent of Wildcat! BBS color codes with
// matching CSS color classes.
func HTMLWildcat(dst *bytes.Buffer, src string) error {
	r := regexp.MustCompile(WildcatMatch)
	x := r.ReplaceAllString(src, `@X$1$2`)
	return parserPCBoard(dst, x)
}

// HTMLWildcat writes to dst the HTML equivalent of WWIV BBS # color codes with
// matching CSS color classes.
func HTMLWHash(dst *bytes.Buffer, src string) error {
	r := regexp.MustCompile(WWIVHashMatch)
	x := r.ReplaceAllString(src, `|0$1`)
	return parserBars(dst, x)
}

// HTMLWildcat writes to dst the HTML equivalent of WWIV BBS ♥ color codes with
// matching CSS color classes.
func HTMLWHeart(dst *bytes.Buffer, src string) error {
	r := regexp.MustCompile(WWIVHeartMatch)
	x := r.ReplaceAllString(src, `|0$1`)
	return parserBars(dst, x)
}

// IsCelerity reports if the bytes contains Celerity BBS color codes.
// The format uses the vertical bar "|" followed by a case sensitive single alphabetic character.
func IsCelerity(b []byte) bool {
	// celerityCodes contains all the character sequences for Celerity.
	for _, code := range []byte(celerityCodes) {
		if bytes.Contains(b, []byte{Celerity.Bytes()[0], code}) {
			return true
		}
	}
	return false
}

// IsPCBoard reports if the bytes contains PCBoard BBS color codes.
// The format uses an "@X" prefix with a background and foreground, 4-bit hexadecimal color value.
func IsPCBoard(b []byte) bool {
	const first, last = 0, 15
	const hexxed = "%X%X"
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf(hexxed, bg, fg))
			subslice = append(PCBoard.Bytes(), subslice...)
			if bytes.Contains(b, subslice) {
				return true
			}
		}
	}
	return false
}

// IsRenegade reports if the bytes contains Renegade BBS color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 00 and 23.
func IsRenegade(b []byte) bool {
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Renegade.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return true
		}
	}
	return false
}

// IsTelegard reports if the bytes contains Telegard BBS color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 00 and 23.
func IsTelegard(b []byte) bool {
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Telegard.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return true
		}
	}
	return false
}

// IsWildcat reports if the bytes contains Wildcat! BBS color codes.
// The format uses an a background and foreground,
// 4-bit hexadecimal color value enclosed by two at "@" characters.
func IsWildcat(b []byte) bool {
	const first, last = 0, 15
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf("%s%X%X%s",
				Wildcat.Bytes(), bg, fg, Wildcat.Bytes()))
			if bytes.Contains(b, subslice) {
				return true
			}
		}
	}
	return false
}

// IsWHash reports if the bytes contains WWIV BBS # (hash or pound) color codes.
// The format uses a vertical bar "|" with the hash "#" characters
// as a prefix with a numeric value between 0 and 9.
func IsWHash(b []byte) bool {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHash.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return true
		}
	}
	return false
}

// IsWHeart reports if the bytes contains WWIV BBS ♥ (heart) color codes.
// The format uses the ETX character as a prefix with a numeric value between 0 and 9.
// In the standard MS-DOS, USA codepage (CP-437), the ETX (end-of-text)
// character is substituted with a heart character.
func IsWHeart(b []byte) bool {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHeart.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return true
		}
	}
	return false
}

// Fields splits the io.Reader around the first instance of one or more consecutive BBS color codes.
// An error is returned if no color codes are found or if ANSI control sequences are first found.
func Fields(src io.Reader) ([]string, BBS, error) {
	var r1 bytes.Buffer
	r2 := io.TeeReader(src, &r1)
	f := Find(r2)
	if !f.Valid() {
		return nil, -1, ErrColorCodes
	}
	b, err := io.ReadAll(&r1)
	if err != nil {
		return nil, -1, err
	}
	switch f {
	case ANSI:
		return nil, -1, ErrANSI
	case Celerity:
		return FieldsCelerity(string(b)), f, nil
	case PCBoard, Telegard, Wildcat:
		return FieldsPCBoard(string(b)), f, nil
	case Renegade, WWIVHash, WWIVHeart:
		return FieldsBars(string(b)), f, nil
	}
	return nil, -1, ErrColorCodes
}

// Find the format of any known BBS color code sequence within the reader.
// If no sequences are found -1 is returned.
func Find(src io.Reader) BBS {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		b := scanner.Bytes()
		ts := bytes.TrimSpace(b)
		if ts == nil {
			continue
		}
		const l = len(ClearCmd)
		if len(ts) > l {
			if bytes.Equal(ts[0:l], []byte(ClearCmd)) {
				b = ts[l:]
			}
		}
		switch {
		case bytes.Contains(b, ANSI.Bytes()):
			return ANSI
		case bytes.Contains(b, Celerity.Bytes()):
			if IsRenegade(b) {
				return Renegade
			}
			if IsCelerity(b) {
				return Celerity
			}
			return -1
		case IsPCBoard(b):
			return PCBoard
		case IsTelegard(b):
			return Telegard
		case IsWildcat(b):
			return Wildcat
		case IsWHash(b):
			return WWIVHash
		case IsWHeart(b):
			return WWIVHeart
		}
	}
	return -1
}

// Bytes returns the BBS color toggle sequence as bytes.
func (b BBS) Bytes() []byte {
	const (
		etx               byte = 3  // CP437 ♥
		esc               byte = 27 // CP437 ←
		hash                   = byte('#')
		atSign                 = byte('@')
		grave                  = byte('`')
		leftSquareBracket      = byte('[')
		verticalBar            = byte('|')
		upperX                 = byte('X')
	)
	switch b {
	case ANSI:
		return []byte{esc, leftSquareBracket}
	case Celerity, Renegade:
		return []byte{verticalBar}
	case PCBoard:
		return []byte{atSign, upperX}
	case Telegard:
		return []byte{grave}
	case Wildcat:
		return []byte{atSign}
	case WWIVHash:
		return []byte{verticalBar, hash}
	case WWIVHeart:
		return []byte{etx}
	default:
		return nil
	}
}

// CSS writes to dst the Cascading Style Sheets classes needed by the HTML.
// The CSS relies on cascading variables.
// See https://developer.mozilla.org/en-US/docs/Web/CSS/Using_CSS_custom_properties for details.
func (b BBS) CSS(dst *bytes.Buffer) error {
	r, err := static.ReadFile("static/css/text_pcboard.css")
	if err != nil {
		return err
	}
	dst.Write(r)
	return nil
}

// HTML writes to dst the HTML equivalent of BBS color codes with matching CSS color classes.
func (b BBS) HTML(dst *bytes.Buffer, src []byte) error {
	x := trimPrefix(string(src))
	switch b {
	case ANSI:
		return ErrANSI
	case Celerity:
		return HTMLCelerity(dst, x)
	case PCBoard:
		return HTMLPCBoard(dst, x)
	case Renegade:
		return HTMLRenegade(dst, x)
	case Telegard:
		return HTMLTelegard(dst, x)
	case Wildcat:
		return HTMLWildcat(dst, x)
	case WWIVHash:
		return HTMLWHash(dst, x)
	case WWIVHeart:
		return HTMLWHeart(dst, x)
	default:
		return ErrColorCodes
	}
}

// trimPrefix removes common PCBoard BBS controls from the string.
func trimPrefix(s string) string {
	r := regexp.MustCompile(`@(CLS|CLS |PAUSE)@`)
	return r.ReplaceAllString(s, "")
}

// Name returns the name of the BBS color format.
func (b BBS) Name() string {
	if !b.Valid() {
		return ""
	}
	return [...]string{
		"ANSI",
		"Celerity",
		"PCBoard",
		"Renegade",
		"Telegard",
		"Wildcat!",
		"WWIV #",
		"WWIV ♥",
	}[b]
}

// Remove the BBS color codes from src and write it to dst.
func (b BBS) Remove(dst *bytes.Buffer, src []byte) error {
	switch b {
	case ANSI:
		return ErrANSI
	case Celerity:
		return remove(dst, src, CelerityMatch)
	case PCBoard:
		return remove(dst, src, PCBoardMatch)
	case Renegade:
		return remove(dst, src, RenegadeMatch)
	case Telegard:
		return remove(dst, src, TelegardMatch)
	case Wildcat:
		return remove(dst, src, WildcatMatch)
	case WWIVHash:
		return remove(dst, src, WWIVHashMatch)
	case WWIVHeart:
		return remove(dst, src, WWIVHeartMatch)
	default:
		return ErrColorCodes
	}
}

func remove(dst *bytes.Buffer, src []byte, expr string) error {
	m := regexp.MustCompile(expr)
	res := m.ReplaceAll(src, []byte(""))
	_, err := dst.Write(res)
	return err
}

// String returns the BBS color format name and toggle sequence.
func (b BBS) String() string {
	if !b.Valid() {
		return ""
	}
	return [...]string{
		"ANSI ←[",
		"Celerity |",
		"PCBoard @X",
		"Renegade |",
		"Telegard `",
		"Wildcat! @@",
		"WWIV |#",
		"WWIV ♥",
	}[b]
}

// Valid reports whether the BBS type is valid.
func (b BBS) Valid() bool {
	switch b {
	case ANSI,
		Celerity,
		PCBoard,
		Renegade,
		Telegard,
		Wildcat,
		WWIVHash,
		WWIVHeart:
		return true
	default:
		return false
	}
}
