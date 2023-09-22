package xhex

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

var ErrBase = fmt.Errorf("base must be 10 or 16")

// Hexadecimal prefix identifiers.
const (
	X    = "X"   // common
	Hash = "#"   // cascading style sheets colors
	Doll = "$"   // retro microcomputers
	U    = "U+"  // unicode
	Zero = "0X"  // linux and unix C syntax
	Esc  = "\\X" // escape
)

const (
	base10 = 10
	base16 = 16
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
// For decimal numbers use base 10.
// For hexadecimal numbers use base 16.
func Parse(base int, vals ...string) []int64 {
	n := make([]int64, len(vals))
	for i, val := range vals {
		x, err := strconv.ParseInt(val, base, 64)
		if err != nil {
			n[i] = -1
			continue
		}
		x = int64(uint64(x)) // remove sign
		n[i] = x
	}
	return n
}

// Dec writes the hexadecimal values of the provided decimal strings.
// If a string is not a hexadecimal number then the value is printed as "invalid".
func Result(w io.Writer, base int, vals ...string) error {
	if w == nil {
		w = io.Discard
	}
	if base != base10 && base != base16 {
		return fmt.Errorf("%w: %d", ErrBase, base)
	}
	const pad = "  "
	b := &strings.Builder{}
	nums := []int64{}
	switch base {
	case base10:
		nums = Parse(base10, vals...)
	case base16:
		nums = Parse(base16, TrimIndents(vals...)...)
	}
	for i, x := range nums {
		s := strings.ToUpper(vals[i])
		if x == -1 {
			fmt.Fprintf(b, "%s = invalid%s", s, pad)
			continue
		}
		switch base {
		case base10:
			fmt.Fprintf(b, "%s = %X%s", s, x, pad)
		case base16:
			fmt.Fprintf(b, "%s = %d%s", s, x, pad)
		}
	}
	fmt.Fprint(w, strings.TrimSpace(b.String()))
	fmt.Fprintln(w)
	return nil
}

// DecRaw writes the hexadecimal values of the provided decimal strings.
// Only the results are written and are separated by a space.
// If a string is not a hexadecimal number then the value is printed as "NaN".
func Raw(w io.Writer, base int, vals ...string) error {
	if w == nil {
		w = io.Discard
	}
	if base != base10 && base != base16 {
		return fmt.Errorf("%w: %d", ErrBase, base)
	}
	const pad = " "
	b := &strings.Builder{}
	nums := []int64{}
	switch base {
	case base10:
		nums = Parse(base10, vals...)
	case base16:
		nums = Parse(base16, TrimIndents(vals...)...)
	}
	for _, x := range nums {
		if x == -1 {
			fmt.Fprintf(b, "NaN%s", pad)
			continue
		}
		switch base {
		case base10:
			fmt.Fprintf(b, "%X%s", x, pad)
		case base16:
			fmt.Fprintf(b, "%d%s", x, pad)
		}
	}
	fmt.Fprint(w, strings.TrimSpace(b.String()))
	fmt.Fprintln(w)
	return nil
}
