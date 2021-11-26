package bbs

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrColorCodes = errors.New("no bbs color codes found in string")
	ErrANSI       = errors.New("ansi escape code found in string")
)

// Bulletin Board System color code format.
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

// IntData template data for integer based color codes.
type IntData struct {
	Background int
	Foreground int
	Content    string
}

// StrData template data for string based color codes.
type StrData struct {
	Background string
	Foreground string
	Content    string
}

const (
	celerityCodes      = "kbgcrmywdBGCRMYWS"
	pcbClear           = `@CLS@`
	verticalBar        = byte('|')
	lf            byte = 10
)

// Valid checks the bbs type is a known BBS.
func (bbs BBS) Valid() bool {
	switch bbs {
	case ANSI, Celerity, PCBoard, Renegade, Telegard, Wildcat, WWIVHash, WWIVHeart:
		return true
	default:
		return false
	}
}

// String returns the BBS color format name and toggle sequence.
func (bbs BBS) String() string {
	if !bbs.Valid() {
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
	}[bbs]
}

// Name returns the BBS color format name.
func (bbs BBS) Name() string {
	if !bbs.Valid() {
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
	}[bbs]
}

// Bytes returns the BBS color code toggle sequence as bytes.
func (bbs BBS) Bytes() []byte {
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
	switch bbs {
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

// HTML transforms a string containing BBS color codes into a
// collection of HTML <i> elements with matching CSS color classes.
func (bbs BBS) HTML(s string) (*bytes.Buffer, error) {
	empty := bytes.Buffer{}
	x := trimPrefix(s)
	switch bbs {
	case ANSI:
		return &empty, ErrANSI
	case Celerity:
		return ParseCelerity(x)
	case PCBoard:
		return ParsePCBoard(x)
	case Renegade:
		return ParseRenegade(x)
	case Telegard:
		return ParseTelegard(x)
	case Wildcat:
		return ParseWildcat(x)
	case WWIVHash:
		return ParseWHash(x)
	case WWIVHeart:
		return ParseWHeart(x)
	default:
		return &empty, ErrColorCodes
	}
}

// trimPrefix removes common PCBoard BBS controls from the string.
func trimPrefix(s string) string {
	r := regexp.MustCompile(`@(CLS|CLS |PAUSE)@`)
	return r.ReplaceAllString(s, "")
}

// Find the format of any known BBS color code sequence within the reader.
// If no sequences are found -1 is returned.
func Find(r io.Reader) BBS {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		ts := bytes.TrimSpace(b)
		if ts == nil {
			continue
		}
		const l = len(pcbClear)
		if len(ts) > l {
			if bytes.Equal(ts[0:l], []byte(pcbClear)) {
				b = ts[l:]
			}
		}
		switch {
		case bytes.Contains(b, ANSI.Bytes()):
			return ANSI
		case bytes.Contains(b, Celerity.Bytes()):
			if f := FindRenegade(b); f == Renegade {
				return Renegade
			}
			if f := FindCelerity(b); f == Celerity {
				return Celerity
			}
			return -1
		case bytes.Contains(b, PCBoard.Bytes()):
			return FindPCBoard(b)
		case bytes.Contains(b, Telegard.Bytes()):
			return FindTelegard(b)
		case bytes.Contains(b, Wildcat.Bytes()):
			return FindWildcat(b)
		case bytes.Contains(b, WWIVHash.Bytes()):
			return FindWWIVHash(b)
		case bytes.Contains(b, WWIVHeart.Bytes()):
			return FindWWIVHeart(b)
		}
	}
	return -1
}

// FindCelerity reports if the bytes contains Celerity BBS color codes.
// The format uses the vertical bar "|" followed by a case sensitive single alphabetic character.
// Returns Celerity when true or -1 if no codes are found.
func FindCelerity(b []byte) BBS {
	// celerityCodes contains all the character sequences for Celerity.
	for _, code := range []byte(celerityCodes) {
		if bytes.Contains(b, []byte{Celerity.Bytes()[0], code}) {
			return Celerity
		}
	}
	return -1
}

// FindPCBoard reports if the bytes contains PCBoard BBS color codes.
// The format uses an "@X" prefix with a background and foreground, 4-bit hexadecimal color value.
// Returns PCBoard when true or -1 if no codes are found.
func FindPCBoard(b []byte) BBS {
	const first, last = 0, 15
	const hexxed = "%X%X"
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf(hexxed, bg, fg))
			subslice = append(PCBoard.Bytes(), subslice...)
			if bytes.Contains(b, subslice) {
				return PCBoard
			}
		}
	}
	return -1
}

// FindRenegade reports if the bytes contains Renegade BBS color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 0 and 23.
// Returns Renegade when true or -1 if no codes are found.
func FindRenegade(b []byte) BBS {
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Renegade.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return Renegade
		}
	}
	return -1
}

// FindTelegard reports if the bytes contains Telegard BBS color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 0 and 23.
// Returns Telegard when true or -1 if no codes are found.
func FindTelegard(b []byte) BBS {
	const first, last = 0, 23
	const leadingZero = "%01d"
	for i := first; i <= last; i++ {
		subslice := []byte(fmt.Sprintf(leadingZero, i))
		subslice = append(Telegard.Bytes(), subslice...)
		if bytes.Contains(b, subslice) {
			return Telegard
		}
	}
	return -1
}

// FindWildcat reports if the bytes contains Wildcat! BBS color codes.
// The format uses an a background and foreground, 4-bit hexadecimal color value enclosed by two at "@" characters.
// Returns Wildcat when true or -1 if no codes are found.
func FindWildcat(b []byte) BBS {
	const first, last = 0, 15
	for bg := first; bg <= last; bg++ {
		for fg := first; fg <= last; fg++ {
			subslice := []byte(fmt.Sprintf("%s%X%X%s",
				Wildcat.Bytes(), bg, fg, Wildcat.Bytes()))
			if bytes.Contains(b, subslice) {
				return Wildcat
			}
		}
	}
	return -1
}

// FindWWIVHash reports if the bytes contains WWIV BBS hash color codes.
// The format uses a vertical bar "|" with the hash "#" characters as a prefix with a numeric value between 0 and 9.
// Returns WWIVHash when true or -1 if no codes are found.
func FindWWIVHash(b []byte) BBS {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHash.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return WWIVHash
		}
	}
	return -1
}

// FindWWIVHeart reports if the bytes contains WWIV BBS heart color codes.
// The format uses the ETX character as a prefix with a numeric value between 0 and 9.
// The ETX (end-of-text) character in the MS-DOS Code Page 437 is substituted with a heart "♥" character.
// Returns WWIVHeart when true or -1 if no codes are found.
func FindWWIVHeart(b []byte) BBS {
	const first, last = 0, 9
	for i := first; i <= last; i++ {
		subslice := append(WWIVHeart.Bytes(), []byte(strconv.Itoa(i))...)
		if bytes.Contains(b, subslice) {
			return WWIVHeart
		}
	}
	return -1
}

// SplitBars slices s into substrings separated by "|" vertical bar codes.
// The first two bytes of each substrings will contain a colour value.
// An empty slice is returned if there are no valid bar code values exist in s.
func SplitBars(s string) []string {
	const sep rune = 65535
	m := regexp.MustCompile(`\|(0[0-9]|1[1-9]|2[0-3])`)
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

// ParserBars parses the string for BBS color codes that use
// vertical bar prefixes to apply a HTML template.
func parserBars(s string) (*bytes.Buffer, error) {
	const idiomaticTpl = `<i class="P{{.Background}} P{{.Foreground}}">{{.Content}}</i>`
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	buf, d := bytes.Buffer{}, IntData{}
	bars := SplitBars(s)
	if len(bars) == 0 {
		fmt.Fprint(&buf, s)
		return &buf, nil
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
		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
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

// SplitCelerity slices s into substrings separated by "|" vertical bar codes.
// The first byte of each substrings will contain a colour value.
// An empty slice is returned if there are no valid bar code values exist in s.
func SplitCelerity(s string) []string {
	// The format uses the vertical bar "|" followed by a case sensitive single alphabetic character.
	const sep rune = 65535
	m := regexp.MustCompile(`\|(k|b|g|c|r|m|y|w|d|B|G|C|R|M|Y|W|S)`)
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

// ParserCelerity parses the string for the unique Celerity BBS color codes
// to apply a HTML template.
func parserCelerity(s string) (*bytes.Buffer, error) {
	const idiomaticTpl, swapCmd = `<i class="PB{{.Background}} PF{{.Foreground}}">{{.Content}}</i>`, "S"
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	buf, background := bytes.Buffer{}, false
	d := StrData{
		Foreground: "w",
		Background: "k",
	}

	bars := SplitCelerity(s)
	if len(bars) == 0 {
		fmt.Fprint(&buf, s)
		return &buf, nil
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
		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

// SplitPCBoard slices s into substrings separated by PCBoard @X codes.
// The first two bytes of each substrings will contain a background
// and foreground hex colour value. An empty slice is returned if
// there are no valid @X code values exist in s.
func SplitPCBoard(s string) []string {
	const sep rune = 65535
	m := regexp.MustCompile("(?i)@X([0-9A-F][0-9A-F])")
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
func parserPCBoard(s string) (*bytes.Buffer, error) {
	const idiomaticTpl = `<i class="PB{{.Background}} PF{{.Foreground}}">{{.Content}}</i>`
	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	buf, d := bytes.Buffer{}, StrData{}
	xcodes := SplitPCBoard(s)
	if len(xcodes) == 0 {
		fmt.Fprint(&buf, s)
		return &buf, nil
	}
	for _, color := range xcodes {
		d.Background = strings.ToUpper(string(color[0]))
		d.Foreground = strings.ToUpper(string(color[1]))
		d.Content = color[2:]
		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

func ParseCelerity(s string) (*bytes.Buffer, error) {
	return parserCelerity(s)
}

func ParseRenegade(s string) (*bytes.Buffer, error) {
	return parserBars(s)
}

func ParsePCBoard(s string) (*bytes.Buffer, error) {
	return parserPCBoard(s)
}

// ParseTelegard parses the string for Telegard BBS color codes.
// It swaps the Telegard color codes with PCBoard @X codes and
// parses those with parserPCBoard.
func ParseTelegard(s string) (*bytes.Buffer, error) {
	r := regexp.MustCompile("`([0-9|A-F])([0-9|A-F])")
	x := r.ReplaceAllString(s, `@X$1$2`)
	return parserPCBoard(x)
}

// ParseWildcat parses the string for Wildcat! BBS color codes.
// It swaps the Wildcat color codes with PCBoard @X codes and
// parses those with parserPCBoard.
func ParseWildcat(s string) (*bytes.Buffer, error) {
	r := regexp.MustCompile(`@([0-9|A-F])([0-9|A-F])@`)
	x := r.ReplaceAllString(s, `@X$1$2`)
	return parserPCBoard(x)
}

// ParseWHash parses the string for WWIV hash color codes.
// It swaps the WWIV color codes with vertical bars and
// parses those with ParserBars.
func ParseWHash(s string) (*bytes.Buffer, error) {
	r := regexp.MustCompile(`\|#(\d)`)
	x := r.ReplaceAllString(s, `|0$1`)
	return parserBars(x)
}

// ParseWHeart parses the string for WWIV heart color codes.
// It swaps the WWIV color codes with vertical bars and
// parses those with ParserBars.
func ParseWHeart(s string) (*bytes.Buffer, error) {
	r := regexp.MustCompile(`\x03(\d)`)
	x := r.ReplaceAllString(s, `|0$1`)
	return parserBars(x)
}
