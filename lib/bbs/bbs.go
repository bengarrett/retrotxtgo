package bbs

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// Bulletin Board System color code format.
type BBS int

const (
	ANSI BBS = iota
	Celerity
	PCBoard
	Renegade
	Telegard
	Wildcat
	WWIVHash
	WWIVHeart
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

// String returns the BBS color format name and toggle characters.
func (bbs BBS) String() string {
	if bbs < ANSI || bbs > WWIVHeart {
		return ""
	}
	return [...]string{
		"ANSI ←[",
		"Celerity |",
		"PCBoard @",
		"Renegade |",
		"Telegard `",
		"Wildcat! @X",
		"WWIV |#",
		"WWIV ♥",
	}[bbs]
}

// Bytes returns the BBS color code toggle characters.
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

// HTML transforms the string containing BBS color codes into HTML <i> elements.
func (bbs BBS) HTML(s string) string {
	x := rmCLS(s)
	switch bbs {
	case ANSI:
		return s
	case Celerity:
		return parseCelerity(x)
	case PCBoard:
		return parsePCBoard(x)
	case Renegade:
		return parseRenegade(x)
	case Telegard:
		return parseTelegard(x)
	case Wildcat:
		return parseWildcat(x)
	case WWIVHash:
		return parseWHash(x)
	case WWIVHeart:
		return parseWHeart(x)
	default:
		return s
	}
}

// rmCLS removes common PCBoard BBS controls from the string.
func rmCLS(s string) string {
	r := regexp.MustCompile(`@(CLS|CLS |PAUSE)@`)
	return r.ReplaceAllString(s, "")
}

// Find the format any known BBS color codes in the reader.
func Find(r io.Reader) BBS {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b := scanner.Bytes()
		ts := bytes.TrimSpace(b)
		if ts == nil {
			return -1
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

// Find Celerity BBS format color codes.
// The format uses the vertical bar "|" followed by a case sensitive single alphabetic character.
// See the celerityCodes const for all the characters.
func FindCelerity(b []byte) BBS {
	for _, code := range []byte(celerityCodes) {
		if bytes.Contains(b, []byte{Celerity.Bytes()[0], code}) {
			return Celerity
		}
	}
	return -1
}

// Find PCBoard BBS format color codes.
// The format uses an "@X" prefix with a background and foreground, 4-bit hexadecimal color value.
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

// Find Renegade BBS format color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 0 and 23.
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

// Find Telegard BBS format color codes.
// The format uses the vertical bar "|" followed by a padded, numeric value between 0 and 23.
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

// Find Wildcat BBS format color codes.
// The format uses an a background and foreground, 4-bit hexadecimal color value enclosed by two at "@" characters.
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

// Find WWIV BBS hash format color codes.
// The format uses a vertical bar "|" with the hash "#" characters as a prefix with a numeric value between 0 and 9.
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

// Find WWIV BBS heart format color codes.
// The format uses the ETX character as a prefix with a numeric value between 0 and 9.
// The ETX (end-of-text) character in the MS-DOS Code Page 437 is substituted with a heart "♥" character.
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

// Validate that the bytes are valid BBS color codes.
// Only Celerity, PCBoard or Renegade types are checked.
func (bbs BBS) validate(b []byte) bool {
	if b == nil {
		return false
	}
	switch bbs {
	case Celerity:
		return validateC(b[0])
	case PCBoard:
		return validateP(b[0])
	case Renegade:
		const min = 2
		if len(b) < min {
			return false
		}
		return validateR([2]byte{b[0], b[1]})
	case ANSI:
	case Telegard:
	case Wildcat:
	case WWIVHash:
	case WWIVHeart:
	}
	return false
}

func validateC(b byte) bool {
	return bytes.Contains([]byte(celerityCodes), []byte{b})
}

func validateP(b byte) bool {
	if b == byte(' ') {
		return false
	}
	const baseHex, bitSize = 16, 64
	i, err := strconv.ParseInt(string(b), baseHex, bitSize)
	if err != nil {
		return false
	}
	if i < 0 || i > 16 {
		return false
	}
	return true
}

func validateR(b [2]byte) bool {
	const bgMin, fgMax = 0, 23
	s := string(b[0]) + string(b[1])
	i, err := strconv.Atoi(s)
	if err != nil {
		return false
	}
	if i < bgMin || i > fgMax {
		return false
	}
	return true
}

// parserBar parses the string for BBS color codes that use vertical bar prefixes to apply a HTML template.
func parserBar(s string) (*bytes.Buffer, error) {
	const idiomaticTpl = `<i class="P{{.Background}},P{{.Foreground}}">{{.Content}}</i>`
	buf := bytes.Buffer{}
	d := IntData{
		Foreground: 0,
		Background: 0,
	}

	subs := strings.Split(s, string(verticalBar))
	if len(subs) <= 1 {
		fmt.Fprint(&buf, s)
		return &buf, nil
	}

	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		if sub == "" {
			continue
		}
		if sub[0] == lf {
			continue
		}
		n, err := strconv.Atoi(sub[0:2])
		if err != nil {
			continue
		}
		if !Renegade.validate([]byte{sub[0], sub[1]}) {
			fmt.Fprint(&buf, string(verticalBar))
			continue
		}

		if n >= 0 && n <= 15 {
			d.Foreground = n
		} else if n >= 16 && n <= 23 {
			d.Background = n
		}
		d.Content = sub[2:]

		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

// parserCelerity parses the string for the unique Celerity BBS color codes to apply a HTML template.
func parserCelerity(s string) (*bytes.Buffer, error) {
	const idiomaticTpl, swap = `<i class="PB{{.Background}},PF{{.Foreground}}">{{.Content}}</i>`, "S"
	buf, background := bytes.Buffer{}, false
	d := StrData{
		Foreground: "w",
		Background: "k",
	}

	subs := strings.Split(s, string(verticalBar))
	if len(subs) <= 1 {
		fmt.Fprint(&buf, s)
		return &buf, nil
	}

	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	for _, sub := range subs {
		if sub == "" {
			continue
		}
		if sub[0] == lf {
			continue
		}
		if !Celerity.validate([]byte{sub[0]}) {
			fmt.Fprint(&buf, string(verticalBar))
			continue
		}
		if sub == swap {
			background = !background
			continue
		}
		if !background {
			d.Foreground = string(sub[0])
		} else {
			d.Background = string(sub[0])
		}
		d.Content = sub[1:]

		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

// parserPCBoard parses the string for the common PCBoard BBS color codes to apply a HTML template.
func parserPCBoard(s string) (*bytes.Buffer, error) {
	const idiomaticTpl = `<i class="PB{{.Background}},PF{{.Foreground}}">{{.Content}}</i>`
	buf := bytes.Buffer{}
	d, b := StrData{}, PCBoard.Bytes()

	codes := strings.Split(s, string(b))
	if len(codes) <= 1 {
		fmt.Fprint(&buf, s)
		return &buf, nil
	}

	tmpl, err := template.New("idomatic").Parse(idiomaticTpl)
	if err != nil {
		return nil, err
	}

	for _, code := range codes {
		if code == "" {
			continue
		}
		if code[0] == lf {
			continue
		}
		if !PCBoard.validate([]byte{code[0]}) || !PCBoard.validate([]byte{code[1]}) {
			fmt.Fprint(&buf, b)
			continue
		}

		d.Background = string(code[0])
		d.Foreground = string(code[1])
		d.Content = code[2:]

		if err := tmpl.Execute(&buf, d); err != nil {
			return nil, err
		}
	}
	return &buf, nil
}

func parseCelerity(s string) string {
	buf, _ := parserCelerity(s)
	return buf.String()
}

func parseRenegade(s string) string {
	buf, _ := parserBar(s)
	return buf.String()
}

func parsePCBoard(s string) string {
	buf, _ := parserPCBoard(s)
	return buf.String()
}

// parseTelegard parses the string for Telegard BBS color codes.
// Internally it replaces all color codes with PCBoard codes and parses those with parserPCBoard.
func parseTelegard(s string) string {
	r := regexp.MustCompile("`([0-9|A-F])([0-9|A-F])")
	x := r.ReplaceAllString(s, `@X$1$2`)
	buf, _ := parserPCBoard(x)
	return buf.String()
}

// parseWildcat parses the string for Wildcat BBS color codes.
// Internally it replaces all color codes with PCBoard codes and parses those with parserPCBoard.
func parseWildcat(s string) string {
	r := regexp.MustCompile(`@([0-9|A-F])([0-9|A-F])@`)
	x := r.ReplaceAllString(s, `@X$1$2`)
	buf, _ := parserPCBoard(x)
	return buf.String()
}

// parseWHash parses the string for WWIV hash color codes.
// Internally it replaces all color codes with vertical bar codes and parses those with parserBar.
func parseWHash(s string) string {
	r := regexp.MustCompile(`\|#(\d)`)
	x := r.ReplaceAllString(s, `|0$1`)
	buf, _ := parserBar(x)
	return buf.String()
}

// parseWHeart parses the string for WWIV heart color codes.
// Internally it replaces all color codes with vertical bar codes and parses those with parserBar.
func parseWHeart(s string) string {
	r := regexp.MustCompile(`\x03(\d)`)
	x := r.ReplaceAllString(s, `|0$1`)
	buf, _ := parserBar(x)
	return buf.String()
}
