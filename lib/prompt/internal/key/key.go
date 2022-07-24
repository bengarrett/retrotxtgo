package key

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/prompt/internal/read"
	"github.com/bengarrett/retrotxtgo/lib/str"
	"github.com/gookit/color"
)

const (
	CursorPreviousLine = "\033[F"
	RemoveChr          = "-"
)

type Keys []string

// gotoInput returns the ANSI cursor previous line control.
func gotoInput() string {
	if color.Enable {
		return CursorPreviousLine
	}
	return ""
}

// Numeric checks input string for a valid int value and returns a matching slice.
func (k Keys) Numeric(input string) string {
	if input == "" {
		return ""
	}
	i, err := strconv.Atoi(input)
	if err != nil {
		return ""
	}
	if i >= len(k) || i < 0 {
		return ""
	}
	sort.Strings(k)
	return k[i]
}

// Prompt parses the reader input for a valid key.
func (k Keys) Prompt(w io.Writer, r io.Reader, setup bool) string {
	prompts := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key := scanner.Text()
		if key == RemoveChr {
			return RemoveChr
		}
		if setup && key == "" {
			return ""
		}
		if n := k.Numeric(key); n != "" {
			return n
		}
		if !k.Validate(key) {
			if err := read.Check(w, prompts); err != nil {
				return key
			}
			continue
		}
		return key
	}
	return ""
}

// ShortPrompt parses the reader input for a valid key or alias of the key.
func (k Keys) ShortPrompt(w io.Writer, r io.Reader) string {
	prompts, scanner := 0, bufio.NewScanner(r)
	for scanner.Scan() {
		prompts++
		key := scanner.Text()
		switch key {
		case "", RemoveChr:
			return key
		}
		if long := k.ShortValidate(key); long != "" {
			return long
		}
		if !k.Validate(key) {
			if err := read.Check(w, prompts); err != nil {
				return ""
			}
			continue
		}
		return key
	}
	return ""
}

// ShortValidate validates the key exists in Keys.
// Both the first letter of the key and the full name of the key are accepted as valid.
// Whenever the key is valid the full key name will be returned otherwise an empty string
// signifies a false result.
func (k Keys) ShortValidate(key string) string {
	if key == "" {
		return ""
	}
	l, x := len(k), strings.ToLower(key)
	sort.Strings(k)
	a, b := make([]string, l), make([]string, l)
	for i, l := range k {
		a[i] = strings.ToLower(string(l[0]))
		b[i] = strings.ToLower(l)
	}
	// match the first letter
	sort.Strings(a)
	i := sort.SearchStrings(a, x)
	if i >= len(a) || a[i] != x {
		// match the whole word
		sort.Strings(b)
		j := sort.SearchStrings(b, x)
		if j >= len(b) || b[j] != x {
			return ""
		}
		return k[j]
	}
	return k[i]
}

// Validate that the key exists in the slice of Keys.
func (k Keys) Validate(key string) bool {
	if key == "" {
		return false
	}
	sort.Strings(k)
	i := sort.SearchStrings(k, key)
	if i >= len(k) || k[i] != key {
		fmt.Printf("%s%s %v\n", gotoInput(), str.Bool(false), key)
		return false
	}
	return true
}
