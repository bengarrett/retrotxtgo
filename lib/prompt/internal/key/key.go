package key

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bengarrett/retrotxtgo/lib/str"
)

const NoChange = "no changes applied"

type Keys []string

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
		if key == "-" {
			return ""
		}
		if setup && key == "" {
			return ""
		}
		if n := k.Numeric(key); n != "" {
			return n
		}
		if !k.Validate(key) {
			check(w, prompts)
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
		case "", "-":
			return key
		}
		if long := k.ShortValidate(key); long != "" {
			return long
		}
		if !k.Validate(key) {
			check(w, prompts)
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
func (k Keys) Validate(key string) (ok bool) {
	if key == "" {
		return false
	}
	sort.Strings(k)
	i := sort.SearchStrings(k, key)
	if i >= len(k) || k[i] != key {
		fmt.Printf("%s %v\n", str.Bool(false), key)
		return false
	}
	return true
}

// Check used in scanner Scans to prompt a user for a valid text input.
func check(w io.Writer, prompts int) {
	const info, max = 2, 4
	switch {
	case prompts == info:
		if i, err := fmt.Fprint(w, "\r  Ctrl+C to keep the existing value\n"); err != nil {
			log.Fatalf("prompt check println at %dB: %s", i, err)
		}
	case prompts >= max:
		fmt.Fprintf(w, "\r  %s", NoChange)
		os.Exit(0) // TODO: return error
	}
}
