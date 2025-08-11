// Package hexa provides rudimental hexadecimal conversion functions.
package hexa

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// Base is the number base system for the conversion.
type Base int

const (
	Base10 Base = 10 // Base10 is the decimal number base system.
	Base16 Base = 16 // Base16 is the hexadecimal number base system.
)

// Hexadecimal prefix identifiers.
const (
	X    = "X"   // common
	Hash = "#"   // cascading style sheets colors
	Doll = "$"   // retro microcomputers
	U    = "U+"  // unicode
	Zero = "0X"  // linux and unix C syntax
	Esc  = "\\X" // escape
)

// TrimIdent trims the hexadecimal prefix identifiers and
// returns the hexadecimal value as an upper case string.
func TrimIdent(s string) string {
	s = strings.ToUpper(s)
	s = strings.TrimPrefix(s, X)
	s = strings.TrimPrefix(s, Hash)
	s = strings.TrimPrefix(s, Doll)
	s = strings.TrimPrefix(s, U)
	s = strings.TrimPrefix(s, Zero)
	s = strings.TrimPrefix(s, Esc)
	return s
}

// TrimNCR trims the numeric character reference prefix and suffix,
// and returns a decimal integer as a string.
// If the string is not in NCR syntax, then the original string is returned.
func TrimNCR(s string) string {
	const dec, suffix = "&#", ";"
	x := strings.ToUpper(s)
	if !strings.HasPrefix(x, dec) || !strings.HasSuffix(x, suffix) {
		return s
	}
	x = strings.TrimPrefix(x, dec)
	x = strings.TrimSuffix(x, suffix)
	if x == "" {
		return s
	}
	return x
}

// TrimIndents trims the hexadecimal prefix identifiers and
// returns the hexadecimal values as upper case strings.
func TrimIndents(vals ...string) []string {
	s := make([]string, len(vals))
	for i, val := range vals {
		s[i] = TrimIdent(val)
	}
	for i, val := range s {
		s[i] = TrimNCR(val)
	}
	return s
}

// Parse converts the provided strings to based unsigned ints.
// If a string is not a valid integer then the value is -1.
func Parse(b Base, vals ...string) []int64 {
	n := make([]int64, len(vals))
	for i, val := range vals {
		x, err := strconv.ParseInt(val, int(b), 64)
		if err != nil {
			n[i] = -1
			continue
		}
		x = int64(math.Abs(float64(x))) // remove sign
		n[i] = x
	}
	return n
}

// Parser writes the hexadecimal values of the provided decimal strings.
// Only the results are written and are separated by a space.
// If a string is not a hexadecimal number then the value is printed as "NaN".
func Parser(w io.Writer, b Base, vals ...string) error {
	if w == nil {
		w = io.Discard
	}
	const pad = " "
	sb := &strings.Builder{}
	nums := []int64{}
	switch b {
	case Base10:
		nums = Parse(b, vals...)
	case Base16:
		nums = Parse(b, TrimIndents(vals...)...)
	}
	for _, x := range nums {
		if x == -1 {
			fmt.Fprintf(sb, "NaN%s", pad)
			continue
		}
		switch b {
		case Base10:
			fmt.Fprintf(sb, "%X%s", x, pad)
		case Base16:
			fmt.Fprintf(sb, "%d%s", x, pad)
		}
	}
	fmt.Fprint(w, strings.TrimSpace(sb.String()))
	fmt.Fprintln(w)
	return nil
}

// Writer the hexadecimal values of the provided decimal strings.
// If a string is not a hexadecimal number then the value is printed as "invalid".
func Writer(w io.Writer, b Base, vals ...string) error {
	if w == nil {
		w = io.Discard
	}
	const pad = "  "
	sb := &strings.Builder{}
	nums := []int64{}
	switch b {
	case Base10:
		nums = Parse(b, vals...)
	case Base16:
		nums = Parse(b, TrimIndents(vals...)...)
	}
	for i, x := range nums {
		if i >= len(vals) {
			break
		}
		s := strings.ToUpper(vals[i])
		if x == -1 {
			fmt.Fprintf(sb, "%s = invalid%s", s, pad)
			continue
		}
		switch b {
		case Base10:
			fmt.Fprintf(sb, "%s = %X%s", s, x, pad)
		case Base16:
			fmt.Fprintf(sb, "%s = %d%s", s, x, pad)
		}
	}
	fmt.Fprint(w, strings.TrimSpace(sb.String()))
	fmt.Fprintln(w)
	return nil
}
